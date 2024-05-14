package closed_gl

import (
	"github.com/EngoEngine/glm"
	"github.com/go-gl/gl/v4.1-core/gl"
)

type Triangle struct {
	shader     *Shader
	projection *glm.Mat4
	Colour     glm.Vec4
	Points     []Vec2 //x,y
	vertices   []float32
	vao, vbo   Buffer
}

func newTriangle(shader *Shader, projection *glm.Mat4, colour glm.Vec4, points []Vec2) Triangle {
	var retTri = Triangle{shader: shader, projection: projection, Colour: colour, Points: points, vertices: make([]float32, 3*(2+4))}
	retTri.createVertices()
	generateBuffers(&retTri.vao, &retTri.vbo, nil, retTri.vertices, 0, nil, []int{2, 4})
	return retTri
}

func (this *Triangle) createVertices() {
	const amountVertices = 3
	const pointStride = 2
	const vertexStride = 6
	const pointEntries = 2
	const colourEntries = 4
	for i := 0; i < amountVertices; i++ {
		for j := 0; j < pointEntries; j++ {
			this.vertices[(i*vertexStride)+j] = this.Points[i][j]
		}
		for j := 0; j < colourEntries; j++ {
			this.vertices[(i*vertexStride)+pointStride+j] = this.Colour[j]
		}
	}
}

func (this *Triangle) Draw() {
	this.shader.use()
	this.shader.setUniformMatrix4("projection", this.projection)
	gl.BindVertexArray(this.vao)
	gl.BindBuffer(gl.ARRAY_BUFFER, this.vbo)
	//	this.createVertices()
	gl.BufferSubData(gl.ARRAY_BUFFER, 0, 4*len(this.vertices), gl.Ptr(this.vertices))
	gl.DrawArrays(gl.TRIANGLES, 0, 6)

}
