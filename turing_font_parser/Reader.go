package turingfontparser

import (
	"math"
	"os"
	"unsafe"

	"github.com/EngoEngine/glm"
	"github.com/Zerou02/closedGL/closedGL"
)

type fword int16

type Reader struct {
	path              string
	file              *os.File
	entries           map[string]*DirEntry
	offset            uint32
	mappings          map[int]int
	smallLocATable    bool
	nrGlyphs          uint16
	loca              []uint32
	ctx               *closedGL.ClosedGLContext
	nrHMetrics        uint16
	horizMetrics      []LongHorizMetric
	remainingBearings []int16
}

type LongHorizMetric struct {
	advanceWidth uint16
	lsb          int16 //left-side-bearing
}

func NewReader(path string, ctx *closedGL.ClosedGLContext) Reader {
	return Reader{
		path:         path,
		entries:      map[string]*DirEntry{},
		offset:       0,
		loca:         []uint32{},
		ctx:          ctx,
		horizMetrics: []LongHorizMetric{},
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

func (this *Reader) init() {
	var f, err = os.Open(this.path)
	if err != nil {
		panic(err.Error())
	}
	this.file = f
	var offsetSubTable = this.parseOffsetSubtable()
	for i := 0; i < int(offsetSubTable.numTable); i++ {
		var e = this.readDirEntry()
		this.entries[e.tag] = &e
	}
	this.parseMaxP()
	this.parseHead()
	this.parseHhea()
	this.parseHmtx()
	this.parseCmap()
	this.parseLocA()
	this.seek(this.entries["glyf"].offset)

}

func (this *Reader) parseCmap() {
	var entry = this.entries["cmap"]
	if this.calcChecksum(entry) != entry.checksum {
		panic("invalid checksum cmap")
	}
	this.seek(entry.offset)
	this.readUint16()
	var subtablesNr = this.readUint16()
	var subTables = [][3]uint32{} //
	for i := 0; i < int(subtablesNr); i++ {
		var subtable = [3]uint32{
			uint32(this.readUint16()),
			uint32(this.readUint16()),
			this.readUint32(),
		}
		subTables = append(subTables, subtable)
	}
	var biggestUnicode = -1
	var unicodeIdx = -1
	for i, x := range subTables {
		if x[0] == 0 {
			if int(x[1]) > biggestUnicode {
				biggestUnicode = int(x[1])
				unicodeIdx = i
			}
		}
	}
	if biggestUnicode == -1 {
		for i, x := range subTables {
			if x[0] == 3 {
				if int(x[1]) > biggestUnicode {
					biggestUnicode = int(x[1])
					unicodeIdx = i
				}
			}
		}
	}
	if biggestUnicode == -1 {
		for i, x := range subTables {
			println("table ", i)
			println("id", x[0])
			println("specID", x[1])
			println("offset", x[2])
		}
		panic("No unicode table found")
	}

	this.seek(entry.offset + subTables[unicodeIdx][2])
	var format = this.readUint16()
	if format == 12 {
		this.parseFormat12()

	} else if format == 4 {
		this.parseFormat4()
	} else {
		println("id", subTables[unicodeIdx][0])
		println("spec", subTables[unicodeIdx][1])
		println("offset", subTables[unicodeIdx][2])
		panic("Format not implemented")
	}

}

func (this *Reader) parseFormat12() {
	this.readUint16()
	this.readUint32()
	this.readUint32()
	var nrGroups = this.readUint32()
	var mappings map[int]int = map[int]int{}
	for i := 0; i < int(nrGroups); i++ {
		var startCode = this.readUint32()
		var endCode = this.readUint32()
		var startGlyphCode = this.readUint32()
		var count = -1
		for j := startCode; j <= endCode; j++ {
			count++
			mappings[int(j)] = int(startGlyphCode) + count
		}

	}
	this.mappings = mappings
}

func (this *Reader) parseFormat4() {
	this.readUint16()
	this.readUint16()
	var segCount = this.readUint16() / 2
	this.readUint16()
	this.readUint16()
	this.readUint16()

	var endCodes = make([]uint16, segCount)
	var startCode = make([]uint16, segCount)
	var idDelta = make([]uint16, segCount)
	var idRangeOffset = make([]uint16, segCount)

	for i := 0; i < int(segCount); i++ {
		endCodes[i] = this.readUint16()
	}
	this.readUint16()
	for i := 0; i < int(segCount); i++ {
		startCode[i] = this.readUint16()
	}
	for i := 0; i < int(segCount); i++ {
		idDelta[i] = this.readUint16()
	}
	for i := 0; i < int(segCount); i++ {
		idRangeOffset[i] = this.readUint16()
	}
	var mappings map[int]int = map[int]int{}
	for i := 0; i < int(segCount); i++ {
		if idRangeOffset[i] == 0 {
			for j := uint32(startCode[i]); j <= uint32(endCodes[i]); j++ {
				var test uint32 = 65536
				mappings[int(j)] = int(uint32((j + uint32(idDelta[i]))) % test)
			}
		} else if idRangeOffset[i] != 0 {
			if idDelta[i] != 0 {
				panic("PROBLEm")
			}
			for j := startCode[i]; j <= endCodes[i]; j++ {
				mappings[int(j)] = int(this.readUint16())
			}
		}
	}
	this.mappings = mappings
}

type CharCodeGroup struct {
	startCode      uint32
	endCode        uint32
	startGlyphCode uint32
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
	this.offset += 4
	return readType[uint32](this.file, 4, true)
}

func (this *Reader) readInt16() int16 {
	this.offset += 2
	return readType[int16](this.file, 2, true)
}

func (this *Reader) readUint8() uint8 {
	this.offset += 1
	return readType[uint8](this.file, 1, false)
}

func (this *Reader) readInt8() int8 {
	this.offset += 1
	return readType[int8](this.file, 1, false)
}

func (this *Reader) readStr(amountChars int) string {
	this.offset += uint32(amountChars)
	return string(readBytes(amountChars, this.file, false))
}

func (this *Reader) readUint16() uint16 {
	this.offset += 2
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
	println("tag", entry.tag)
	println("checksum", entry.checksum)
	println("offset", entry.offset)
	println("len", entry.length)
}

func (this *Reader) seek(offset uint32) {
	var newOffset, _ = this.file.Seek(int64(offset), 0)
	this.offset = uint32(newOffset)
}

func (this *Reader) calcChecksum(e *DirEntry) uint32 {
	this.seek(e.offset)
	var sum uint32 = 0
	var nrLongs = (e.length + 3) / 4
	for nrLongs > 0 {
		sum += this.readUint32()
		nrLongs--
	}
	return sum
}

func (this *Reader) readGlyfHeader() GlyfHeader {
	var nrContours = this.readInt16()
	if nrContours != 0 {
		return GlyfHeader{
			nrContours: nrContours,
			xMin:       float32(this.readInt16()),
			yMin:       float32(this.readInt16()),
			xMax:       float32(this.readInt16()),
			yMax:       float32(this.readInt16()),
		}
	} else {
		return GlyfHeader{
			nrContours: nrContours,
			xMin:       float32(this.readInt16()),
			yMin:       float32(this.readInt16()),
			xMax:       float32(this.readInt16()),
			yMax:       float32(this.readInt16()),
		}
	}
}

func (this *Reader) readSimpleGlyph(header GlyfHeader) SimpleGlyf {
	var body = SimpleGlyfBody{}

	body.endOfContours = make([]uint16, header.nrContours)
	for i := 0; i < len(body.endOfContours); i++ {
		body.endOfContours[i] = this.readUint16()
	}
	var nrPoints = 0
	for _, x := range body.endOfContours {
		nrPoints = int(math.Max(float64(x+1), float64(nrPoints)))
	}
	body.instructionLength = this.readUint16()
	body.instructions = make([]uint8, body.instructionLength)
	for i := 0; i < len(body.instructions); i++ {
		body.instructions[i] = this.readUint8()
	}
	body.flags = make([]uint8, nrPoints)
	for i := 0; i < int(nrPoints); i++ {
		body.flags[i] = this.readUint8()
		//repeat
		if getBit(body.flags[i], 3) == 1 {
			var repetitions = this.readUint8()
			var base = body.flags[i]
			for j := 0; j < int(repetitions); j++ {
				i++
				body.flags[i] = base
			}
		}
	}
	var xCoordinates = make([]int16, nrPoints)
	var yCoordinates = make([]int16, nrPoints)
	this.parsePoints(&body.flags, &xCoordinates, true)
	this.parsePoints(&body.flags, &yCoordinates, false)
	body.Points = [][]glm.Vec2{}
	var glyfPoints = [][]GlyfPoints{}

	var currContour = []GlyfPoints{}
	for i, x := range xCoordinates {
		var endP = closedGL.Contains(&body.endOfContours, uint16(i))
		var OnCurve = getBit(body.flags[i], 0) == 1
		var cartPos = glm.Vec2{float32(x), float32(yCoordinates[i])}
		var newP = GlyfPoints{
			OnCurve:  OnCurve,
			EndPoint: endP,
			Pos:      cartPos,
		}
		currContour = append(currContour, newP)
		if endP {
			glyfPoints = append(glyfPoints, currContour)
			currContour = []GlyfPoints{}
		}
	}

	body.Points = this.transformPoints2(glyfPoints, header.nrContours)

	return SimpleGlyf{
		body: body,
	}
}

func printBezierPoints(points []glm.Vec2) {
	for i, x := range points {
		if (i)%3 == 0 {
			println("--------------")
		}
		closedGL.PrintlnVec2(x)
	}
}

func (this *Reader) transformPoints2(points [][]GlyfPoints, nrContours int16) [][]glm.Vec2 {

	//add implicit points
	var newPoints = [][]glm.Vec2{}
	for i := 0; i < len(points); i++ {
		var currPoints = []glm.Vec2{}
		var currContour = points[i]
		for j := 0; j < len(currContour)-1; j++ {
			var first = currContour[j]
			var second = currContour[j+1]
			currPoints = append(currPoints, first.Pos)
			if first.OnCurve == second.OnCurve {
				var control = this.createOnCurveMiddlePoint(first, second).Pos
				currPoints = append(currPoints, control)
			}
		}
		var last = currContour[len(currContour)-1]
		if last.OnCurve {
			currPoints = append(currPoints, last.Pos)
		} else {
			currPoints = append(currPoints, last.Pos)
			currPoints = append(currPoints, this.createOnCurveMiddlePoint(last, currContour[0]).Pos)
		}
		newPoints = append(newPoints, currPoints)
	}
	//make splines
	var bezierSpline = [][]glm.Vec2{}
	for _, x := range newPoints {
		var newSpline = []glm.Vec2{}
		for i := 1; i < len(x)-1; i += 2 {
			var first = x[i-1]
			var control = x[i]
			var second = x[i+1]
			newSpline = append(newSpline, first, control, second)
		}
		newSpline = append(newSpline, x[len(x)-1], x[len(x)-2], x[0])
		/* println("splite")
		printBezierPoints(newSpline) */
		bezierSpline = append(bezierSpline, newSpline)
	}
	return bezierSpline
}

func (this *Reader) createOnCurveMiddlePoint(p1, p2 GlyfPoints) GlyfPoints {
	return GlyfPoints{
		Pos:      closedGL.LerpVec2(p1.Pos, p2.Pos, 0.5),
		OnCurve:  false,
		EndPoint: false,
	}
}

func (this *Reader) transformPoints(points []GlyfPoints) []GlyfPoints {
	var start = true
	var startP GlyfPoints
	var newPoints = []GlyfPoints{}
	for _, x := range points {
		if start {
			startP = x
			start = false
		}
		newPoints = append(newPoints, x)
		if x.EndPoint {
			newPoints = append(newPoints, startP)
			start = true
		}
	}

	points = newPoints
	newPoints = []GlyfPoints{}
	for i := 0; i < len(points); i++ {
		newPoints = append(newPoints, points[i])
		if i < len(points)-1 && !points[i].OnCurve && !points[i+1].OnCurve {
			var newPos = closedGL.LerpVec2(points[i].Pos, points[i+1].Pos, 0.5)
			var newP = GlyfPoints{
				OnCurve:  true,
				EndPoint: points[i].EndPoint || points[i+1].EndPoint,
				Pos:      newPos,
			}
			newPoints = append(newPoints, newP)
		}
	}
	points = newPoints
	newPoints = []GlyfPoints{}
	for i := 0; i < len(points); i++ {
		newPoints = append(newPoints, points[i])
		if i < len(points)-1 && points[i].OnCurve && points[i+1].OnCurve {
			var newPos = closedGL.LerpVec2(points[i].Pos, points[i+1].Pos, 0.5)
			var newP = GlyfPoints{
				OnCurve:  false,
				EndPoint: points[i].EndPoint || points[i+1].EndPoint,
				Pos:      newPos,
			}
			newPoints = append(newPoints, newP)
		}
	}
	newPoints = append(newPoints, newPoints[len(newPoints)-1])

	return newPoints
}

func (this *Reader) parsePoints(flags *[]uint8, destArr *[]int16, isX bool) {
	var lastVal int16 = 0
	var offset uint8 = 0
	if !isX {
		offset++
	}
	for i, x := range *flags {
		var short = getBit(x, 1+offset)
		var same = getBit(x, 4+offset)
		if same == 0 && short == 0 {
			lastVal += this.readInt16()
		} else if same == 0 && short == 1 {
			lastVal -= int16(this.readUint8())
		} else if same == 1 && short == 0 {
			lastVal = lastVal
		} else if same == 1 && short == 1 {
			lastVal += int16(this.readUint8())
		}
		(*destArr)[i] = lastVal
	}
}

func (this *Reader) printGlfyHeader(h GlyfHeader) {
	println(h.nrContours)
	println("xMin", h.xMin)
	println("yMin", h.yMin)
	println("xMax", h.xMax)
	println("yMax", h.yMax)
}

func getBit[T uint32 | uint16 | uint8](val T, bit T) T {
	return (val >> bit) & 1
}

/* func (this *Reader) printGlyfBody(b SimpleGlyfBody) {
	println("body")
	println("instLen", b.instructionLength)
	println(len(b.instructions))
	println("points")
	for _, x := range b.Points {
		closedGL.PrintlnVec2(x.Pos)
	}
} */

func (this *Reader) readGlyf(unicodeVal uint32) Glyf {
	var entry = this.entries["glyf"]
	var glyphId = this.mappings[int(unicodeVal)]
	this.seek(entry.offset + this.loca[glyphId])
	var ret = newGlyf()
	ret.AdvanceWidth = float32(this.horizMetrics[glyphId].advanceWidth)
	ret.header = this.readGlyfHeader()
	var isZeroLenGlyf = this.loca[glyphId-1] == this.loca[glyphId]
	if isZeroLenGlyf {
		return ret
	}
	if ret.header.nrContours > 0 {
		var g = this.readSimpleGlyph(ret.header)
		ret.SimpleGlyfs = append(ret.SimpleGlyfs, &g)
	} else if ret.header.nrContours < 0 {
		var comp = this.readCompundGlyf(ret.header)
		for _, x := range comp.compundDescr {
			var g = this.readGlyfIdx(uint32(x.glyfIdx))
			if x.flags&0x02 == 0x02 {
				g.AddOffset(glm.Vec2{float32(x.arg1), float32(x.arg2)})
			} else {
				panic("Kind of flag not yet supported")
			}
			ret.SimpleGlyfs = append(ret.SimpleGlyfs, &g)
		}
	} else {
		println("nr contours 0")
	}
	return ret
}

func (this *Reader) readGlyfIdx(idx uint32) SimpleGlyf {
	var entry = this.entries["glyf"]
	this.seek(entry.offset + this.loca[int(idx)])
	var ret SimpleGlyf
	var header = this.readGlyfHeader()
	if header.nrContours > 0 {
		ret = this.readSimpleGlyph(header)
	} else if header.nrContours < 0 {
		//var comp = this.readCompundGlyf(header)
		panic("nested comp glyf not yet supported")
	} else {
		panic("nr contours 0")
	}
	return ret
}

func (this *SimpleGlyf) GetBody() *SimpleGlyfBody {
	return &this.body
}

func (this *Reader) readCompundGlyf(header GlyfHeader) CompoundGlyf {
	var bodies = []CompoundBody{}
	bodies = append(bodies, this.readCompundGlyfBody())
	for getBit(bodies[len(bodies)-1].flags, 5) == 1 {
		bodies = append(bodies, this.readCompundGlyfBody())
	}
	return CompoundGlyf{
		header:       header,
		compundDescr: bodies,
	}
}

func (this *Reader) readCompundGlyfBody() CompoundBody {
	var flags = this.readUint16()
	var idx = this.readUint16()
	var arg1, arg2 int32
	if getBit(flags, 1) == 1 {
		if getBit(flags, 0) == 1 {
			arg1 = int32(this.readInt16())
			arg2 = int32(this.readInt16())
		} else {
			arg1 = int32(this.readInt8())
			arg2 = int32(this.readInt8())
		}
	} else {
		panic("Points!!")
		if getBit(flags, 0) == 1 {
			arg1 = int32(this.readUint16())
			arg2 = int32(this.readUint16())
		} else {
			arg1 = int32(this.readUint8())
			arg2 = int32(this.readUint8())
		}
	}
	var weHaveAScale = getBit(flags, 3) == 1
	var weHaveAXAndYScale = getBit(flags, 6) == 1
	var weHaveATwoByTwo = getBit(flags, 7) == 1
	var a, b, c, d uint16 = 1, 0, 0, 1

	if weHaveAScale {
		a = this.readUint16()
		d = a
	} else if weHaveAXAndYScale {
		a = this.readUint16()
		d = this.readUint16()
	} else if weHaveATwoByTwo {
		a = this.readUint16()
		b = this.readUint16()
		c = this.readUint16()
		d = this.readUint16()
	}
	return CompoundBody{
		flags:   flags,
		glyfIdx: idx,
		arg1:    arg1,
		arg2:    arg2,
		a:       a,
		b:       b,
		c:       c,
		d:       d,
	}
}

func (this *Reader) parseHead() {
	var entry = this.entries["head"]
	this.seek(entry.offset)
	this.readUint32()
	this.readUint32()
	this.readUint32()
	this.readUint32()
	this.readUint16()
	this.readUint16()
	//date
	this.readUint32()
	this.readUint32()
	this.readUint32()
	this.readUint32()
	//xMin
	this.readInt16()
	this.readInt16()
	this.readInt16()
	this.readInt16()
	this.readUint16()
	this.readUint16()
	this.readInt16()
	this.smallLocATable = this.readInt16() == 0

}

func (this *Reader) parseMaxP() {
	var entry = this.entries["maxp"]
	if this.calcChecksum(entry) != entry.checksum {
		panic("checksum")
	}
	this.seek(entry.offset)
	this.readUint32()
	this.nrGlyphs = this.readUint16()

}

func (this *Reader) parseLocA() {
	var entry = this.entries["loca"]
	if this.calcChecksum(entry) != entry.checksum {
		panic("checksum")
	}
	this.seek(entry.offset)
	this.loca = make([]uint32, this.nrGlyphs+1)
	for i := 0; i < len(this.loca); i++ {
		if this.smallLocATable {
			this.loca[i] = uint32(this.readUint16() * 2)
		} else {
			this.loca[i] = uint32(this.readUint32())
		}
	}
}

func (this *Reader) parseHhea() {
	var entry = this.entries["hhea"]
	if this.calcChecksum(entry) != entry.checksum {
		panic("wrong checksum in hhea")
	}
	this.seek(entry.offset)
	this.readUint16()
	this.readUint16()

	this.readInt16()
	this.readInt16()
	this.readInt16()

	this.readUint16()

	this.readInt16()
	this.readInt16()
	this.readInt16()

	this.readInt16()
	this.readInt16()
	this.readInt16()

	this.readInt16()
	this.readInt16()
	this.readInt16()
	this.readInt16()

	this.readInt16()
	this.nrHMetrics = this.readUint16()
}

func (this *Reader) parseHmtx() {
	var entry = this.entries["hmtx"]
	if this.calcChecksum(entry) != entry.checksum {
		panic("wrong checksum")
	}
	this.seek(entry.offset)
	this.horizMetrics = make([]LongHorizMetric, this.nrHMetrics)
	this.remainingBearings = make([]int16, this.nrGlyphs-this.nrHMetrics)

	for i := 0; i < int(this.nrHMetrics); i++ {
		this.horizMetrics[i] = LongHorizMetric{
			advanceWidth: this.readUint16(),
			lsb:          this.readInt16(),
		}
	}
	for i := 0; i < len(this.remainingBearings); i++ {
		this.remainingBearings[i] = this.readInt16()
	}

}
