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
	strides          []uint32
	buffer           []*IBuffer
	indices          []uint16
	dirty            bool
	update           bool
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

func (this *SimpleMesh) InitBuffer(indices []uint16, vertexAttribBytes [][]int, divisorValues [][]int, types []uint32) {
	this.indices = indices
	this.strides = make([]uint32, len(types))
	this.buffer = make([]*IBuffer, len(types))
	var amountEntries = 0
	for i := 0; i < len(types); i++ {
		this.strides[i] = uint32(ArrSum(&vertexAttribBytes[i]))
		var b = genInterleavedBufferGen(this.vao, amountEntries, vertexAttribBytes[i], divisorValues[i], types[i])
		this.buffer[i] = &b
		amountEntries += len(vertexAttribBytes[i])
	}
}

/***
* Set dirty, um im nächsten Frame Updates zu machen
 */
func (this *SimpleMesh) Draw() {
	if this.dirty {
		this.CopyToGPU()
	}
	if this.amountElements < 1 {
		return
	}

	gl.Disable(gl.DEPTH_TEST)
	this.shader.use()
	this.shader.setUniformMatrix4("projection", &this.projection)
	this.shader.setUniformMatrix4("view", &this.view)

	gl.BindVertexArray(this.vao)
	if len(this.buffer) > 1 {
		if len(this.indices) > 0 {
			gl.DrawElementsInstanced(gl.TRIANGLES, int32(len(this.indices)), gl.UNSIGNED_SHORT, gl.Ptr(this.indices), int32(this.amountElements))
		} else {
			println("not impl")
		}
	} else {
		if len(this.indices) > 0 {

			gl.DrawElements(gl.LINES, int32(len(this.indices)), gl.UNSIGNED_SHORT, gl.Ptr(this.indices))
		} else {
			println("not impl")
		}
	}

	gl.Enable(gl.DEPTH_TEST)
	if this.update {
		this.update = false
	}
	if this.dirty {
		this.update = true
		this.dirty = false
	}
}

func (this *SimpleMesh) Clear() {
	for _, x := range this.buffer {
		(*x).clear()
	}
	this.amountElements = 0
	this.CopyToGPU()
}

func (this *SimpleMesh) CopyToGPU() {
	gl.BindVertexArray(this.vao)
	for _, x := range this.buffer {
		(*x).copyToGPU()
	}
}

func (this *SimpleMesh) setDirty() {
	this.dirty = true
}

func (this *SimpleMesh) isUpdate() bool {
	return this.update
}

func addVertices(mesh *SimpleMesh, vertices []*[]any, indices *[]uint16) {
	for i, x := range mesh.buffer {
		var vertices = (vertices)[i]
		if len(*vertices) == 0 {
			continue
		}
		(*x).addVertices(vertices, mesh.strides[i], mesh.amountElements)
	}
	mesh.indices = append(mesh.indices, *indices...)
	mesh.amountElements++
	mesh.setDirty()
}
