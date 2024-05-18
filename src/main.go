package main

import (
	closedGL "closed_gl/src/closedGL"
	_ "image/png"
	"runtime"

	"github.com/EngoEngine/glm"
	"github.com/go-gl/glfw/v3.2/glfw"
)

func main() {
	runtime.LockOSThread()

	var openGL = closedGL.InitClosedGL(800, 600)

	openGL.Window.SetScrollCallback(openGL.Camera.ScrollCb)
	openGL.Window.SetCursorPosCallback(openGL.Camera.MouseCallback)
	//openGL.Window.SetInputMode(glfw.CursorMode, glfw.CursorDisabled)
	var isWireframeMode = false

	var tri = openGL.Factory.NewTriangle([][2]float32{{130, 130}, {100, 500}, {550, 400}}, glm.Vec4{0, 0.5, 1, 1})
	glfw.SwapInterval(0)

	openGL.Camera.CameraPos = glm.Vec3{0, 0, 0}

	var delta float64

	var incenter = closedGL.CalculateIncenter(&tri)
	var circle = openGL.Factory.NewCircle(glm.Vec4{0, 1, 0, 1}, glm.Vec4{1, 1, 0, 1}, 10, incenter, 10)
	var fsms = createIncenterFSM(&tri, &openGL, &delta)

	for !openGL.Window.ShouldClose() {
		delta = openGL.FPSCounter.Delta
		for _, x := range fsms {
			x.Process()
		}

		closedGL.ClearBG()
		if openGL.KeyBoardManager.IsPressed(glfw.KeyF) {
			isWireframeMode = !isWireframeMode
			closedGL.SetWireFrameMode(isWireframeMode)
		}
		/* for _, x := range chunks {
		x.Draw()
		} */
		openGL.DrawFPS(0, 0)
		tri.Draw()
		for _, x := range fsms {
			x.DrawFn()
		}

		circle.Draw()

		openGL.Process()
	}
	openGL.Free()
}

func createIncenterFSM(tri *closedGL.Triangle, ctx *closedGL.ClosedGLContext, delta *float64) []Fsm {
	var fsms = []Fsm{}
	for i := 0; i < 3; i++ {
		fsms = append(fsms, createSimpleFSM(tri, i, delta, ctx))
	}
	return fsms
}

func createSimpleFSM(tri *closedGL.Triangle, startPointIdx int, delta *float64, ctx *closedGL.ClosedGLContext) Fsm {

	var fsm = newFsm()
	var line = ctx.Factory.NewLine(closedGL.Vec2{0, 0}, closedGL.Vec2{0, 0}, glm.Vec3{1, 0, 0}, glm.Vec3{1, 0, 0})
	var journey = newJourneyCircle(closedGL.Vec2{0, 0}, 10, *ctx)
	fsm.DrawFn = func() { line.Draw(); journey.draw() }
	journey.Position = tri.Points[startPointIdx]
	var allIdx = [3]int{0, 1, 2}
	allIdx[startPointIdx] = -1
	var otherIdx = []int{}
	for _, x := range allIdx {
		if x != -1 {
			otherIdx = append(otherIdx, x)
		}
	}
	var currPoint = tri.Points[startPointIdx]
	var mp = closedGL.MiddlePoint(tri.Points[otherIdx[0]], tri.Points[otherIdx[1]])
	var dx = mp[0] - currPoint[0]
	var dy = mp[1] - currPoint[1]
	var vec = glm.Vec2{dx, dy}

	vec.Normalize()
	var state = newState(func() {
		journey.Position[0] += vec[0] * float32(*delta) * 70
		journey.Position[1] += vec[1] * float32(*delta) * 70
		var dist = closedGL.Dist(journey.Position, currPoint)
		journey.Radius = dist
		line.Points[0] = tri.Points[startPointIdx][0]
		line.Points[1] = tri.Points[startPointIdx][1]
		line.Points[1*5+0] = journey.CentreCircle.Centre[0]
		line.Points[1*5+1] = journey.CentreCircle.Centre[1]
	})
	state.addTransition(func() bool {
		return closedGL.CircleLineIntersection(tri.Points[otherIdx[0]], tri.Points[otherIdx[1]], 1, journey.Position)
	}, 1, func() {
		println("a")
		journey.MakeInvis()
	})
	fsm.States = append(fsm.States, state)
	return fsm
}

func createFSM(tri *closedGL.Triangle, animManager *AnimationManager, delta *float64, ctx *closedGL.ClosedGLContext) Fsm {
	var line = ctx.Factory.NewLine(closedGL.Vec2{0, 0}, closedGL.Vec2{0, 0}, glm.Vec3{1, 0, 0}, glm.Vec3{1, 0, 0})
	var journey = newJourneyCircle(closedGL.Vec2{0, 0}, 10, *ctx)

	var intersection = closedGL.IntersectionOfLines(tri.Points[0], tri.Points[1], journey.Position, closedGL.Vec2{0, journey.Position[1]})
	animManager.addAnim("infly", &journey.Position[0], 2, journey.Position[0], intersection[0], true)
	var state0 = newState(func() {})
	state0.addTransition(func() bool {
		return journey.Position[0] <= intersection[0]
	}, 1, func() {
		var iSDown = closedGL.IntersectionOfLines(tri.Points[2], closedGL.Vec2{tri.Points[1][0], tri.Points[1][1] + 1}, journey.Position, closedGL.Vec2{journey.Position[0] + 1, 600})
		animManager.cancelAnim("infly")
		animManager.addAnim("down", &journey.Position[1], 2, journey.Position[1], iSDown[1], true)
	})
	var state1 = newState(func() { journey.Radius += float32(*delta) * 20 })
	state1.addTransition(func() bool { return journey.Position[1]+journey.Radius >= tri.Points[1][1] }, 2, func() {
		animManager.cancelAnim("down")
		journey.Position[1]--
	})
	var state2 = newState(func() { journey.Position[0] += float32(*delta) * 80 })
	state2.addTransition(func() bool {
		return closedGL.CircleLineIntersection(tri.Points[0], tri.Points[2], journey.Radius, journey.Position)
	}, 3, func() {
		journey.Position[0]--
	})
	var state3 = newState(func() {
		var lowerRight = tri.Points[2]
		var upper = tri.Points[0]
		var dx = lowerRight[0] - upper[0]
		var dy = lowerRight[1] - upper[1]

		journey.Radius += float32(*delta) * 10
		journey.Position[0] -= float32(*delta) * 10

		var vec = glm.Vec2{dx, dy}
		vec.Normalize()
		journey.Position[0] -= vec[0] * float32(*delta) * 70
		journey.Position[1] -= vec[1] * float32(*delta) * 70
	})
	state3.addTransition(func() bool {
		return closedGL.CircleLineIntersection(tri.Points[0], tri.Points[1], journey.Radius, journey.Position)
	}, 4, func() {})
	var state4 = newState(func() {
		var lowerLeft = tri.Points[1]
		var upper = tri.Points[0]
		var dx = lowerLeft[0] - upper[0]
		var dy = lowerLeft[1] - upper[1]

		if journey.Radius+journey.Position[1] >= lowerLeft[1] {
			var newRad = lowerLeft[1] - journey.Position[1]
			journey.Position[0] -= 1.6 * (journey.Radius - newRad)
			journey.Radius = newRad
		}
		var vec = glm.Vec2{dx, dy}
		vec.Normalize()
		journey.Position[0] += vec[0] * float32(*delta) * 70
		journey.Position[1] += vec[1] * float32(*delta) * 70
	})
	state4.addTransition(func() bool {
		var lowerLeft = tri.Points[1]
		return journey.Position[0] <= lowerLeft[0]
	}, 5, func() {})

	var state5 = newState(func() {
		var lowerRight = tri.Points[2]
		var lowerLeft = tri.Points[1]
		var upper = tri.Points[0]
		var mp = closedGL.MiddlePoint(lowerRight, upper)
		var dx = mp[0] - lowerLeft[0]
		var dy = mp[1] - lowerLeft[1]

		var vec = glm.Vec2{dx, dy}
		vec.Normalize()
		journey.Position[0] += vec[0] * float32(*delta) * 70
		journey.Position[1] += vec[1] * float32(*delta) * 70
		journey.Radius = lowerRight[1] - journey.Position[1]
		line.Points[0] = lowerLeft[0]
		line.Points[1] = lowerLeft[1]

		line.Points[1*5+0] = journey.CentreCircle.Centre[0]
		line.Points[1*5+1] = journey.CentreCircle.Centre[1]
	})
	state5.addTransition(func() bool {
		var upper = tri.Points[0]
		return journey.Position[1]-journey.Radius < upper[1]
	}, 6, func() {
		journey.MakeInvis()
	})
	var fsm = newFsm()
	fsm.DrawFn = func() {
		line.Draw()
		journey.draw()
	}
	fsm.States = append(fsm.States, state0)
	fsm.States = append(fsm.States, state1)
	fsm.States = append(fsm.States, state2)
	fsm.States = append(fsm.States, state3)
	fsm.States = append(fsm.States, state4)
	fsm.States = append(fsm.States, state5)
	return fsm
}
