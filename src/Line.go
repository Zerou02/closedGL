package main

import (
	"github.com/EngoEngine/glm"
	"github.com/go-gl/gl/v4.1-core/gl"
)

type Line struct {
	shader     *Shader
	points     []float32
	vao, vbo   uint32
	projection *glm.Mat4
}

func newLine(shader *Shader, projection *glm.Mat4) Line {
	var line = Line{shader: shader, projection: projection}
	line.generateBuffers()
	return line
}

func (this *Line) deleteBuffers() {
	gl.DeleteBuffers(1, &this.vao)
	gl.DeleteBuffers(1, &this.vbo)
}

func (this *Line) generateBuffers() {
	this.deleteBuffers()
	generateBuffers(&this.vao, &this.vbo, nil, this.points, 0, nil, []int{2, 3})
}

func (this *Line) draw() {
	if len(this.points)/5 < 1 {
		return
	}
	this.shader.use()
	this.shader.setUniformMatrix4("projection", this.projection)
	gl.BindVertexArray(this.vao)
	gl.BindBuffer(gl.ARRAY_BUFFER, this.vbo)

	gl.DrawArrays(gl.LINES, 0, int32(len(this.points)/5))
}

func (this *Line) addPoint(p Point) {
	this.points = append(this.points, p.pos[0])
	this.points = append(this.points, p.pos[1])
	this.points = append(this.points, p.colour[0])
	this.points = append(this.points, p.colour[1])
	this.points = append(this.points, p.colour[2])
	this.generateBuffers()
}
