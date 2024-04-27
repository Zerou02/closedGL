package main

import (
	"github.com/EngoEngine/glm"
	"github.com/go-gl/gl/v4.1-core/gl"
)

type Rectangle struct {
	dim        glm.Vec4
	shader     *Shader
	projection *glm.Mat4
	colour     glm.Vec3
	visible    bool
	vao, vbo   uint32
	vertices   []float32
}

func newRect(shader *Shader, projection *glm.Mat4, dim glm.Vec4, colour glm.Vec3) Rectangle {
	var rect = Rectangle{shader: shader, projection: projection, dim: dim, colour: colour}
	rect.vertices = make([]float32, 6*2*3*4)
	rect.vertices = []float32{
		dim[0] + dim[2], dim[1], colour[0], colour[1], colour[2], //top r
		dim[0] + dim[2], dim[1] + dim[3], colour[0], colour[1], colour[2], // bottom r
		dim[0], dim[1], colour[0], colour[1], colour[2], // top l
		dim[0] + dim[2], dim[1] + dim[3], colour[0], colour[1], colour[2], // bottom r
		dim[0], dim[1] + dim[3], colour[0], colour[1], colour[2], // bottom l
		dim[0], dim[1], colour[0], colour[1], colour[2], // top l
	}
	generateBuffers(&rect.vao, &rect.vbo, nil, nil, 6*2*3*4, nil, []VertexInfo{{2, 0}, {3, 8}})
	return rect
}

func (this *Rectangle) deleteBuffers() {
	gl.DeleteBuffers(1, &this.vao)
	gl.DeleteBuffers(1, &this.vbo)
}

func (this *Rectangle) createVertices() {
	var colour = this.colour
	if !this.visible {
		colour = glm.Vec3{0, 0, 0}
	}
	for i := 0; i < 6; i++ {
		this.vertices[(i*5)+2] = colour[0]
		this.vertices[(i*5)+3] = colour[1]
		this.vertices[(i*5)+4] = colour[2]
	}

}

func (this *Rectangle) draw() {

	this.shader.use()
	this.shader.setUniformMatrix4("projection", this.projection)
	gl.BindVertexArray(this.vao)
	gl.BindBuffer(gl.ARRAY_BUFFER, this.vbo)
	this.createVertices()
	gl.BufferSubData(gl.ARRAY_BUFFER, 0, 4*len(this.vertices), gl.Ptr(this.vertices))
	gl.DrawArrays(gl.TRIANGLES, 0, 6)
}
