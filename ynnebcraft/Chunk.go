package ynnebcraft

import (
	"github.com/EngoEngine/glm"
	"github.com/Zerou02/closedGL/closedGL"
)

type Chunk struct {
	origin, size glm.Vec3
	ctx          *closedGL.ClosedGLContext
}

func NewChunk(origin, size glm.Vec3, ctx *closedGL.ClosedGLContext) Chunk {
	return Chunk{origin: origin, size: size, ctx: ctx}

}

func (this *Chunk) Draw() {
	for i := 0; i < int(this.size[0]); i++ {
		for j := 0; j < int(this.size[1]); j++ {
			for k := 0; k < int(this.size[2]); k++ {
				this.ctx.DrawCube(glm.Vec3{float32(i) + this.origin[0], float32(j) + this.origin[1], float32(k) + this.origin[2]}, "./assets/sprites/fence.png", 1)
			}
		}
	}
}
