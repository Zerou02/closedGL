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
	generateBuffers(&line.vao, &line.vbo, nil, nil, 2*4*5, nil, []VertexInfo{{2, 0}, {3, 8}})
	return line
}

func (this *Line) draw() {
	if len(this.points)/5 < 1 {
		return
	}
	this.shader.use()
	this.shader.setUniformMatrix4("projection", this.projection)
	gl.BindVertexArray(this.vao)
	gl.BindBuffer(gl.ARRAY_BUFFER, this.vbo)

	for i := 0; i < len(this.points)/5-1; i++ {
		var slice = this.points[i*5 : i*5+5]
		gl.BufferSubData(gl.ARRAY_BUFFER, 0, 2*4*5, gl.Ptr(slice))
		gl.DrawArrays(gl.LINES, 0, 2)
	}
}

func (this *Line) addPoint(p Point) {
	this.points = append(this.points, p.pos[0])
	this.points = append(this.points, p.pos[1])
	this.points = append(this.points, p.colour[0])
	this.points = append(this.points, p.colour[1])
	this.points = append(this.points, p.colour[2])
}
