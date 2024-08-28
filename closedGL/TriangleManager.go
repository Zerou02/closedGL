package closedGL

import (
	"github.com/EngoEngine/glm"
	"github.com/go-gl/gl/v4.6-core/gl"
)

type TriangleManager struct {
	shader     *Shader
	projection *glm.Mat4
	vao        uint32
	//vec2-pos;vec4-colour
	dataBuffer      BufferFloat
	amountTriangles int
}

func newTriangleManager(shader *Shader, projection *glm.Mat4) TriangleManager {
	var retTri = TriangleManager{shader: shader, projection: projection}
	retTri.generateBuffers()
	return retTri
}

func (this *TriangleManager) generateBuffers() {
	this.vao = genVAO()
	this.dataBuffer = BufferFloat{
		buffer:     generateInterleavedVBOFloat(this.vao, 0, []int{2, 4}, []int{0, 0}),
		bufferSize: 0,
		cpuArr:     []float32{},
	}
}

func (this *TriangleManager) beginDraw() {
	this.dataBuffer.clear()
	this.amountTriangles = 0
}

func (this *TriangleManager) createVertices(pos [3]glm.Vec2, colour glm.Vec4) {

	const stride = 6
	const amountVertices = 3
	this.dataBuffer.resizeCPUData((this.amountTriangles + 1) * stride * amountVertices)
	for i := 0; i < amountVertices; i++ {
		this.dataBuffer.cpuArr[this.amountTriangles*stride+0+i*stride] = pos[i][0]
		this.dataBuffer.cpuArr[this.amountTriangles*stride+1+i*stride] = pos[i][1]
		this.dataBuffer.cpuArr[this.amountTriangles*stride+2+i*stride] = colour[0]
		this.dataBuffer.cpuArr[this.amountTriangles*stride+3+i*stride] = colour[1]
		this.dataBuffer.cpuArr[this.amountTriangles*stride+4+i*stride] = colour[2]
		this.dataBuffer.cpuArr[this.amountTriangles*stride+5+i*stride] = colour[3]
	}
	this.amountTriangles++
}

func (this *TriangleManager) draw() {
	this.shader.use()
	gl.Disable(gl.CULL_FACE)
	gl.Disable(gl.DEPTH_TEST)

	this.shader.setUniformMatrix4("projection", this.projection)
	gl.BindVertexArray(this.vao)
	this.dataBuffer.copyToGPU()
	gl.DrawArrays(gl.TRIANGLES, 0, int32(6*this.amountTriangles))
	gl.Enable(gl.CULL_FACE)
	gl.Enable(gl.DEPTH_TEST)
}
