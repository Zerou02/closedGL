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

func hasSameKantenzug(poly []glm.Vec2) bool {
	for i := 0; i < len(poly); i += 2 {
		if poly[i].Equal(&poly[i+1]) {
			println("I", i)
			return true
		}
	}
	return false
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

	/* var ret = []glm.Vec2{}
	for _, x := range contours {
		ret = append(ret, extractPolyFromContour(x)...)
	}
	*/
	var outer = extractPolyFromContour(contours[0])
	var outerContours = [][]glm.Vec2{}
	var innerContours = [][]glm.Vec2{}
	for i, x := range contours {
		var poly = extractPolyFromContour(x)
		if i == 0 {
			outerContours = append(outerContours, poly)
		} else {
			var isInner = isInsideOfOtherPolyogn(outerContours[0], poly)
			if !isInner {
				panic("multiple outer!")
			}
			innerContours = append(innerContours, poly)
		}
	}

	println("inner-outer", len(innerContours), len(outerContours))
	//Löcher stopfen
	for _, x := range innerContours {
		var rightMost = findRightmostPoint(x)
		var baseLine = closedGL.CalculateLine(glm.Vec2{0, rightMost[1]}, rightMost)
		var found = false

		closedGL.PrintlnVec2(rightMost)
		for i := 0; i < len(outer)-1; i += 2 {
			if found {
				break
			}
			var p1 = outer[i]
			var p2 = outer[i+1]
			if p1[0] < rightMost[0] && p2[0] < rightMost[0] {
				continue
			}
			var line = closedGL.CalculateLine(p1, p2)
			var cp, succ = baseLine.GetIntersection(line)
			if succ {
				if line.IsOnLine(cp) {
					var metVertex = p1[1] == rightMost[1] || p2[1] == rightMost[1]
					println("met vertex", metVertex)
					if !metVertex {
						var idx = closedGL.FindIdx(outer, p1)
						//remove old kantenzug
						outer = closedGL.RemoveAt(outer, idx+1)
						outer = closedGL.RemoveAt(outer, idx+1)

						idx = closedGL.FindIdx(outer, p1)
						if idx == -1 {
							panic("error")
						}

						var new = glm.Vec2{cp[0], cp[1]}
						//new doppelt als Endpunkt für geschnittenen
						var firstKantenzug = []glm.Vec2{p1, new, new, rightMost}
						outer = closedGL.InsertArrAt(outer, firstKantenzug, idx+1)

						var innerIdx = closedGL.FindIdx(x, rightMost)
						var newInner = []glm.Vec2{}
						for j := innerIdx + 1; j < len(x); j++ {
							newInner = append(newInner, x[j])
						}
						for j := 0; j < innerIdx; j++ {
							newInner = append(newInner, x[j])
						}
						newInner = append(newInner, newInner[0])
						newInner = append(newInner, rightMost, new, new, p2)
						outer = closedGL.InsertArrAt(outer, newInner, idx+5)
						found = true
					} else {
						var met = p1
						if p2[1] == rightMost[1] {
							met = p2
						}

						var idx = closedGL.FindIdx(outer, met)
						for idx != -1 {
							outer[idx] = glm.Vec2{outer[idx][0], outer[idx][1] - 1}
							idx = closedGL.FindIdx(outer, met)
						}
						i = 0
						/* 	var met = p1
						var other = p2
						if rightMost[1] == p2[1] {
							met = p2
							other = p1
						}
						var idx = closedGL.FindIdx(outer, met)
						_ = other
						outer = closedGL.RemoveAt(outer, idx)
						outer = closedGL.RemoveAt(outer, idx)
						var first = []glm.Vec2{met, met, rightMost}
						outer = closedGL.InsertArrAt(outer, first, idx)

						printlnPoly(outer)
						println(hasSameKantenzug(outer))
						panic("show")
						found = true */
						//				panic("met vertex. Please investigate. might be harmless")
					}
				}
			}
		}
	}

	//genauere Approximation
	var ret = outer
	if hasSameKantenzug(ret) {
		panic("has same kantenzug")
	}
	printlnPoly(ret)
	println("beofre approx")

	var foundOffPoints = []glm.Vec2{}
	for _, x := range contours {
		var found = 0
		for i, y := range x {

			if !y.OnCurve && isPointInPolygon(y.Pos, ret) && !closedGL.Contains(&foundOffPoints, y.Pos) {
				var lastOnPoint = x[i-2].Pos
				var nextOnPoint = x[i-1].Pos

				var line = closedGL.CalculateLine(lastOnPoint, nextOnPoint)

				var evalX, _ = line.EvalX(y.Pos[0])
				var evalY, _ = line.EvalY(y.Pos[1])

				evalX[0] = math.Floor(evalX[0])
				evalX[1] = math.Floor(evalX[1])

				evalY[0] = math.Floor(evalY[0])
				evalY[1] = math.Floor(evalY[1])

				//no straight lines
				if evalX.Equal(&y.Pos) || evalY.Equal(&y.Pos) {
					continue
				}
				if lastOnPoint[0] == nextOnPoint[0] && nextOnPoint[0] == y.Pos[0] {
					continue
				}
				if lastOnPoint[1] == nextOnPoint[1] && nextOnPoint[1] == y.Pos[1] {
					continue
				}
				//Sonderfall beim kleinen b
				if lastOnPoint.Equal(&y.Pos) {
					continue
				}
				foundOffPoints = append(foundOffPoints, y.Pos)
				var amonut = closedGL.FindAmount(ret, lastOnPoint)
				if amonut != 2 {
					println("amo!!")
					println(amonut)
					closedGL.PrintlnVec2(lastOnPoint)
					continue
					panic("amo")
				}
				var polyIdx = closedGL.FindIdx(ret, lastOnPoint)
				var isAtEnd = polyIdx == len(ret)-3

				var startKantenzug = polyIdx == 2
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
				if found == 1 {

					println("evalX")
					closedGL.PrintlnVec2(evalX)
					println("evalY")
					closedGL.PrintlnVec2(evalY)
					println("on")
					closedGL.PrintlnVec2(lastOnPoint)
					println("next on")
					closedGL.PrintlnVec2(nextOnPoint)
					println("off")
					closedGL.PrintlnVec2(y.Pos)
					println("start", startKantenzug)
					println("end", isAtEnd)
					println("i", i)
				}
				if hasSameKantenzug(ret) {
					printlnPoly(ret)

					println("evalX")
					closedGL.PrintlnVec2(evalX)
					println("on")
					closedGL.PrintlnVec2(lastOnPoint)
					println("next on")
					closedGL.PrintlnVec2(nextOnPoint)
					println("off")
					closedGL.PrintlnVec2(y.Pos)
					println("start", startKantenzug)
					println("end", isAtEnd)
					println("i", i)

					panic("ERRROR, ident kantenzug")
				}
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

func glyfToPoints(pointsGlyf []turingfontparser.GlyfPoints, poly []glm.Vec2, triMesh *closedGL.TriangleMesh) {
	var points = extractPoints(pointsGlyf)
	points = pointsToSS(points)
	for i := 0; i < len(points); i += 3 {
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

type VertexData struct {
	p0, p1, tip   glm.Vec2
	isConvex      bool
	isEar         bool
	interiorAngle float32
}

func findNeighboursInPoly(x glm.Vec2, poly []glm.Vec2) [2]glm.Vec2 {
	var idx = closedGL.FindIdx(poly, x)
	var prev glm.Vec2
	if idx == 0 {
		prev = poly[len(poly)-2]
	} else {
		prev = poly[idx-1]
	}
	var next = poly[idx+2]
	return [2]glm.Vec2{prev, next}
}

func getAngleData(poly []glm.Vec2) map[glm.Vec2]VertexData {
	var retMap = map[glm.Vec2]VertexData{}
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
		var vec1 = p0.Sub(&tip)
		var vec2 = p1.Sub(&tip)
		var positiveY = glm.Vec2{0, 1}
		var nVec1 = vec1.Normalized()
		var nVec2 = vec2.Normalized()
		var rotAngle = closedGL.AngleTo(positiveY, nVec1)
		var rot1 = closedGL.Rotate(rotAngle, nVec1)
		var rot2 = closedGL.Rotate(rotAngle, nVec2)
		if math.Abs(rot1[0]) > 0.1 {
			closedGL.PrintlnVec2(p0)
			closedGL.PrintlnVec2(tip)
			closedGL.PrintlnVec2(p1)
			println("--")
			closedGL.PrintlnVec2(rot1)
			//		panic("hm")
		}
		if rot1[1] > 0 {
			rot1[1] *= -1
			rot2[0] *= -1
			rot2[1] *= -1
		}
		var rightTurn = rot2[0] > glm.Epsilon

		var angle = closedGL.AngleTo(vec1, vec2)
		var convex = rightTurn || glm.FloatEqual(angle, math.Pi)
		if !rightTurn {
			angle = 2*math.Pi - angle
		}

		tip = closedGL.CartesianToSS(tip, 800)
		p0 = closedGL.CartesianToSS(p0, 800)
		p1 = closedGL.CartesianToSS(p1, 800)

		retMap[tip] = VertexData{
			p0:            p0,
			p1:            p1,
			tip:           tip,
			isConvex:      convex,
			isEar:         true,
			interiorAngle: angle,
		}
	}
	for k, x := range retMap {
		var data = x
		if !x.isConvex {
			data.isEar = false
			retMap[k] = data
			continue
		}
		for _, p := range poly {
			if p.Equal(&data.p0) || p.Equal(&data.p1) || p.Equal(&data.tip) {
				continue
			}
			if closedGL.PointInTriangle(p, data.p0, data.tip, data.p1) {
				if !retMap[p].isConvex {
					data.isEar = false
					break
				}
			}
		}
		retMap[k] = data
	}

	return retMap
}

func removeVertexFromPoly(x glm.Vec2, poly []glm.Vec2) []glm.Vec2 {
	var idx = closedGL.FindIdx(poly, x)
	if idx == 0 {
		println("???")
		//???
		poly = closedGL.RemoveAt(poly, 0)
		poly = closedGL.RemoveAt(poly, 0)
		poly[len(poly)-1] = poly[0]
	} else {
		poly = closedGL.Remove(poly, x)
		poly = closedGL.Remove(poly, x)
	}
	return poly
}

func findDoubleVertices(poly []glm.Vec2) []glm.Vec2 {
	var countMap = map[glm.Vec2]int{}
	for _, x := range poly {
		countMap[x] += 1
	}
	var ret = []glm.Vec2{}
	for k, x := range countMap {
		if x >= 3 {
			ret = append(ret, k)
		}
	}
	return ret
}

func triangulatePoly(poly []glm.Vec2, mesh *closedGL.TriangleMesh, pMesh *closedGL.PixelMesh) {
	var list = findDoubleVertices(poly)
	println("double?", len(list))
	println(len(poly))
	for _, x := range list {
		var idx = closedGL.FindIdx(poly, x)
		println(idx)
		var vertex = poly[idx]
		var newVertx = glm.Vec2{vertex[0] - 1, vertex[1] - 1}
		poly[idx] = newVertx
		poly[idx+1] = newVertx
	}
	list = findDoubleVertices(poly)
	println("still double?", len(list))

	var i = 0
	var ret = 0
	var amountIter = 0
	for len(poly) >= 7 {
		amountIter++
		if amountIter == len(poly)*2 {
			printlnPoly(poly)
			println(len(poly))
			var vertices []glm.Vec2 = []glm.Vec2{poly[0], poly[1], poly[3], poly[5]}

			var found = false
			for _, x := range vertices {
				for _, y := range vertices {
					var dist = y.Sub(&x)
					if !x.Equal(&y) && math.Abs(dist.Len()) < 5 {
						vertices = closedGL.Remove(vertices, y)
						found = true
						break
					}
				}
				if found {
					break
				}
			}
			println("AAAA")
			mesh.AddTri([3]glm.Vec2{vertices[0], vertices[1], vertices[2]}, [3]glm.Vec2{{1, 1}, {1, 1}, {1, 1}}, 1)
			return
		}
		var x = poly[i%len(poly)]
		var neighbours = findNeighboursInPoly(x, poly)
		var isValidTri = true
		var line = closedGL.CalculateLine(neighbours[0], neighbours[1])
		for _, sp := range line.SamplePointsOnLine(10) {
			if sp.Equal(&x) || sp.Equal(&neighbours[0]) || sp.Equal(&neighbours[1]) {
				continue
			}
			if !isPointInPolygon(sp, poly) {
				isValidTri = false
			}
		}

		if isValidTri {
			mesh.AddTri([3]glm.Vec2{x, neighbours[0], neighbours[1]}, [3]glm.Vec2{{1, 1}, {1, 1}, {1, 1}}, 1)
			poly = removeVertexFromPoly(x, poly)
			ret++
			amountIter = 0
		}
		i++
	}
	mesh.AddTri([3]glm.Vec2{poly[0], poly[1], poly[3]}, [3]glm.Vec2{{1, 1}, {1, 1}, {1, 1}}, 1)
}

func StartTTF() {
	var opengl = closedGL.InitClosedGL(800, 800, "comic")

	opengl.LimitFPS(true)
	var p = turingfontparser.NewTuringFont("./assets/font/jetbrains_mono_medium.ttf", &opengl)
	var points = []turingfontparser.GlyfPoints{}
	var offset float32 = 0
	//input
	for _, x := range "h" {
		_ = x
		var newPoints = p.ParseGlyf(uint32(x), 0.25).AddXOffset(offset)
		var biggestX float32 = offset
		for i, y := range newPoints {
			var offset = glm.Vec2{0, 50}
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
	glyfToPoints(points, poly, &tri)
	for i := 0; i < len(poly); i += 2 {
		lines.AddLine(poly[i], poly[i+1], glm.Vec4{1, 1, 1, 1}, glm.Vec4{1, 1, 1, 1})
	}

	lines.Copy()
	polyMesh.SetPixelSize(1)

	debugMesh.SetPixelSize(5)
	leaveMesh.SetPixelSize(1)
	for _, x := range poly {
		debugMesh.AddPixel(x, glm.Vec4{1, 0, 0, 1})
	}
	triangulatePoly(poly, &tri, &debugMesh)
	tri.Copy()

	var breakAt = 1
	var print = false
	var amountPixel = 1

	tri.Copy()
	debugMesh.Copy()

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

		opengl.SetWireFrameMode(opengl.IsKeyDown(glfw.KeyF))
		opengl.ClearBG(glm.Vec4{0, 0, 0, 0})
		tri.Draw()
		/* 		var currMode = opengl.GetWireFrameMode()
		   		opengl.SetWireFrameMode(!currMode)
		   		tri.Draw()
		   		opengl.SetWireFrameMode(!currMode) */
		//	closedGL.SetWireFrameMode(true)
		//	tri.Draw()
		lines.Draw()
		//	debugMesh.Draw()

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
