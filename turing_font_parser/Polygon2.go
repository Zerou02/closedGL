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
	polygons []Polygon1
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

func NewPolygon2(glyf SimpleGlyf, ctx *closedGL.ClosedGLContext, triMesh *closedGL.TriangleMesh, lineMesh *closedGL.LineMesh, pMesh *closedGL.PixelMesh) Polygon2 {
	var retPolygon = Polygon2{
		glyf: glyf,
	}
	var points = glyf.GetPoints()
	for _, x := range points {
		var vertices = []glm.Vec2{}
		var edges = []uint32{}
		x = pointsToSS(x, ctx.Window.Wh)
		for i := 0; i < len(x); i += 3 {
			var first = x[i]
			var control = x[i+1]
			var second = x[i+2]

			first[1] -= 50
			control[1] -= 50
			second[1] -= 50

			var line = closedGL.CalculateLine(first, second)
			var isInside = retPolygon.isPointRightOfLine(first, second, control, ctx)
			vertices = append(vertices, first)

			if line.LineType == "normal" {
				vertices = append(vertices, control)
			}
			triMesh.AddTri([3]glm.Vec2{first, control, second}, [3]glm.Vec2{{0, 0}, {1.0 / 2.0, 0}, {1, 1}}, closedGL.Ternary(!isInside, float32(1), float32(-1)))
		}
		for i := 0; i < len(vertices); i++ {
			edges = append(edges, uint32(i))
		}
		edges = append(edges, 0)

		for i := 0; i < len(edges)-1; i++ {
			var a = vertices[edges[i]]
			var b = vertices[edges[i+1]]
			pMesh.AddPixel(a, glm.Vec4{1, 0, 0, 1})
			pMesh.AddPixel(b, glm.Vec4{1, 0, 0, 1})
			lineMesh.AddPath([]glm.Vec2{a, b}, []glm.Vec4{{0, 1, 0, 1}, {0, 1, 0, 1}})
		}
		retPolygon.polygons = append(retPolygon.polygons, Polygon1{
			vertices: vertices,
			edges:    edges,
		})
	}
	return retPolygon
}
