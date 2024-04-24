package main

import (
	"github.com/EngoEngine/glm"
	"github.com/go-gl/gl/v4.1-core/gl"
)

type Rectangle struct {
	dim           glm.Vec4
	shader        *Shader
	projection    *glm.Mat4
	colour        glm.Vec3
	visible       bool
	vao, vbo, ebo uint32
}

func newRect(shader *Shader, projection *glm.Mat4, dim glm.Vec4, colour glm.Vec3) Rectangle {
	var rect = Rectangle{shader: shader, projection: projection, dim: dim}
	var vertices = []float32{
		dim[0] + dim[2], dim[1], colour[0], colour[1], colour[2], //top r
		dim[0] + dim[2], dim[1] + dim[3], colour[0], colour[1], colour[2], // bottom r
		dim[0], dim[1] + dim[3], colour[0], colour[1], colour[2], // bottom l
		dim[0], dim[1], colour[0], colour[1], colour[2], // top l
	}
	generateBuffers(&rect.vao, &rect.vbo, &rect.ebo, vertices, 0, indicesQuad, []VertexInfo{{2, 0}, {3, 8}})
	return rect
}

func (this *Rectangle) draw() {
	if !this.visible {
		return
	}
	this.shader.use()
	this.shader.setUniformMatrix4("projection", this.projection)
	gl.BindVertexArray(this.vao)
	//gl.BindBuffer(gl.ARRAY_BUFFER, this.vbo)
	gl.DrawElements(gl.TRIANGLES, 6, gl.UNSIGNED_INT, nil)
}
