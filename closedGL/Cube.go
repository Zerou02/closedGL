package closedGL

import (
	"github.com/EngoEngine/glm"
	"github.com/go-gl/gl/v4.1-core/gl"
)

type CubeMesh struct {
	baseBuffer     BufferFloat
	instanceBuffer BufferFloat
	amountCubes    uint32
	vao            uint32
}

type Cube struct {
	shader      *Shader
	camera      *Camera
	projection  *glm.Mat4
	currMesh    CubeMesh
	textureMane TextureMane
	indices     []byte
}

func NewCube(shader *Shader, camera *Camera, projection *glm.Mat4, pos glm.Vec3) Cube {
	var retCube = Cube{shader: shader, camera: camera, projection: projection, textureMane: newTextureMane()}
	retCube.initMesh()
	retCube.indices = []byte{
		0, 1, 2, 2, 3, 0,
		4, 5, 6, 6, 7, 4,
		8, 9, 10, 10, 11, 8,
		12, 13, 14, 14, 15, 12,
		16, 17, 18, 18, 19, 16,
		20, 21, 22, 22, 23, 20,
	}

	return retCube
}

func (this *Cube) initMesh() {
	var vao = genVAO()
	var baseBuffer = generateInterleavedVBOFloat2(vao, 0, []int{3, 2})
	baseBuffer.cpuArr = cube
	baseBuffer.copyToGPU()
	var instanceBuffer = generateInterleavedVBOFloat2(vao, 2, []int{3, 2})
	gl.BindBuffer(gl.ARRAY_BUFFER, instanceBuffer.buffer)
	gl.VertexAttribDivisor(2, 1)

	this.currMesh = CubeMesh{
		vao:            vao,
		baseBuffer:     baseBuffer,
		instanceBuffer: instanceBuffer,
		amountCubes:    this.currMesh.amountCubes,
	}
}

func (this *Cube) beginDraw() {
	this.currMesh.amountCubes = 0
}

func (this *Cube) createVertices(pos glm.Vec3, texPath string, ctx *ClosedGLContext) {
	var stride uint32 = 5

	this.currMesh.instanceBuffer.resizeCPUData(int(this.currMesh.amountCubes+1) * int(stride))

	this.currMesh.instanceBuffer.cpuArr[this.currMesh.amountCubes*stride+0] = pos[0]
	this.currMesh.instanceBuffer.cpuArr[this.currMesh.amountCubes*stride+1] = pos[1]
	this.currMesh.instanceBuffer.cpuArr[this.currMesh.amountCubes*stride+2] = pos[2]

	this.textureMane.loadTex(texPath)

	var handle = this.textureMane.getHandle(texPath)
	var lower uint32 = uint32(handle & 0xffff_ffff)
	var upper uint32 = uint32((handle >> 32) & 0xffff_ffff)
	this.currMesh.instanceBuffer.cpuArr[this.currMesh.amountCubes*stride+3] = float32(lower)
	this.currMesh.instanceBuffer.cpuArr[this.currMesh.amountCubes*stride+4] = float32(upper)

	this.currMesh.amountCubes++
}

func (this *Cube) copyCurrMesh() CubeMesh {
	gl.BindVertexArray(this.currMesh.vao)
	this.currMesh.instanceBuffer.copyToGPU()
	this.currMesh.instanceBuffer.cpuArr = []float32{}
	this.currMesh.baseBuffer.cpuArr = []float32{}
	var retMesh = CubeMesh{
		baseBuffer:     this.currMesh.baseBuffer.copy(),
		instanceBuffer: this.currMesh.instanceBuffer.copy(),
		amountCubes:    this.currMesh.amountCubes,
		vao:            this.currMesh.vao,
	}
	this.currMesh = CubeMesh{}
	return retMesh
}

func (this *Cube) draw() {
	if this.currMesh.amountCubes == 0 {
		return
	}
	this.shader.use()
	this.textureMane.makeResident()

	this.shader.setUniformMatrix4("view", &this.camera.lookAtMat)
	this.shader.setUniformMatrix4("projection", &this.camera.perspective)

	gl.BindVertexArray(this.currMesh.vao)
	this.currMesh.instanceBuffer.copyToGPU()
	gl.DrawArraysInstanced(gl.TRIANGLES, 0, 36, int32(this.currMesh.amountCubes))
	this.textureMane.makeNonResident()
}

func (this *Cube) drawMesh(mesh *CubeMesh, ctx *ClosedGLContext) {
	if mesh.amountCubes == 0 {
		return
	}
	this.shader.use()
	this.textureMane.makeResident()

	this.shader.setUniformMatrix4("view", &this.camera.lookAtMat)
	this.shader.setUniformMatrix4("projection", &this.camera.perspective)

	gl.BindVertexArray(mesh.vao)

	gl.DrawArraysInstanced(gl.TRIANGLES, 0, 36, int32(mesh.amountCubes))
	this.textureMane.makeNonResident()

}
