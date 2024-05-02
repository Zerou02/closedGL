package main

import (
	"github.com/EngoEngine/glm"
	"github.com/go-gl/gl/v4.1-core/gl"
)

var cumulatedVertices = make([]float32, 36*(3+3+2)*16*16*16)

type CubeVertex struct {
	vertices []float32
	isInner  bool
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
				var vertices = make([]float32, (3+2+3)*36)
				for i := 0; i < 36; i++ {
					//pos
					vertices[i*8] = cube[i*5]
					vertices[i*8+1] = cube[i*5+1]
					vertices[i*8+2] = cube[i*5+2]
					//tex
					vertices[i*8+3] = cube[i*5+3]
					vertices[i*8+4] = cube[i*5+4]
					//model
					vertices[i*8+5] = float32(x) + chunk.pos[0]
					vertices[i*8+6] = float32(y) + chunk.pos[1]
					vertices[i*8+7] = float32(z) + chunk.pos[2]
				}
				var newVertex = CubeVertex{isInner: false, vertices: vertices}
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
	var pos = glm.Vec3{cube.vertices[5], cube.vertices[6], cube.vertices[7]}
	var first = this.cubes[0]
	var basePos = glm.Vec3{first.vertices[5], first.vertices[6], first.vertices[7]}

	for i := 0; i < len(offsets); i += 3 {
		var newX, newY, newZ = pos[0] - basePos[0] + offsets[i], pos[1] - basePos[1] + offsets[i+1], pos[2] - basePos[2] + offsets[i+2]

		var idx = pos3ToIdx(int(newX), int(newY), int(newZ), int(this.dim[0]), int(this.dim[1]), int(this.dim[2]))
		if idx >= 0 && idx < len(this.cubes) {
			retAmount += 1
		}
	}
	return retAmount
}

func (this *Chunk) createVBO() {
	this.calculateInnerBlocks()
	var vboSize = 0
	for i := 0; i < len(this.cubes); i++ {
		var c = this.cubes[i]
		if !c.isInner {
			for j := 0; j < len(c.vertices); j++ {
				cumulatedVertices[vboSize*288+j] = c.vertices[j]
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
	var pos = glm.Vec3{cube.vertices[5], cube.vertices[6], cube.vertices[7]}
	var first = this.cubes[0]
	var basePos = glm.Vec3{first.vertices[5], first.vertices[6], first.vertices[7]}
	var posX = pos[0] - basePos[0]
	var posY = pos[1] - basePos[1]
	var posZ = pos[2] - basePos[2]

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
	gl.DrawArrays(gl.TRIANGLES, 0, int32(this.amountOuter)*36)
}
