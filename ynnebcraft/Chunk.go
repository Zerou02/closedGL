package ynnebcraft

/*

import (
	"github.com/EngoEngine/glm"
	"github.com/Zerou02/closedGL/closedGL"
)

type Chunk struct {
	origin, size glm.Vec3
	iSize        [3]int
	ctx          *closedGL.ClosedGLContext
	//little-endian: ,1bit transparency,6bit faceMask(little oben,vorne,...)
	cubes []uint16
	mesh  closedGL.CubeMesh
}

func NewChunk(origin, size glm.Vec3, ctx *closedGL.ClosedGLContext, mesher *GreedyMesher) Chunk {
	var amountCubes = int(size[0] * size[1] * size[2])
	var cubeArr = make([]uint16, amountCubes)

	var ret = Chunk{origin: origin, size: size, ctx: ctx, cubes: cubeArr,
		iSize: [3]int{int(size[0]), int(size[1]), int(size[2])},
	}
	ret.setTransparency(0, true)
	ret.setTransparency(1, true)
	ret.setTransparency(2, true)
	ret.setTransparency(3, true)
	ret.setTransparency(4, true)
	ret.setTransparency(closedGL.Pos3ToIdx(1, 1, 1, int(size[0]), int(size[1]), int(size[2])), true)
	ret.setTransparency(closedGL.Pos3ToIdx(1, 2, 1, int(size[0]), int(size[1]), int(size[2])), true)
	ret.setTransparency(closedGL.Pos3ToIdx(2, 2, 2, int(size[0]), int(size[1]), int(size[2])), true)

	ret.CreateMesh(mesher)

	return ret
}

func (this *Chunk) CreateMesh(mesher *GreedyMesher) {
	this.ctx.Logger.Start("meshing")
	var buffer = mesher.mesh(this)
	this.ctx.Logger.End("meshing")

	this.ctx.InitCubeMesh(this.origin, 1)
	for i := 0; i < len(*buffer); i++ {
		var f = &(*buffer)[i]
		this.ctx.DrawCube(glm.Vec3{float32(f.pos[0]), float32(f.pos[1]), float32(f.pos[2])}, glm.Vec3{float32(f.size[0]), float32(f.size[1]), float32(f.size[2])}, "./assets/sprites/sheet1.png", f.side, 1, 0+int(f.side), 1)

	}
	this.mesh = this.ctx.CopyCurrCubeMesh(1)
}

func (this *Chunk) Draw() {
	this.ctx.DrawCubeMesh(&this.mesh, 1)
}

func (this *Chunk) isTransparent(cube uint16) bool {
	return (cube>>6)&1 == 1
}

func (this *Chunk) setTransparency(idx int, val bool) {
	var a uint16 = 1
	if !val {
		a = 0
	}
	this.cubes[idx] |= a << 6
}
*/
