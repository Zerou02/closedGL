package closedGL

/*

import (
	"github.com/EngoEngine/glm"
	"github.com/go-gl/gl/v4.1-core/gl"
)

type Button struct {
	pos, size  glm.Vec2
	text       string
	BaseTex    Sprite2D
	OutlineTex Sprite2D
	isHover    bool
	ctx        *ClosedGLContext
}

func NewBtn(pos, size glm.Vec2, basePath, outlinePath string, text string, ctx *ClosedGLContext) Button {
	var ret = Button{
		pos:        pos,
		size:       size,
		text:       text,
		BaseTex:    factory.NewSprite2D(pos, size, glm.Vec4{1, 1, 1, 1}, basePath),
		OutlineTex: factory.NewSprite2D(pos, size, glm.Vec4{1, 1, 1, 1}, outlinePath),
		isHover:    true,
		ctx:        ctx,
	}
	return ret
}

func (this *Button) Draw() {
	gl.Disable(gl.DEPTH_TEST)
	this.BaseTex.Draw()
	if this.isHover {
		this.OutlineTex.Draw()
	}
	gl.Enable(gl.DEPTH_TEST)

	text.DrawText(int(this.pos[0]), int(this.pos[1]+0.3*this.size[1]), this.text)
}

func (this *Button) Process() {
	var x, y = this.ctx.Window.Window.GetCursorPos()
	if float32(x) >= this.pos[0] && float32(x) <= this.pos[0]+this.size[0] &&
		float32(y) >= this.pos[1] && float32(y) <= this.pos[1]+this.size[1] {
		this.isHover = true
	} else {
		this.isHover = false
	}
}
*/
