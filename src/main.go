package main

import (
	closedGL "closed_gl/src/closedGL"
	closed_gl "closed_gl/src/closedGL"
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
	var circle = openGL.Factory.NewCircle(glm.Vec4{1, 0, 0, 1}, glm.Vec4{0, 1, 0, 1}, 20, [2]float32{800, 230}, 6)
	var tri = openGL.Factory.NewTriangle([][2]float32{{430, 130}, {100, 400}, {550, 400}}, glm.Vec4{0, 0.5, 1, 1})

	glfw.SwapInterval(0)

	openGL.Camera.CameraPos = glm.Vec3{0, 0, 0}

	var anim = newAnimation(&circle.BorderThickness, 1, circle.BorderThickness, 5, true)
	var anim2 = newAnimation(&circle.Radius, 1, circle.Radius, 12, true)
	var intersection = closedGL.IntersectionOfLines(tri.Points[0], tri.Points[1], [2]float32{0, circle.Centre[1]}, [2]float32{circle.Centre[0], circle.Centre[1] + 1})

	var si = closed_gl.CartesianToSS(intersection)
	closedGL.PrintlnFloat(si[0])
	closedGL.PrintlnFloat(si[1])

	//circle.Centre[0] = x
	var anim3 = newAnimation(&circle.CentreColour[1], 0.5, 0, 1, true)
	var anim4 = newAnimation(&circle.Centre[0], 10, circle.Centre[0], 0, false)
	var anim5 = newAnimation(&circle.Centre[1], 1, circle.Centre[1], 600, true)
	anim5.Stopped = true

	for !openGL.Window.ShouldClose() {
		var delta = openGL.FPSCounter.Delta
		_ = delta
		anim.process(float32(delta))
		anim2.process(float32(delta))
		anim3.process(float32(delta))
		anim4.process(float32(delta))
		anim5.process(float32(delta))
		if circle.Centre[0] <= si[0] {
			anim4.Stopped = true
			anim5.Stopped = false

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
		circle.Draw()
		openGL.Process()
	}
	openGL.Free()
}
