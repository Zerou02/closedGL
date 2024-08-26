package closedGL

import (
	"github.com/EngoEngine/glm"
	"github.com/go-gl/gl/v4.6-core/gl"
)

type PixelMesh struct {
	shader         *Shader
	vao            uint32
	amountQuads    int32
	projection     glm.Mat4
	View           glm.Mat4
	instanceBuffer BufferFloat
	indices        []byte
}

func newPixelMesh(shader *Shader, defaultProj glm.Mat4) PixelMesh {
	var indices = []byte{0, 1, 2, 2, 1, 3}

	var mesh = PixelMesh{
		shader:     shader,
		View:       glm.Ident4(),
		projection: defaultProj,
		indices:    indices,
	}
	mesh.initBuffer()
	return mesh
}

func (this *PixelMesh) initBuffer() {
	this.vao = genVAO()
	this.instanceBuffer = BufferFloat{
		buffer:     generateInterleavedVBOFloat(this.vao, 0, []int{2, 4}),
		bufferSize: 0,
		cpuArr:     []float32{},
	}
	gl.BindBuffer(gl.ARRAY_BUFFER, this.instanceBuffer.buffer)
}

func (this *PixelMesh) Copy() {
	gl.BindVertexArray(this.vao)
	this.instanceBuffer.copyToGPU()
}

func (this *PixelMesh) Clear() {
	this.instanceBuffer.clear()
	this.amountQuads = 0
}

func (this *PixelMesh) AddPixel(pos glm.Vec2, colour glm.Vec4) {
	var stride int32 = 6

	this.instanceBuffer.resizeCPUData(int(this.amountQuads+1) * int(stride))

	this.instanceBuffer.cpuArr[this.amountQuads*stride+0] = pos[0]
	this.instanceBuffer.cpuArr[this.amountQuads*stride+1] = pos[1]
	this.instanceBuffer.cpuArr[this.amountQuads*stride+2] = colour[0]
	this.instanceBuffer.cpuArr[this.amountQuads*stride+3] = colour[1]
	this.instanceBuffer.cpuArr[this.amountQuads*stride+4] = colour[2]
	this.instanceBuffer.cpuArr[this.amountQuads*stride+5] = colour[3]

	this.amountQuads++
}

func (this *PixelMesh) Draw() {

	gl.Enable(gl.PROGRAM_POINT_SIZE)
	gl.PointSize(1)

	this.shader.use()
	this.shader.setUniformMatrix4("projection", &this.projection)
	this.shader.setUniformMatrix4("view", &this.View)
	gl.Disable(gl.DEPTH_TEST)
	gl.BindVertexArray(this.vao)
	gl.DrawArrays(gl.POINTS, 0, this.amountQuads)
	gl.Enable(gl.DEPTH_TEST)

}
