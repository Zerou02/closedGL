package closedGL

import (
	"github.com/EngoEngine/glm"
	"github.com/go-gl/gl/v4.6-core/gl"
)

type CircleMesh struct {
	mesh SimpleMesh
}

func newCircleMesh(shader *Shader, projection glm.Mat4) CircleMesh {
	var indices = []uint16{0, 1, 2, 2, 1, 3}
	var mesh = CircleMesh{
		mesh: newSimpleMesh(shader, projection, glm.Ident4()),
	}
	mesh.mesh.InitBuffer(indices, [][]int{{2}, {4, 4, 4}}, [][]int{{0}, {1, 1, 1}}, []uint32{gl.FLOAT, gl.FLOAT})
	var one float32 = 1.0
	var zero float32 = 0.0
	var first = []any{one, zero}
	var second = []any{zero, zero}
	var third = []any{one, one}
	var fourth = []any{zero, one}

	addVertices(&mesh.mesh, []*[]any{&first, {}}, &[]uint16{})
	addVertices(&mesh.mesh, []*[]any{&second, {}}, &[]uint16{})
	addVertices(&mesh.mesh, []*[]any{&third, {}}, &[]uint16{})
	addVertices(&mesh.mesh, []*[]any{&fourth, {}}, &[]uint16{})

	return mesh
}

func (this *CircleMesh) CopyToGPU() {
	this.mesh.CopyToGPU()
}

func (this *CircleMesh) Clear() {
	this.mesh.Clear()
}

func (this *CircleMesh) AddCircle(centre glm.Vec2, colour, borderColour glm.Vec4, radius, borderThickness float32) {
	var values = []any{
		centre[0], centre[1],
		radius, borderThickness,
		colour[0], colour[1], colour[2], colour[3],
		borderColour[0], borderColour[1], borderColour[2], borderColour[3],
	}
	addVertices(&this.mesh, []*[]any{{}, &values}, &[]uint16{})
}

func (this *CircleMesh) Draw() {
	this.mesh.Draw()
}
