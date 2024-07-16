package closedGL

/*
import (
	"github.com/EngoEngine/glm"
	"github.com/go-gl/gl/v4.1-core/gl"
)

type Sprite2D struct {
	shader    *Shader
	texture   *uint32
	pos, size glm.Vec2
	Tint      glm.Vec4
	vao, vbo  uint32
	Visible   bool
}

func newSprite2D(s *Shader, pos, size glm.Vec2, tint glm.Vec4, texturePath string) Sprite2D {
	var tex = LoadImage(texturePath, gl.RGBA)
	var sprite = Sprite2D{shader: s, texture: tex, pos: pos, size: size, Tint: tint, Visible: true}
	generateBuffers(&sprite.vao, &sprite.vbo, nil, fullQuad, 0, indicesQuad, []int{3, 2})
	return sprite
}

func (this *Sprite2D) Process(delta float32) {
	 var a = this.velocity.Mul(delta)
	 var b = a.Mul(this.speed)
	 this.move(b)
}

func (this *Sprite2D) Free() {
	gl.DeleteBuffers(1, &this.vao)
	gl.DeleteBuffers(1, &this.vbo)
	gl.DeleteTextures(1, this.texture)

}
func (this *Sprite2D) Draw() {
	if !this.Visible {
		return
	}
	this.shader.use()
	var model = createTransformation(glm.Vec3{0, 0, 0}, glm.Vec3{this.pos.X(), this.pos.Y(), 0}, glm.Vec3{this.size.X(), this.size.Y(), 0})
	this.shader.setUniformMatrix4("model", &model)
	this.shader.setUniformVec4("colour", &this.Tint)
	gl.ActiveTexture(gl.TEXTURE0)
	gl.BindTexture(gl.TEXTURE_2D, *this.texture)
	this.shader.setUniform1i("tex", 0)
	gl.BindVertexArray(this.vao)
	gl.BindBuffer(gl.ARRAY_BUFFER, this.vbo)
	gl.DrawArrays(gl.TRIANGLES, 0, 6)
}

func (this *Sprite2D) Move(vec glm.Vec2) {
	this.pos.AddWith(&vec)
}

func (this *Sprite2D) MoveTo(vec glm.Vec2) {
	this.pos = vec
}

func (this *Sprite2D) colAabb(s Sprite2D) bool {
	return aabbAabbCol(glm.Vec4{this.pos[0], this.pos[1], this.size[0], this.size[1]}, glm.Vec4{s.pos[0], s.pos[1], s.size[0], s.size[1]})
}

func (this *Sprite2D) colCircle(s Sprite2D) (bool, Direction, glm.Vec2) {
	return aabbCircleCol(glm.Vec3{s.pos[0], s.pos[1], s.size[1] - 25}, glm.Vec4{this.pos[0], this.pos[1], this.size[0], this.size[1]})
}
*/
