package turingfontparser

import (
	"os"
	"strconv"
	"unsafe"
)

type fword int16

type Reader struct {
	path string
	file *os.File
}

type GlyfHeader struct {
	nrContours int16
	xMin       fword
	yMin       fword
	xMax       fword
	yMax       fword
}

type SimpleGlyfBody struct {
	endOfContours     []uint16
	instructionLength uint16
	instructions      []uint8
	flags             []uint8
	xCoordinates      []int16
	yCoordinates      []int16
}

func NewReader(path string) Reader {
	return Reader{
		path: path,
	}
}

type DirEntry struct {
	tag      string
	checksum uint32
	offset   uint32
	length   uint32
}

type OffsetSubTable struct {
	scaler        uint32
	numTable      uint16
	searchRange   uint16
	entrySelector uint16
	rangeShift    uint16
}

func (this *Reader) parse() {
	var f, err = os.Open(this.path)
	if err != nil {
		panic(err.Error())
	}
	this.file = f
	var offsetSubTable = this.parseOffsetSubtable()
	var glyfEntry DirEntry
	for i := 0; i < int(offsetSubTable.numTable); i++ {
		var e = this.readDirEntry()
		if e.tag == "glyf" {
			glyfEntry = e
		}
	}
	this.printEntry(&glyfEntry)
	this.readEnty(&glyfEntry)
}

func (this *Reader) parseOffsetSubtable() OffsetSubTable {
	return OffsetSubTable{
		scaler:        this.readUint32(),
		numTable:      this.readUint16(),
		searchRange:   this.readUint16(),
		entrySelector: this.readUint16(),
		rangeShift:    this.readUint16(),
	}
}

func reverseByteArr(arr *[]byte) []byte {
	var newArr = make([]byte, len((*arr)))
	var j = 0
	for i := len(*arr) - 1; i >= 0; i-- {
		newArr[j] = (*arr)[i]
		j++
	}
	return newArr
}
func readBytes(amount int, f *os.File, bigEndian bool) []byte {
	var b = make([]byte, amount)
	f.Read(b)
	if bigEndian {
		return reverseByteArr(&b)
	} else {
		return b
	}
}

// trust me,lad
// @Safe
func readType[T any](f *os.File, amountBytes int, bigEndian bool) T {
	var bytes = readBytes(amountBytes, f, bigEndian)
	var test = (*T)(unsafe.Pointer(&bytes[0]))
	bytes = nil
	return *test
}

func (this *Reader) readUint32() uint32 {
	return readType[uint32](this.file, 4, true)
}

func (this *Reader) readInt16() int16 {
	return readType[int16](this.file, 2, true)
}

func (this *Reader) readUint8() uint8 {
	return readType[uint8](this.file, 1, true)
}

func (this *Reader) readStr(amountChars int) string {
	return string(readBytes(amountChars, this.file, false))
}

func (this *Reader) readUint16() uint16 {
	return readType[uint16](this.file, 2, true)
}

func (this *Reader) readDirEntry() DirEntry {
	return DirEntry{
		tag:      this.readStr(4),
		checksum: this.readUint32(),
		offset:   this.readUint32(),
		length:   this.readUint32(),
	}
}

func (this *Reader) printEntry(entry *DirEntry) {
	println(entry.tag)
	println(entry.checksum)
	println(entry.offset)
	println(entry.length)
}

func (this *Reader) calcChecksum(e *DirEntry) uint32 {
	this.file.Seek(int64(e.offset), 0)
	var sum uint32 = 0
	var nrLongs = (e.length + 3) / 4
	for nrLongs > 0 {
		sum += this.readUint32()
		nrLongs--
	}
	return sum
}

func (this *Reader) readGlyfHeader() GlyfHeader {
	return GlyfHeader{
		nrContours: this.readInt16(),
		xMin:       fword(this.readInt16()),
		yMin:       fword(this.readInt16()),
		xMax:       fword(this.readInt16()),
		yMax:       fword(this.readInt16()),
	}
}

func (this *Reader) readSimpleGlyph(amountContours int16) SimpleGlyfBody {
	var body = SimpleGlyfBody{}
	var endPtsOfContours = make([]uint16, amountContours)
	for i := 0; i < len(endPtsOfContours); i++ {
		endPtsOfContours[i] = this.readUint16()
	}
	body.endOfContours = endPtsOfContours
	var nrPoints = body.endOfContours[amountContours-1] + 1
	body.instructionLength = this.readUint16()
	var instructions = make([]uint8, body.instructionLength)
	for i := 0; i < len(instructions); i++ {
		instructions[i] = this.readUint8()
	}
	body.instructions = instructions
	var flags = make([]uint8, nrPoints)
	for i := 0; i < int(nrPoints); i++ {
		flags[i] = this.readUint8()
		//repeat
		if flags[i]&0b1000 == 1 {
			println("rep")
			var repetitions = this.readUint8()
			var base = flags[i]
			for j := 0; j < int(repetitions); j++ {
				i++
				flags[i] = base
			}
		}
	}
	body.flags = flags
	var xCoordinates = make([]int16, nrPoints)
	var yCoordinates = make([]int16, nrPoints)
	for i := 0; i < len(xCoordinates); i++ {
		if (flags[i]>>1)&1 == 1 {
			xCoordinates[i] = int16(this.readUint8())
		} else {
			xCoordinates[i] = this.readInt16()
		}
	}
	for i := 0; i < len(yCoordinates); i++ {
		if (flags[i]>>2)&1 == 1 {
			println("FFLL", i)
			yCoordinates[i] = int16(this.readUint8())
		} else {
			yCoordinates[i] = this.readInt16()
		}
	}
	body.xCoordinates = xCoordinates
	body.yCoordinates = yCoordinates
	return body
}

func (this *Reader) printGlfyHeader(h GlyfHeader) {
	println(h.nrContours)
	println(h.xMin)
	println(h.yMin)
	println(h.xMax)
	println(h.yMax)
}

func (this *Reader) printGlyfBody(b SimpleGlyfBody) {
	println(b.endOfContours)
	println("amountPoints", b.endOfContours[len(b.endOfContours)-1]+1)
	println(b.instructionLength)
	println(b.instructions)
	println(b.flags)
	println(b.xCoordinates)
	println(b.yCoordinates)
	println("endPoints")
	for i := 0; i < len(b.endOfContours); i++ {
		println(b.endOfContours[i])
	}
	println("points")
	for i := 0; i < len(b.xCoordinates); i++ {
		println("flag", strconv.FormatInt(int64(b.flags[i]), 2))
		println("points", b.xCoordinates[i], ",", b.yCoordinates[i])
	}
	var test uint8 = 3
	println(int16(test))
}
func (this *Reader) readEnty(e *DirEntry) {
	println("---read entry--")
	if e.checksum != this.calcChecksum(e) {
		panic("Wrong checksum in table" + e.tag)
	}
	this.file.Seek(int64(e.offset), 0)
	var header = this.readGlyfHeader()
	this.printGlfyHeader(header)
	this.printGlyfBody(this.readSimpleGlyph(header.nrContours))
}
