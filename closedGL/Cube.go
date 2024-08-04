package closedGL

import (
	"github.com/EngoEngine/glm"
	"github.com/go-gl/gl/v4.1-core/gl"
)

type CubeMesh struct {
	baseBuffer     BufferFloat
	instanceBuffer BufferFloat
	uintBuffer     BufferU32
	amountCubes    uint32
	textureMane    TextureMane
	ssbo           SSBOU32
	vao            uint32
}

type Cube struct {
	shader     *Shader
	camera     *Camera
	projection *glm.Mat4
	currMesh   CubeMesh
	indices    []byte
}

func NewCube(shader *Shader, camera *Camera, projection *glm.Mat4, pos glm.Vec3) Cube {
	var retCube = Cube{shader: shader, camera: camera, projection: projection}
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
	var baseBuffer = generateInterleavedVBOFloat2(vao, 1, []int{3, 2})
	baseBuffer.cpuArr = cube
	baseBuffer.copyToGPU()
	var instanceBuffer = generateInterleavedVBOFloat2(vao, 3, []int{3})

	var uintBuffer = generateInterleavedVBOU32(vao, 0, []int{2})
	var ssbo = genSSBOU32(1)
	gl.BindBuffer(gl.ARRAY_BUFFER, instanceBuffer.buffer)
	gl.VertexAttribDivisor(3, 1)
	gl.BindBuffer(gl.ARRAY_BUFFER, uintBuffer.buffer)
	gl.VertexAttribDivisor(0, 1)

	gl.VertexAttribDivisor(3, 1)

	this.currMesh = CubeMesh{
		vao:            vao,
		baseBuffer:     baseBuffer,
		instanceBuffer: instanceBuffer,
		amountCubes:    this.currMesh.amountCubes,
		ssbo:           ssbo,
		uintBuffer:     uintBuffer,
		textureMane:    newTextureMane(),
	}
}

func (this *Cube) beginDraw() {
	this.currMesh.amountCubes = 0
	//this.instanceBuffer.clear()
	//this.ssbo.clear()
}

func (this *Cube) createVertices(pos glm.Vec3, texPath string, ctx *ClosedGLContext) {
	var stride uint32 = 3
	var ssboStride uint32 = 2

	this.currMesh.instanceBuffer.resizeCPUData(int(this.currMesh.amountCubes+1) * int(stride))
	this.currMesh.ssbo.resizeCPUData(int(this.currMesh.amountCubes+1) * int(ssboStride))
	this.currMesh.uintBuffer.resizeCPUData(int(this.currMesh.amountCubes+1) * int(ssboStride))

	this.currMesh.instanceBuffer.cpuArr[this.currMesh.amountCubes*stride+0] = pos[0]
	this.currMesh.instanceBuffer.cpuArr[this.currMesh.amountCubes*stride+1] = pos[1]
	this.currMesh.instanceBuffer.cpuArr[this.currMesh.amountCubes*stride+2] = pos[2]

	this.currMesh.textureMane.loadTex(texPath)

	var handle = this.currMesh.textureMane.getHandle(texPath)
	var lower uint32 = uint32(handle & 0xffff_ffff)
	var upper uint32 = uint32((handle >> 32) & 0xffff_ffff)
	this.currMesh.ssbo.cpuArr[this.currMesh.amountCubes*ssboStride+0] = lower
	this.currMesh.ssbo.cpuArr[this.currMesh.amountCubes*ssboStride+1] = upper
	this.currMesh.uintBuffer.cpuArr[this.currMesh.amountCubes*ssboStride+0] = lower
	this.currMesh.uintBuffer.cpuArr[this.currMesh.amountCubes*ssboStride+1] = upper

	this.currMesh.amountCubes++
}

func (this *Cube) copyCurrMesh() CubeMesh {
	gl.BindVertexArray(this.currMesh.vao)
	this.currMesh.instanceBuffer.copyToGPU()
	this.currMesh.ssbo.copyToGPU()
	this.currMesh.uintBuffer.copyToGPU()
	this.currMesh.instanceBuffer.cpuArr = []float32{}
	return this.currMesh
}

func (this *Cube) draw() {
	if this.currMesh.amountCubes == 0 {
		return
	}
	this.shader.use()
	this.currMesh.textureMane.makeResident()

	this.shader.setUniformMatrix4("view", &this.camera.lookAtMat)
	this.shader.setUniformMatrix4("projection", &this.camera.perspective)

	gl.BindVertexArray(this.currMesh.vao)
	this.currMesh.instanceBuffer.copyToGPU()
	this.currMesh.ssbo.copyToGPU()
	gl.DrawArraysInstanced(gl.TRIANGLES, 0, 36, int32(this.currMesh.amountCubes))
	this.currMesh.textureMane.makeNonResident()
}

func (this *Cube) drawMesh(mesh *CubeMesh, ctx *ClosedGLContext) {
	if mesh.amountCubes == 0 {
		return
	}
	//
	println(mesh.ssbo.cpuArr[0])
	println(mesh.ssbo.cpuArr[1])
	println(mesh.uintBuffer.cpuArr[0])
	println(mesh.uintBuffer.cpuArr[1])

	this.shader.use()
	ctx.Logger.Start("resident")
	mesh.textureMane.makeResident()
	ctx.Logger.End("resident")

	ctx.Logger.Start("uniform")

	this.shader.setUniformMatrix4("view", &this.camera.lookAtMat)
	this.shader.setUniformMatrix4("projection", &this.camera.perspective)

	ctx.Logger.End("uniform")
	ctx.Logger.Start("bind")
	gl.BindVertexArray(this.currMesh.vao)
	mesh.ssbo.copyToGPU()
	ctx.Logger.End("bind")

	gl.DrawArraysInstanced(gl.TRIANGLES, 0, 36, int32(mesh.amountCubes))
	mesh.textureMane.makeNonResident()

}
