package closedGL

import (
	"github.com/EngoEngine/glm"
	"github.com/go-gl/gl/v4.6-core/gl"
)

type TriangleMesh struct {
	shader      *Shader
	vao         uint32
	amountQuads uint32
	projection  glm.Mat4
	view        glm.Mat4
	buffer      Buffer[float32]
}

func newTriMesh(shader *Shader, defaultProj glm.Mat4) TriangleMesh {
	var mesh = TriangleMesh{
		shader:     shader,
		view:       glm.Ident4(),
		projection: defaultProj,
	}
	mesh.initBuffer()
	return mesh
}

func (this *TriangleMesh) initBuffer() {
	this.vao = genVAO()
	this.buffer = genInterleavedBuffer[float32](this.vao, 0, []int{2, 2, 1}, []int{0, 0, 0}, gl.FLOAT)
}

func (this *TriangleMesh) Copy() {
	gl.BindVertexArray(this.vao)
	this.buffer.copyToGPU()
}

func (this *TriangleMesh) Clear() {
	this.buffer.clear()
	this.amountQuads = 0
}

func (this *TriangleMesh) AddTri(dim, uv [3]glm.Vec2, sign float32) {
	var stride uint32 = 15

	this.buffer.resizeCPUData(int(this.amountQuads+1) * int(stride))

	this.buffer.cpuBuffer[this.amountQuads*stride+0] = dim[0][0]
	this.buffer.cpuBuffer[this.amountQuads*stride+1] = dim[0][1]
	this.buffer.cpuBuffer[this.amountQuads*stride+2] = uv[0][0]
	this.buffer.cpuBuffer[this.amountQuads*stride+3] = uv[0][1]
	this.buffer.cpuBuffer[this.amountQuads*stride+4] = sign

	this.buffer.cpuBuffer[this.amountQuads*stride+5] = dim[1][0]
	this.buffer.cpuBuffer[this.amountQuads*stride+6] = dim[1][1]
	this.buffer.cpuBuffer[this.amountQuads*stride+7] = uv[1][0]
	this.buffer.cpuBuffer[this.amountQuads*stride+8] = uv[1][1]
	this.buffer.cpuBuffer[this.amountQuads*stride+9] = sign

	this.buffer.cpuBuffer[this.amountQuads*stride+10] = dim[2][0]
	this.buffer.cpuBuffer[this.amountQuads*stride+11] = dim[2][1]
	this.buffer.cpuBuffer[this.amountQuads*stride+12] = uv[2][0]
	this.buffer.cpuBuffer[this.amountQuads*stride+13] = uv[2][1]
	this.buffer.cpuBuffer[this.amountQuads*stride+14] = sign

	this.amountQuads++
}

func (this *TriangleMesh) Draw() {
	this.shader.use()
	this.shader.setUniformMatrix4("projection", &this.projection)
	this.shader.setUniformMatrix4("view", &this.view)
	gl.Disable(gl.DEPTH_TEST)
	gl.Disable(gl.CULL_FACE)
	gl.BindVertexArray(this.vao)
	gl.DrawArrays(gl.TRIANGLES, 0, int32(this.amountQuads)*3)
	gl.Enable(gl.DEPTH_TEST)
	gl.Enable(gl.CULL_FACE)

}
