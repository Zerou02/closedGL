package closedGL

import (
	"github.com/EngoEngine/glm"
	"github.com/go-gl/gl/v4.1-core/gl"
)

type SpriteManager struct {
	shader         *Shader
	projection     *glm.Mat4
	vao            uint32
	amountQuads    uint32
	instanceBuffer BufferFloat
	baseVBO        BufferFloat
	indices        []byte
	tex            *uint32
	ssbo           SSBOU64
	textures       []uint32
	handles        []uint64
}

func newSpriteMane(shader *Shader, projection *glm.Mat4) SpriteManager {
	var indices = []byte{0, 1, 2, 2, 1, 3}
	var rect = SpriteManager{shader: shader, projection: projection, amountQuads: 0, indices: indices}

	rect.vao = genVAO()
	rect.baseVBO = genSingularBufferFloat(rect.vao, 0, 2, gl.FLOAT, false, 0)
	rect.instanceBuffer = BufferFloat{
		buffer:     generateInterleavedVBOFloat(rect.vao, 1, []int{4, 4, 4}),
		bufferSize: 0,
		cpuArr:     []float32{},
	}
	gl.BindBuffer(gl.ARRAY_BUFFER, rect.instanceBuffer.buffer)
	gl.VertexAttribDivisor(1, 1)
	gl.VertexAttribDivisor(2, 1)
	gl.VertexAttribDivisor(3, 1)

	rect.ssbo = genSSBOU64(1)

	var cArr = []uint64{0}
	rect.ssbo.cpuArr = cArr

	rect.tex = LoadImage("./assets/sprites/fence.png", gl.RGBA)

	var quadBaseData = []float32{
		1.0, 0.0, //top r
		0.0, 0.0, // top l
		1.0, 1.0, // bottom r
		0.0, 1.0, // bottom l,
	}
	rect.baseVBO.cpuArr = quadBaseData
	rect.baseVBO.copyToGPU()

	var textureData = [32 * 32 * 3]byte{}
	for i := 0; i < len(textureData); i++ {
		textureData[i] = 1
	}

	var numInstances = 10

	for i := 0; i < numInstances; i++ {
		var texture = *LoadImage("./assets/sprites/fence.png", gl.RGBA)
		var handle = gl.GetTextureHandleARB(texture)
		if handle == 0 {
			panic("123")
		}
		rect.handles = append(rect.handles, handle)
	}
	rect.ssbo.cpuArr = rect.handles
	return rect
}

func (this *SpriteManager) beginDraw() {
	this.amountQuads = 0
	this.instanceBuffer.clear()
}

func (this *SpriteManager) deleteBuffers() {
	gl.DeleteBuffers(1, &this.vao)
	gl.DeleteBuffers(1, &this.baseVBO.buffer)
	gl.DeleteBuffers(1, &this.instanceBuffer.buffer)

}

func (this *SpriteManager) createVertices(dim, colour glm.Vec4) {
	var stride uint32 = 12

	this.instanceBuffer.resizeCPUData(int(this.amountQuads+1) * int(stride))

	this.instanceBuffer.cpuArr[this.amountQuads*stride+0] = dim[0]
	this.instanceBuffer.cpuArr[this.amountQuads*stride+1] = dim[1]
	this.instanceBuffer.cpuArr[this.amountQuads*stride+2] = dim[2]
	this.instanceBuffer.cpuArr[this.amountQuads*stride+3] = dim[3]

	this.instanceBuffer.cpuArr[this.amountQuads*stride+4] = 1
	this.instanceBuffer.cpuArr[this.amountQuads*stride+5] = 0
	this.instanceBuffer.cpuArr[this.amountQuads*stride+6] = 1
	this.instanceBuffer.cpuArr[this.amountQuads*stride+7] = 0

	this.instanceBuffer.cpuArr[this.amountQuads*stride+8] = 0
	this.instanceBuffer.cpuArr[this.amountQuads*stride+9] = 0
	this.instanceBuffer.cpuArr[this.amountQuads*stride+10] = 1
	this.instanceBuffer.cpuArr[this.amountQuads*stride+11] = 1

	this.amountQuads++

}

func (this *SpriteManager) draw() {
	this.shader.use()

	for i := 0; i < len(this.handles); i++ {
		gl.MakeTextureHandleResidentARB(this.handles[i])
	}
	this.shader.setUniformMatrix4("projection", this.projection)
	gl.Disable(gl.DEPTH_TEST)
	gl.BindVertexArray(this.vao)
	this.ssbo.copyToGPU()
	this.instanceBuffer.copyToGPU()
	gl.DrawElementsInstanced(gl.TRIANGLES, 6, gl.UNSIGNED_BYTE, gl.Ptr(this.indices), 1)
	gl.Enable(gl.DEPTH_TEST)

	for i := 0; i < len(this.handles); i++ {
		gl.MakeTextureHandleNonResidentARB(this.handles[i])
	}
}
