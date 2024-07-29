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
	ssbo           SSBOU64
	textures       []uint32
	handleMap      map[string]uint64
}

func newSpriteMane(shader *Shader, projection *glm.Mat4) SpriteManager {
	var indices = []byte{0, 1, 2, 2, 1, 3}
	var rect = SpriteManager{shader: shader, projection: projection, amountQuads: 0, indices: indices, handleMap: map[string]uint64{}}

	rect.vao = genVAO()
	rect.baseVBO = generateInterleavedVBOFloat2(rect.vao, 0, []int{4})
	gl.BindBuffer(gl.ARRAY_BUFFER, rect.baseVBO.buffer)
	gl.VertexAttribDivisor(0, 0)
	rect.instanceBuffer = generateInterleavedVBOFloat2(rect.vao, 1, []int{4})

	gl.BindBuffer(gl.ARRAY_BUFFER, rect.instanceBuffer.buffer)
	gl.VertexAttribDivisor(1, 1)

	rect.ssbo = genSSBOU64(1)

	var cArr = []uint64{0}
	rect.ssbo.cpuArr = cArr

	var quadBaseData = []float32{
		//pos,uv
		1.0, 0.0, 1.0, 0.0, //top r
		0.0, 0.0, 0.0, 0.0, // top l
		1.0, 1.0, 1.0, 1.0, // bottom r
		0.0, 1.0, 0.0, 1.0, // bottom l,
	}
	rect.baseVBO.cpuArr = quadBaseData
	rect.baseVBO.copyToGPU()

	return rect
}

func (this *SpriteManager) loadTex(path string) {
	if this.handleMap[path] != 0 {
		return
	}
	var texture = *LoadImage(path, gl.RGBA)
	var handle = gl.GetTextureHandleARB(texture)
	if handle == 0 {
		panic("invalid textureHandle")
	}
	this.textures = append(this.textures, texture)
	this.handleMap[path] = handle
}

func (this *SpriteManager) beginDraw() {
	this.amountQuads = 0
	this.instanceBuffer.clear()
	this.ssbo.clear()
}

func (this *SpriteManager) deleteBuffers() {
	gl.DeleteBuffers(1, &this.vao)
	gl.DeleteBuffers(1, &this.baseVBO.buffer)
	gl.DeleteBuffers(1, &this.instanceBuffer.buffer)

}

func (this *SpriteManager) createVertices(dim glm.Vec4, texPath string) {
	var stride uint32 = 4

	this.instanceBuffer.resizeCPUData(int(this.amountQuads+1) * int(stride))
	this.ssbo.resizeCPUData(int(this.amountQuads+1) * 1)

	this.instanceBuffer.cpuArr[this.amountQuads*stride+0] = dim[0]
	this.instanceBuffer.cpuArr[this.amountQuads*stride+1] = dim[1]
	this.instanceBuffer.cpuArr[this.amountQuads*stride+2] = dim[2]
	this.instanceBuffer.cpuArr[this.amountQuads*stride+3] = dim[3]

	this.loadTex(texPath)
	this.ssbo.cpuArr[this.amountQuads] = this.handleMap[texPath]

	this.amountQuads++

}

func (this *SpriteManager) draw() {
	this.shader.use()

	for _, v := range this.handleMap {
		gl.MakeTextureHandleResidentARB(v)
	}

	this.shader.setUniformMatrix4("projection", this.projection)
	gl.Disable(gl.DEPTH_TEST)
	gl.BindVertexArray(this.vao)
	this.ssbo.copyToGPU()
	this.instanceBuffer.copyToGPU()
	gl.DrawElementsInstanced(gl.TRIANGLES, 6, gl.UNSIGNED_BYTE, gl.Ptr(this.indices), int32(this.amountQuads))
	gl.Enable(gl.DEPTH_TEST)

	for _, v := range this.handleMap {
		gl.MakeTextureHandleNonResidentARB(v)
	}
}
