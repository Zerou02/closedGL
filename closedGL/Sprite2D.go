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

	rect.tex = LoadImage("./assets/sprites/fence.png", gl.RGBA)

	var quadBaseData = []float32{
		1.0, 0.0, //top r
		0.0, 0.0, // top l
		1.0, 1.0, // bottom r
		0.0, 1.0, // bottom l,
	}
	rect.baseVBO.cpuArr = quadBaseData
	rect.baseVBO.copyToGPU()

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
	for i := 0; i < 5000; i++ {
		this.beginDraw()
		this.createVertices(glm.Vec4{float32(i), float32(i), 100, 200}, glm.Vec4{1, 1, 1, 1})
		this.shader.use()
		this.shader.setUniformMatrix4("projection", this.projection)
		gl.ActiveTexture(gl.TEXTURE0)
		gl.BindTexture(gl.TEXTURE_2D, *this.tex)
		this.shader.setUniform1i("tex", 0)
		gl.Disable(gl.DEPTH_TEST)
		gl.BindVertexArray(this.vao)

		this.instanceBuffer.copyToGPU()
		gl.DrawElementsInstanced(gl.TRIANGLES, 6, gl.UNSIGNED_BYTE, gl.Ptr(this.indices), 1)
		gl.Enable(gl.DEPTH_TEST)
	}
}
