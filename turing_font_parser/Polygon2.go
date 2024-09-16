package turingfontparser

import (
	"github.com/EngoEngine/glm"
	"github.com/Zerou02/closedGL/closedGL"
)

type Polygon1 struct {
	vertices []glm.Vec2
	edges    [][2]uint32
}

type Polygon2 struct {
	polygons []Polygon1
	glyf     SimpleGlyf
}

func NewPolygon2(glyf SimpleGlyf, ctx *closedGL.ClosedGLContext, triMesh *closedGL.TriangleMesh, lineMesh *closedGL.LineMesh) Polygon2 {
	var retPolygon = Polygon2{
		glyf: glyf,
	}
	var points = glyf.GetPoints()
	for _, x := range points {
		var vertices = []glm.Vec2{}
		var edges = [][2]uint32{}
		x = pointsToSS(x, ctx.Window.Wh)
		for i := 0; i < len(x); i += 3 {
			var first = x[i]
			var control = x[i+1]
			var second = x[i+2]
			if !closedGL.Contains(&vertices, first) {
				vertices = append(vertices, first)
			}
			if !closedGL.Contains(&vertices, control) {
				vertices = append(vertices, control)
			}
			if !closedGL.Contains(&vertices, second) {
				vertices = append(vertices, second)
			}

			triMesh.AddTri([3]glm.Vec2{first, control, second}, [3]glm.Vec2{{0, 0}, {1.0 / 2.0, 0}, {1, 1}}, 1)
		}
		for i := 0; i < len(vertices)-1; i++ {
			edges = append(edges, [2]uint32{uint32(i), uint32(i + 1)})
		}
		edges = append(edges, [2]uint32{uint32(len(vertices) - 1), 0})
		for _, x := range edges {
			lineMesh.AddPath([]glm.Vec2{vertices[x[0]], vertices[x[1]]}, []glm.Vec4{{0, 1, 0, 1}, {0, 1, 0, 1}})
		}
		retPolygon.polygons = append(retPolygon.polygons, Polygon1{
			vertices: vertices,
			edges:    edges,
		})
	}
	return retPolygon
}
