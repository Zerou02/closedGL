package closedGL

import (
	"github.com/EngoEngine/glm"
	"github.com/go-gl/gl/v4.6-core/gl"
)

type SpriteManager struct {
	shader         *Shader
	projection     *glm.Mat4
	vao            uint32
	amountQuads    uint32
	instanceBuffer Buffer[float32]
	baseVBO        Buffer[float32]
	ssbo           SSBO[uint64]
	indices        []byte
	textureMane    TextureMane
}

func newSpriteMane(shader *Shader, projection *glm.Mat4) SpriteManager {
	var indices = []byte{0, 1, 2, 2, 1, 3}
	var rect = SpriteManager{shader: shader, projection: projection, amountQuads: 0, indices: indices, textureMane: newTextureMane()}

	rect.vao = genVAO()
	rect.baseVBO = genInterleavedBuffer[float32](rect.vao, 0, []int{4}, []int{0}, gl.FLOAT)
	rect.instanceBuffer = genInterleavedBuffer[float32](rect.vao, 1, []int{4, 4, 2}, []int{1, 1, 1}, gl.FLOAT)

	rect.ssbo = genSSBOGen[uint64](1, gl.UNSIGNED_INT64_ARB)

	var cArr = []uint64{0}
	rect.ssbo.cpuBuffer = cArr

	var quadBaseData = []float32{
		//pos,uv
		1.0, 0.0, 1.0, 0.0, //top r
		0.0, 0.0, 0.0, 0.0, // top l
		1.0, 1.0, 1.0, 1.0, // bottom r
		0.0, 1.0, 0.0, 1.0, // bottom l,
	}
	rect.baseVBO.cpuBuffer = quadBaseData
	rect.baseVBO.copyToGPU()

	return rect
}

func (this *SpriteManager) beginDraw() {
	this.amountQuads = 0
	this.instanceBuffer.clear()
	this.ssbo.clear()
}

func (this *SpriteManager) deleteBuffers() {
	gl.DeleteBuffers(1, &this.vao)
	gl.DeleteBuffers(1, &this.baseVBO.gpuBuffer)
	gl.DeleteBuffers(1, &this.instanceBuffer.gpuBuffer)

}

func (this *SpriteManager) createVertices(dim glm.Vec4, texPath string, uv glm.Vec4, cellSpriteSize glm.Vec2) {
	var stride uint32 = 10

	this.instanceBuffer.resizeCPUData(int(this.amountQuads+1) * int(stride))
	this.ssbo.resizeCPUData(int(this.amountQuads+1) * 1)

	this.instanceBuffer.cpuBuffer[this.amountQuads*stride+0] = dim[0]
	this.instanceBuffer.cpuBuffer[this.amountQuads*stride+1] = dim[1]
	this.instanceBuffer.cpuBuffer[this.amountQuads*stride+2] = dim[2]
	this.instanceBuffer.cpuBuffer[this.amountQuads*stride+3] = dim[3]
	//uvData
	this.instanceBuffer.cpuBuffer[this.amountQuads*stride+4] = uv[0]
	this.instanceBuffer.cpuBuffer[this.amountQuads*stride+5] = uv[1]
	this.instanceBuffer.cpuBuffer[this.amountQuads*stride+6] = uv[2]
	this.instanceBuffer.cpuBuffer[this.amountQuads*stride+7] = uv[3]
	//cellSprite
	this.instanceBuffer.cpuBuffer[this.amountQuads*stride+8] = cellSpriteSize[0]
	this.instanceBuffer.cpuBuffer[this.amountQuads*stride+9] = cellSpriteSize[1]

	this.textureMane.loadTex(texPath)
	this.ssbo.cpuBuffer[this.amountQuads] = this.textureMane.getHandle(texPath)

	this.amountQuads++
}

func (this *SpriteManager) draw() {
	if len(this.ssbo.cpuBuffer) == 0 {
		return
	}
	this.shader.use()
	this.textureMane.makeResident()
	this.shader.setUniformMatrix4("projection", this.projection)
	gl.Disable(gl.DEPTH_TEST)
	gl.BindVertexArray(this.vao)

	this.ssbo.copyToGPU()
	this.instanceBuffer.copyToGPU()
	gl.DrawElementsInstanced(gl.TRIANGLES, 6, gl.UNSIGNED_BYTE, gl.Ptr(this.indices), int32(this.amountQuads))
	gl.Enable(gl.DEPTH_TEST)
	this.textureMane.makeNonResident()
}
