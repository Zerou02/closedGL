package ynnebcraft

import (
	"github.com/EngoEngine/glm"
	"github.com/Zerou02/closedGL/closedGL"
)

type Chunk struct {
	origin, size glm.Vec3
	ctx          *closedGL.ClosedGLContext
	cubes        []uint16
	mesh         closedGL.CubeMesh
}

func NewChunk(origin, size glm.Vec3, ctx *closedGL.ClosedGLContext) Chunk {
	var amountCubes = int(size[0] * size[1] * size[2])
	var cubeArr = make([]uint16, amountCubes)
	var ret = Chunk{origin: origin, size: size, ctx: ctx, cubes: cubeArr}
	for i := 0; i < len(cubeArr); i++ {
		if ret.isInnerBlock(i) {
			cubeArr[i] = 1
		} else {
			cubeArr[i] = 0
		}
	}
	ret.CreateMesh()
	return ret
}

func (this *Chunk) CreateMesh() {
	this.ctx.InitCubeMesh(1)
	for i := 0; i < len(this.cubes); i++ {
		var c = this.cubes[i]
		if c&1 == 0 {
			var x, y, z = closedGL.IdxToPos3(i, int(this.size[0]), int(this.size[1]), int(this.size[2]))
			this.ctx.DrawCube(glm.Vec3{float32(x) + this.origin[0], float32(y) + this.origin[1], float32(z) + this.origin[2]}, "./assets/sprites/fence_small.png", 1)
		}
	}
	this.mesh = this.ctx.CopyCurrCubeMesh(1)
}

func (this *Chunk) Draw() {
	this.ctx.DrawCubeMesh(&this.mesh, 1)
}

// surrounded on all sides
func (this *Chunk) isInnerBlock(idx int) bool {
	var neighbours = this.getAmountNeighbours(idx)
	var posX, posY, posZ = closedGL.IdxToPos3(idx, int(this.size[0]), int(this.size[1]), int(this.size[2]))

	var isInner = posX > 0 && posX < int(this.size[0])-1 && posY > 0 && posY < int(this.size[1])-1 && posZ > 0 && posZ < int(this.size[2])-1
	return isInner && neighbours >= 4
}

// No Diagonal Neighbours
func (this *Chunk) getAmountNeighbours(idx int) int {
	var retAmount = 0
	var offsets = []int{
		-1, 0, 0,
		1, 0, 0,
		0, 1, 0,
		0, -1, 0,
		0, 0, 1,
		0, 0, -1,
	}

	var posX, posY, posZ = closedGL.IdxToPos3(idx, int(this.size[0]), int(this.size[1]), int(this.size[2]))

	for i := 0; i < len(offsets); i += 3 {
		var newX, newY, newZ = posX + offsets[i], posY + offsets[i+1], posZ + offsets[i+2]

		var idx = closedGL.Pos3ToIdx(int(newX), int(newY), int(newZ), int(this.size[0]), int(this.size[1]), int(this.size[2]))
		if idx >= 0 && idx < len(this.cubes) {
			retAmount += 1
		}
	}
	return retAmount
}
