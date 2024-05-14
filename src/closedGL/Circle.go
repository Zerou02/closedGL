package closed_gl

import (
	"github.com/EngoEngine/glm"
	"github.com/go-gl/gl/v4.1-core/gl"
)

type Circle struct {
	Radius          float32
	Centre          Vec2 //x,,y
	BorderThickness float32
	CentreColour    glm.Vec4
	BorderColour    glm.Vec4
	shader          *Shader
	vertices        []float32
	projection      *glm.Mat4
	vao, vbo        Buffer
}

func newCircle(shader *Shader, projection *glm.Mat4, radius float32, centre Vec2, centreColour, borderColour glm.Vec4, borderThickness float32) Circle {
	var circle = Circle{shader: shader, projection: projection, Radius: radius, Centre: centre, CentreColour: centreColour, BorderColour: borderColour, vertices: make([]float32, 6*2), BorderThickness: borderThickness}
	circle.createVertices()
	generateBuffers(&circle.vao, &circle.vbo, nil, circle.vertices, 0, nil, []int{2})
	return circle
}

func (this *Circle) createVertices() {
	this.vertices = []float32{
		this.Centre[0] + this.Radius + this.BorderThickness, this.Centre[1] - this.Radius - this.BorderThickness, //top r
		this.Centre[0] + this.Radius + this.BorderThickness, this.Centre[1] + this.Radius + this.BorderThickness, // bottom r
		this.Centre[0] - this.Radius - this.BorderThickness, this.Centre[1] - this.Radius - this.BorderThickness, // top l
		this.Centre[0] + this.Radius + this.BorderThickness, this.Centre[1] + this.Radius + this.BorderThickness, // bottom r
		this.Centre[0] - this.Radius - this.BorderThickness, this.Centre[1] + this.Radius + this.BorderThickness, // bottom l
		this.Centre[0] - this.Radius - this.BorderThickness, this.Centre[1] - this.Radius - this.BorderThickness, // top l
	}
}

func (this *Circle) Draw() {

	this.shader.use()
	this.shader.setUniformMatrix4("projection", this.projection)
	this.shader.setUniform1f("radius", this.Radius)
	this.shader.setUniform1f("borderThickness", this.BorderThickness)
	this.shader.setUniform2fv("centre", this.Centre)
	this.shader.setUniformVec4("centreColour", &this.CentreColour)
	this.shader.setUniformVec4("borderColour", &this.BorderColour)

	gl.BindVertexArray(this.vao)
	gl.BindBuffer(gl.ARRAY_BUFFER, this.vbo)
	this.createVertices()
	gl.BufferSubData(gl.ARRAY_BUFFER, 0, 4*len(this.vertices), gl.Ptr(this.vertices))
	gl.DrawArrays(gl.TRIANGLES, 0, 6)
}
