package turingfontparser

import (
	"github.com/EngoEngine/glm"
	"github.com/EngoEngine/math"
	"github.com/Zerou02/closedGL/closedGL"
)

type Polygon1 struct {
	vertices []glm.Vec2
	edges    []uint32
}

type Polygon2 struct {
	polygons []*Polygon1
	glyf     SimpleGlyf
}

// points in ss
func (this *Polygon2) isPointRightOfLine(p1, p2 glm.Vec2, p glm.Vec2, ctx *closedGL.ClosedGLContext) bool {

	var first = p1
	var second = p2
	var control = p
	//Ist CP rechts der Linie?
	first = closedGL.SsToCartesian(first, ctx.Window.Wh)
	control = closedGL.SsToCartesian(control, ctx.Window.Wh)
	second = closedGL.SsToCartesian(second, ctx.Window.Wh)
	var vec1 = second.Sub(&first)
	var vec2 = control.Sub(&first)
	vec1.Normalize()
	vec2.Normalize()
	var positiveY = glm.Vec2{0, 1}
	var angle = closedGL.AngleTo(vec1, positiveY)
	var other = 2*math.Pi - angle
	var base = vec1
	vec1 = closedGL.Rotate(angle, base)

	if math.Abs(vec1[0]) < 0.1 {
		vec2 = closedGL.Rotate(angle, vec2)
	} else {
		vec1 = closedGL.Rotate(other, base)
		vec2 = closedGL.Rotate(other, vec2)
	}
	return vec2[0] > glm.Epsilon
}

func (this *Polygon2) formLine(first, second, third glm.Vec2) bool {
	var line = closedGL.CalculateLine(first, third)
	return line.LineType != "normal" || line.IsOnLine(second)
}
func (this *Polygon2) removePointsOnLine(vertices []glm.Vec2) []glm.Vec2 {
	var newPoints = []glm.Vec2{}
	for i := 0; i < len(vertices)-2; i++ {
		var first = vertices[i]
		var second = vertices[i+1]
		var third = vertices[i+2]
		newPoints = append(newPoints, first)
		if this.formLine(first, second, third) {
			i++
		}
	}
	if this.formLine(vertices[len(vertices)-1], vertices[len(vertices)-2], vertices[len(vertices)-3]) {
		newPoints = append(newPoints, vertices[len(vertices)-1])
	} else {
		newPoints = append(newPoints, vertices[len(vertices)-2])
		newPoints = append(newPoints, vertices[len(vertices)-1])

	}
	return newPoints
}
func NewPolygon2(glyf SimpleGlyf, ctx *closedGL.ClosedGLContext, triMesh *closedGL.TriangleMesh, lineMesh *closedGL.LineMesh, pMesh *closedGL.PixelMesh) Polygon2 {
	var retPolygon = Polygon2{
		glyf: glyf,
	}
	var points = glyf.GetPoints()
	for _, x := range points {
		var vertices = []glm.Vec2{}
		x = pointsToSS(x, ctx.Window.Wh)
		for i := 0; i < len(x); i += 3 {
			var first = x[i]
			var control = x[i+1]
			var second = x[i+2]

			var isInside = retPolygon.isPointRightOfLine(first, second, control, ctx)
			vertices = append(vertices, first)
			if isInside {
				vertices = append(vertices, control)
			}
			var uvCoords = [3]glm.Vec2{{0, 0}, {1.0 / 2.0, 0}, {1, 1}}
			var sign = closedGL.Ternary(!isInside, float32(1), float32(-1))
			triMesh.AddTri([3]glm.Vec2{first, control, second}, uvCoords, sign)
		}
		for _, x := range vertices {
			pMesh.AddPixel(x, glm.Vec4{0, 0, 1, 1})
		}
		vertices = retPolygon.removePointsOnLine(vertices)

		var poly = Polygon1{
			vertices: vertices,
		}
		poly.enumerateVertices()
		retPolygon.polygons = append(retPolygon.polygons, &poly)
	}

	retPolygon.mergeContours()
	for _, x := range retPolygon.polygons {
		x.triangulateThisShit(triMesh)
		x.showVerticesAndLines(pMesh, lineMesh)
	}

	return retPolygon
}

func (this *Polygon2) mergeContours() {
	var first = this.polygons[0]
	var newPolys = []*Polygon1{first}
	for i := 1; i < len(this.polygons); i++ {
		var x = this.polygons[i]
		if first.isOtherInsideOfThis(x) {
			first.mergeOtherInThisToFillHoles(x)
		} else {
			newPolys = append(newPolys, x)
		}
	}
	this.polygons = newPolys
}
func (this *Polygon1) polyFill(pMesh *closedGL.PixelMesh) {
	for i := 0; i < 800; i++ {
		for j := 0; j < 800; j++ {
			var p = glm.Vec2{float32(i), float32(j)}
			_ = p
			if this.isPointInPolygon(p) {
				var gray = glm.Vec4{0.5, 0.5, 0.5, 1}
				var white = glm.Vec4{1, 1, 1, 1}
				_, _ = gray, white
				pMesh.AddPixel(p, gray)
			}
		}
	}
}

func (this *Polygon1) isVertex(p glm.Vec2) bool {
	return closedGL.Contains(&this.vertices, p)
}
func (this *Polygon1) isOnVertexVerticalLine(p glm.Vec2) bool {
	for _, x := range this.vertices {
		if math.Abs(x[1]-p[1]) < glm.Epsilon {
			return true
		}
	}
	return false
}

func (this *Polygon1) getIntersectionPointsInPolygon(p glm.Vec2) []glm.Vec2 {
	var ret = []glm.Vec2{}
	if this.isOnVertexVerticalLine(p) {
		return this.getIntersectionPointsInPolygon(glm.Vec2{p[0], p[1] - 0.01})

	}
	var ray = closedGL.CalculateLine(glm.Vec2{0, p[1]}, p)
	for i := 0; i < len(this.edges)-1; i++ {
		var first = this.vertices[this.edges[i]]
		var second = this.vertices[this.edges[i+1]]
		var line = closedGL.CalculateLine(first, second)
		var cp, valid = line.GetIntersection(ray)
		if !valid {
			continue
		}
		if line.IsOnLine(cp) && cp[0] <= p[0] {
			ret = append(ret, cp)
		}
	}
	return ret
}

func (this *Polygon1) isPointInPolygon(p glm.Vec2) bool {
	return len(this.getIntersectionPointsInPolygon(p))%2 == 1
}

func (this *Polygon1) isOtherInsideOfThis(p *Polygon1) bool {
	var allInside = true
	for _, x := range p.vertices {
		if !this.isPointInPolygon(x) {
			allInside = false
		}
	}
	return allInside
}

func (this *Polygon1) findRightmostPoint() (glm.Vec2, int) {
	var rightMost = this.vertices[0]
	var idx = 0
	for i, x := range this.vertices {
		if x[0] > rightMost[0] {
			rightMost = x
			idx = i
		}
	}
	return rightMost, idx
}

func (this *Polygon1) enumerateVertices() {
	this.edges = []uint32{}
	for i := 0; i < len(this.vertices); i++ {
		this.edges = append(this.edges, uint32(i))
	}
	this.edges = append(this.edges, 0)
}

func (this *Polygon1) mergeOtherInThisToFillHoles(inner *Polygon1) {
	var rightMost, rightMostIdx = inner.findRightmostPoint()
	var ray = closedGL.CalculateLine(rightMost, glm.Vec2{800, rightMost[1]})
	var newVertices = []glm.Vec2{}
	for i := 0; i < len(this.vertices)-1; i++ {
		var first = this.vertices[i]
		var second = this.vertices[i+1]
		var line = closedGL.CalculateLine(first, second)
		var cp, valid = line.GetIntersection(ray)
		if this.isVertex(cp) {
			panic("met vertex")
		}
		if !valid {
			continue
		}
		if line.IsOnLine(cp) && cp[0] > rightMost[0] {
			//not met vertex
			var newVertexUpper = cp
			var newVertexLower = cp
			var rightMostUpper = rightMost
			var rightMostLower = rightMost

			newVertexUpper[1] -= 0.1
			newVertexLower[1] += 0.1
			rightMostUpper[1] -= 0.1
			rightMostLower[1] += 0.1

			_ = rightMostIdx
			var outerIdx = i
			for j := 0; j < outerIdx+1; j++ {
				newVertices = append(newVertices, this.vertices[j])
			}
			newVertices = append(newVertices, newVertexUpper)

			var newInner = []glm.Vec2{newVertexUpper, rightMostUpper}
			for j := rightMostIdx + 1; j < len(inner.vertices); j++ {
				newInner = append(newInner, inner.vertices[j])
			}
			for j := 0; j < rightMostIdx; j++ {
				newInner = append(newInner, inner.vertices[j])
			}

			newVertices = append(newVertices, newInner...)
			newVertices = append(newVertices, rightMostLower, newVertexLower)

			for j := outerIdx + 1; j < len(this.vertices); j++ {
				newVertices = append(newVertices, this.vertices[j])
			}
			break
		}
	}
	this.vertices = newVertices
	this.enumerateVertices()
}

func (this *Polygon1) triangulateThisShit(triMesh *closedGL.TriangleMesh) {
	var i = 0
	var edges = this.edges
	var nrFound = 0
	for len(edges) > 3 {
		var edgeLen = len(edges)
		i++
		var first = this.vertices[edges[(i-1)%edgeLen]]
		var tip = this.vertices[edges[i%edgeLen]]
		var second = this.vertices[edges[(i+1)%edgeLen]]

		var isInsidePoly = true
		//var containsOtherVertex = false
		var line = closedGL.CalculateLine(first, second)
		if first[1] == tip[1] && tip[1] == second[1] || first[0] == tip[0] && first[0] == second[0] {
			continue
		}
		for _, sp := range line.SamplePointsOnLine(10) {
			if sp.Equal(&first) || sp.Equal(&tip) || sp.Equal(&second) {
				continue
			}
			if !this.isPointInPolygon(sp) {
				isInsidePoly = false
			}
		}
		if isInsidePoly {
			triMesh.AddTri([3]glm.Vec2{first, tip, second}, [3]glm.Vec2{{1, 1}, {1, 1}, {1, 1}}, 1)
			edges = closedGL.RemoveAt(edges, i%edgeLen)
			println("found:", nrFound)
		}

	}
	triMesh.AddTri([3]glm.Vec2{this.vertices[edges[0]], this.vertices[edges[1]], this.vertices[edges[2]]}, [3]glm.Vec2{{1, 1}, {1, 1}, {1, 1}}, 1)
}

func (this *Polygon1) showVerticesAndLines(pMesh *closedGL.PixelMesh, lMesh *closedGL.LineMesh) {
	for i := 0; i < len(this.edges)-1; i++ {
		var a = this.vertices[this.edges[i]]
		var b = this.vertices[this.edges[i+1]]
		pMesh.AddPixel(a, glm.Vec4{1, 0, 0, 1})
		pMesh.AddPixel(b, glm.Vec4{1, 0, 0, 1})
		lMesh.AddPath([]glm.Vec2{a, b}, []glm.Vec4{{0, 1, 0, 1}, {0, 1, 0, 1}})
	}
}
