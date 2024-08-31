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
		if highY > rightMost[1] && lowY < rightMost[1] {
			if p1[0] > rightMost[0] || p2[0] > rightMost[0] {
				println("AAA")
				closedGL.FindIdx(outer, p1)
				println("i", i)
				closedGL.PrintlnVec2(cp)
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

func findInteriorAngle(p glm.Vec2, poly []glm.Vec2, mesh *closedGL.PixelMesh) float32 {
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
	/* 	p = closedGL.SsToCartesian(p, 800)
	   	previous = closedGL.SsToCartesian(previous, 800)
	   	next = closedGL.SsToCartesian(next, 800) */

	//so dass die Schwänze übereinander liegen
	var firstVec = previous.Sub(&p)
	var secondVec = next.Sub(&p)
	firstVec.Normalize()
	secondVec.Normalize()
	println("first")
	closedGL.PrintlnVec2(firstVec)
	println("secundo")
	closedGL.PrintlnVec2(secondVec)

	var angle = closedGL.AngleTo(firstVec, secondVec)
	var testEq = closedGL.CalcLinearEquation(p, next)
	var dx = next[0] - p[0]
	dx *= 0.1
	var newPoint = closedGL.EvalLinEq(testEq, p[0]+dx)

	var otherDirAngle = 2*math.Pi - angle
	closedGL.PrintlnVec2(newPoint)
	var halfAngle = angle * 0.5
	var halfOtherAngle = 2*math.Pi - halfAngle
	_ = halfOtherAngle

	var rest = 2*math.Pi - otherDirAngle
	var first = closedGL.RotateAroundPoint(otherDirAngle*0.5, newPoint, p)
	var second = closedGL.RotateAroundPoint(otherDirAngle+0.5*rest, newPoint, p)
	var retAngle float32
	if pointInPolygon(first, poly) {
		println("cw in poly")
		retAngle = otherDirAngle
	}
	if pointInPolygon(second, poly) {
		println("ccw in poly")
		retAngle = angle
	}
	mesh.AddPixel(newPoint, glm.Vec4{1, 0.5, 0.5, 1})
	mesh.AddPixel(first, glm.Vec4{1, 0.5, 0.5, 1})
	mesh.AddPixel(second, glm.Vec4{1, 0.75, 0.5, 1})
	mesh.AddPixel(next, glm.Vec4{1, 0.75, 0.5, 1})

	println("angle", closedGL.RadToDeg(angle))
	println("otherangle", closedGL.RadToDeg(otherDirAngle))

	println("retangle", closedGL.RadToDeg(retAngle))
	closedGL.PrintlnVec2(newPoint)
	return retAngle
}

func triangulatePolygon(points []glm.Vec2, mesh *closedGL.PixelMesh) [][3]glm.Vec2 {
	var base = points
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

		for i := 1; i < steps-2; i++ {
			var newX = p0[0] + float32(i)*stepX
			var newP = closedGL.EvalLinEq(eq, newX)
			if pointInPolygon(newP, points) || p0[0] == p1[0] {
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
				println("inPoly", inPoly)

				println("p0")
				closedGL.PrintlnVec2(p0)
				println("tip")
				closedGL.PrintlnVec2(tip)
				println("p1")
				closedGL.PrintlnVec2(p1)

				findInteriorAngle(tip, base, mesh)
				/* 		for _, x := range base {
					var angle = findInteriorAngle(x, base, mesh)
					var c = glm.Vec4{1, 0, 0, 1}
					if angle < math.Pi {
						c = glm.Vec4{0, 1, 0, 1}
					}
					mesh.AddPixel(x, c)
				} */

				break
			}
			i = 0
		}
		i++
	}

	//	ret = append(ret, [3]glm.Vec2{points[0], points[1], points[2]})

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
	for _, x := range "a" {
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

	var test = triangulatePolygon(poly, &leaveMesh)
	leaveMesh.Copy()
	for _, x := range test {
		_ = x
		tri.AddTri(x, [3]glm.Vec2{{1, 1}, {1, 1}, {1, 1}}, -1)
	}

	lines.Copy()
	tri.Copy()
	polyMesh.SetPixelSize(10)
	for _, x := range poly {
		polyMesh.AddPixel(x, glm.Vec4{1, 1, 1, 1})
	}
	polyMesh.Copy()

	debugMesh.SetPixelSize(10)
	leaveMesh.SetPixelSize(10)

	for !opengl.WindowShouldClose() {
		if opengl.IsKeyPressed(glfw.KeyL) {
			closedGL.PrintlnVec2(opengl.GetMousePos())
		}
		closedGL.SetWireFrameMode(!opengl.IsKeyDown(glfw.KeyF))
		debugMesh.Clear()
		var mouse = opengl.GetMousePos()
		//mouse = closedGL.SsToCartesian(mouse, 800)
		debugMesh.AddPixel(mouse, glm.Vec4{0, 0.5, 0.5, 1})
		var inter = getInterSectionPointsInPolygon(mouse, poly)
		for _, x := range inter {
			debugMesh.AddPixel(x, glm.Vec4{1, 1, 0, 1})
		}
		debugMesh.AddPixel(glm.Vec2{448, 495}, glm.Vec4{1, 0, 1, 1})
		debugMesh.Copy()

		opengl.BeginDrawing()
		opengl.ClearBG(glm.Vec4{0, 0, 0, 0})
		lines.Draw()
		tri.Draw()
		polyMesh.Draw()
		debugMesh.Draw()
		leaveMesh.Draw()
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
