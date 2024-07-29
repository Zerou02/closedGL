package closedGL

import (
	"github.com/EngoEngine/glm"
	"github.com/go-gl/gl/v4.1-core/gl"
)

type Cube struct {
	shader         *Shader
	camera         *Camera
	projection     *glm.Mat4
	position       glm.Vec3
	vao            uint32
	isInner        bool
	baseBuffer     BufferFloat
	instanceBuffer BufferFloat
	amountCubes    uint32
	textureMane    TextureMane
	ssbo           SSBOU64
	indices        []byte
}

func NewCube(shader *Shader, camera *Camera, projection *glm.Mat4, pos glm.Vec3) Cube {
	var retCube = Cube{shader: shader, camera: camera, projection: projection, position: pos, isInner: false, amountCubes: 0, textureMane: newTextureMane()}
	//TODO:Fix
	retCube.vao = genVAO()
	retCube.baseBuffer = generateInterleavedVBOFloat2(retCube.vao, 0, []int{3, 2})
	retCube.baseBuffer.cpuArr = cube
	retCube.baseBuffer.copyToGPU()
	retCube.instanceBuffer = generateInterleavedVBOFloat2(retCube.vao, 2, []int{3})
	retCube.ssbo = genSSBOU64(1)
	gl.BindBuffer(gl.ARRAY_BUFFER, retCube.instanceBuffer.buffer)
	gl.VertexAttribDivisor(2, 1)
	retCube.indices = []byte{
		0, 1, 2, 2, 3, 0,
		4, 5, 6, 6, 7, 4,
		8, 9, 10, 10, 11, 8,
		12, 13, 14, 14, 15, 12,
		16, 17, 18, 18, 19, 16,
		20, 21, 22, 22, 23, 20,
	}

	retCube.baseBuffer.copyToGPU()
	return retCube
}

func (this *Cube) beginDraw() {
	this.amountCubes = 0
	this.instanceBuffer.clear()
	this.ssbo.clear()
}

func (this *Cube) createVertices(pos glm.Vec3, texPath string) {
	var stride uint32 = 3

	this.instanceBuffer.resizeCPUData(int(this.amountCubes+1) * int(stride))
	this.ssbo.resizeCPUData(int(this.amountCubes+1) * 1)

	this.instanceBuffer.cpuArr[this.amountCubes*stride+0] = pos[0]
	this.instanceBuffer.cpuArr[this.amountCubes*stride+1] = pos[1]
	this.instanceBuffer.cpuArr[this.amountCubes*stride+2] = pos[2]

	this.textureMane.loadTex(texPath)
	this.ssbo.cpuArr[this.amountCubes] = this.textureMane.getHandle(texPath)

	this.amountCubes++
}

func (this *Cube) draw() {
	this.shader.use()
	this.textureMane.makeResident()

	this.shader.setUniformMatrix4("view", &this.camera.lookAtMat)
	this.shader.setUniformMatrix4("projection", &this.camera.perspective)

	gl.BindVertexArray(this.vao)
	this.instanceBuffer.copyToGPU()
	this.ssbo.copyToGPU()
	gl.DrawElementsInstanced(gl.TRIANGLES, 6*6, gl.UNSIGNED_BYTE, gl.Ptr(this.indices), int32(this.amountCubes))
	this.textureMane.makeNonResident()
}
