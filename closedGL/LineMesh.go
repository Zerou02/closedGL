package closedGL

import (
	"github.com/EngoEngine/glm"
	"github.com/go-gl/gl/v4.6-core/gl"
)

type LineMesh struct {
	lineShader       *Shader
	projection, view glm.Mat4
	indices          []uint16
	vao              uint32
	amountPoints     int
	buffer           BufferFloat
}

func newLineMesh(shader *Shader, projection glm.Mat4) LineMesh {
	var LineMesh = LineMesh{lineShader: shader, projection: projection, vao: 0, view: glm.Ident4()}
	LineMesh.generateBuffers()
	return LineMesh
}

func (this *LineMesh) beginDraw() {
	this.indices = []uint16{}

	this.amountPoints = 0
}

func (this *LineMesh) generateBuffers() {
	this.vao = genVAO()
	this.buffer = BufferFloat{
		buffer:     generateInterleavedVBOFloat(this.vao, 0, []int{2, 4}, []int{0, 0}),
		bufferSize: 0,
		cpuArr:     []float32{},
	}
}

func (this *LineMesh) Draw() {
	if this.amountPoints < 1 {
		return
	}

	this.lineShader.use()
	this.lineShader.setUniformMatrix4("projection", &this.projection)
	this.lineShader.setUniformMatrix4("view", &this.view)

	gl.BindVertexArray(this.vao)
	gl.DrawElements(gl.LINES, int32(len(this.indices)), gl.UNSIGNED_SHORT, gl.Ptr(this.indices))
}

func (this *LineMesh) AddPoint(pos glm.Vec2, colour glm.Vec4) {
	this.buffer.resizeCPUData((this.amountPoints + 1) * 6)

	this.indices = append(this.indices, uint16(this.amountPoints))

	this.buffer.cpuArr[this.amountPoints*6+0] = pos[0]
	this.buffer.cpuArr[this.amountPoints*6+1] = pos[1]
	this.buffer.cpuArr[this.amountPoints*6+2] = colour[0]
	this.buffer.cpuArr[this.amountPoints*6+3] = colour[1]
	this.buffer.cpuArr[this.amountPoints*6+4] = colour[2]
	this.buffer.cpuArr[this.amountPoints*6+5] = colour[3]
	this.amountPoints++
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
		this.indices = append(this.indices, uint16(this.amountPoints))
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

func (this *LineMesh) Copy() {
	gl.BindVertexArray(this.vao)
	this.buffer.copyToGPU()
}
