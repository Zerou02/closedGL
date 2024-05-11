package closed_gl

import (
	"github.com/EngoEngine/glm"
	"github.com/go-gl/gl/v4.1-core/gl"
)

type LineArr struct {
	lineShader *Shader
	projection *glm.Mat4
	points     []float32
	vao, vbo   uint32
}

func newLineArr(shader *Shader, projection *glm.Mat4) LineArr {
	var lineArr = LineArr{lineShader: shader, points: []float32{}, projection: projection}
	return lineArr
}

func (this *LineArr) generateBuffers() {
	this.deleteBuffers()
	generateBuffers(&this.vao, &this.vbo, nil, this.points, 5*4, nil, []int{2, 3})
}

func (this *LineArr) deleteBuffers() {
	gl.DeleteBuffers(1, &this.vao)
	gl.DeleteBuffers(1, &this.vbo)
}

func (this *LineArr) addLine(line *Line) {
	this.points = append(this.points, line.points...)
	this.generateBuffers()
}

func (this *LineArr) draw() {
	if len(this.points)/5 < 1 {
		return
	}
	this.lineShader.use()
	this.lineShader.setUniformMatrix4("projection", this.projection)
	gl.BindVertexArray(this.vao)
	gl.BindBuffer(gl.ARRAY_BUFFER, this.vbo)
	gl.DrawArrays(gl.LINES, 0, int32(len(this.points)/5))
}

func (this *LineArr) addPoint(pos glm.Vec2, colour glm.Vec4) {
	this.points = append(this.points, pos[0])
	this.points = append(this.points, pos[1])
	this.points = append(this.points, colour[0])
	this.points = append(this.points, colour[1])
	this.points = append(this.points, colour[2])
	this.generateBuffers()
}
