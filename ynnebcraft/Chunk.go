package ynnebcraft

import (
	"github.com/EngoEngine/glm"
	"github.com/Zerou02/closedGL/closedGL"
)

type CubeFace struct {
	id            uint
	pos           glm.Vec3
	sizeX, sizeY  uint
	side          byte
	alreadyMeshed bool
}

type Chunk struct {
	origin, size glm.Vec3
	ctx          *closedGL.ClosedGLContext
	//little-endian: ,1bit transparency,6bit faceMask(little oben,vorne,...)
	cubes       []uint16
	mesh        closedGL.CubeMesh
	columns     map[string]*[32 * 32][32]CubeFace
	faceBuffer  []CubeFace
	upBuffer    [32][1024]CubeFace
	downBuffer  [32][1024]CubeFace
	leftBuffer  [32][1024]CubeFace
	rightBuffer [32][1024]CubeFace
	frontBuffer [32][1024]CubeFace
	backBuffer  [32][1024]CubeFace
}

func NewChunk(origin, size glm.Vec3, ctx *closedGL.ClosedGLContext) Chunk {
	var amountCubes = int(size[0] * size[1] * size[2])
	var cubeArr = make([]uint16, amountCubes)
	var ret = Chunk{origin: origin, size: size, ctx: ctx, cubes: cubeArr,
		columns:     map[string]*[32 * 32][32]CubeFace{},
		faceBuffer:  []CubeFace{},
		upBuffer:    [32][1024]CubeFace{},
		leftBuffer:  [32][1024]CubeFace{},
		rightBuffer: [32][1024]CubeFace{},
		frontBuffer: [32][1024]CubeFace{},
		backBuffer:  [32][1024]CubeFace{},
	}
	var keys = []string{"up", "front", "left", "right", "back", "down"}
	for _, x := range keys {
		var arr = [32 * 32][32]CubeFace{}
		ret.columns[x] = &arr
	}

	ret.setTransparency(0, true)
	ret.setTransparency(1, true)
	ret.setTransparency(2, true)
	ret.setTransparency(3, true)
	ret.setTransparency(4, true)
	ret.setTransparency(closedGL.Pos3ToIdx(1, 1, 1, int(size[0]), int(size[1]), int(size[2])), true)
	ret.setTransparency(closedGL.Pos3ToIdx(1, 2, 1, int(size[0]), int(size[1]), int(size[2])), true)

	ret.setTransparency(closedGL.Pos3ToIdx(2, 2, 2, int(size[0]), int(size[1]), int(size[2])), true)
	ret.faceCullCubes()

	//	ret.greedyMesh2dMesh()

	for _, x := range ret.frontBuffer[0] {
		if x.pos[0] != 0 && x.pos[1] != 0 && x.pos[2] != 0 {

			closedGL.PrintlnFloat(x.pos[0])
			closedGL.PrintlnFloat(x.pos[1])
			closedGL.PrintlnFloat(x.pos[2])
		}
	}
	//var up = ret.createUpDownVertices("up")
	//var down = ret.createUpDownVertices("down")
	//var left = ret.createLeftRightVertices("left")
	//var right = ret.createLeftRightVertices("right")

	//var front = ret.createFrontBackVertices()
	for i := 0; i < 32; i++ {
		//	ret.greedyMesh2dPlane(up[i])
		//	ret.greedyMesh2dPlane(down[i])
		//	ret.greedyMesh2dPlane(down[i])
	}
	ret.greedyMesh2dPlane(ret.frontBuffer[0])

	/* for _, x := range down[0] {
		if x.pos[0] != 0 && x.pos[1] != 0 && x.pos[2] != 0 {
			println(x.pos[0])
			println(x.pos[1])
			println(x.pos[2])
		}
	} */

	println("len", len(ret.faceBuffer))
	for _, x := range ret.faceBuffer {
		println(x.pos[0])
		println(x.pos[1])
		println(x.pos[2])

		println(x.sizeX)
		println(x.sizeY)

	}
	ret.CreateMesh()

	return ret
}

func (this *Chunk) createUpDownVertices(dir string) [32][32 * 32]CubeFace {
	var faces = [32][32 * 32]CubeFace{}
	for y := 0; y < 32; y++ {
		for i := 0; i < 32; i++ {
			for x := 0; x < 32; x++ {
				faces[y][i*32+x] = this.columns[dir][closedGL.GridPosToIdx(x, i, 32)][y]
			}
		}
	}
	return faces
}

func (this *Chunk) createLeftRightVertices(dir string) [32][32 * 32]CubeFace {
	var faces = [32][32 * 32]CubeFace{}
	for i := 0; i < 32; i++ {
		for y := 0; y < 32; y++ {
			for x := 0; x < 32; x++ {
				faces[i][y*32+x] = this.columns[dir][y*32+i][x]
			}
		}
	}
	return faces
}

func (this *Chunk) createFrontBackVertices() [32][32 * 32]CubeFace {
	//i = z;
	//j = y
	//k = x
	//oder so
	var faces = [32][32 * 32]CubeFace{}
	for i := 0; i < 32; i++ {
		for j := 0; j < 32; j++ {
			for k := 0; k < 32; k++ {
				faces[j][j*32+k] = this.columns["back"][i*32+k][j]
			}
		}
	}
	return faces
}

func (this *Chunk) greedyMesh2dPlane(plane [32 * 32]CubeFace) {
	var currType uint = 0

	var x, z = -1, 0
	var startX = 0
	var finished = false
	for !finished {
		x++
		var entry = plane[closedGL.GridPosToIdx(x, z, 32)]
		if currType == 0 && entry.id != 0 && !entry.alreadyMeshed {
			currType = entry.id
			startX = x
		}
		//mesh
		if (x == 31 || entry.alreadyMeshed || entry.id != currType) && currType != 0 {
			println(startX, x, z)
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
			this.faceBuffer = append(this.faceBuffer, CubeFace{
				id:            currType,
				pos:           glm.Vec3{float32(startX), 0, float32(z)},
				sizeX:         uint(xSteps),
				sizeY:         uint(j),
				side:          0,
				alreadyMeshed: true,
			})
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

func (this *Chunk) greedyMesh2dMesh() {
	this.faceBuffer = []CubeFace{}
	var sides = this.columns["up"]
	var currType = 0
	var x, y, z = -1, 0, 0
	var stepsX = 0

	var finished = false
	for !finished {
		x++
		stepsX++

		var next = &sides[closedGL.GridPosToIdx(x, z, 32)][y]
		println(x, z, y, currType, next.alreadyMeshed)

		if currType == 0 && next.id != 0 && !next.alreadyMeshed {
			currType = int(next.id)
			stepsX = 0
			x = int(next.pos[0])
			z = int(next.pos[2])
			println("curr", x, z, y)
			continue
		}
		if z == 31 && x == 31 {
			if currType != 0 {
				next.alreadyMeshed = true

				this.faceBuffer = append(this.faceBuffer, CubeFace{
					id:            uint(currType),
					pos:           glm.Vec3{float32(x), float32(y), float32(z)},
					sizeX:         uint(stepsX),
					sizeY:         uint(1),
					side:          0,
					alreadyMeshed: true,
				})
			}
			if y == 31 {
				finished = true
			} else {
				y++
				x = -1
				z = 0
				stepsX = 0
				currType = 0
			}
		}
		if (int(next.id) != currType || x == 31 || next.alreadyMeshed) && currType != 0 {
			println("i", stepsX, z)
			var i = 0 //stepsZ
			var valid = true
			var min = stepsX
			for valid {
				var innerValid = true
				for j := 0; j < min; j++ {
					if int(sides[closedGL.GridPosToIdx(x-stepsX+j, z+i, 32)][y].id) != currType {
						innerValid = false
					}
				}
				if innerValid {
					println("inner", i)
					for j := 0; j <= stepsX; j++ {
						sides[closedGL.GridPosToIdx(x-stepsX+j, z+i, 32)][y].alreadyMeshed = true
					}
				}
				println(i, min)
				i++
				if !innerValid || i == 32 || i == min {
					valid = false
					break
				}

			}
			this.faceBuffer = append(this.faceBuffer, CubeFace{
				id:            uint(currType),
				pos:           glm.Vec3{float32(x - stepsX), float32(y), float32(z)},
				sizeX:         uint(stepsX),
				sizeY:         uint(i),
				side:          0,
				alreadyMeshed: true,
			})
			x = -1
			currType = 0
			stepsX = 0
			z++
			/* 			println(z, i)
			   			if z+i >= 32 {
			   				finished = true
			   			} */
		}
		if x == 31 {
			x = -1
			z++
		}
	}

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
		var bitToFaceMap = []string{"up", "front", "left", "right", "back", "down"}
		for i := 0; i < len(offsets); i += 3 {
			var newX, newY, newZ = posX + offsets[i], posY + offsets[i+1], posZ + offsets[i+2]
			var isOuter = (newX < 0 || newX >= int(this.size[0])) || (newY < 0 || newY >= int(this.size[1])) || (newZ < 0 || newZ >= int(this.size[2]))
			var newIdx = closedGL.Pos3ToIdx(newX, newY, newZ, int(this.size[0]), int(this.size[1]), int(this.size[2]))
			if isOuter {
				allowedFaceMask |= (uint16(1) << (i / 3))
				var entry = bitToFaceMap[i/3]
				var sideColumns = &(*(this.columns[bitToFaceMap[i/3]]))
				var column = &sideColumns[closedGL.GridPosToIdx(posX, posZ, 32)]
				var face = CubeFace{
					pos:           glm.Vec3{float32(posX), float32(posY), float32(posZ)},
					id:            1,
					side:          byte((uint16(1) << (i / 3))),
					alreadyMeshed: false,
				}
				column[posY] = face
				if entry == "up" {
					this.upBuffer[posY][closedGL.GridPosToIdx(posX, posZ, 32)] = face
				} else if entry == "down" {
					this.downBuffer[posY][closedGL.GridPosToIdx(posX, posZ, 32)] = face
				} else if entry == "left" {
					this.leftBuffer[posX][closedGL.GridPosToIdx(posZ, posY, 32)] = face
				} else if entry == "right" {
					this.rightBuffer[posX][closedGL.GridPosToIdx(posZ, posY, 32)] = face
				} else if entry == "front" {
					this.frontBuffer[posZ][closedGL.GridPosToIdx(posX, posY, 32)] = face
				} else if entry == "back" {
					this.backBuffer[posZ][closedGL.GridPosToIdx(posX, posY, 32)] = face
				}
			} else {
				_ = newIdx
				var c = this.cubes[newIdx]
				if this.isTransparent(c) {
					allowedFaceMask |= (uint16(1) << (i / 3))
					var entry = bitToFaceMap[i/3]
					var sideColumns = &(*(this.columns[bitToFaceMap[i/3]]))
					var column = &sideColumns[closedGL.GridPosToIdx(posX, posZ, 32)]
					var face = CubeFace{
						pos:           glm.Vec3{float32(posX), float32(posY), float32(posZ)},
						id:            1,
						side:          byte((uint16(1) << (i / 3))),
						alreadyMeshed: false,
					}
					column[posY] = face
					if entry == "up" {
						this.upBuffer[posY][closedGL.GridPosToIdx(posX, posZ, 32)] = face
					} else if entry == "down" {
						this.downBuffer[posY][closedGL.GridPosToIdx(posX, posZ, 32)] = face
					} else if entry == "left" {
						this.leftBuffer[posX][closedGL.GridPosToIdx(posZ, posY, 32)] = face
					} else if entry == "right" {
						this.rightBuffer[posX][closedGL.GridPosToIdx(posZ, posY, 32)] = face
					} else if entry == "front" {
						this.frontBuffer[posZ][closedGL.GridPosToIdx(posX, posY, 32)] = face
					} else if entry == "back" {
						this.backBuffer[posZ][closedGL.GridPosToIdx(posX, posY, 32)] = face
					}
				}
			}
		}
		this.cubes[i] <<= 6
		this.cubes[i] |= allowedFaceMask
	}
}
