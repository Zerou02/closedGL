package ynnebcraft

import (
	"github.com/EngoEngine/glm"
	"github.com/Zerou02/closedGL/closedGL"
)

type Chunk struct {
	origin, size glm.Vec3
	ctx          *closedGL.ClosedGLContext
	//little-endian: ,1bit transparency,6bit faceMask(little oben,vorne,...)
	cubes []uint16
	mesh  closedGL.CubeMesh
}

func NewChunk(origin, size glm.Vec3, ctx *closedGL.ClosedGLContext) Chunk {
	var amountCubes = int(size[0] * size[1] * size[2])
	var cubeArr = make([]uint16, amountCubes)
	var ret = Chunk{origin: origin, size: size, ctx: ctx, cubes: cubeArr}

	closedGL.PrintlnFloat(origin[0])
	ret.setTransparency(0, true)
	ret.setTransparency(1, true)
	ret.setTransparency(2, true)
	ret.setTransparency(3, true)
	ret.setTransparency(4, true)
	ret.setTransparency(closedGL.Pos3ToIdx(1, 1, 1, int(size[0]), int(size[1]), int(size[2])), true)
	ret.setTransparency(closedGL.Pos3ToIdx(2, 2, 2, int(size[0]), int(size[1]), int(size[2])), true)
	ret.faceCullCubes()
	ret.CreateMesh()
	return ret
}

func (this *Chunk) CreateMesh() {
	this.ctx.InitCubeMesh(this.origin, 1)
	for i := 0; i < len(this.cubes); i++ {
		var c = this.cubes[i]
		if !this.isTransparent(c) {
			var x, y, z = closedGL.IdxToPos3(i, int(this.size[0]), int(this.size[1]), int(this.size[2]))
			var faceMask = c & 63
			for j := 0; j < 6; j++ {
				if (faceMask>>j)&1 == 1 {
					this.ctx.DrawCube(glm.Vec3{float32(x), float32(y), float32(z)}, "./assets/sprites/sheet1.png", byte(j), 1, 0+j, 1)
				}
			}
		}
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
func (this *Chunk) faceCullCubes() {
	for i := 0; i < len(this.cubes); i++ {
		if this.isTransparent(this.cubes[i]) {
			continue
		}
		var dimX = int(this.size[0])
		var dimY = int(this.size[1])
		var dimZ = int(this.size[2])

		var allowedFaceMask uint16 = 0
		var posX, posY, posZ = closedGL.IdxToPos3(i, dimX, dimY, dimZ)
		var offsets = []int{
			0, 1, 0,
			0, 0, 1,
			-1, 0, 0,
			1, 0, 0,
			0, 0, -1,
			0, -1, 0,
		}
		for i := 0; i < len(offsets); i += 3 {
			var newX, newY, newZ = posX + offsets[i], posY + offsets[i+1], posZ + offsets[i+2]
			var isOuter = (newX < 0 || newX >= int(this.size[0])) || (newY < 0 || newY >= int(this.size[1])) || (newZ < 0 || newZ >= int(this.size[2]))
			var newIdx = closedGL.Pos3ToIdx(newX, newY, newZ, int(this.size[0]), int(this.size[1]), int(this.size[2]))
			if isOuter {
				allowedFaceMask |= uint16(1) << (i / 3)
			} else {
				_ = newIdx
				var c = this.cubes[newIdx]
				if this.isTransparent(c) {
					allowedFaceMask |= uint16(1) << (i / 3)
				}
			}
		}
		this.cubes[i] <<= 6
		this.cubes[i] |= allowedFaceMask
	}
}
