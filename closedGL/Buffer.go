package closedGL

import (
	"github.com/go-gl/gl/v4.6-core/gl"
)

type IBuffer interface {
	clear()
	copyToGPU()
	resizeCPUData(newLen int)
	//copy() *IBuffer
	addVertices(vertices *[]any, stride uint32, amountElements uint32)
}
type GLType interface {
	float32 | uint32 | uint64 | int
}

type Buffer[T GLType] struct {
	cpuBuffer         []T
	gpuBuffer         uint32
	bufferSize        int
	sizeOfTypeInBytes int32
}

type SSBO[T GLType] struct {
	gpuBuffer        uint32
	cpuBuffer        []T
	bindingPoint     uint32
	elementsPerEntry int32
}

// dataType = gl.unsigned_short bspw.
func genSingularBuffer[T GLType](vao, index uint32, elementsPerEntry int32, dataType uint32, normalized bool, instanceCount uint32) Buffer[T] {
	var buf = Buffer[T]{
		cpuBuffer:         []T{},
		gpuBuffer:         genSingularVBO(vao, index, elementsPerEntry, dataType, normalized, instanceCount),
		bufferSize:        0,
		sizeOfTypeInBytes: getByteLenOfGLDataType(dataType),
	}
	buf.copyToGPU()
	return buf
}

func genInterleavedBufferGen(vao uint32, startIdx int, vertexAttribBytes []int, divisorValues []int, glType uint32) IBuffer {
	if glType == gl.FLOAT {
		var b = genInterleavedBuffer[float32](vao, startIdx, vertexAttribBytes, divisorValues, glType)
		return &b
	} else if glType == gl.UNSIGNED_INT {
		var b = genInterleavedBuffer[uint32](vao, startIdx, vertexAttribBytes, divisorValues, glType)
		return &b
	} else {
		panic("Not implemented")
	}
}
func genInterleavedBuffer[T GLType](vao uint32, startIdx int, vertexAttribBytes []int, divisorValues []int, glType uint32) Buffer[T] {
	var buf = Buffer[T]{
		cpuBuffer:         []T{},
		gpuBuffer:         generateInterleavedVBO(vao, startIdx, vertexAttribBytes, divisorValues, glType),
		bufferSize:        0,
		sizeOfTypeInBytes: getByteLenOfGLDataType(glType),
	}
	buf.copyToGPU()
	return buf
}

func getByteLenOfGLDataType(glDataType uint32) int32 {
	var dataSizes map[uint32]int32 = map[uint32]int32{}
	dataSizes[gl.UNSIGNED_BYTE] = 1
	dataSizes[gl.UNSIGNED_SHORT] = 2
	dataSizes[gl.UNSIGNED_INT] = 4
	dataSizes[gl.INT] = 4
	dataSizes[gl.FLOAT] = 4
	dataSizes[gl.UNSIGNED_INT64_ARB] = 8
	return dataSizes[glDataType]
}

func isGLIntType(dataType uint32) bool {
	return dataType == gl.UNSIGNED_SHORT || dataType == gl.UNSIGNED_BYTE || dataType == gl.UNSIGNED_INT || dataType == gl.INT || dataType == gl.UNSIGNED_INT64_ARB

}

// dataType = gl.unsigned_short bspw.
func genSingularVBOGeneric(vao, index uint32, elementsPerEntry int32, dataType uint32, normalized bool, instanceCount uint32) uint32 {
	var vbo uint32 = 0
	gl.BindVertexArray(vao)

	gl.GenBuffers(1, &vbo)
	gl.EnableVertexAttribArray(index)
	gl.BindBuffer(gl.ARRAY_BUFFER, vbo)

	var dataSize = getByteLenOfGLDataType(dataType)
	var isInt = isGLIntType(dataType)

	if isInt {
		gl.VertexAttribIPointerWithOffset(index, elementsPerEntry, dataType, dataSize*elementsPerEntry, 0)
	} else {
		gl.VertexAttribPointerWithOffset(index, elementsPerEntry, dataType, normalized, dataSize*elementsPerEntry, 0)
	}
	gl.VertexAttribDivisor(index, instanceCount)
	return vbo
}

// Kein Plan, ob das hier funktioniert...
// glType = gl.float,...
// startIdx = shader layout locations
// vertexBytes = elementsPerEntry?
func generateInterleavedVBO(vao uint32, startIdx int, vertexAttribBytes []int, divisorValues []int, glType uint32) uint32 {
	gl.BindVertexArray(vao)
	var vbo uint32
	gl.GenBuffers(1, &vbo)
	gl.BindBuffer(gl.ARRAY_BUFFER, vbo)

	var stride int32 = 0
	for i := 0; i < len(vertexAttribBytes); i++ {
		stride += int32(vertexAttribBytes[i])
	}

	var currOffset int32 = 0
	var isInt = isGLIntType(glType)

	var typeByteSize = getByteLenOfGLDataType(glType)
	for i := startIdx; i < startIdx+len(vertexAttribBytes); i++ {
		gl.EnableVertexAttribArray(uint32(i))
		if isInt {
			gl.VertexAttribIPointerWithOffset(uint32(i), int32(vertexAttribBytes[i-startIdx]), glType, stride*typeByteSize, uintptr(currOffset*typeByteSize))
		} else {
			gl.VertexAttribPointerWithOffset(uint32(i), int32(vertexAttribBytes[i-startIdx]), glType, false, stride*typeByteSize, uintptr(currOffset*typeByteSize))
		}
		gl.VertexAttribDivisor(uint32(i), uint32(divisorValues[i-startIdx]))
		currOffset += int32(vertexAttribBytes[i-startIdx])
	}
	return vbo
}

func setVerticesInVboGen[T GLType](vertices *[]T, vboSizeEntries *int, vbo uint32, dataTypeSizeInBytes int) {
	if len(*vertices) == 0 {
		return
	}
	gl.BindBuffer(gl.ARRAY_BUFFER, vbo)
	if len(*vertices)*dataTypeSizeInBytes >= *vboSizeEntries {
		*vboSizeEntries = len(*vertices) * dataTypeSizeInBytes
		gl.BufferData(gl.ARRAY_BUFFER, *vboSizeEntries, gl.Ptr(*vertices), gl.DYNAMIC_DRAW)
	} else {
		gl.BufferSubData(gl.ARRAY_BUFFER, 0, len(*vertices)*dataTypeSizeInBytes, gl.Ptr(*vertices))
	}
}
func extendArrayGen[T GLType](arr *[]T, newLenEntries int) {
	if newLenEntries != 0 && len(*arr) == 0 {
		*arr = make([]T, 1)
	}
	for newLenEntries >= len(*arr) {
		var newArr = make([]T, len(*arr)*2)
		copy(newArr, *arr)
		*arr = newArr
	}
}

// u64 = gl.unsigned_int64_arb
func genSSBOGen[T GLType](bindingPoint uint32, dType uint32) SSBO[T] {

	return SSBO[T]{
		gpuBuffer:        genSSBO(),
		cpuBuffer:        []T{},
		bindingPoint:     bindingPoint,
		elementsPerEntry: getByteLenOfGLDataType(dType),
	}
}

func (this *SSBO[T]) clear() {
	this.cpuBuffer = []T{}
}

func (this *SSBO[T]) resizeCPUData(newLenEntries int) {
	extendArrayGen(&this.cpuBuffer, newLenEntries)
}

func (this *SSBO[T]) copyToGPU() {
	gl.BindBufferBase(gl.SHADER_STORAGE_BUFFER, this.bindingPoint, this.gpuBuffer)
	gl.BufferData(gl.SHADER_STORAGE_BUFFER, len(this.cpuBuffer)*int(this.elementsPerEntry), gl.Ptr(this.cpuBuffer), gl.DYNAMIC_DRAW)
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
func genSSBO() uint32 {
	var ssbo uint32 = 0
	gl.CreateBuffers(1, &ssbo)
	return ssbo
}

func genVAO() uint32 {
	var vao uint32 = 0
	gl.GenVertexArrays(1, &vao)
	gl.BindVertexArray(vao)
	return vao
}

func (this *Buffer[T]) resizeCPUData(newLenEntries int) {
	extendArrayGen[T](&this.cpuBuffer, newLenEntries)
}

func (this *Buffer[T]) copyToGPU() {
	setVerticesInVboGen(&this.cpuBuffer, &this.bufferSize, this.gpuBuffer, int(this.sizeOfTypeInBytes))
}

func (this *Buffer[T]) clear() {
	this.cpuBuffer = []T{}
	this.copyToGPU()
}

/* func (this *Buffer[T]) copy() IBuffer {
	var newArr = make([]T, len(this.cpuBuffer))
	copy(newArr, this.cpuBuffer)
	var buf = Buffer[T]{
		cpuBuffer:         newArr,
		gpuBuffer:         this.gpuBuffer,
		bufferSize:        this.bufferSize,
		sizeOfTypeInBytes: this.sizeOfTypeInBytes,
	}
	return buf
} */

func (this *Buffer[T]) addVertices(vertices *[]any, stride uint32, amountElements uint32) {
	this.resizeCPUData(int((amountElements + 1) * stride))
	for i := 0; i < len(*vertices); i++ {
		this.cpuBuffer[int(amountElements*stride)+i] = (*vertices)[i].(T) //type assertion
	}
}
