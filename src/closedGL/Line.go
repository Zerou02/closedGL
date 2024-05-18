package closed_gl

import (
	"github.com/EngoEngine/glm"
	"github.com/go-gl/gl/v4.1-core/gl"
)

type Line struct {
	shader     *Shader
	Points     []float32
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
	generateBuffers(&this.vao, &this.vbo, nil, this.Points, 0, nil, []int{2, 3})
}

func (this *Line) Draw() {
	if len(this.Points)/5 < 1 {
		return
	}
	this.shader.use()
	this.shader.setUniformMatrix4("projection", this.projection)
	gl.BindVertexArray(this.vao)
	gl.BindBuffer(gl.ARRAY_BUFFER, this.vbo)
	gl.BufferSubData(gl.ARRAY_BUFFER, 0, len(this.Points)*4, gl.Ptr(this.Points))

	gl.DrawArrays(gl.LINES, 0, int32(len(this.Points)/5))
}

func (this *Line) Rehydrate() {
	this.generateBuffers()
}

func (this *Line) addPoint(pos Vec2, colour glm.Vec3) {
	this.Points = append(this.Points, pos[0])
	this.Points = append(this.Points, pos[1])
	this.Points = append(this.Points, colour[0])
	this.Points = append(this.Points, colour[1])
	this.Points = append(this.Points, colour[2])
	this.generateBuffers()
}
