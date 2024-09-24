package nimblemossbilliard

import (
	"github.com/EngoEngine/glm"
	"github.com/Zerou02/closedGL/closedGL"
)

type Circle struct {
	pos    glm.Vec2
	vel    glm.Vec2
	radius float32
	colour glm.Vec4
}

func (this *Circle) process(delta float32, speed float32, grav float32) {
	this.move(delta, speed, grav)
}

func (this *Circle) move(delta float32, speed float32, grav float32) {
	this.vel[1] += grav * delta
	this.pos.AddScaledVec(delta*speed, &this.vel)

}

func (this *Circle) drawInto(mesh *closedGL.CircleMesh) {
	mesh.AddCircle(this.pos, this.colour, this.colour, this.radius, 0)
}

func (this *Circle) getCentre() glm.Vec2 {
	return glm.Vec2{this.pos[0] + this.radius/2, this.pos[1] + this.radius/2}
}

func newCircle(pos, vel glm.Vec2, radius float32, colour glm.Vec4) Circle {
	return Circle{
		pos:    pos,
		vel:    vel,
		radius: radius,
		colour: colour,
	}
}
