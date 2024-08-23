package closedGL

/* package closedGL

import (
	"math"

	"github.com/EngoEngine/glm"
	"github.com/go-gl/gl/v4.3-core/gl"
)

//pos:x,y,z

var cubeVertexStride = 3
var cummVertexStride = 8

var bytesPerVertex = 4
var verticesPerCube = 36
var dim = 16
var cumulatedVertices2 = make([]uint32, verticesPerCube*(cummVertexStride)*dim*dim*dim)

// 10 bit TexOffset
// big-endian: xxxxx yyyyy
type CubeVertex = uint16

type Chunk struct {
	shader         *Shader
	camera         *Camera
	projection     *glm.Mat4
	tex            *Texture
	dim, pos       glm.Vec3
	cubes          []CubeVertex
	amountVertices uint32
	vao, vbo       uint32
}

func NewChunk(dim, pos glm.Vec3, tex *Texture, camera *Camera, projection *glm.Mat4, shader *Shader) Chunk {
	var chunk = Chunk{dim: dim, pos: pos, camera: camera, projection: projection, tex: tex, shader: shader}
	var amountCubes = int(dim[0] * dim[1] * dim[2])
	var cubeArr = make([]CubeVertex, amountCubes)
	var count = 0
	var texId = 0
	var texX, texY = IdxToGridPos(texId, 32, 32)
	for y := 0; y < int(dim[1]); y++ {
		for z := 0; z < int(dim[2]); z++ {
			for x := 0; x < int(dim[0]); x++ {

				var newVertex uint16 = 0
				newVertex |= uint16(texX)
				newVertex <<= 5
				newVertex |= uint16(texY)
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

func (this *Chunk) createVertices() {
	var blockStride = 6
	var idx = 0
	var amountVertices = 0
	for i := 0; i < len(this.cubes); i++ {
		var c = this.cubes[i]
		var baseY = int(c & 31)
		c >>= 5
		var baseX = int(c & 31)

		var posX, posY, posZ = idxToPos3(i, int(this.dim[0]), int(this.dim[1]), int(this.dim[2]))
		//vertexNr
		var allowFaces = this.faceCullCube2(i)
		for l := 0; l < len(allowFaces); l++ {
			var j = allowFaces[l]
			for k := 0; k < 6; k++ {
				//copy pos-3bit
				var vertex uint32 = 0
				if cubeVertices[(j*6+k)*5+0] == 1.0 {
					vertex |= 0b100
				}
				if cubeVertices[(j*6+k)*5+1] == 1.0 {
					vertex |= 0b010
				}
				if cubeVertices[(j*6+k)*5+2] == 1.0 {
					vertex |= 0b001
				}
				vertex <<= 5
				//copy tex-10bit
				var texX = j + int(cubeVertices[(j*6+k)*5+3]) + baseX*blockStride
				var texY = int(cubeVertices[(j*6+k)*5+4]) + baseY*blockStride
				_, _, _ = blockStride, baseX, baseY
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
				cumulatedVertices2[idx*6+k] = vertex
				amountVertices += 1
			}
			idx += 1
		}
	}
	this.amountVertices = uint32(amountVertices)
}
func (this *Chunk) createVBO() {
	this.createVertices()
	this.delete()
	//TODO: Fix
	//generateBuffersCopy2(&this.vao, &this.vbo, nil, cumulatedVertices2, 0, nil, []int{1})
}

// surrounded on all sides
func (this *Chunk) isInnerBlock(idx int) bool {
	var neighbours = this.getAmountNeighbours(idx)
	var posX, posY, posZ = idxToPos3(idx, int(this.dim[0]), int(this.dim[1]), int(this.dim[2]))

	var isInner = posX > 0 && posX < int(this.dim[0])-1 && posY > 0 && posY < int(this.dim[1])-1 && posZ > 0 && posZ < int(this.dim[2])-1
	return isInner && neighbours >= 4
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
		this.pos.Add(&glm.Vec3{0, this.dim[1], 0}),
		this.pos.Add(&glm.Vec3{this.dim[0], this.dim[1], 0}),
		this.pos.Add(&glm.Vec3{this.dim[0], this.dim[1], this.dim[2]}),
		this.pos.Add(&glm.Vec3{0, this.dim[1], this.dim[2]}),
		//Seiten
		this.pos.Add(&glm.Vec3{this.dim[0] / 2, this.dim[1] / 2, 0}),
		this.pos.Add(&glm.Vec3{0, this.dim[1] / 2, this.dim[2] / 2}),
		this.pos.Add(&glm.Vec3{0, this.dim[1] / 2, this.dim[2] / 2}),
		this.pos.Add(&glm.Vec3{this.dim[0] / 2, this.dim[1] / 2, this.dim[2]}),
		this.pos.Add(&glm.Vec3{this.dim[0] / 2, 0, this.dim[2] / 2}),
		this.pos.Add(&glm.Vec3{this.dim[0] / 2, this.dim[1] / 2, this.dim[2] / 2}),
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

func (this *Chunk) Draw() {
	if !this.isVisible() {
		return
	}

	this.shader.use()
	this.shader.setUniformMatrix4("view", &this.camera.lookAtMat)
	this.shader.setUniformMatrix4("projection", this.projection)
	this.shader.setUniformVec3("chunkOrigin", &this.pos)
	this.shader.setUniform1i("tex", 0)

	gl.BindVertexArray(this.vao)
	gl.BindBuffer(gl.ARRAY_BUFFER, this.vbo)

	gl.ActiveTexture(gl.TEXTURE0)
	gl.BindTexture(gl.TEXTURE_2D, *this.tex)
	gl.DrawArrays(gl.TRIANGLES, 0, int32(this.amountVertices)*3)

	gl.BindVertexArray(0)
	gl.BindBuffer(gl.ARRAY_BUFFER, 0)
	gl.BindTexture(gl.TEXTURE_2D, 0)

}

func (this *Chunk) faceCullCube(posCentre glm.Vec3) []int {

	var ca = this.camera
	var vec = posCentre.Sub(&ca.CameraPos)
	//oben,vorne,links,...,unten
	var normals = []glm.Vec3{{0, 1, 0}, {0, 0, 1}, {-1, 0, 0}, {-0, 0, -1}, {1, 0, 0}, {0, -1, 0}}
	var retVec = []int{}
	for j, x := range normals {
		for i := 0; i < 3; i++ {
			if x[i] == 0 {
				continue
			}
			var camSign = vec[i] > 0
			var cubeSign = x[i] > 0
			if camSign != cubeSign && math.Abs(float64(vec[i])) > 0.2 {
				retVec = append(retVec, j)
				break
			}
		}
	}
	return retVec
}

func (this *Chunk) faceCullCube2(cubeIdx int) []int {
	var dimX = int(this.dim[0])
	var dimY = int(this.dim[1])
	var dimZ = int(this.dim[2])

	var posX, posY, posZ = idxToPos3(cubeIdx, dimX, dimY, dimZ)
	var offsets = []int{
		0, 1, 0,
		0, 0, 1,
		-1, 0, 0,
		0, 0, -1,
		1, 0, 0,
		0, -1, 0,
	}

	var allowedFaces = []int{}
	for i := 0; i < len(offsets); i += 3 {
		var newX, newY, newZ = posX + offsets[i], posY + offsets[i+1], posZ + offsets[i+2]
		var isOuter = (newX < 0 || newX >= int(this.dim[0])) || (newY < 0 || newY >= int(this.dim[1])) || (newZ < 0 || newZ >= int(this.dim[2]))
		if isOuter {
			allowedFaces = append(allowedFaces, i/3)
		}
	}
	return allowedFaces
}
*/
