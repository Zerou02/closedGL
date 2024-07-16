package closedGL

import (
	"github.com/EngoEngine/glm"
	"github.com/go-gl/gl/v4.1-core/gl"
)

type LineArr struct {
	lineShader        *Shader
	projection        *glm.Mat4
	indices           []uint16
	vao               uint32
	amountPoints      int
	masterVBO         BufferFloat
	currLastColourIdx int
}

func NewLineArr(shader *Shader, projection *glm.Mat4) LineArr {
	var lineArr = LineArr{lineShader: shader, projection: projection, vao: 0, currLastColourIdx: 0}
	lineArr.generateBuffers()
	return lineArr
}

func (this *LineArr) beginDraw() {
	this.indices = []uint16{}

	this.amountPoints = 0
	this.currLastColourIdx = 0
}

func (this *LineArr) generateBuffers() {
	this.vao = genVAO()
	this.masterVBO = BufferFloat{
		buffer:     generateInterleavedVBOFloat(this.vao, 0, []int{2, 4}),
		bufferSize: 0,
		cpuArr:     []float32{},
	}
}

func (this *LineArr) deleteBuffers() {
	gl.DeleteBuffers(1, &this.vao)
	gl.DeleteBuffers(1, &this.masterVBO.buffer)

}

func (this *LineArr) Draw() {
	if this.amountPoints < 1 {
		return
	}

	this.lineShader.use()
	this.lineShader.setUniformMatrix4("projection", this.projection)

	gl.BindVertexArray(this.vao)
	this.masterVBO.copyToGPU()

	gl.DrawElements(gl.LINES, int32(len(this.indices)), gl.UNSIGNED_SHORT, gl.Ptr(this.indices))
}

func (this *LineArr) addPoint(pos glm.Vec2, colour glm.Vec4) {
	this.masterVBO.resizeCPUData((this.amountPoints + 1) * 6)

	this.indices = append(this.indices, uint16(this.amountPoints))

	this.masterVBO.cpuArr[this.amountPoints*6+0] = pos[0]
	this.masterVBO.cpuArr[this.amountPoints*6+1] = pos[1]
	this.masterVBO.cpuArr[this.amountPoints*6+2] = colour[0]
	this.masterVBO.cpuArr[this.amountPoints*6+3] = colour[1]
	this.masterVBO.cpuArr[this.amountPoints*6+4] = colour[2]
	this.masterVBO.cpuArr[this.amountPoints*6+5] = colour[3]
	this.amountPoints++
}

func (this *LineArr) addLine(pos1, pos2 glm.Vec2, colour1, colour2 glm.Vec4) {
	this.addPoint(pos1, colour1)
	this.addPoint(pos2, colour2)
}

func (this *LineArr) AddPath(pos []glm.Vec2, colours []glm.Vec4) {
	if len(pos) != len(colours) {
		return
	}
	this.addPoint(pos[0], colours[0])
	for i := 1; i < len(pos)-1; i++ {
		this.indices = append(this.indices, uint16(this.amountPoints))
		this.addPoint(pos[i], colours[i])
	}
	this.addPoint(pos[len(pos)-1], colours[len(pos)-1])
}

func (this *LineArr) AddQuadraticBezier(p1, p2, controlPoint glm.Vec2, colour glm.Vec4) {
	var path = []glm.Vec2{}
	var colours = []glm.Vec4{}
	var maxPoints = 20
	for i := 0; i < maxPoints; i++ {
		var t = float32(i) / float32(maxPoints-1)
		path = append(path, BezierLerp(p1, p2, controlPoint, t))
		colours = append(colours, colour)
	}
	this.AddPath(path, colours)
}

func (this *LineArr) AddQuadraticBezierLerp(p1, p2, controlPoint glm.Vec2, colour1, colour2 glm.Vec4) {
	var path = []glm.Vec2{}
	var colours = []glm.Vec4{}
	var maxPoints = 20
	for i := 0; i < maxPoints; i++ {
		var t = float32(i) / float32(maxPoints-1)
		path = append(path, BezierLerp(p1, p2, controlPoint, t))
		colours = append(colours, glm.Vec4{
			Lerp(colour1[0], colour2[0], t),
			Lerp(colour1[1], colour2[1], t),
			Lerp(colour1[2], colour2[2], t),
			Lerp(colour1[3], colour2[3], t),
		})
		this.AddPath(path, colours)
	}
}
