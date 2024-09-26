package closedGL

import (
	"github.com/EngoEngine/glm"
	"github.com/EngoEngine/math"
	"github.com/go-gl/gl/v4.6-core/gl"
)

type BezierShader struct {
	shader         *Shader
	projection     *glm.Mat4
	vao            uint32
	baseBuffer     Buffer[float32]
	instanceBuffer Buffer[float32]
	amountLines    uint32
	indices        []byte
}

func newBezier(shader *Shader, projection *glm.Mat4) BezierShader {
	var indices = []byte{0, 1, 2, 2, 1, 3}
	var b = BezierShader{shader: shader, projection: projection, vao: genVAO(), amountLines: 0, indices: indices}
	b.genBuffers()
	return b
}

func (this *BezierShader) genBuffers() {
	gl.BindVertexArray(this.vao)
	var quadBaseData = []float32{
		1.0, 0.0, //top r
		0.0, 0.0, // top l
		1.0, 1.0, // bottom r
		0.0, 1.0, // bottom l,
	}
	this.baseBuffer = genSingularBuffer[float32](this.vao, 0, 2, gl.FLOAT, false, 0)
	this.baseBuffer.cpuBuffer = quadBaseData
	this.baseBuffer.copyToGPU()
	this.instanceBuffer = genInterleavedBuffer[float32](this.vao, 1, []int{4, 2, 2, 2}, []int{1, 1, 1, 1}, gl.FLOAT)
	gl.BindBuffer(gl.ARRAY_BUFFER, this.instanceBuffer.gpuBuffer)

}

func (this *BezierShader) beginDraw() {
	this.amountLines = 0
	this.instanceBuffer.clear()
}

func (this *BezierShader) deleteBuffers() {

}

func (this *BezierShader) createVertices(p1, p2, cp glm.Vec2) {
	var stride uint32 = 4 + 3*2
	var minX = math.Min(math.Min(p1[0], p2[0]), cp[0])
	var minY = math.Min(math.Min(p1[1], p2[1]), cp[1])
	var maxX = math.Max(math.Max(p1[0], p2[0]), cp[0])
	var maxY = math.Max(math.Max(p1[1], p2[1]), cp[1])

	var dim = glm.Vec4{minX, minY, maxX - minX, maxY - minY}

	this.instanceBuffer.resizeCPUData(int(this.amountLines+1) * int(stride))

	this.instanceBuffer.cpuBuffer[this.amountLines*stride+0] = dim[0]
	this.instanceBuffer.cpuBuffer[this.amountLines*stride+1] = dim[1]
	this.instanceBuffer.cpuBuffer[this.amountLines*stride+2] = dim[2]
	this.instanceBuffer.cpuBuffer[this.amountLines*stride+3] = dim[3]
	this.instanceBuffer.cpuBuffer[this.amountLines*stride+4] = p1[0]
	this.instanceBuffer.cpuBuffer[this.amountLines*stride+5] = p1[1]
	this.instanceBuffer.cpuBuffer[this.amountLines*stride+6] = p2[0]
	this.instanceBuffer.cpuBuffer[this.amountLines*stride+7] = p2[1]
	this.instanceBuffer.cpuBuffer[this.amountLines*stride+8] = cp[0]
	this.instanceBuffer.cpuBuffer[this.amountLines*stride+9] = cp[1]

	this.amountLines++
}

func (this *BezierShader) draw() {
	this.shader.use()
	this.shader.setUniformMatrix4("projection", this.projection)
	gl.Disable(gl.DEPTH_TEST)
	gl.BindVertexArray(this.vao)

	this.instanceBuffer.copyToGPU()

	gl.DrawElementsInstanced(gl.TRIANGLES, 6, gl.UNSIGNED_BYTE, gl.Ptr(this.indices), int32(this.amountLines))
	gl.Enable(gl.DEPTH_TEST)
}
