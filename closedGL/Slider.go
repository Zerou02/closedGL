package closedGL

import "github.com/EngoEngine/glm"

type Quad struct {
	dim    glm.Vec4
	colour glm.Vec4
}
type Slider struct {
	ctx  *ClosedGLContext
	mesh *RectangleMesh
	bg   Quad
	knob Quad
	dim  glm.Vec4
	drag bool
}

func NewSlider(mesh *RectangleMesh, ctx *ClosedGLContext, dim glm.Vec4) Slider {
	var slider = Slider{
		ctx:  ctx,
		mesh: mesh,
		dim:  dim,
	}
	slider.bg = Quad{dim: dim, colour: glm.Vec4{0.5, 0.5, 0.5, 1}}
	slider.knob = Quad{dim: glm.Vec4{dim[0] + dim[2]/2, dim[1], 10, dim[3]}, colour: glm.Vec4{1, 1, 1, 1}}
	mesh.AddQuad(&slider.bg)
	mesh.AddQuad(&slider.knob)
	return slider
}

func (this *Slider) Process() {
	if this.mesh.isUpdate() {
		this.mesh.AddQuad(&this.bg)
		this.mesh.AddQuad(&this.knob)
	}
	if this.ctx.IsMouseDown() && this.ctx.IsMouseInRect(this.dim) {
		this.drag = true
	}
	if this.drag {
		var newX = this.ctx.GetMousePos()[0]
		newX = Clamp(this.dim[0], this.dim[0]+this.dim[2]-this.knob.dim[2], newX)
		this.knob.dim = glm.Vec4{newX, this.knob.dim[1], this.knob.dim[2], this.knob.dim[3]}
		this.mesh.setDirty()
	}
	if !this.ctx.IsMouseDown() && this.drag {
		this.drag = false
	}
}

func (this *Slider) GetPercentage() float32 {
	return CalcPercentage(this.bg.dim[0], this.bg.dim[0]+this.bg.dim[2], this.knob.dim[0])
}
