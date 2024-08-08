package ynnebcraft

import "github.com/Zerou02/closedGL/closedGL"

type GreedyMeshFace struct {
	id            uint
	alreadyMeshed bool
}

type BufferHolder struct {
	//"up", "front", "left", "right", "back", "down"
	buffer [6][32][1024]GreedyMeshFace
}

type CubeFace struct {
	id        uint
	pos, size [3]int
	side      byte
}

type GreedyMesher struct {
	faceBuffer   []CubeFace
	bufferHolder *BufferHolder
}

func NewGreedyMesher() GreedyMesher {
	var bufferHolder = BufferHolder{
		buffer: [6][32][1024]GreedyMeshFace{},
	}
	return GreedyMesher{
		faceBuffer:   []CubeFace{},
		bufferHolder: &bufferHolder,
	}
}

func (this *GreedyMesher) mesh(chunk *Chunk) []CubeFace {
	this.faceCullCubes(chunk)
	this.greedyMesh()
	return this.faceBuffer
}

func (this *GreedyMesher) faceCullCubes(chunk *Chunk) {

	for i := 0; i < len(chunk.cubes); i++ {
		if chunk.isTransparent(chunk.cubes[i]) {
			continue
		}
		var dimX = int(chunk.size[0])
		var dimY = int(chunk.size[1])
		var dimZ = int(chunk.size[2])

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
			var isOuter = (newX < 0 || newX >= chunk.iSize[0]) || (newY < 0 || newY >= chunk.iSize[1]) || (newZ < 0 || newZ >= chunk.iSize[2])
			var newIdx = closedGL.Pos3ToIdx(newX, newY, newZ, chunk.iSize[0], chunk.iSize[1], chunk.iSize[2])
			var otherTransparent = false
			if !isOuter {
				var c = chunk.cubes[newIdx]
				otherTransparent = chunk.isTransparent(c)
			}

			var face GreedyMeshFace

			if isOuter || otherTransparent {
				allowedFaceMask |= (uint16(1) << (i / 3))
				face = GreedyMeshFace{
					id:            1,
					alreadyMeshed: false,
				}
			}

			this.bufferHolder.buffer[i/3][nb[i/3][0]][closedGL.GridPosToIdx(nb[i/3][1], nb[i/3][2], 32)] = face
		}
		chunk.cubes[i] <<= 6
		chunk.cubes[i] |= allowedFaceMask
	}
}

func (this *GreedyMesher) greedyMesh() {
	for b := 0; b < 6; b++ {
		for i := 0; i < 32; i++ {
			this.greedyMesh2dPlane(&this.bufferHolder.buffer[b][i], i, b)
		}
	}
}

func (this *GreedyMesher) greedyMesh2dPlane(plane *[32 * 32]GreedyMeshFace, sliceID int, side int) {
	var currType uint = 0

	var x, z = -1, 0
	var startX = 0
	var finished = false

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
