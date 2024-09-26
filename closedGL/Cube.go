package closedGL

import (
	"github.com/EngoEngine/glm"
	"github.com/go-gl/gl/v4.6-core/gl"
)

type CubeMesh struct {
	instanceBuffer Buffer[uint32]
	ssbo           SSBO[uint32]
	amountCubes    uint32
	vao            uint32
}

type Cube struct {
	shader          *Shader
	camera          *Camera
	projection      *glm.Mat4
	currMesh        CubeMesh
	textureMane     TextureMane
	indices         []byte
	baseMeshSSBO    SSBO[uint32]
	textureContains []string
}

func NewCube(shader *Shader, camera *Camera, projection *glm.Mat4) Cube {
	var retCube = Cube{shader: shader, camera: camera, projection: projection, textureMane: newTextureMane(), textureContains: []string{}}
	retCube.initMesh(glm.Vec3{0, 0, 0})
	retCube.indices = []byte{
		0, 1, 2, 2, 3, 0,
		4, 5, 6, 6, 7, 4,
		8, 9, 10, 10, 11, 8,
		12, 13, 14, 14, 15, 12,
		16, 17, 18, 18, 19, 16,
		20, 21, 22, 22, 23, 20,
	}

	var baseMeshSSBO = genSSBOGen[uint32](0, gl.UNSIGNED_INT)
	for _, x := range CompressedCubeVertices {
		baseMeshSSBO.cpuBuffer = append(baseMeshSSBO.cpuBuffer, uint32(x))
	}
	retCube.baseMeshSSBO = baseMeshSSBO

	return retCube
}

func (this *Cube) initMesh(anchor glm.Vec3) {
	var vao = genVAO()
	var ssbo = genSSBOGen[uint32](1, gl.UNSIGNED_INT)
	var meshData = genInterleavedBuffer[uint32](vao, 0, []int{2}, []int{1}, gl.UNSIGNED_INT)
	ssbo.cpuBuffer = append(ssbo.cpuBuffer, uint32(anchor[0]))
	ssbo.cpuBuffer = append(ssbo.cpuBuffer, uint32(anchor[1]))
	ssbo.cpuBuffer = append(ssbo.cpuBuffer, uint32(anchor[2]))

	var instanceBuffer = genInterleavedBuffer[float32](vao, 2, []int{2}, []int{1}, gl.FLOAT)
	_ = instanceBuffer

	this.currMesh = CubeMesh{
		vao:            vao,
		amountCubes:    this.currMesh.amountCubes,
		instanceBuffer: meshData,
		ssbo:           ssbo,
	}
}

func (this *Cube) beginDraw() {
	this.currMesh.amountCubes = 0
}

// side: 0 up, front,left,right,back,down
func (this *Cube) createVertices(pos glm.Vec3, size glm.Vec3, texPath string, side byte, texIdX, texIdY int) {
	var stride = 2
	this.currMesh.instanceBuffer.resizeCPUData(int(this.currMesh.amountCubes+1) * stride)

	this.textureMane.loadTex(texPath)
	var handle = this.textureMane.getHandle(texPath)

	var lower uint32 = uint32(handle & 0xffff_ffff)
	var upper uint32 = uint32((handle >> 32) & 0xffff_ffff)
	if !this.doesSSBOContainTex(texPath) {
		this.baseMeshSSBO.cpuBuffer = append(this.baseMeshSSBO.cpuBuffer, lower)
		this.baseMeshSSBO.cpuBuffer = append(this.baseMeshSSBO.cpuBuffer, upper)
		this.textureContains = append(this.textureContains, texPath)
	}

	//4texID,5u,5v,5x,5y,5z,3side
	var entry uint32 = 0
	//tex
	entry |= this.findIdxOfTex(texPath)
	//u
	entry <<= 5
	entry |= uint32(texIdX)
	//v
	entry <<= 5
	entry |= uint32(texIdY)
	//x
	entry <<= 5
	entry |= uint32(pos[0])
	entry <<= 5
	entry |= uint32(pos[1])
	entry <<= 5
	entry |= uint32(pos[2])
	entry <<= 3
	entry |= uint32(side)
	this.currMesh.instanceBuffer.cpuBuffer[this.currMesh.amountCubes*uint32(stride)+0] = entry

	var entry2 uint32 = 0
	entry2 |= uint32(size[0])
	entry2 <<= 6
	entry2 |= uint32(size[1])
	entry2 <<= 6
	entry2 |= uint32(size[2])
	this.currMesh.instanceBuffer.cpuBuffer[this.currMesh.amountCubes*uint32(stride)+1] = entry2

	this.currMesh.amountCubes++
}

func (this *Cube) doesSSBOContainTex(path string) bool {
	var retVal = false
	for _, x := range this.textureContains {
		if x == path {
			retVal = true
		}
	}
	return retVal
}

func (this *Cube) findIdxOfTex(path string) uint32 {
	var idx = 0
	for i, x := range this.textureContains {
		if x == path {
			idx = i
			break
		}
	}
	return uint32(idx)
}

func (this *Cube) copyCurrMesh() CubeMesh {

	gl.BindVertexArray(this.currMesh.vao)
	this.currMesh.instanceBuffer.copyToGPU()
	this.currMesh.instanceBuffer.cpuBuffer = []uint32{}

	var retMesh = CubeMesh{
		amountCubes: this.currMesh.amountCubes,
		vao:         this.currMesh.vao,
		//	instanceBuffer: this.currMesh.instanceBuffer.copy(),
		ssbo: this.currMesh.ssbo.copy(),
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
	gl.DrawArraysInstanced(gl.TRIANGLES, 0, 36, int32(this.currMesh.amountCubes))
	this.textureMane.makeNonResident()
}

func (this *Cube) drawMesh(mesh *CubeMesh) {

	if mesh.amountCubes == 0 {
		return
	}

	this.baseMeshSSBO.copyToGPU()
	mesh.ssbo.copyToGPU()
	this.shader.use()
	this.textureMane.makeResident()

	this.shader.setUniformMatrix4("view", &this.camera.lookAtMat)
	this.shader.setUniformMatrix4("projection", &this.camera.perspective)

	gl.BindVertexArray(mesh.vao)

	gl.DrawArraysInstanced(gl.TRIANGLES, 0, 6, int32(mesh.amountCubes))
	this.textureMane.makeNonResident()

}

func (this *SSBO[T]) copy() SSBO[T] {
	var newArr = make([]T, len(this.cpuBuffer))
	copy(newArr, this.cpuBuffer)
	return SSBO[T]{
		cpuBuffer:        newArr,
		gpuBuffer:        this.gpuBuffer,
		bindingPoint:     this.bindingPoint,
		elementsPerEntry: this.elementsPerEntry,
	}
}
