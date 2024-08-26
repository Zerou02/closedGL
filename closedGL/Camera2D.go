package closedGL

import (
	"github.com/EngoEngine/glm"
	"github.com/go-gl/glfw/v3.2/glfw"
)

type Camera2D struct {
	projMat, ViewMat glm.Mat4
	pos              glm.Vec2
	ctx              *ClosedGLContext
	Speed            float32
}

func newCamera2D(ww, wh float32, ctx *ClosedGLContext) Camera2D {
	return Camera2D{
		projMat: glm.Ortho2D(0, ww, wh, 0),
		ViewMat: glm.Ident4(),
		ctx:     ctx,
		Speed:   100,
	}
}

func (this *Camera2D) Process(delta float32) {
	if this.ctx.IsKeyDown(glfw.KeyD) {
		this.pos[0] += this.Speed * delta
	}
	if this.ctx.IsKeyDown(glfw.KeyA) {
		this.pos[0] -= this.Speed * delta
	}
	if this.ctx.IsKeyDown(glfw.KeyW) {
		this.pos[1] -= this.Speed * delta
	}
	if this.ctx.IsKeyDown(glfw.KeyS) {
		this.pos[1] += this.Speed * delta
	}
	this.ViewMat = glm.Translate3D(-this.pos[0], -this.pos[1], 0)
}
