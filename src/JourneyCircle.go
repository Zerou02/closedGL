package main

import (
	closed_gl "closed_gl/src/closedGL"

	"github.com/EngoEngine/glm"
)

type JourneyCircle struct {
	Position closed_gl.Vec2
	Radius   float32

	CentreCircle closed_gl.Circle
	RadiusCircle closed_gl.Circle
}

func newJourneyCircle(pos closed_gl.Vec2, radius float32, ctx closed_gl.ClosedGLContext) JourneyCircle {
	return JourneyCircle{Position: pos, Radius: radius,
		CentreCircle: ctx.Factory.NewCircle(glm.Vec4{1, 0, 0, 1}, glm.Vec4{1, 0, 0, 0}, 10, pos, 5),
		RadiusCircle: ctx.Factory.NewCircle(glm.Vec4{0, 0, 0, 0}, glm.Vec4{1, 0, 0, 0}, radius, pos, 10),
	}
}

func (this *JourneyCircle) MakeInvis() {
	this.CentreCircle.BorderColour = glm.Vec4{0, 0, 0, 0}
	this.CentreCircle.CentreColour = glm.Vec4{0, 0, 0, 0}
	this.RadiusCircle.BorderColour = glm.Vec4{0, 0, 0, 0}
	this.RadiusCircle.CentreColour = glm.Vec4{0, 0, 0, 0}
	this.CentreCircle.BorderThickness = 0
	this.RadiusCircle.BorderThickness = 0

	this.Radius = 0
}

func (this *JourneyCircle) process() {
	this.CentreCircle.Centre = this.Position
	this.RadiusCircle.Centre = this.Position
	this.RadiusCircle.Radius = this.Radius

}
func (this *JourneyCircle) draw() {
	this.process()
	this.CentreCircle.Draw()
	this.RadiusCircle.Draw()
}
