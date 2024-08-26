package closedGL

import (
	"github.com/EngoEngine/glm"
	"github.com/go-gl/gl/v4.6-core/gl"
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

type SSBOU8 struct {
	buffer       uint32
	cpuArr       []uint8
	bindingPoint uint32
}

type SSBOU64 struct {
	buffer       uint32
	cpuArr       []uint64
	bindingPoint uint32
}

type SSBOU32 struct {
	buffer       uint32
	cpuArr       []uint32
	bindingPoint uint32
}

type BufferFloat struct {
	buffer     uint32
	bufferSize int
	cpuArr     []float32
}

type BufferU8 struct {
	buffer     uint32
	bufferSize int
	cpuArr     []uint8
}

type BufferU16 struct {
	buffer     uint32
	bufferSize int
	cpuArr     []uint16
}

type BufferU32 struct {
	buffer     uint32
	bufferSize int
	cpuArr     []uint32
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

func genSSBOU64(bindingPoint uint32) SSBOU64 {
	return SSBOU64{
		buffer:       genSSBO(),
		cpuArr:       []uint64{},
		bindingPoint: bindingPoint,
	}
}

func (this *SSBOU64) clear() {
	this.cpuArr = []uint64{}
}

func (this *SSBOU64) resizeCPUData(newLenEntries int) {
	extendArrayU64(&this.cpuArr, newLenEntries)
}

func (this *SSBOU64) copyToGPU() {
	gl.BindBufferBase(gl.SHADER_STORAGE_BUFFER, this.bindingPoint, this.buffer)
	gl.BufferData(gl.SHADER_STORAGE_BUFFER, len(this.cpuArr)*8, gl.Ptr(this.cpuArr), gl.DYNAMIC_DRAW)
}

func genSSBOU32(bindingPoint uint32) SSBOU32 {
	return SSBOU32{
		buffer:       genSSBO(),
		cpuArr:       []uint32{},
		bindingPoint: bindingPoint,
	}
}

func (this *SSBOU32) clear() {
	this.cpuArr = []uint32{}
}

func (this *SSBOU32) resizeCPUData(newLenEntries int) {
	extendArrayU32(&this.cpuArr, newLenEntries)
}

func (this *SSBOU32) copyToGPU() {
	gl.BindBufferBase(gl.SHADER_STORAGE_BUFFER, this.bindingPoint, this.buffer)
	gl.BufferData(gl.SHADER_STORAGE_BUFFER, len(this.cpuArr)*4, gl.Ptr(this.cpuArr), gl.DYNAMIC_DRAW)
}

func (this *SSBOU32) copy() SSBOU32 {
	var newArr = make([]uint32, len(this.cpuArr))
	copy(newArr, this.cpuArr)
	return SSBOU32{
		buffer:       this.buffer,
		bindingPoint: this.bindingPoint,
		cpuArr:       newArr,
	}
}

func genSSBOU8(bindingPoint uint32) SSBOU8 {
	return SSBOU8{
		buffer:       genSSBO(),
		cpuArr:       []uint8{},
		bindingPoint: bindingPoint,
	}
}

func (this *SSBOU8) clear() {
	this.cpuArr = []uint8{}
}

func (this *SSBOU8) resizeCPUData(newLenEntries int) {
	extendArrayU8(&this.cpuArr, newLenEntries)
}

func (this *SSBOU8) copyToGPU() {
	gl.BindBufferBase(gl.SHADER_STORAGE_BUFFER, this.bindingPoint, this.buffer)
	gl.BufferData(gl.SHADER_STORAGE_BUFFER, len(this.cpuArr)*1, gl.Ptr(this.cpuArr), gl.DYNAMIC_DRAW)
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

func (this *BufferFloat) copy() BufferFloat {
	var newArr = make([]float32, len(this.cpuArr))
	copy(newArr, this.cpuArr)
	return BufferFloat{
		buffer:     this.buffer,
		bufferSize: this.bufferSize,
		cpuArr:     newArr,
	}
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

func (this *BufferU32) resizeCPUData(newLenEntries int) {
	extendArrayU32(&this.cpuArr, newLenEntries)
}

func (this *BufferU32) copyToGPU() {
	setVerticesInVboU32(&this.cpuArr, &this.bufferSize, this.buffer)
}
func (this *BufferU32) clear() {
	this.cpuArr = []uint32{}
}

func (this *BufferU32) copy() BufferU32 {
	var newArr = make([]uint32, len(this.cpuArr))
	copy(newArr, this.cpuArr)
	return BufferU32{
		buffer:     this.buffer,
		bufferSize: this.bufferSize,
		cpuArr:     newArr,
	}
}
func (this *BufferU8) resizeCPUData(newLenEntries int) {
	extendArrayU8(&this.cpuArr, newLenEntries)
}

func (this *BufferU8) copyToGPU() {
	setVerticesInVboU8(&this.cpuArr, &this.bufferSize, this.buffer)
}

func (this *BufferU8) clear() {
	this.cpuArr = []uint8{}
}

func (this *BufferU8) copy() BufferU8 {
	var newArr = make([]uint8, len(this.cpuArr))
	copy(newArr, this.cpuArr)
	return BufferU8{
		buffer:     this.buffer,
		bufferSize: this.bufferSize,
		cpuArr:     newArr,
	}
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

func generateInterleavedVBOFloat2(vao uint32, startIdx int, vertexAttribBytes []int) BufferFloat {
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
	return BufferFloat{
		buffer:     vbo,
		bufferSize: 0,
		cpuArr:     []float32{},
	}
}

func generateInterleavedVBOU32(vao uint32, startIdx int, vertexAttribBytes []int) BufferU32 {
	gl.BindVertexArray(vao)
	var vbo uint32
	gl.GenBuffers(1, &vbo)
	gl.BindBuffer(gl.ARRAY_BUFFER, vbo)

	var stride = 0
	for i := 0; i < len(vertexAttribBytes); i++ {
		stride += int(vertexAttribBytes[i])
	}

	var bytePerType = 4
	var currOffset = 0
	for i := startIdx; i < startIdx+len(vertexAttribBytes); i++ {
		gl.EnableVertexAttribArray(uint32(i))
		gl.VertexAttribIPointerWithOffset(uint32(i), int32(vertexAttribBytes[i-startIdx]), gl.UNSIGNED_INT, int32(stride*bytePerType), uintptr(currOffset*bytePerType))
		currOffset += vertexAttribBytes[i-startIdx]
	}
	return BufferU32{
		buffer:     vbo,
		bufferSize: 0,
		cpuArr:     []uint32{},
	}
}

func generateInterleavedVBOU8(vao uint32, startIdx int, vertexAttribBytes []int) BufferU8 {
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
		gl.VertexAttribIPointerWithOffset(uint32(i), int32(vertexAttribBytes[i-startIdx]), gl.UNSIGNED_BYTE, int32(stride*1), uintptr(currOffset)*1)
		currOffset += vertexAttribBytes[i-startIdx]
	}
	return BufferU8{
		buffer:     vbo,
		bufferSize: 0,
		cpuArr:     []uint8{},
	}
}

func setVerticesInVbo(vertices *[]float32, vboSizeEntries *int, vbo uint32) {
	if len(*vertices) == 0 {
		return
	}
	gl.BindBuffer(gl.ARRAY_BUFFER, vbo)
	if len(*vertices)*4 >= *vboSizeEntries {
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
	if len(*vertices)*bytesPerEntry >= *vboSizeEntries {
		*vboSizeEntries = len(*vertices) * bytesPerEntry
		gl.BufferData(gl.ARRAY_BUFFER, *vboSizeEntries, gl.Ptr(*vertices), gl.DYNAMIC_DRAW)
	} else {
		gl.BufferSubData(gl.ARRAY_BUFFER, 0, len(*vertices)*bytesPerEntry, gl.Ptr(*vertices))
	}
}

func setVerticesInVboU32(vertices *[]uint32, vboSizeEntries *int, vbo uint32) {
	var bytesPerEntry = 4
	if len(*vertices) == 0 {
		return
	}
	gl.BindBuffer(gl.ARRAY_BUFFER, vbo)
	if len(*vertices)*bytesPerEntry >= *vboSizeEntries {
		*vboSizeEntries = len(*vertices) * bytesPerEntry
		gl.BufferData(gl.ARRAY_BUFFER, *vboSizeEntries, gl.Ptr(*vertices), gl.DYNAMIC_DRAW)
	} else {
		gl.BufferSubData(gl.ARRAY_BUFFER, 0, len(*vertices)*bytesPerEntry, gl.Ptr(*vertices))
	}
}

func setVerticesInVboU8(vertices *[]uint8, vboSizeEntries *int, vbo uint32) {
	var bytesPerEntry = 1
	if len(*vertices) == 0 {
		return
	}
	gl.BindBuffer(gl.ARRAY_BUFFER, vbo)
	if len(*vertices)*bytesPerEntry >= *vboSizeEntries {
		*vboSizeEntries = len(*vertices) * bytesPerEntry
		gl.BufferData(gl.ARRAY_BUFFER, *vboSizeEntries, gl.Ptr(*vertices), gl.DYNAMIC_DRAW)
	} else {
		gl.BufferSubData(gl.ARRAY_BUFFER, 0, len(*vertices)*bytesPerEntry, gl.Ptr(*vertices))
	}
}
