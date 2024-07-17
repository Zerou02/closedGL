package closedGL

import (
	"github.com/EngoEngine/glm"
	"github.com/go-gl/gl/v4.1-core/gl"
)

type CircleManager struct {
	shader         *Shader
	projection     *glm.Mat4
	vao            uint32
	amountCircles  uint32
	baseVBO        BufferFloat
	instanceBuffer BufferFloat
	indices        []byte
}

func newCircleManger(shader *Shader, projection *glm.Mat4) CircleManager {
	var indices = []byte{0, 1, 2, 2, 1, 3}
	var circleMan = CircleManager{shader: shader, projection: projection, amountCircles: 0, indices: indices}
	circleMan.generateBuffers()

	return circleMan
}

func (this *CircleManager) generateBuffers() {
	this.vao = genVAO()
	gl.BindVertexArray(0)
	this.baseVBO = genSingularBufferFloat(this.vao, 0, 2, gl.FLOAT, false, 0)
	this.instanceBuffer = BufferFloat{
		buffer: generateInterleavedVBOFloat(this.vao, 1, []int{4, 4, 4}), //centre,colour
	}
	gl.VertexAttribDivisor(1, 1)
	gl.VertexAttribDivisor(2, 1)
	gl.VertexAttribDivisor(3, 1)

	var quadBaseData = []float32{
		1.0, 0.0, //top r
		0.0, 0.0, // top l
		1.0, 1.0, // bottom r
		0.0, 1.0, // bottom l,
	}
	this.baseVBO.cpuArr = quadBaseData
	this.baseVBO.copyToGPU()
}

func (this *CircleManager) beginDraw() {
	this.amountCircles = 0
	this.instanceBuffer.clear()
}

func (this *CircleManager) deleteBuffers() {
	gl.DeleteBuffers(1, &this.vao)
	gl.DeleteBuffers(1, &this.baseVBO.buffer)
	gl.DeleteBuffers(1, &this.instanceBuffer.buffer)

}

func (this *CircleManager) createVertices(centre glm.Vec2, colour, borderColour glm.Vec4, radius, borderThickness float32) {
	var stride uint32 = 12

	this.instanceBuffer.resizeCPUData((int(this.amountCircles) + 1) * int(stride))

	this.instanceBuffer.cpuArr[this.amountCircles*stride+0] = centre[0]
	this.instanceBuffer.cpuArr[this.amountCircles*stride+1] = centre[1]
	this.instanceBuffer.cpuArr[this.amountCircles*stride+2] = radius
	this.instanceBuffer.cpuArr[this.amountCircles*stride+3] = borderThickness

	this.instanceBuffer.cpuArr[this.amountCircles*stride+4] = colour[0]
	this.instanceBuffer.cpuArr[this.amountCircles*stride+5] = colour[1]
	this.instanceBuffer.cpuArr[this.amountCircles*stride+6] = colour[2]
	this.instanceBuffer.cpuArr[this.amountCircles*stride+7] = colour[3]

	this.instanceBuffer.cpuArr[this.amountCircles*stride+8] = borderColour[0]
	this.instanceBuffer.cpuArr[this.amountCircles*stride+9] = borderColour[1]
	this.instanceBuffer.cpuArr[this.amountCircles*stride+10] = borderColour[2]
	this.instanceBuffer.cpuArr[this.amountCircles*stride+11] = borderColour[3]
	this.amountCircles++
}

func (this *CircleManager) draw() {
	this.shader.use()
	this.shader.setUniformMatrix4("projection", this.projection)
	gl.Disable(gl.DEPTH_TEST)
	gl.BindVertexArray(this.vao)

	this.instanceBuffer.copyToGPU()

	gl.DrawElementsInstanced(gl.TRIANGLES, 6, gl.UNSIGNED_BYTE, gl.Ptr(this.indices), int32(this.amountCircles))
	gl.Enable(gl.DEPTH_TEST)
}
