package closedGL

import (
	"github.com/EngoEngine/glm"
	"github.com/go-gl/gl/v4.6-core/gl"
)

type RectangleMesh struct {
	shader         *Shader
	vao            uint32
	amountQuads    uint32
	projection     glm.Mat4
	view           glm.Mat4
	instanceBuffer BufferFloat
	baseVBO        BufferFloat
	indices        []byte
}

func newRectMesh(shader *Shader, defaultProj glm.Mat4) RectangleMesh {
	var indices = []byte{0, 1, 2, 2, 1, 3}

	var mesh = RectangleMesh{
		shader:     shader,
		view:       glm.Ident4(),
		projection: defaultProj,
		indices:    indices,
	}
	mesh.initBuffer()
	return mesh
}

func (this *RectangleMesh) initBuffer() {
	this.vao = genVAO()
	this.baseVBO = genSingularBufferFloat(this.vao, 0, 2, gl.FLOAT, false, 0)
	this.instanceBuffer = BufferFloat{
		buffer:     generateInterleavedVBOFloat(this.vao, 1, []int{4, 4}),
		bufferSize: 0,
		cpuArr:     []float32{},
	}
	gl.BindBuffer(gl.ARRAY_BUFFER, this.instanceBuffer.buffer)
	gl.VertexAttribDivisor(1, 1)
	gl.VertexAttribDivisor(2, 1)

	var quadBaseData = []float32{
		1.0, 0.0, //top r
		0.0, 0.0, // top l
		1.0, 1.0, // bottom r
		0.0, 1.0, // bottom l,
	}
	this.baseVBO.cpuArr = quadBaseData
	this.baseVBO.copyToGPU()
}

func (this *RectangleMesh) Copy() {
	gl.BindVertexArray(this.vao)
	this.baseVBO.copyToGPU()
	this.instanceBuffer.copyToGPU()
}

func (this *RectangleMesh) Clear() {
	this.instanceBuffer.clear()
	this.amountQuads = 0
}

func (this *RectangleMesh) AddRect(dim, colour glm.Vec4) {
	var stride uint32 = 8

	this.instanceBuffer.resizeCPUData(int(this.amountQuads+1) * int(stride))

	this.instanceBuffer.cpuArr[this.amountQuads*stride+0] = dim[0]
	this.instanceBuffer.cpuArr[this.amountQuads*stride+1] = dim[1]
	this.instanceBuffer.cpuArr[this.amountQuads*stride+2] = dim[2]
	this.instanceBuffer.cpuArr[this.amountQuads*stride+3] = dim[3]
	this.instanceBuffer.cpuArr[this.amountQuads*stride+4] = colour[0]
	this.instanceBuffer.cpuArr[this.amountQuads*stride+5] = colour[1]
	this.instanceBuffer.cpuArr[this.amountQuads*stride+6] = colour[2]
	this.instanceBuffer.cpuArr[this.amountQuads*stride+7] = colour[3]

	this.amountQuads++
}

func (this *RectangleMesh) Draw() {
	this.shader.use()
	this.shader.setUniformMatrix4("projection", &this.projection)
	//	this.shader.setUniformMatrix4("view", &this.view)
	gl.Disable(gl.DEPTH_TEST)
	gl.BindVertexArray(this.vao)
	gl.DrawElementsInstanced(gl.TRIANGLES, 6, gl.UNSIGNED_BYTE, gl.Ptr(this.indices), int32(this.amountQuads))
	gl.Enable(gl.DEPTH_TEST)
}

type RectangleManager struct {
	shader         *Shader
	projection     *glm.Mat4
	vao            uint32
	amountQuads    uint32
	instanceBuffer BufferFloat
	baseVBO        BufferFloat
	indices        []byte
}

func newRect(shader *Shader, projection *glm.Mat4) RectangleManager {
	var indices = []byte{0, 1, 2, 2, 1, 3}
	var rect = RectangleManager{shader: shader, projection: projection, amountQuads: 0, indices: indices}

	rect.vao = genVAO()
	rect.baseVBO = genSingularBufferFloat(rect.vao, 0, 2, gl.FLOAT, false, 0)
	rect.instanceBuffer = BufferFloat{
		buffer:     generateInterleavedVBOFloat(rect.vao, 1, []int{4, 4}),
		bufferSize: 0,
		cpuArr:     []float32{},
	}
	gl.BindBuffer(gl.ARRAY_BUFFER, rect.instanceBuffer.buffer)
	gl.VertexAttribDivisor(1, 1)
	gl.VertexAttribDivisor(2, 1)

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

func (this *RectangleManager) beginDraw() {
	this.amountQuads = 0
	this.instanceBuffer.clear()
}

func (this *RectangleManager) deleteBuffers() {
	gl.DeleteBuffers(1, &this.vao)
	gl.DeleteBuffers(1, &this.baseVBO.buffer)
	gl.DeleteBuffers(1, &this.instanceBuffer.buffer)

}

func (this *RectangleManager) createVertices(dim, colour glm.Vec4) {
	var stride uint32 = 8

	this.instanceBuffer.resizeCPUData(int(this.amountQuads+1) * int(stride))

	this.instanceBuffer.cpuArr[this.amountQuads*stride+0] = dim[0]
	this.instanceBuffer.cpuArr[this.amountQuads*stride+1] = dim[1]
	this.instanceBuffer.cpuArr[this.amountQuads*stride+2] = dim[2]
	this.instanceBuffer.cpuArr[this.amountQuads*stride+3] = dim[3]
	this.instanceBuffer.cpuArr[this.amountQuads*stride+4] = colour[0]
	this.instanceBuffer.cpuArr[this.amountQuads*stride+5] = colour[1]
	this.instanceBuffer.cpuArr[this.amountQuads*stride+6] = colour[2]
	this.instanceBuffer.cpuArr[this.amountQuads*stride+7] = colour[3]

	this.amountQuads++
}

func (this *RectangleManager) draw() {
	this.shader.use()
	this.shader.setUniformMatrix4("projection", this.projection)
	gl.Disable(gl.DEPTH_TEST)
	gl.BindVertexArray(this.vao)

	this.instanceBuffer.copyToGPU()
	gl.DrawElementsInstanced(gl.TRIANGLES, 6, gl.UNSIGNED_BYTE, gl.Ptr(this.indices), int32(this.amountQuads))
	gl.Enable(gl.DEPTH_TEST)
}
