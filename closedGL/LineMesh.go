package closedGL

import (
	"github.com/EngoEngine/glm"
	"github.com/go-gl/gl/v4.6-core/gl"
)

type LineMesh struct {
	mesh SimpleMesh
}

func newLineMesh(shader *Shader, projection glm.Mat4) LineMesh {
	var LineMesh = LineMesh{
		mesh: newSimpleMesh(shader, projection, glm.Ident4()),
	}
	LineMesh.mesh.InitBuffer([]uint16{}, [][]int{{2, 4}}, [][]int{{0, 0}}, []uint32{gl.FLOAT})
	return LineMesh
}

func (this *LineMesh) Draw() {
	this.mesh.Draw()
}

func (this *LineMesh) AddPoint(pos glm.Vec2, colour glm.Vec4) {
	var vertices = []any{pos[0], pos[1], colour[0], colour[1], colour[2], colour[3]}
	var indices = []uint16{uint16(this.mesh.amountElements)}
	var test = []*[]any{&vertices}
	addVertices(&this.mesh, test, &indices)
}

func (this *LineMesh) AddLine(pos1, pos2 glm.Vec2, colour1, colour2 glm.Vec4) {
	this.AddPoint(pos1, colour1)
	this.AddPoint(pos2, colour2)

}

func (this *LineMesh) AddPath(pos []glm.Vec2, colours []glm.Vec4) {
	if len(pos) != len(colours) {
		return
	}
	this.AddPoint(pos[0], colours[0])
	for i := 1; i < len(pos)-1; i++ {
		this.mesh.indices = append(this.mesh.indices, uint16(this.mesh.amountElements))
		this.AddPoint(pos[i], colours[i])
	}
	this.AddPoint(pos[len(pos)-1], colours[len(pos)-1])
}

func (this *LineMesh) AddQuadraticBezier(p1, p2, controlPoint glm.Vec2, colour glm.Vec4) {
	var path = []glm.Vec2{}
	var colours = []glm.Vec4{}
	var maxPoints = 20
	for i := 0; i < maxPoints; i++ {
		var t = float32(i) / float32(maxPoints-1)
		path = append(path, BezierLerp(p1, p2, controlPoint, t))
		colours = append(colours, colour)
	}
	this.AddPath(path, colours)
}

func (this *LineMesh) AddQuadraticBezierLerp(p1, p2, controlPoint glm.Vec2, colour1, colour2 glm.Vec4) {
	var path = []glm.Vec2{}
	var colours = []glm.Vec4{}
	var maxPoints = 20
	for i := 0; i < maxPoints; i++ {
		var t = float32(i) / float32(maxPoints-1)
		path = append(path, BezierLerp(p1, p2, controlPoint, t))
		colours = append(colours, glm.Vec4{
			Lerp(colour1[0], colour2[0], t),
			Lerp(colour1[1], colour2[1], t),
			Lerp(colour1[2], colour2[2], t),
			Lerp(colour1[3], colour2[3], t),
		})
		this.AddPath(path, colours)
	}
}

func (this *LineMesh) CopyToGPU() {
	this.mesh.CopyToGPU()
}

func (this *LineMesh) Clear() {
	this.mesh.Clear()
}
