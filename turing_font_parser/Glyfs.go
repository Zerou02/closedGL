package turingfontparser

import (
	"github.com/EngoEngine/glm"
	"github.com/EngoEngine/math"
)

type Glyf struct {
	header       GlyfHeader
	SimpleGlyfs  []SimpleGlyf
	AdvanceWidth float32
}

func (this *Glyf) GetHeader() *GlyfHeader {
	return &this.header
}
func newGlyf() Glyf {
	return Glyf{SimpleGlyfs: []SimpleGlyf{}}
}

func (this *GlyfHeader) GetMaxX() float32 {
	return float32(this.xMax)
}

type SimpleGlyf struct {
	body SimpleGlyfBody
}
type GlyfHeader struct {
	nrContours int16
	xMin       float32
	yMin       float32
	xMax       float32
	yMax       float32
}

type CompoundBody struct {
	flags      uint16
	glyfIdx    uint16
	arg1, arg2 int32
	//todo: change to f16
	a, b, c, d uint16
}
type CompoundGlyf struct {
	header       GlyfHeader
	compundDescr []CompoundBody
	points       []GlyfPoints
}

type SimpleGlyfBody struct {
	endOfContours     []uint16
	instructionLength uint16
	instructions      []uint8
	flags             []uint8
	Points            []GlyfPoints
}
type GlyfPoints struct {
	Pos      glm.Vec2
	OnCurve  bool
	EndPoint bool
}

func (this SimpleGlyf) GetPoints() []GlyfPoints {
	return this.body.Points
}

func (this CompoundGlyf) AddOffset(y glm.Vec2) {
	for i, x := range this.points {
		this.points[i].Pos = x.Pos.Add(&y)
	}
}

func (this SimpleGlyf) AddOffset(y glm.Vec2) {
	for i, x := range this.body.Points {
		this.body.Points[i].Pos = x.Pos.Add(&y)
	}
}

func (this CompoundGlyf) GetPoints() []GlyfPoints {
	return this.points
}

func (this *Glyf) Normalize() {
	var xMin = math.Abs(this.header.xMin)
	var yMin = math.Abs(this.header.yMin)
	this.header.xMin += xMin
	this.header.yMin += yMin
	this.header.xMax += xMin
	this.header.yMax += yMin
	for i := 0; i < len(this.SimpleGlyfs); i++ {
		for j := 0; j < len(this.SimpleGlyfs[i].body.Points); j++ {
			this.SimpleGlyfs[i].body.Points[j].Pos[0] += xMin
			this.SimpleGlyfs[i].body.Points[j].Pos[1] += yMin
		}
	}
}

func (this *Glyf) CalcScaleFactor(newHeight float32) float32 {
	return newHeight / this.header.yMax
}

func (this *Glyf) ScaleToHeight(newHeight float32) {
	this.Scale(this.CalcScaleFactor(newHeight))
}

func (this *Glyf) Scale(scale float32) {
	this.header.yMax *= scale
	this.header.xMax *= scale
	this.AdvanceWidth *= scale
	for i := 0; i < len(this.SimpleGlyfs); i++ {
		for j := 0; j < len(this.SimpleGlyfs[i].body.Points); j++ {
			this.SimpleGlyfs[i].body.Points[j].Pos[0] *= scale
			this.SimpleGlyfs[i].body.Points[j].Pos[1] *= scale
		}
	}
}
