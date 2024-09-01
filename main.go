package main

import (
	_ "image/png"
	"strconv"

	"github.com/EngoEngine/glm"
	"github.com/EngoEngine/math"
	"github.com/Zerou02/closedGL/closedGL"
	turingfontparser "github.com/Zerou02/closedGL/turing_font_parser"
	"github.com/go-gl/glfw/v3.2/glfw"
)

func main() {
	StartTTF()
	//StartClosedGL()
	//StartTuwuing()
}

func findRightmostPoint(points []glm.Vec2) glm.Vec2 {
	var maxX = points[0][0]
	var k = 0
	for i, x := range points {
		if x[0] > maxX {
			maxX = x[0]
			k = i
		}
	}
	return points[k]
}

func extractPolyFromContour(contour []turingfontparser.GlyfPoints) []glm.Vec2 {
	var newPoints = []glm.Vec2{}
	for _, x := range contour {
		if x.OnCurve && !closedGL.Contains(&newPoints, x.Pos) {
			newPoints = append(newPoints, x.Pos)
		}
	}
	return newPoints
}

func glyfToPolygon(points []turingfontparser.GlyfPoints, mesh *closedGL.PixelMesh) []glm.Vec2 {

	var contours = [][]turingfontparser.GlyfPoints{}
	var single = []turingfontparser.GlyfPoints{}
	for _, x := range points {
		var a = x
		a.Pos = closedGL.CartesianToSS(a.Pos, 800)
		single = append(single, a)

		if x.EndPoint {
			contours = append(contours, single)
			single = []turingfontparser.GlyfPoints{}
		}

	}
	var outer = extractPolyFromContour(contours[0])
	var inner = extractPolyFromContour(contours[1])

	var rightMost = findRightmostPoint(inner)
	closedGL.PrintlnVec2(rightMost)
	var baseLine = closedGL.CalcLinearEquation(glm.Vec2{0, rightMost[1]}, rightMost)
	for i := 0; i < len(outer)-1; i++ {
		var p1 = outer[i]
		var p2 = outer[i+1]
		var line = closedGL.CalcLinearEquation(p1, p2)
		var cp = closedGL.CalcCrossingPoint(baseLine, line)
		var highY = math.Max(p1[1], p2[1])
		var lowY = math.Min(p1[1], p2[1])
		var dy float32 = 13
		if highY >= rightMost[1] && lowY <= rightMost[1] {
			println("AA")
			if p1[0] > rightMost[0] || p2[0] > rightMost[0] {
				closedGL.FindIdx(outer, p1)
				var firstNew = glm.Vec2{cp[0], cp[1] - dy}
				var secondNew = glm.Vec2{cp[0], cp[1] + dy}
				outer = closedGL.InsertAt(outer, firstNew, i+1)
				var innerIdx = closedGL.FindIdx(inner, rightMost)
				var newInner = []glm.Vec2{}
				for j := innerIdx; j < len(inner); j++ {
					newInner = append(newInner, inner[j])
				}
				for j := 0; j < innerIdx; j++ {
					newInner = append(newInner, inner[j])
				}
				var newInnerP = glm.Vec2{rightMost[0], rightMost[1] + dy}
				newInner = append(newInner, newInnerP)
				newInner = append(newInner, secondNew)
				outer = closedGL.InsertArrAt(outer, newInner, i+2)
				break
			}
		}
	}
	return outer

}

// ausgelegt für SS-Polys,einzelner Kantenzug, nichts Dopplung
func getInterSectionPointsInPolygon(p glm.Vec2, points []glm.Vec2) []glm.Vec2 {
	var ret = []glm.Vec2{}
	for i := 0; i < len(points)-1; i++ {
		var p1 = points[i]
		var p2 = points[i+1]
		var highY = math.Max(p1[1], p2[1])
		var lowY = math.Min(p1[1], p2[1])

		if highY > p[1] && lowY < p[1] {
			var line = closedGL.CalcLinearEquation(p1, p2)
			var x = (p[1] - line[1]) / line[0]
			if x < p[0] {
				ret = append(ret, glm.Vec2{x, p[1]})
			}
		}
	}
	return ret
}

func pointInPolygon(p glm.Vec2, points []glm.Vec2) bool {
	return len(getInterSectionPointsInPolygon(p, points))%2 == 1
}

func printlnPoints(points []glm.Vec2) {
	for _, x := range points {
		closedGL.PrintlnVec2(x)
		println("--")
	}
}

func findInteriorAngle(p glm.Vec2, poly []glm.Vec2, mesh *closedGL.PixelMesh) (float32, bool) {
	var idx = closedGL.FindIdx(poly, p)
	if idx == -1 {
		panic("dmae")
	}
	var previous, next glm.Vec2
	if idx == 0 {
		previous = poly[len(poly)-1]
	} else {
		previous = poly[idx-1]
	}
	if idx == len(poly)-1 {
		next = poly[0]
	} else {
		next = poly[idx+1]
	}
	//so dass die Schwänze übereinander liegen
	p = closedGL.SsToCartesian(p, 800)
	previous = closedGL.SsToCartesian(previous, 800)
	next = closedGL.SsToCartesian(next, 800)

	var firstVec = previous.Sub(&p)
	var secondVec = next.Sub(&p)
	firstVec.Normalize()
	secondVec.Normalize()

	var isInner bool
	var angle = closedGL.AngleTo(firstVec, secondVec)

	var eq = closedGL.CalcLinearEquation(p, next)
	var dx = next[0] - p[0]
	dx *= 0.1
	var test = closedGL.EvalLinEq(eq, p[0]+dx)
	test = closedGL.RotateAroundPoint(angle, test, p)
	isInner = !pointInPolygon(test, poly)

	return angle, isInner
}

func triangulatePolygon(points []glm.Vec2, mesh *closedGL.PixelMesh, breakAt int, print bool) [][3]glm.Vec2 {
	var base = points
	_ = base
	var i = 1
	var ret = [][3]glm.Vec2{}

	for len(points) > 3 {
		var length = len(points)
		var p0 = points[(i-1)%length]
		var tip = points[i%length]
		var p1 = points[(i+1)%length]

		var inTri = false

		var dx = p1[0] - p0[0]
		var steps = 10
		var stepX = math.Floor(dx / float32(steps))
		var eq = closedGL.CalcLinearEquation(p0, p1)
		var inPoly = false

		for j := 1; j < steps-2; j++ {
			var newX = p0[0] + float32(j)*stepX
			var newP = closedGL.EvalLinEq(eq, newX)
			if pointInPolygon(newP, points) || p0[0] == p1[0] || newP[0] == p0[0] {
				if len(ret)+1 == breakAt {

					if print {
						println("step I", i, j, "newP")
						closedGL.PrintlnVec2(newP)
						println("tip")
						closedGL.PrintlnVec2(tip)
						println("p0")
						closedGL.PrintlnVec2(p0)
						println("p1")
						closedGL.PrintlnVec2(p1)
					}
					mesh.AddPixel(newP, glm.Vec4{0, 1, 0, 1})
				}

				inPoly = true

				break
			}
		}
		if !inPoly {
			i++
			continue
		}
		for j, x := range points {
			if j != i && j != (i-1)%length && j != (i+1)%length {
				if closedGL.PointInTriangle(x, p0, tip, p1) {
					inTri = true
				}
			}
		}
		if !inTri {
			ret = append(ret, [3]glm.Vec2{p0, tip, p1})
			points = closedGL.Remove(points, tip)
			if len(ret) == 1 {

				mesh.AddPixel(tip, glm.Vec4{0.5, 0.5, 1, 1})

				for _, x := range points {
					mesh.AddPixel(x, glm.Vec4{1, 0, 0, 1})
				}
				mesh.AddPixel(p0, glm.Vec4{0, 1, 1, 1})
				mesh.AddPixel(p1, glm.Vec4{0, 1, 1, 1})
				break
			}
			i = 0
		}
		i++
	}

	ret = append(ret, [3]glm.Vec2{points[0], points[1], points[2]})

	return ret
}

func pointsToSS(points []glm.Vec2) []glm.Vec2 {
	var ret = []glm.Vec2{}
	for _, x := range points {
		ret = append(ret, closedGL.CartesianToSS(x, 800))
	}
	return ret
}

func extractPoints(glyfPoints []turingfontparser.GlyfPoints) []glm.Vec2 {
	var points = []glm.Vec2{}
	for _, x := range glyfPoints {
		points = append(points, x.Pos)
	}
	return points
}

func glyfToPoints(pointsGlyf []turingfontparser.GlyfPoints, poly []glm.Vec2, lineMesh *closedGL.LineMesh, triMesh *closedGL.TriangleMesh) {
	var points = extractPoints(pointsGlyf)
	points = pointsToSS(points)
	for i := 0; i < len(points); i += 3 {
		lineMesh.AddQuadraticBezier(points[i], points[i+1], points[i+2], glm.Vec4{1, 1, 1, 1})
		var sign float32 = 1
		if pointInPolygon(points[i+2], poly) {
			sign = -1
		}
		_ = sign
		triMesh.AddTri([3]glm.Vec2{points[i], points[i+2], points[i+1]}, [3]glm.Vec2{{0, 0}, {1.0 / 2.0, 0}, {1, 1}}, sign)
	}
}

func StartTTF() {
	var opengl = closedGL.InitClosedGL(800, 800, "comic")
	opengl.LimitFPS(false)
	var p = turingfontparser.NewTuringFont("./assets/font/jetbrains_mono_medium.ttf", &opengl)
	var points = []turingfontparser.GlyfPoints{}
	var offset float32 = 0
	for _, x := range "R" {
		_ = x
		var newPoints = p.ParseGlyf(uint32(x), 1).AddXOffset(offset)
		var biggestX float32 = offset
		for i, y := range newPoints {
			var offset = glm.Vec2{0, 0}
			newPoints[i].Pos = newPoints[i].Pos.Add(&offset)
			if y.Pos[0] > biggestX {
				biggestX = y.Pos[0]
			}
		}
		offset += (biggestX - offset)
		points = append(points, newPoints...)
	}
	var lines = opengl.CreateLineMesh()

	var leaveMesh = opengl.CreatePixelMesh()
	var debugMesh = opengl.CreatePixelMesh()

	var poly = glyfToPolygon(points, &debugMesh)

	var tri = opengl.CreateTriMesh()
	var polyMesh = opengl.CreatePixelMesh()
	glyfToPoints(points, poly, &lines, &tri)

	lines.Copy()
	tri.Copy()
	polyMesh.SetPixelSize(10)
	for _, x := range poly {
		polyMesh.AddPixel(x, glm.Vec4{1, 1, 1, 1})
	}
	polyMesh.Copy()

	debugMesh.SetPixelSize(10)
	leaveMesh.SetPixelSize(10)

	var breakAt = 1
	var print = false

	testLines()
	for !opengl.WindowShouldClose() {
		if opengl.IsKeyPressed(glfw.KeyL) {
			closedGL.PrintlnVec2(opengl.GetMousePos())
		}
		if opengl.IsKeyPressed(glfw.KeyO) {
			breakAt++
		}
		if opengl.IsKeyPressed(glfw.KeyP) {
			breakAt--
		}
		if opengl.IsKeyPressed(glfw.KeyU) {
			print = !print
		}

		closedGL.SetWireFrameMode(!opengl.IsKeyDown(glfw.KeyF))
		opengl.BeginDrawing()
		opengl.ClearBG(glm.Vec4{0, 0, 0, 0})
		tri.Draw()
		polyMesh.Draw()
		lines.Draw()
		opengl.DrawFPS(600, 0, 1)
		opengl.EndDrawing()
	}
	opengl.Free()
}

/* func StartTuwuing() {
	var openGL = closedGL.InitClosedGL(1400, 800, "demo")

	openGL.LimitFPS(true)
	var complete = tuwuing_complete.NewTuwuingComplete(&openGL)

	for !openGL.WindowShouldClose() {
		complete.Process()

		openGL.ClearBG(glm.Vec4{0, 0, 0, 0})
		openGL.BeginDrawing()
		complete.Draw()

		openGL.DrawFPS(500, 0, 1)

		openGL.EndDrawing()
	}
	openGL.Free()
} */

func StartClosedGL() {
	var openGL = closedGL.InitClosedGL(800, 600, "demo")
	openGL.LimitFPS(false)
	glfw.GetCurrentContext().SetInputMode(glfw.CursorMode, glfw.CursorDisabled)

	openGL.Window.SetScrollCallback(openGL.Camera.ScrollCb)
	openGL.Window.SetCursorPosCallback(openGL.Camera.MouseCallback)

	var val = true

	//	var chunks = []ynnebcraft.Chunk{}
	//	var mesher = ynnebcraft.NewGreedyMesher()

	/* 	for i := 0; i < 20; i++ {
		for j := 0; j < 10; j++ {
			chunks = append(chunks, ynnebcraft.NewChunk(glm.Vec3{float32(i) * 32, 0, float32(j) * 32}, glm.Vec3{32, 32, 32}, &openGL, &mesher))
		}
	} */
	openGL.Logger.Enabled = true

	var pxTest = openGL.CreatePixelMesh()
	var lines = openGL.CreateLineMesh()

	for i := 0; i < 800; i++ {
		for j := 0; j < 600; j++ {
			pxTest.AddPixel(glm.Vec2{float32(i), float32(j)}, glm.Vec4{float32(j) / 600, float32(i) / 800, 1, 1})
		}
	}
	pxTest.Copy()
	lines.AddPath([]glm.Vec2{{0, 0}, {200, 200}}, []glm.Vec4{{1, 1, 1, 1}, glm.Vec4{1, 1, 1, 1}})
	lines.AddQuadraticBezier(glm.Vec2{0, 0}, glm.Vec2{200, 200}, glm.Vec2{0, 200}, glm.Vec4{0, 0.54, 0.57, 1})
	lines.Copy()

	var cam = openGL.NewCam2D()
	openGL.LimitFPS(false)
	for !openGL.WindowShouldClose() {
		cam.Process(float32(openGL.FPSCounter.Delta))
		pxTest.View = cam.ViewMat
		if openGL.IsKeyPressed(glfw.Key(glfw.KeyF)) {
			val = !val
			closedGL.SetWireFrameMode(val)
		}
		openGL.BeginDrawing()
		openGL.ClearBG(glm.Vec4{0, 0, 0, 0})
		openGL.DrawFPS(500, 0, 1)

		//openGL.DrawSprite(glm.Vec4{0, 0, 20, 20}, "./assets/sprites/fence_small.png", glm.Vec4{0, 0, 1, 1}, glm.Vec2{32, 1024}, 1)

		openGL.Text.DrawText(500, 50, "x:"+strconv.FormatFloat(float64(openGL.Camera.CameraFront[0]), 'f', 2, 64), 1)
		openGL.Text.DrawText(500, 75, "y:"+strconv.FormatFloat(float64(openGL.Camera.CameraFront[1]), 'f', 2, 64), 1)
		openGL.Text.DrawText(500, 100, "z:"+strconv.FormatFloat(float64(openGL.Camera.CameraFront[2]), 'f', 2, 64), 1)

		openGL.Text.DrawText(600, 50, "x:"+strconv.FormatFloat(float64(openGL.Camera.CameraPos[0]), 'f', 2, 64), 1)
		openGL.Text.DrawText(600, 75, "y:"+strconv.FormatFloat(float64(openGL.Camera.CameraPos[1]), 'f', 2, 64), 1)
		openGL.Text.DrawText(600, 100, "z:"+strconv.FormatFloat(float64(openGL.Camera.CameraPos[2]), 'f', 2, 64), 1)

		/* for i := 0; i < len(chunks); i++ {
			chunks[i].Draw()
		} */
		//	lines.Draw()
		pxTest.Draw()
		//openGL.DrawRect(glm.Vec4{0, 0, 100, 100}, glm.Vec4{0, 1, 0, 1}, 1)
		openGL.EndDrawing()
	}
	openGL.Free()
}

func assert(val bool) {
	if !val {
		panic("assert failed")
	}
}

func testLines() {
	var p1 = glm.Vec2{0, 0}
	var p2 = glm.Vec2{100, 100}
	var p3 = glm.Vec2{100, 0}
	var p4 = glm.Vec2{0, 100}
	//normal-normal
	var l1 = closedGL.CalculateLine(p1, p2)
	var l2 = closedGL.CalculateLine(p3, p4)
	var poin, _ = l1.GetIntersection(l2)
	var point2, _ = l2.GetIntersection(l1)
	assert(point2.Equal(&poin))
	println("expect 50,50")
	closedGL.PrintlnVec2(poin)

	//vert-horiz
	var l3 = closedGL.CalculateLine(p2, p3)
	var l4 = closedGL.CalculateLine(p1, p3)
	var cp2, _ = l3.GetIntersection(l4)
	var cp22, _ = l4.GetIntersection(l3)
	assert(cp2.Equal(&cp22))
	println("expect 0,100")
	closedGL.PrintlnVec2(cp2)

	//vert-normal
	var l5 = closedGL.CalculateLine(p1, p4)
	var l6 = closedGL.CalculateLine(p4, p3)
	var cp3, _ = l5.GetIntersection(l6)
	var cp32, _ = l6.GetIntersection(l5)

	assert(cp3.Equal(&cp32))
	println("expect 0,100")
	closedGL.PrintlnVec2(cp3)

	//vert-vert
	var l7 = closedGL.CalculateLine(p2, p3)
	var l8 = closedGL.CalculateLine(p1, p4)
	var cp4, succ = l7.GetIntersection(l8)
	var cp42, _ = l8.GetIntersection(l7)
	assert(cp4.Equal(&cp42))
	println("expect 0,0,false")
	print(succ, ",")
	closedGL.PrintlnVec2(cp4)

	//horiz,horiz
	var l9 = closedGL.CalculateLine(p2, p4)
	var l10 = closedGL.CalculateLine(p1, p3)
	var cp5, succ2 = l9.GetIntersection(l10)
	var cp52, _ = l10.GetIntersection(l9)
	assert(cp5.Equal(&cp52))
	println("expect 0,0,false")
	print(succ2, ",")
	closedGL.PrintlnVec2(cp5)

	//horiz,normal
	var l11 = closedGL.CalculateLine(p1, p2)
	var l12 = closedGL.CalculateLine(p1, p3)
	var cp6, succ3 = l11.GetIntersection(l12)
	var cp62, _ = l12.GetIntersection(l11)
	assert(cp62.Equal(&cp6))
	println("expect 0,0,true")
	print(succ3, ",")
	closedGL.PrintlnVec2(cp6)
}
