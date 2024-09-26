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
	instanceBuffer Buffer[float32]
	baseVBO        Buffer[float32]
	indices        []byte
	dirty          bool
	update         bool
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
	this.baseVBO = genSingularBuffer[float32](this.vao, 0, 2, gl.FLOAT, false, 0)
	this.instanceBuffer = genInterleavedBuffer[float32](this.vao, 1, []int{4, 4}, []int{1, 1}, gl.FLOAT)

	var quadBaseData = []float32{
		1.0, 0.0, //top r
		0.0, 0.0, // top l
		1.0, 1.0, // bottom r
		0.0, 1.0, // bottom l,
	}
	this.baseVBO.cpuBuffer = quadBaseData
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

func (this *RectangleMesh) AddQuad(q *Quad) {
	this.AddRect(q.dim, q.colour)
}
func (this *RectangleMesh) AddRect(dim, colour glm.Vec4) {
	var stride uint32 = 8

	this.instanceBuffer.resizeCPUData(int(this.amountQuads+1) * int(stride))

	this.instanceBuffer.cpuBuffer[this.amountQuads*stride+0] = dim[0]
	this.instanceBuffer.cpuBuffer[this.amountQuads*stride+1] = dim[1]
	this.instanceBuffer.cpuBuffer[this.amountQuads*stride+2] = dim[2]
	this.instanceBuffer.cpuBuffer[this.amountQuads*stride+3] = dim[3]
	this.instanceBuffer.cpuBuffer[this.amountQuads*stride+4] = colour[0]
	this.instanceBuffer.cpuBuffer[this.amountQuads*stride+5] = colour[1]
	this.instanceBuffer.cpuBuffer[this.amountQuads*stride+6] = colour[2]
	this.instanceBuffer.cpuBuffer[this.amountQuads*stride+7] = colour[3]

	this.amountQuads++
	this.dirty = true
}

func (this *RectangleMesh) Draw() {
	if this.dirty {
		this.Copy()
	}
	this.shader.use()
	this.shader.setUniformMatrix4("projection", &this.projection)
	//	this.shader.setUniformMatrix4("view", &this.view)
	gl.Disable(gl.DEPTH_TEST)
	gl.BindVertexArray(this.vao)
	gl.DrawElementsInstanced(gl.TRIANGLES, 6, gl.UNSIGNED_BYTE, gl.Ptr(this.indices), int32(this.amountQuads))
	gl.Enable(gl.DEPTH_TEST)
	if this.update {
		this.update = false
	}
	if this.dirty {
		this.update = true
		this.dirty = false
		this.Clear()
	}
}

func (this *RectangleMesh) setDirty() {
	this.dirty = true
}

func (this *RectangleMesh) isUpdate() bool {
	return this.update
}
