package ynnebcraft

import (
	"github.com/EngoEngine/glm"
	"github.com/Zerou02/closedGL/closedGL"
)

type GreedyMeshFace struct {
	id            uint
	alreadyMeshed bool
}

type CubeFace struct {
	id        uint
	pos, size [3]int
	side      byte
}

type BufferHolder struct {
	//"up", "front", "left", "right", "back", "down"
	buffer [6][32][1024]GreedyMeshFace
}

type Chunk struct {
	origin, size glm.Vec3
	iSize        [3]int
	ctx          *closedGL.ClosedGLContext
	//little-endian: ,1bit transparency,6bit faceMask(little oben,vorne,...)
	cubes        []uint16
	mesh         closedGL.CubeMesh
	faceBuffer   []CubeFace
	bufferHolder *BufferHolder
}

func NewChunk(origin, size glm.Vec3, ctx *closedGL.ClosedGLContext) Chunk {
	var amountCubes = int(size[0] * size[1] * size[2])
	var cubeArr = make([]uint16, amountCubes)

	var ret = Chunk{origin: origin, size: size, ctx: ctx, cubes: cubeArr,
		faceBuffer:   []CubeFace{},
		bufferHolder: nil,
		iSize:        [3]int{int(size[0]), int(size[1]), int(size[2])},
	}
	ret.setTransparency(0, true)
	ret.setTransparency(1, true)
	ret.setTransparency(2, true)
	ret.setTransparency(3, true)
	ret.setTransparency(4, true)
	ret.setTransparency(closedGL.Pos3ToIdx(1, 1, 1, int(size[0]), int(size[1]), int(size[2])), true)
	ret.setTransparency(closedGL.Pos3ToIdx(1, 2, 1, int(size[0]), int(size[1]), int(size[2])), true)

	ret.setTransparency(closedGL.Pos3ToIdx(2, 2, 2, int(size[0]), int(size[1]), int(size[2])), true)

	ret.CreateMesh()

	return ret
}

func (this *Chunk) greedyMesh() {
	var bufferHolder = BufferHolder{
		buffer: [6][32][1024]GreedyMeshFace{},
	}
	this.bufferHolder = &bufferHolder
	this.greedyMeshPrep()
	for b := 0; b < 6; b++ {
		for i := 0; i < 32; i++ {
			this.greedyMesh2dPlane(&this.bufferHolder.buffer[b][i], i, b)
		}
	}

}
func (this *Chunk) greedyMesh2dPlane(plane *[32 * 32]GreedyMeshFace, sliceID int, side int) {
	var currType uint = 0

	var x, z = -1, 0
	var startX = 0
	var finished = false
	var sideMap = map[string]byte{}
	sideMap["up"] = 0
	sideMap["front"] = 1
	sideMap["left"] = 2
	sideMap["right"] = 3
	sideMap["back"] = 4
	sideMap["down"] = 5
	for !finished {
		x++
		var entry = &plane[closedGL.GridPosToIdx(x, z, 32)]
		if currType == 0 && entry.id != 0 && !entry.alreadyMeshed {
			currType = entry.id
			startX = x
		}
		//mesh
		if (x == 31 || entry.alreadyMeshed || entry.id != currType) && currType != 0 {
			//extend rightward
			//off-by-one hack. Don't know why, don't care
			if x == 31 {
				x++
			}
			var xSteps = x - startX
			var valid = true
			var j = 0
			for valid && j+z < 32 {
				var allSameType = true
				for i := 0; i < xSteps; i++ {
					if plane[closedGL.GridPosToIdx(startX+i, z+j, 32)].id != currType {
						allSameType = false
					}
				}
				valid = allSameType
				if allSameType {
					for i := 0; i < xSteps; i++ {
						plane[closedGL.GridPosToIdx(startX+i, z+j, 32)].alreadyMeshed = true
					}
				}
				if valid {
					j++
				}
			}
			//"up", "front", "left", "right", "back", "down"
			var size = [6][3]int{
				{xSteps, 1, j},
				{xSteps, j, 1},
				{1, j, xSteps},
				{1, j, xSteps},
				{xSteps, j, 1},
				{xSteps, 1, j},
			}
			var pos = [6][3]int{
				{startX, sliceID, z},
				{startX, z, sliceID},
				{sliceID, z, startX},
				{sliceID, z, startX},
				{startX, z, sliceID},
				{startX, sliceID, z},
			}
			var face = CubeFace{
				id:   currType,
				pos:  pos[side],
				size: size[side],
				side: byte(side),
			}
			this.faceBuffer = append(this.faceBuffer, face)
			currType = 0
			x = -1
			z = 0
		}
		if x == 31 {
			x = -1
			z++
		}
		if z == 32 {
			finished = true
		}
	}
}

func (this *Chunk) CreateMesh() {
	this.ctx.InitCubeMesh(this.origin, 1)

	this.ctx.Logger.Start("greedyMesh")
	this.greedyMesh()
	this.ctx.Logger.End("greedyMesh")
	for i := 0; i < len(this.faceBuffer); i++ {
		var f = &this.faceBuffer[i]
		this.ctx.DrawCube(glm.Vec3{float32(f.pos[0]), float32(f.pos[1]), float32(f.pos[2])}, glm.Vec3{float32(f.size[0]), float32(f.size[1]), float32(f.size[2])}, "./assets/sprites/sheet1.png", f.side, 1, 0+int(f.side), 1)

	}
	this.mesh = this.ctx.CopyCurrCubeMesh(1)
	this.faceBuffer = make([]CubeFace, 0)
	this.bufferHolder = nil

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
func (this *Chunk) greedyMeshPrep() {

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
		var nb = [6][3]int{
			{posY, posX, posZ},
			{posZ, posX, posY},
			{posX, posZ, posY},
			{posX, posZ, posY},
			{posZ, posX, posY},
			{posY, posX, posZ},
		}
		for i := 0; i < len(offsets); i += 3 {
			var newX, newY, newZ = posX + offsets[i], posY + offsets[i+1], posZ + offsets[i+2]
			var isOuter = (newX < 0 || newX >= this.iSize[0]) || (newY < 0 || newY >= this.iSize[1]) || (newZ < 0 || newZ >= this.iSize[2])
			var newIdx = closedGL.Pos3ToIdx(newX, newY, newZ, this.iSize[0], this.iSize[1], this.iSize[2])
			var otherTransparent = false
			if !isOuter {
				var c = this.cubes[newIdx]
				otherTransparent = this.isTransparent(c)
			}

			if isOuter || otherTransparent {
				allowedFaceMask |= (uint16(1) << (i / 3))
				var face = GreedyMeshFace{
					id:            1,
					alreadyMeshed: false,
				}
				this.bufferHolder.buffer[i/3][nb[i/3][0]][closedGL.GridPosToIdx(nb[i/3][1], nb[i/3][2], 32)] = face
			}
		}
		this.cubes[i] <<= 6
		this.cubes[i] |= allowedFaceMask
	}
}
