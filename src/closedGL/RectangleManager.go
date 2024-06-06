package closed_gl

import (
	"github.com/EngoEngine/glm"
	"github.com/go-gl/gl/v4.1-core/gl"
)

type RectangleManager struct {
	shader     *Shader
	projection *glm.Mat4
	vao, vbo   uint32
	vboLen     uint32
	vertices   []float32
}

func newRect(shader *Shader, projection *glm.Mat4) RectangleManager {
	var rect = RectangleManager{shader: shader, projection: projection, vboLen: 1, vertices: []float32{}}
	generateBuffers(&rect.vao, &rect.vbo, nil, nil, int(rect.vboLen), nil, []int{2, 4})
	return rect
}

func (this *RectangleManager) deleteBuffers() {
	gl.DeleteBuffers(1, &this.vao)
	gl.DeleteBuffers(1, &this.vbo)
}

func (this *RectangleManager) clearVertices() {
	this.vertices = make([]float32, 0)
}
func (this *RectangleManager) createVertices(dim, colour glm.Vec4) {
	var newVertices = []float32{
		dim[0] + dim[2], dim[1], colour[0], colour[1], colour[2], colour[3], //top r
		dim[0], dim[1], colour[0], colour[1], colour[2], colour[3], // top l
		dim[0] + dim[2], dim[1] + dim[3], colour[0], colour[1], colour[2], colour[3], // bottom r
		dim[0] + dim[2], dim[1] + dim[3], colour[0], colour[1], colour[2], colour[3], // bottom r
		dim[0], dim[1], colour[0], colour[1], colour[2], colour[3], // top l
		dim[0], dim[1] + dim[3], colour[0], colour[1], colour[2], colour[3], // bottom l
	}
	this.vertices = append(this.vertices, newVertices...)

}

func (this *RectangleManager) Draw() {
	this.shader.use()
	this.shader.setUniformMatrix4("projection", this.projection)
	gl.Disable(gl.DEPTH_TEST)
	gl.BindVertexArray(this.vao)
	gl.BindBuffer(gl.ARRAY_BUFFER, this.vbo)
	setVerticesInVbo(&this.vertices, &this.vboLen, this.vbo)
	gl.DrawArrays(gl.TRIANGLES, 0, int32(len(this.vertices)))
	gl.Enable(gl.DEPTH_TEST)
}
