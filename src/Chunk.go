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
var cumulatedVertices = make([]float32, verticesPerCube*(cummVertexStride)*dim*dim*dim)

type CubeVertex struct {
	vertices []byte
	isInner  bool
	texId    int
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
	var chunk = Chunk{dim: dim, pos: pos, camera: camera, projection: projection, tex: tex, shader: shader}
	var amountCubes = int(dim[0] * dim[1] * dim[2])
	var cubeArr = make([]CubeVertex, amountCubes)
	var count = 0
	for y := 0; y < int(dim[1]); y++ {
		for z := 0; z < int(dim[2]); z++ {
			for x := 0; x < int(dim[0]); x++ {
				var vertices = make([]byte, (cubeVertexStride)*verticesPerCube)
				for i := 0; i < verticesPerCube; i++ {
					//pos
					vertices[i*cubeVertexStride] = byte(x)
					vertices[i*cubeVertexStride+1] = byte(y)
					vertices[i*cubeVertexStride+2] = byte(z)
				}
				var newVertex = CubeVertex{isInner: false, vertices: vertices, texId: 1}
				cubeArr[count] = newVertex
				count += 1
			}
		}
	}
	chunk.cubes = cubeArr
	//copy into vbo
	chunk.createVBO()
	return chunk
}

// No Diagonal Neighbours
func (this *Chunk) getAmountNeighbours(cube *CubeVertex) int {
	var retAmount = 0
	var offsets = []float32{
		-1, 0, 0,
		1, 0, 0,
		0, 1, 0,
		0, -1, 0,
		0, 0, 1,
		0, 0, -1,
	}
	var pos = glm.Vec3{float32(cube.vertices[0]), float32(cube.vertices[1]), float32(cube.vertices[2])}

	for i := 0; i < len(offsets); i += 3 {
		var newX, newY, newZ = pos[0] + offsets[i], pos[1] + offsets[i+1], pos[2] + offsets[i+2]

		var idx = pos3ToIdx(int(newX), int(newY), int(newZ), int(this.dim[0]), int(this.dim[1]), int(this.dim[2]))
		if idx >= 0 && idx < len(this.cubes) {
			retAmount += 1
		}
	}
	return retAmount
}

func (this *Chunk) createVBO() {

	var blockStride = 6
	this.calculateInnerBlocks()
	var vboSize = 0
	for i := 0; i < len(this.cubes); i++ {
		var c = this.cubes[i]
		var baseX, baseY = idxToGridPos(c.texId, 32, 32)
		if !c.isInner {
			//vertexNr
			for j := 0; j < verticesPerCube; j++ {
				//copy pos
				for k := 0; k < 3; k++ {
					cumulatedVertices[vboSize*288+j*8+k] = cubeVertices[j*5+k]
				}
				//copy tex
				cumulatedVertices[vboSize*288+j*8+3] = float32((j / 6)) + cubeVertices[j*5+3] + float32(baseX)*float32(blockStride)
				cumulatedVertices[vboSize*288+j*8+4] = cubeVertices[j*5+4] + float32(baseY)*float32(blockStride)
				//copy model
				for k := 0; k < 3; k++ {
					cumulatedVertices[vboSize*288+j*8+5+k] = float32(c.vertices[k]) + this.pos[k]
				}
			}
			vboSize += 1
		}
	}
	this.amountOuter = uint32(vboSize)
	gl.DeleteBuffers(1, &this.vao)
	gl.DeleteBuffers(1, &this.vbo)
	generateBuffersCopy(&this.vao, &this.vbo, nil, cumulatedVertices, 0, nil, []int{3, 2, 3})
}

// surrounded on all sides
func (this *Chunk) isInnerBlock(cube *CubeVertex) bool {
	var neighbours = this.getAmountNeighbours(cube)
	var posX = float32(cube.vertices[0])
	var posY = float32(cube.vertices[1])
	var posZ = float32(cube.vertices[2])

	var isInner = posX > 0 && posX < this.dim[0]-1 && posY > 0 && posY < this.dim[1]-1 && posZ > 0 && posZ < this.dim[2]-1
	return isInner && neighbours >= 4
}

func (this *Chunk) calculateInnerBlocks() {
	for i := 0; i < len(this.cubes); i++ {
		this.cubes[i].isInner = this.isInnerBlock(&this.cubes[i])
	}
}

func (this *Chunk) draw() {
	this.shader.use()
	this.shader.setUniformMatrix4("view", &this.camera.lookAtMat)
	this.shader.setUniformMatrix4("projection", this.projection)
	this.shader.setUniform1i("tex", 0)

	gl.BindVertexArray(this.vao)
	gl.BindBuffer(gl.ARRAY_BUFFER, this.vbo)
	gl.ActiveTexture(gl.TEXTURE0)
	gl.BindTexture(gl.TEXTURE_2D, *this.tex)
	gl.DrawArrays(gl.TRIANGLES, 0, int32(this.amountOuter)*int32(verticesPerCube))
}
