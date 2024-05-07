package main

import (
	"github.com/EngoEngine/glm"
	"github.com/go-gl/gl/v4.1-core/gl"
)

//pos:x,y,z

var cubeVertexStride = 3
var cummVertexStride = 8

var verticesPerCube = 36
var dim = 16
var cumulatedVertices2 = make([]uint32, verticesPerCube*(cummVertexStride)*dim*dim*dim)

type CubeVertex struct {
	isInner bool
	texId   int
}

type Chunk struct {
	shader      *Shader
	camera      *Camera
	projection  *glm.Mat4
	tex         *Texture
	dim, pos    glm.Vec3
	cubes       []CubeVertex
	amountOuter uint32
	vao, vbo    uint32
}

func newChunk(dim, pos glm.Vec3, tex *Texture, camera *Camera, projection *glm.Mat4, shader *Shader) Chunk {
	profiler.startTime("init")
	var chunk = Chunk{dim: dim, pos: pos, camera: camera, projection: projection, tex: tex, shader: shader}
	var amountCubes = int(dim[0] * dim[1] * dim[2])
	var cubeArr = make([]CubeVertex, amountCubes)
	var count = 0
	var texId = 0
	for y := 0; y < int(dim[1]); y++ {
		for z := 0; z < int(dim[2]); z++ {
			for x := 0; x < int(dim[0]); x++ {
				var newVertex = CubeVertex{isInner: false, texId: texId}
				cubeArr[count] = newVertex
				count += 1
			}
		}
	}
	chunk.cubes = cubeArr
	profiler.endTime("init")
	profiler.startTime("vbo")
	//copy into vbo
	chunk.createVBO()
	profiler.endTime("vbo")

	return chunk
}

func (this *Chunk) delete() {
	gl.DeleteBuffers(1, &this.vao)
	gl.DeleteBuffers(1, &this.vbo)
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

	var posX, posY, posZ = idxToPos3(idx, int(this.dim[0]), int(this.dim[1]), int(this.dim[2]))

	for i := 0; i < len(offsets); i += 3 {
		var newX, newY, newZ = posX + offsets[i], posY + offsets[i+1], posZ + offsets[i+2]

		var idx = pos3ToIdx(int(newX), int(newY), int(newZ), int(this.dim[0]), int(this.dim[1]), int(this.dim[2]))
		if idx >= 0 && idx < len(this.cubes) {
			retAmount += 1
		}
	}
	return retAmount
}

func (this *Chunk) createVBO() {
	profiler.startTime("vertexCreation")
	var blockStride = 6
	this.calculateInnerBlocks()
	var vboSize = 0
	var idx = 0
	for i := 0; i < len(this.cubes); i++ {
		var c = this.cubes[i]
		var baseX, baseY = idxToGridPos(c.texId, 64, 64)
		var posX, posY, posZ = idxToPos3(i, int(this.dim[0]), int(this.dim[1]), int(this.dim[2]))
		if !c.isInner {
			//vertexNr
			var colour = (vboSize % 2)
			colour = 1
			_ = colour
			for j := 0; j < verticesPerCube; j++ {
				//copy pos-3bit
				var vertex uint32 = 0
				if cubeVertices[j*5+0] == 1.0 {
					vertex |= 0b100
				}
				if cubeVertices[j*5+1] == 1.0 {
					vertex |= 0b010
				}
				if cubeVertices[j*5+2] == 1.0 {
					vertex |= 0b001
				}
				vertex <<= 5
				//copy tex-10bit
				var texX = (j / 6) + int(cubeVertices[j*5+3]) + baseX*blockStride
				var texY = int(cubeVertices[j*5+4]) + baseY*blockStride
				_, _, _ = blockStride, baseX, baseY
				//var texX, texY = colour, colour
				vertex |= uint32(texX)
				vertex <<= 5
				vertex |= uint32(texY)
				vertex <<= 5
				//copy model-15bit
				vertex |= uint32(posX)
				vertex <<= 5
				vertex |= uint32(posY)
				vertex <<= 5
				vertex |= uint32(posZ)
				cumulatedVertices2[vboSize*36+j] = vertex

				idx += 1
			}
			vboSize += 1
		}
	}
	profiler.endTime("vertexCreation")
	this.amountOuter = uint32(vboSize)
	profiler.startTime("buffer")
	this.delete()
	profiler.startTime("genBufs")
	generateBuffersCopy2(&this.vao, &this.vbo, nil, cumulatedVertices2, 0, nil, []int{1})
	profiler.endTime("genBufs")

	profiler.endTime("buffer")

}

// surrounded on all sides
func (this *Chunk) isInnerBlock(idx int) bool {
	var neighbours = this.getAmountNeighbours(idx)
	var posX, posY, posZ = idxToPos3(idx, int(this.dim[0]), int(this.dim[1]), int(this.dim[2]))

	var isInner = posX > 0 && posX < int(this.dim[0])-1 && posY > 0 && posY < int(this.dim[1])-1 && posZ > 0 && posZ < int(this.dim[2])-1
	return isInner && neighbours >= 4
}

func (this *Chunk) calculateInnerBlocks() {
	for i := 0; i < len(this.cubes); i++ {
		this.cubes[i].isInner = this.isInnerBlock(i)
	}
}

func (this *Chunk) isVisible() bool {
	var points = []glm.Vec3{
		//Mittelpunk
		this.pos.Add(&glm.Vec3{this.dim[0] / 2, this.dim[1] / 2, this.dim[2] / 2}),
		//Ecken
		this.pos,

		this.pos.Add(&glm.Vec3{this.dim[0], 0, 0}),
		this.pos.Add(&glm.Vec3{this.dim[0], 0, this.dim[2]}),
		this.pos.Add(&glm.Vec3{0, 0, this.dim[2]}),
		this.pos.Add(&glm.Vec3{0, -this.dim[1], 0}),
		this.pos.Add(&glm.Vec3{this.dim[0], -this.dim[1], 0}),
		this.pos.Add(&glm.Vec3{this.dim[0], -this.dim[1], this.dim[2]}),
		this.pos.Add(&glm.Vec3{0, -this.dim[1], this.dim[2]}),
		//Seiten
		this.pos.Add(&glm.Vec3{this.dim[0] / 2, -this.dim[1] / 2, 0}),
		this.pos.Add(&glm.Vec3{0, -this.dim[1] / 2, this.dim[2] / 2}),
		this.pos.Add(&glm.Vec3{0, -this.dim[1] / 2, this.dim[2] / 2}),
		this.pos.Add(&glm.Vec3{this.dim[0] / 2, -this.dim[1] / 2, this.dim[2]}),
		this.pos.Add(&glm.Vec3{this.dim[0] / 2, 0, this.dim[2] / 2}),
		this.pos.Add(&glm.Vec3{this.dim[0] / 2, -this.dim[1] / 2, this.dim[2] / 2}),
	}

	var retVal = false
	for _, x := range points {
		if isPointInFrustum(this.camera, x) {

			retVal = true
			break
		}
	}
	return retVal
}

func (this *Chunk) draw() {
	//this.createVBO()
	this.shader.use()
	this.shader.setUniformMatrix4("view", &this.camera.lookAtMat)
	this.shader.setUniformMatrix4("projection", this.projection)
	this.shader.setUniformVec3("chunkOrigin", &this.pos)
	this.shader.setUniform1i("tex", 0)

	gl.BindVertexArray(this.vao)
	gl.BindBuffer(gl.ARRAY_BUFFER, this.vbo)
	//gl.BufferSubData(gl.ARRAY_BUFFER, 0, int(this.amountOuter)*36, gl.Ptr(cumulatedVertices2))
	gl.ActiveTexture(gl.TEXTURE0)
	gl.BindTexture(gl.TEXTURE_2D, *this.tex)
	gl.DrawArrays(gl.TRIANGLES, 0, int32(this.amountOuter)*int32(verticesPerCube))

	gl.BindVertexArray(0)
	gl.BindBuffer(gl.ARRAY_BUFFER, 0)
	gl.BindTexture(gl.TEXTURE_2D, 0)

}
