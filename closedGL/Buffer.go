package closedGL

import (
	"github.com/EngoEngine/glm"
	"github.com/go-gl/gl/v4.1-core/gl"
)

type BufferCompose interface {
	clear()
	copyToGPU()
	resizeCPUData(newLenEntries int)
}

type SSBOVec4 struct {
	buffer       uint32
	cpuArr       []glm.Vec4
	bindingPoint uint32
}

type BufferFloat struct {
	buffer     uint32
	bufferSize int
	cpuArr     []float32
}
type BufferU16 struct {
	buffer     uint32
	bufferSize int
	cpuArr     []uint16
}

func genVAO() uint32 {
	var vao uint32 = 0
	gl.GenVertexArrays(1, &vao)
	gl.BindVertexArray(vao)
	return vao
}

func genSSBO() uint32 {
	var ssbo uint32 = 0
	gl.CreateBuffers(1, &ssbo)
	return ssbo
}

func genSSBOVec4(bindingPoint uint32) SSBOVec4 {
	return SSBOVec4{
		buffer:       genSSBO(),
		cpuArr:       []glm.Vec4{},
		bindingPoint: bindingPoint,
	}
}

func (this *SSBOVec4) clear() {
	this.cpuArr = []glm.Vec4{}
}

func (this *SSBOVec4) resizeCPUData(newLenEntries int) {
	extendArrayVec4(&this.cpuArr, newLenEntries)
}

func (this *SSBOVec4) copyToGPU() {
	gl.BindBufferBase(gl.SHADER_STORAGE_BUFFER, this.bindingPoint, this.buffer)
	gl.BufferData(gl.SHADER_STORAGE_BUFFER, len(this.cpuArr)*4*4, gl.Ptr(this.cpuArr), gl.DYNAMIC_DRAW)

}

func genSingularVBO(vao, index uint32, elementsPerEntry int32, dataType uint32, normalized bool, instanceCount uint32) uint32 {
	var vbo uint32 = 0
	gl.BindVertexArray(vao)

	gl.GenBuffers(1, &vbo)
	gl.EnableVertexAttribArray(index)
	gl.BindBuffer(gl.ARRAY_BUFFER, vbo)
	var dataSizes map[uint32]int32 = map[uint32]int32{}
	dataSizes[gl.UNSIGNED_BYTE] = 1
	dataSizes[gl.UNSIGNED_SHORT] = 2
	dataSizes[gl.FLOAT] = 4
	var dataSize = dataSizes[dataType]

	var isInt = dataType == gl.UNSIGNED_SHORT || dataType == gl.UNSIGNED_BYTE

	if isInt {
		gl.VertexAttribIPointerWithOffset(index, elementsPerEntry, dataType, dataSize*elementsPerEntry, 0)
	} else {
		gl.VertexAttribPointerWithOffset(index, elementsPerEntry, dataType, normalized, dataSize*elementsPerEntry, 0)
	}
	gl.VertexAttribDivisor(index, instanceCount)
	return vbo
}

func genSingularBufferFloat(vao, index uint32, elementsPerEntry int32, dataType uint32, normalized bool, instanceCount uint32) BufferFloat {
	return BufferFloat{
		buffer:     genSingularVBO(vao, index, elementsPerEntry, dataType, normalized, instanceCount),
		bufferSize: 0,
		cpuArr:     []float32{},
	}
}

func genSingularBufferU16(vao, index uint32, elementsPerEntry int32, dataType uint32, normalized bool, instanceCount uint32) BufferU16 {
	return BufferU16{
		buffer:     genSingularVBO(vao, index, elementsPerEntry, dataType, normalized, instanceCount),
		bufferSize: 0,
		cpuArr:     []uint16{},
	}
}

func (this *BufferFloat) resizeCPUData(newLenEntries int) {
	extendArray(&this.cpuArr, newLenEntries)
}

func (this *BufferFloat) copyToGPU() {
	setVerticesInVbo(&this.cpuArr, &this.bufferSize, this.buffer)
}
func (this *BufferFloat) clear() {
	this.cpuArr = []float32{}
}

func (this *BufferU16) resizeCPUData(newLenEntries int) {
	extendArrayU16(&this.cpuArr, newLenEntries)
}

func (this *BufferU16) copyToGPU() {
	setVerticesInVboU16(&this.cpuArr, &this.bufferSize, this.buffer)
}
func (this *BufferU16) clear() {
	this.cpuArr = []uint16{}
}

func generateInterleavedVBOFloat(vao uint32, startIdx int, vertexAttribBytes []int) uint32 {
	gl.BindVertexArray(vao)
	var vbo uint32
	gl.GenBuffers(1, &vbo)
	gl.BindBuffer(gl.ARRAY_BUFFER, vbo)

	var stride = 0
	for i := 0; i < len(vertexAttribBytes); i++ {
		stride += int(vertexAttribBytes[i])
	}

	var currOffset = 0
	for i := startIdx; i < startIdx+len(vertexAttribBytes); i++ {
		gl.EnableVertexAttribArray(uint32(i))
		gl.VertexAttribPointerWithOffset(uint32(i), int32(vertexAttribBytes[i-startIdx]), gl.FLOAT, false, int32(stride*4), uintptr(currOffset)*4)
		currOffset += vertexAttribBytes[i-startIdx]
	}
	return vbo
}

func setVerticesInVbo(vertices *[]float32, vboSizeEntries *int, vbo uint32) {
	if len(*vertices) == 0 {
		return
	}
	gl.BindBuffer(gl.ARRAY_BUFFER, vbo)
	if len(*vertices) >= *vboSizeEntries {
		*vboSizeEntries = len(*vertices) * 4
		gl.BufferData(gl.ARRAY_BUFFER, *vboSizeEntries, gl.Ptr(*vertices), gl.DYNAMIC_DRAW)

	} else {
		gl.BufferSubData(gl.ARRAY_BUFFER, 0, len(*vertices)*4, gl.Ptr(*vertices))
	}
}

func setVerticesInVboU16(vertices *[]uint16, vboSizeEntries *int, vbo uint32) {
	var bytesPerEntry = 2
	if len(*vertices) == 0 {
		return
	}
	gl.BindBuffer(gl.ARRAY_BUFFER, vbo)
	if len(*vertices) >= *vboSizeEntries {
		*vboSizeEntries = len(*vertices) * bytesPerEntry
		gl.BufferData(gl.ARRAY_BUFFER, *vboSizeEntries, gl.Ptr(*vertices), gl.DYNAMIC_DRAW)
	} else {
		gl.BufferSubData(gl.ARRAY_BUFFER, 0, len(*vertices)*bytesPerEntry, gl.Ptr(*vertices))
	}
}
