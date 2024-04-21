package main

import (
	"github.com/EngoEngine/glm"
	"github.com/go-gl/gl/v4.1-core/gl"
)

type Sprite2D struct {
	shader        *Shader
	texture       *uint32
	pos, size     glm.Vec2
	tint          glm.Vec4
	vao, vbo, ebo uint32
	visible       bool
}

func newSprite2D(s *Shader, tex *uint32, pos, size glm.Vec2, tint glm.Vec4) Sprite2D {
	var sprite = Sprite2D{shader: s, texture: tex, pos: pos, size: size, tint: tint, visible: true}
	generateBuffers(&sprite.vao, &sprite.vbo, &sprite.ebo, fullQuad, 0, indicesQuad, []VertexInfo{{3, 0}, {2, 12}})
	return sprite
}

func (this *Sprite2D) process(delta float32) {
	/* var a = this.velocity.Mul(delta) */
	/* var b = a.Mul(this.speed) */
	/* this.move(b) */
}

func (this *Sprite2D) draw() {
	this.shader.use()
	var model = createTransformation(glm.Vec3{0, 0, 0}, glm.Vec3{this.pos.X(), this.pos.Y(), 0}, glm.Vec3{this.size.X(), this.size.Y(), 0})
	this.shader.setUniformMatrix4("model", &model)
	this.shader.setUniformVec4("colour", &this.tint)
	gl.ActiveTexture(gl.TEXTURE0)
	gl.BindTexture(gl.TEXTURE_2D, *this.texture)
	this.shader.setUniform1i("tex", 0)

	gl.BindVertexArray(this.vao)
	//gl.BindBuffer(gl.ARRAY_BUFFER, this.vbo)
	//gl.BindBuffer(gl.ELEMENT_ARRAY_BUFFER, this.ebo)
	gl.DrawElements(gl.TRIANGLES, 6, gl.UNSIGNED_INT, nil)
}

func (this *Sprite2D) move(vec glm.Vec2) {
	this.pos.AddWith(&vec)
}

func (this *Sprite2D) moveTo(vec glm.Vec2) {
	this.pos = vec
}

func (this *Sprite2D) colAabb(s Sprite2D) bool {
	return aabbAabbCol(glm.Vec4{this.pos[0], this.pos[1], this.size[0], this.size[1]}, glm.Vec4{s.pos[0], s.pos[1], s.size[0], s.size[1]})
}

func (this *Sprite2D) colCircle(s Sprite2D) (bool, Direction, glm.Vec2) {
	return aabbCircleCol(glm.Vec3{s.pos[0], s.pos[1], s.size[1] - 25}, glm.Vec4{this.pos[0], this.pos[1], this.size[0], this.size[1]})
}
