package main

import (
	_ "image/png"

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
	//duplicate Points to Kantenzug
	var duplicateTrain = []glm.Vec2{}
	for i := 0; i < len(newPoints)-1; i++ {
		duplicateTrain = append(duplicateTrain, newPoints[i], newPoints[i+1])
	}
	duplicateTrain = append(duplicateTrain, duplicateTrain[len(duplicateTrain)-1])
	duplicateTrain = append(duplicateTrain, duplicateTrain[0])
	return duplicateTrain
}

func isInsideOfOtherPolyogn(base, polyToTest []glm.Vec2) bool {
	var minX, maxX = base[0][0], base[0][0]
	var minY, maxY = base[0][1], base[0][1]
	for _, x := range base {
		minX = math.Min(x[0], minX)
		maxX = math.Max(x[0], maxX)
		minY = math.Min(x[1], minY)
		maxY = math.Max(x[1], maxY)
	}
	var isInside = false
	for _, x := range polyToTest {
		if isPointInPolygon(x, base) {
			isInside = true
			break
		}
	}
	return isInside
}

func glyfToPolygon(points []turingfontparser.GlyfPoints, mesh *closedGL.PixelMesh) []glm.Vec2 {

	var contours = [][]turingfontparser.GlyfPoints{}
	var single = []turingfontparser.GlyfPoints{}
	for _, x := range points {
		var a = x
		a.Pos = closedGL.CartesianToSS(a.Pos, 800)
		a.Pos[0] = math.Floor(a.Pos[0])
		a.Pos[1] = math.Floor(a.Pos[1])
		single = append(single, a)
		if x.EndPoint {
			contours = append(contours, single)
			single = []turingfontparser.GlyfPoints{}
		}

	}

	var ret = []glm.Vec2{}
	for _, x := range contours {
		ret = append(ret, extractPolyFromContour(x)...)
	}
	var foundOffPoints = []glm.Vec2{}
	for _, x := range contours {
		var found = 0
		for i, y := range x {
			if !y.OnCurve && isPointInPolygon(y.Pos, ret) && !closedGL.Contains(&foundOffPoints, y.Pos) {
				var lastOnPoint = x[i-2].Pos
				var nextOnPoint = x[i-1].Pos
				if lastOnPoint.Equal(&y.Pos) {
					continue
				}
				foundOffPoints = append(foundOffPoints, y.Pos)
				var polyIdx = closedGL.FindIdx(ret, lastOnPoint)
				var isAtEnd = polyIdx == len(ret)-3

				var startKantenzug = i == 2
				var newKantenzug = []glm.Vec2{lastOnPoint, y.Pos, y.Pos, nextOnPoint}
				if isAtEnd {
					newKantenzug = []glm.Vec2{lastOnPoint, y.Pos, y.Pos, ret[len(ret)-1]}
				}

				if startKantenzug {
					ret = closedGL.RemoveAt(ret, polyIdx)
					ret = closedGL.RemoveAt(ret, polyIdx)
					ret = closedGL.InsertArrAt(ret, newKantenzug, polyIdx)
				} else {
					ret = closedGL.RemoveAt(ret, polyIdx+1)
					ret = closedGL.RemoveAt(ret, polyIdx+1)
					ret = closedGL.InsertArrAt(ret, newKantenzug, polyIdx+1)
				}

				found++
			}
		}
	}
	return ret

}

func isVertexOfPolygon(p glm.Vec2, points []glm.Vec2) bool {
	for _, x := range points {
		if x.Equal(&p) {
			return true
		}
	}
	return false
}

func isOnVertexLineOfPolygon(p glm.Vec2, points []glm.Vec2) bool {
	for _, x := range points {
		if x[1] == p[1] {
			return true
		}
	}
	return false
}

// ausgelegt für SS-Polys,mit Kantenzügen
func getInterSectionPointsInPolygon2(p glm.Vec2, points []glm.Vec2) ([]glm.Vec2, bool) {
	if isOnVertexLineOfPolygon(p, points) {
		return getInterSectionPointsInPolygon2(glm.Vec2{p[0], p[1] - 1}, points)
	}

	var ray = closedGL.CalculateLine(glm.Vec2{0, p[1]}, p)
	_ = ray
	var retMap = map[glm.Vec2]glm.Vec2{}
	for i := 0; i < len(points); i += 2 {
		var p1 = points[i]
		var p2 = points[i+1]
		var line = closedGL.CalculateLine(p1, p2)
		var cp, _ = line.GetIntersection(ray)
		if line.IsOnLine(cp) && cp[0] <= p[0] {
			retMap[cp] = cp
		}
	}
	var ret = []glm.Vec2{}
	for _, x := range retMap {
		ret = append(ret, x)
	}
	return ret, false
}

func isPointInPolygon(p glm.Vec2, points []glm.Vec2) bool {
	var intersections, isOnOutline = getInterSectionPointsInPolygon2(p, points)
	return len(intersections)%2 == 1 || isOnOutline
}

func printlnPoints(points []glm.Vec2) {
	for _, x := range points {
		closedGL.PrintlnVec2(x)
		println("--")
	}
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
		var off = points[i+2]
		off[0] = math.Floor(off[0])
		off[1] = math.Floor(off[1])
		if closedGL.Contains(&poly, off) {
			sign = -1
		}
		_ = sign
		triMesh.AddTri([3]glm.Vec2{points[i], points[i+2], points[i+1]}, [3]glm.Vec2{{0, 0}, {1.0 / 2.0, 0}, {1, 1}}, sign)
	}
}

func pixelFillPoly(poly []glm.Vec2, mesh *closedGL.PixelMesh) {
	for i := 0; i < 800; i++ {
		for j := 0; j < 800; j++ {
			var p = glm.Vec2{float32(j), float32(i)}
			if isPointInPolygon(p, poly) {
				mesh.AddPixel(p, glm.Vec4{1, 1, 1, 1})
			}
		}
	}
}

func printlnPoly(poly []glm.Vec2) {
	for i := 0; i < len(poly); i++ {
		closedGL.PrintlnVec2(poly[i])
		if i%2 == 1 {
			println("----", i)
		}
	}
}

func findInteriorAnglesOfPoly(poly []glm.Vec2, dbMesh *closedGL.PixelMesh) {
	for i := 1; i < len(poly); i += 2 {
		var tip = poly[i]
		var p0 = poly[i-1]
		var p1 = poly[1]
		//nur am Ende anders
		if i != len(poly)-1 {
			p1 = poly[i+2]
		}
		tip = closedGL.SsToCartesian(tip, 800)
		p0 = closedGL.SsToCartesian(p0, 800)
		p1 = closedGL.SsToCartesian(p1, 800)
		var line = closedGL.CalculateLine(tip, p0)
		var line2 = closedGL.CalculateLine(tip, p1)

		var vec1 = p0.Sub(&tip)
		var vec2 = p1.Sub(&tip)

		var positiveY = glm.Vec2{0, 1}
		var nVec1 = vec1.Normalized()
		var nVec2 = vec2.Normalized()

		var rotAngle = closedGL.AngleTo(positiveY, nVec1)
		var rot1 = closedGL.Rotate(rotAngle, nVec1)
		var rot2 = closedGL.Rotate(rotAngle, nVec2)

		if math.Abs(rot1[0]) > 0.1 {
			closedGL.PrintlnVec2(rot1)
			panic("hm")
		}
		if rot1[1] > 0 {
			rot1[1] *= -1
			rot2[0] *= -1
			rot2[1] *= -1
		}
		var rightTurn = rot2[0] > glm.Epsilon

		var angle = closedGL.AngleTo(vec1, vec2)
		var p = line.LerpPointOnLine(0.1)
		var p2 = line2.LerpPointOnLine(0.1)
		var rotated = closedGL.RotateAroundPoint(angle*0.5, p2, tip)
		_, _ = p, p2
		var angleIsInsideAngle = false
		_ = angleIsInsideAngle
		rotated = closedGL.CartesianToSS(rotated, 800)
		var isInside = isPointInPolygon(rotated, poly)
		//	var p2 = line2.LerpPointOnLine(0.1)
		//		dbMesh.AddPixel(p2, glm.Vec4{0, 1, 0.5, 1})
		//		dbMesh.AddPixel(rotated, glm.Vec4{0.5, 0.71, 1, 1})

		dbMesh.AddPixel(closedGL.CartesianToSS(tip, 800), glm.Vec4{1, 0, 0, 1})
		dbMesh.AddPixel(rotated, glm.Vec4{0, 0, 1, 1})
		_, _ = angle, isInside

		if rightTurn {
			dbMesh.AddPixel(closedGL.CartesianToSS(tip, 800), glm.Vec4{0, 1, 0, 1})

		}
	}

}

func StartTTF() {
	var opengl = closedGL.InitClosedGL(800, 800, "comic")

	opengl.LimitFPS(true)
	var p = turingfontparser.NewTuringFont("./assets/font/jetbrains_mono_medium.ttf", &opengl)
	var points = []turingfontparser.GlyfPoints{}
	var offset float32 = 0
	//input
	for _, x := range "I" {
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
	_ = poly

	var tri = opengl.CreateTriMesh()
	var polyMesh = opengl.CreatePixelMesh()
	glyfToPoints(points, poly, &lines, &tri)

	lines.Copy()
	tri.Copy()
	polyMesh.SetPixelSize(1)

	debugMesh.SetPixelSize(5)
	leaveMesh.SetPixelSize(1)
	for _, x := range poly {
		debugMesh.AddPixel(x, glm.Vec4{1, 0, 0, 1})
	}
	findInteriorAnglesOfPoly(poly, &debugMesh)

	var breakAt = 1
	var print = false
	var amountPixel = 1

	tri.Copy()
	debugMesh.Copy()
	var pixelFillMesh = opengl.CreatePixelMesh()
	pixelFillPoly(poly, &pixelFillMesh)
	pixelFillMesh.SetPixelSize(1)
	pixelFillMesh.Copy()

	for !opengl.WindowShouldClose() {
		opengl.BeginDrawing()

		if opengl.IsKeyPressed(glfw.KeyL) {
			closedGL.PrintlnVec2(opengl.GetMousePos())
		}
		if opengl.IsKeyPressed(glfw.KeyO) {
			breakAt++
		}
		if opengl.IsKeyPressed(glfw.KeyP) {
			breakAt--
		}
		if opengl.IsKeyPressed(glfw.KeyT) {
			amountPixel -= 2
		}
		if opengl.IsKeyPressed(glfw.KeyY) {
			amountPixel += 2
		}
		print = false
		if opengl.IsKeyPressed(glfw.KeyU) {
			print = !print
		}
		if opengl.IsKeyPressed(glfw.KeyI) {
			closedGL.PrintlnVec2(opengl.GetMousePos())
		}

		polyMesh.Clear()
		polyMesh.Copy()

		closedGL.SetWireFrameMode(!opengl.IsKeyDown(glfw.KeyF))
		opengl.ClearBG(glm.Vec4{0, 0, 0, 0})
		tri.Draw()
		polyMesh.Draw()
		lines.Draw()
		pixelFillMesh.Draw()
		debugMesh.Draw()

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
/*
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
	}
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
		}
		//	lines.Draw()
		pxTest.Draw()
		//openGL.DrawRect(glm.Vec4{0, 0, 100, 100}, glm.Vec4{0, 1, 0, 1}, 1)
		openGL.EndDrawing()
	}
	openGL.Free()
}
*/

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

func isOnLineTest() {
	var p1 = glm.Vec2{0, 0}     //ul
	var p2 = glm.Vec2{100, 100} //or
	var p3 = glm.Vec2{100, 0}   //ur
	var p4 = glm.Vec2{0, 100}   //ol

	var vertical = closedGL.CalculateLine(p1, p4)
	var horizontal = closedGL.CalculateLine(p2, p4)
	var normal = closedGL.CalculateLine(p3, p4)
	assert(vertical.IsOnLine(glm.Vec2{0, 50}))
	assert(!vertical.IsOnLine(glm.Vec2{0, 150}))
	assert(horizontal.IsOnLine(glm.Vec2{50, 100}))
	assert(!horizontal.IsOnLine(glm.Vec2{120, 100}))
	assert(normal.IsOnLine(glm.Vec2{50, 50}))
	assert(!normal.IsOnLine(glm.Vec2{60, 60}))

}
