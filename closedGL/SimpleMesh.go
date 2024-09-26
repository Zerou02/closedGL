package closedGL

import (
	"github.com/EngoEngine/glm"
	"github.com/go-gl/gl/v4.6-core/gl"
)

type SimpleMesh struct {
	shader           *Shader
	vao              uint32
	amountElements   uint32
	projection, view glm.Mat4
	//	instanceBuffer   Buffer[float32]
	buffer  Buffer[float32]
	indices []uint16
}

func newSimpleMesh(shader *Shader, projection, view glm.Mat4) SimpleMesh {
	var mesh = SimpleMesh{
		shader:         shader,
		vao:            genVAO(),
		amountElements: 0,
		projection:     projection,
		view:           view,
	}
	return mesh
}

func (this *SimpleMesh) InitBuffer(indices []uint16, vertexAttribBytes []int, divisorValues []int) {
	this.indices = indices
	this.buffer = genInterleavedBuffer[float32](this.vao, 0, vertexAttribBytes, divisorValues, gl.FLOAT)
}
