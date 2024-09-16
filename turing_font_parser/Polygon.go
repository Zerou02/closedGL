package turingfontparser

import (
	"github.com/EngoEngine/glm"
	"github.com/EngoEngine/math"

	"github.com/Zerou02/closedGL/closedGL"
)

type Polygon struct {
	edges []glm.Vec2 //tuple of points
	glyf  SimpleGlyf
}

func NewPolygon(edges []glm.Vec2) Polygon {
	return Polygon{
		edges: edges,
	}
}

func (this *Polygon) FindDoubleVertices() []glm.Vec2 {
	var countMap = map[glm.Vec2]int{}
	for _, x := range this.edges {
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

func (this *Polygon) RemoveVertex(x glm.Vec2) {
	var idx = closedGL.FindIdx(this.edges, x)
	if idx == 0 {
		//???
		this.edges = closedGL.RemoveAt(this.edges, 0)
		this.edges = closedGL.RemoveAt(this.edges, 0)
		this.edges[len(this.edges)-1] = this.edges[0]
	} else {
		this.edges = closedGL.Remove(this.edges, x)
		this.edges = closedGL.Remove(this.edges, x)
	}
}

func (this *Polygon) getVertices() []glm.Vec2 {
	var retMap = map[glm.Vec2]glm.Vec2{}
	for _, x := range this.edges {
		retMap[x] = x
	}
	var retVec = []glm.Vec2{}
	for _, x := range retMap {
		retVec = append(retVec, x)
	}
	return retVec
}

func (this *Polygon) handleDoubleVertices() {
	var list = this.FindDoubleVertices()
	for _, x := range list {
		var idx = closedGL.FindIdx(this.edges, x)
		var vertex = this.edges[idx]

		var newVertx = glm.Vec2{vertex[0] + 0.1, vertex[1] + 0.1}

		this.edges[idx] = newVertx
		this.edges[idx+1] = newVertx
	}
}
func (this *Polygon) Triangulate(mesh *closedGL.TriangleMesh) {
	this.handleDoubleVertices()
	var i = 0
	var ret = 0
	var amountIter = 0
	for len(this.edges) >= 7 {
		amountIter++
		if amountIter == len(this.edges)*2 {
			var vertices []glm.Vec2 = this.getVertices()

			for len(vertices) > 3 {
				for i := 0; i < len(vertices); i++ {
					for j := 0; j < len(vertices); j++ {
						for k := 0; k < len(vertices); k++ {
							if k == j || k == i || j == i {
								continue
							}
							var tip = vertices[i]
							var first = vertices[j]
							var second = vertices[k]
							var line = closedGL.CalculateLine(first, second)
							var points = line.SamplePointsOnLine(10)
							var isInPoly = true
							for _, x := range points {
								if !this.isPointInPolygon(x) {
									isInPoly = false
									break
								}
							}
							if isInPoly {
								mesh.AddTri([3]glm.Vec2{first, tip, second}, [3]glm.Vec2{{1, 1}, {1, 1}, {1, 1}}, 1)
							}
						}
					}
				}
				vertices = closedGL.RemoveAt(vertices, 0)
			}
			return
		}
		var x = this.edges[i%len(this.edges)]
		var neighbours = this.findNeighboursInPoly(x)
		var isValidTri = true
		var line = closedGL.CalculateLine(neighbours[0], neighbours[1])
		for _, sp := range line.SamplePointsOnLine(10) {
			if sp.Equal(&x) || sp.Equal(&neighbours[0]) || sp.Equal(&neighbours[1]) {
				continue
			}
			if !this.isPointInPolygon(sp) {
				isValidTri = false
			}
		}

		if isValidTri {
			mesh.AddTri([3]glm.Vec2{x, neighbours[0], neighbours[1]}, [3]glm.Vec2{{1, 1}, {1, 1}, {1, 1}}, 1)
			this.RemoveVertex(x)
			ret++
			amountIter = 0
		}
		i++
	}
	mesh.AddTri([3]glm.Vec2{this.edges[0], this.edges[1], this.edges[3]}, [3]glm.Vec2{{1, 1}, {1, 1}, {1, 1}}, 1)
}

func (this *Polygon) findNeighboursInPoly(x glm.Vec2) [2]glm.Vec2 {
	var idx = closedGL.FindIdx(this.edges, x)
	var prev glm.Vec2
	if idx == 0 {
		prev = this.edges[len(this.edges)-2]
	} else {
		prev = this.edges[idx-1]
	}
	var next = this.edges[idx+2]
	return [2]glm.Vec2{prev, next}
}

func (this *Polygon) println() {
	for i := 0; i < len(this.edges); i++ {
		closedGL.PrintlnVec2(this.edges[i])
		if i%2 == 1 {
			println("----", i)
		}
	}
}

func NewPolyFromGlyf(glyf SimpleGlyf) Polygon {
	return Polygon{
		glyf: glyf,
	}
	/* 	var contours = [][]GlyfPoints{}
	   	var single = []GlyfPoints{}
	   	for _, x := range glyf.GetPoints() {
	   		var a = x
	   		a.Pos = closedGL.CartesianToSS(a.Pos, 800)
	   		 		a.Pos[0] = math.Ceil(a.Pos[0])
	   		   		a.Pos[1] = math.Ceil(a.Pos[1])
	   		single = append(single, a)
	   		if x.EndPoint {
	   			contours = append(contours, single)
	   			single = []GlyfPoints{}
	   		}
	   	}

	   	var outer = NewPolyFromContour(contours[0])
	   	var outerPolys = []Polygon{}
	   	var innerPolys = []Polygon{}
	   	for i, x := range contours {
	   		var poly = NewPolyFromContour(x)
	   		if i == 0 {
	   			outerPolys = append(outerPolys, poly)
	   		} else {
	   			var isInner = outer.isOtherInThis(poly.edges)
	   			if !isInner {
	   				panic("multiple outer!")
	   			}

	   			innerPolys = append(innerPolys, poly)
	   		}
	   	}

	   	//Löcher stopfen
	   	var dy float32 = 0.1
	   	for _, x := range innerPolys {
	   		var rightMost = x.FindRightmostPoint()
	   		var baseLine = closedGL.CalculateLine(glm.Vec2{0, rightMost[1]}, rightMost)
	   		var found = false

	   		for i := 0; i < len(outer.edges)-1; i += 2 {
	   			if found {
	   				break
	   			}
	   			var p1 = outer.edges[i]
	   			var p2 = outer.edges[i+1]
	   			if p1[0] < rightMost[0] && p2[0] < rightMost[0] {
	   				continue
	   			}
	   			var line = closedGL.CalculateLine(p1, p2)
	   			var cp, succ = baseLine.GetIntersection(line)
	   			if succ {
	   				if line.IsOnLine(cp) {
	   					var metVertex = p1[1] == rightMost[1] || p2[1] == rightMost[1]
	   					if !metVertex {
	   						var idx = closedGL.FindIdx(outer.edges, p1)
	   						//remove old kantenzug
	   						outer.edges = closedGL.RemoveAt(outer.edges, idx+1)
	   						outer.edges = closedGL.RemoveAt(outer.edges, idx+1)

	   						idx = closedGL.FindIdx(outer.edges, p1)
	   						if idx == -1 {
	   							panic("error")
	   						}

	   						var new = glm.Vec2{cp[0], cp[1]}
	   						//new doppelt als Endpunkt für geschnittenen
	   						var firstKantenzug = []glm.Vec2{p1, new, new, rightMost}
	   						outer.edges = closedGL.InsertArrAt(outer.edges, firstKantenzug, idx+1)

	   						var innerIdx = closedGL.FindIdx(x.edges, rightMost)
	   						var newInner = []glm.Vec2{}
	   						for j := innerIdx + 1; j < len(x.edges); j++ {
	   							newInner = append(newInner, x.edges[j])
	   						}
	   						for j := 0; j < innerIdx; j++ {
	   							newInner = append(newInner, x.edges[j])
	   						}
	   						newInner = append(newInner, newInner[0])
	   						newInner = append(newInner, rightMost, new, new, p2)
	   						outer.edges = closedGL.InsertArrAt(outer.edges, newInner, idx+5)
	   						found = true
	   					} else {
	   						var met = p1
	   						if p2[1] == rightMost[1] {
	   							met = p2
	   						}

	   						var idx = closedGL.FindIdx(outer.edges, met)
	   						for idx != -1 {
	   							outer.edges[idx] = glm.Vec2{outer.edges[idx][0], outer.edges[idx][1] - dy}
	   							idx = closedGL.FindIdx(outer.edges, met)
	   						}
	   						i = 0
	   					}
	   				}
	   			}
	   		}
	   	} */

	//genauere Approximation
	/* 	var ret = outer */
	/*
		if ret.hasSameKantenzug() {
			panic("has same kantenzug")
		}

		var foundOffPoints = []glm.Vec2{}
		for _, x := range contours {
			for i, y := range x {

				if !y.OnCurve && ret.isPointInPolygon(y.Pos) && !closedGL.Contains(&foundOffPoints, y.Pos) {
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
					var amount = closedGL.FindAmount(ret.edges, lastOnPoint)
					var polyIdx = closedGL.FindIdx(ret.edges, lastOnPoint)
					if amount != 2 {
						ret.handleDoubleVertices()
						i = -1
						println("special case: amount not 2!!", polyIdx)
						continue
					}
					var isAtEnd = polyIdx == len(ret.edges)-3

					var startKantenzug = polyIdx == 2
					var newKantenzug = []glm.Vec2{lastOnPoint, y.Pos, y.Pos, nextOnPoint}
					if isAtEnd {
						newKantenzug = []glm.Vec2{lastOnPoint, y.Pos, y.Pos, ret.edges[len(ret.edges)-1]}
					}

					if startKantenzug {
						ret.edges = closedGL.RemoveAt(ret.edges, polyIdx)
						ret.edges = closedGL.RemoveAt(ret.edges, polyIdx)
						ret.edges = closedGL.InsertArrAt(ret.edges, newKantenzug, polyIdx)
					} else {
						ret.edges = closedGL.RemoveAt(ret.edges, polyIdx+1)
						ret.edges = closedGL.RemoveAt(ret.edges, polyIdx+1)
						ret.edges = closedGL.InsertArrAt(ret.edges, newKantenzug, polyIdx+1)
					}
					if ret.hasSameKantenzug() {
						panic("ERRROR, ident kantenzug")
					}
				}
			}
		}
	*/
	/* 	ret.glyf = glyf
	   	return ret */
	return Polygon{}
}

func (this *Polygon) FindRightmostPoint() glm.Vec2 {
	var maxX = this.edges[0][0]
	var k = 0
	for i, x := range this.edges {
		if x[0] > maxX {
			maxX = x[0]
			k = i
		}
	}
	return this.edges[k]
}

func NewPolyFromContour(contour []GlyfPoints) Polygon {
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

	var p = Polygon{
		edges: duplicateTrain,
	}
	return p

}

func (this *Polygon) isOtherInThis(other []glm.Vec2) bool {
	var minX, maxX = this.edges[0][0], this.edges[0][0]
	var minY, maxY = this.edges[0][1], this.edges[0][1]
	for _, x := range this.edges {
		minX = math.Min(x[0], minX)
		maxX = math.Max(x[0], maxX)
		minY = math.Min(x[1], minY)
		maxY = math.Max(x[1], maxY)
	}
	var isInside = false
	for _, x := range other {
		if this.isPointInPolygon(x) {
			isInside = true
			break
		}
	}
	return isInside
}

func (this *Polygon) hasSameKantenzug() bool {
	for i := 0; i < len(this.edges); i += 2 {
		if this.edges[i].Equal(&this.edges[i+1]) {
			println("I", i)
			return true
		}
	}
	return false
}

func (this *Polygon) isVertexOfPolygon(p glm.Vec2) bool {
	for _, x := range this.edges {
		if x.Equal(&p) {
			return true
		}
	}
	return false
}

func (this *Polygon) isOnVertexLineOfPolygon(p glm.Vec2) bool {
	for _, x := range this.edges {
		if x[1] == p[1] {
			return true
		}
	}
	return false
}

// ausgelegt für SS-Polys,mit Kantenzügen
func (this *Polygon) getInterSectionPoints(p glm.Vec2) ([]glm.Vec2, bool) {
	var dy float32 = 0.1
	if this.isOnVertexLineOfPolygon(p) {
		return this.getInterSectionPoints(glm.Vec2{p[0], p[1] - dy})
	}

	var ray = closedGL.CalculateLine(glm.Vec2{0, p[1]}, p)
	_ = ray
	var retMap = map[glm.Vec2]glm.Vec2{}
	for i := 0; i < len(this.edges)-1; i += 2 {
		var p1 = this.edges[i]
		var p2 = this.edges[i+1]
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

func (this *Polygon) isPointInPolygon(p glm.Vec2) bool {
	var intersections, _ = this.getInterSectionPoints(p)
	return len(intersections)%2 == 1
}

func (this *Polygon) pixelFillPoly(mesh *closedGL.PixelMesh) {
	for i := 0; i < 800; i++ {
		for j := 0; j < 800; j++ {
			var p = glm.Vec2{float32(j), float32(i)}
			if this.isPointInPolygon(p) {
				mesh.AddPixel(p, glm.Vec4{1, 1, 1, 1})
			}
		}
	}
}

func (this *Polygon) GetEdges() []glm.Vec2 {
	return this.edges
}

func pointsToSS(points []glm.Vec2, wh float32) []glm.Vec2 {
	var ret = []glm.Vec2{}
	for _, x := range points {
		ret = append(ret, closedGL.CartesianToSS(x, wh))
	}
	return ret
}

func extractPoints(glyfPoints []GlyfPoints) []glm.Vec2 {
	var points = []glm.Vec2{}
	for _, x := range glyfPoints {
		points = append(points, x.Pos)
	}
	return points
}

func (this *Polygon) FillMeshes(triMesh *closedGL.TriangleMesh, lineMesh *closedGL.LineMesh, pixelMesh *closedGL.PixelMesh, wh float32) {

	/* 	var points = this.glyf.GetPoints()
	   	for _, x := range points {
	   		x = pointsToSS(x, wh)
	   		for i := 0; i < len(x); i += 3 {
	   			var first = x[i]
	   			var control = x[i+1]
	   			var second = x[i+2]
	   			triMesh.AddTri([3]glm.Vec2{first, control, second}, [3]glm.Vec2{{0, 0}, {1.0 / 2.0, 0}, {1, 1}}, 1)
	   			lineMesh.AddPath([]glm.Vec2{first, control, second}, []glm.Vec4{{0, 1, 0, 1}, {0, 1, 0, 1}, {0, 1, 0, 1}})
	   		}

	   	} */

	/* 	var points = extractPoints(this.glyf.GetPoints())
	   	points = pointsToSS(points, wh)
	   	for _, x := range points {
	   		closedGL.PrintlnVec2(x)
	   		pixelMesh.AddPixel(x, glm.Vec4{1, 0, 0, 1})
	   	}
	   	println(len(points))
	   	for i := 0; i < len(points); i += 3 {
	   		var sign float32 = 1
	   		var off = points[i+2]
	   		 off[0] = math.Floor(off[0])
	   		off[1] = math.Floor(off[1])
	   		if closedGL.Contains(&this.edges, off) {
	   			sign = -1
	   		}
	   		triMesh.AddTri([3]glm.Vec2{points[i], points[i+2], points[i+1]}, [3]glm.Vec2{{0, 0}, {1.0 / 2.0, 0}, {1, 1}}, sign)
	   	}
	   	for i := 0; i < len(this.edges); i += 2 {
	   		lineMesh.AddLine(this.edges[i], this.edges[i+1], glm.Vec4{1, 1, 1, 1}, glm.Vec4{1, 1, 1, 1})
	   	} */

}
