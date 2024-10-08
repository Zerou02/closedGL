package closedGL

import (
	"github.com/EngoEngine/glm"
	"github.com/EngoEngine/math"
)

type Line struct {
	p1, p2, lineData glm.Vec2 //m,n
	LineType         string
}

func (this *Line) IsParallel(line Line) bool {
	return this.lineData[0] == line.lineData[0]
}

func (this *Line) GetIntersection(line Line) (glm.Vec2, bool) {
	if this.IsParallel(line) {
		return glm.Vec2{}, false
	}
	if this.LineType == "horizontal" {
		if line.LineType == "vertical" {
			return glm.Vec2{line.p1[0], this.p1[1]}, true
		} else if line.LineType == "normal" {
			return line.EvalY(this.p1[1])
		} else {
			panic("unreachable")
		}
	} else if this.LineType == "vertical" {
		if line.LineType == "horizontal" {
			return glm.Vec2{this.p1[0], line.p1[1]}, true
		} else if line.LineType == "normal" {
			return line.EvalX(this.p1[0])
		} else {
			panic("unreachable")
		}
	} else if this.LineType == "normal" {
		if line.LineType == "horizontal" {
			return this.EvalY(line.p1[1])
		} else if line.LineType == "vertical" {
			return this.EvalX(line.p1[0])
		} else {
			var x = (this.lineData[1] - line.lineData[1]) / (line.lineData[0] - this.lineData[0])
			return glm.Vec2{x, this.lineData[0]*x + this.lineData[1]}, true
		}
	}
	return glm.Vec2{}, false
}

// funktioniert wahrscheinlich nur für SS
func (this *Line) IsOnLine(p glm.Vec2) bool {
	var eps float32 = 0.005
	if this.LineType == "vertical" {
		var minY = math.Min(this.p1[1], this.p2[1])
		var maxY = math.Max(this.p1[1], this.p2[1])
		return this.p1[0] == p[0] && minY <= p[1] && p[1] <= maxY
	} else if this.LineType == "horizontal" {
		var minX = math.Min(this.p1[0], this.p2[0])
		var maxX = math.Max(this.p1[0], this.p2[0])
		return this.p1[1] == p[1] && minX <= p[0] && p[0] <= maxX
	} else {
		var y, _ = this.EvalX(p[0])
		var minX = math.Min(this.p1[0], this.p2[0])
		var maxX = math.Max(this.p1[0], this.p2[0])
		var minY = math.Min(this.p1[1], this.p2[1])
		var maxY = math.Max(this.p1[1], this.p2[1])
		return math.Abs(y[1]-p[1]) <= eps && IsPointInRect(p, glm.Vec4{minX, minY, maxX - minX, maxY - minY})
	}
}

func (this *Line) EvalX(newX float32) (glm.Vec2, bool) {
	if this.LineType == "vertical" {
		return glm.Vec2{}, false
	} else if this.LineType == "horizontal" {
		return glm.Vec2{newX, this.p1[1]}, true
	} else {
		return glm.Vec2{newX, this.lineData[0]*newX + this.lineData[1]}, true
	}
}

func (this *Line) EvalY(newY float32) (glm.Vec2, bool) {
	if this.LineType == "vertical" {
		return glm.Vec2{}, false
	} else if this.LineType == "horizontal" {
		return glm.Vec2{}, false
	} else {
		return glm.Vec2{(newY - this.lineData[1]) / this.lineData[0], newY}, true
	}
}

// start,end exklusiv
func (this *Line) SamplePointsOnLine(amount float32) []glm.Vec2 {
	var ret = []glm.Vec2{}
	if this.LineType == "vertical" {
		var dy = this.p2[1] - this.p1[1]
		dy /= (amount + 1)
		for i := 1; i < int(amount); i++ {
			ret = append(ret, glm.Vec2{this.p1[0], this.p1[1] + dy*float32(i)})
		}
	} else if this.LineType == "normal" {
		var dx = this.p2[0] - this.p1[0]
		for i := 1; i < int(amount); i++ {
			var newP, _ = this.EvalX(this.p1[0] + dx*(float32(i)/amount))
			ret = append(ret, newP)
		}
	} else if this.LineType == "horizontal" {
		var dx = this.p2[0] - this.p1[0]
		for i := 1; i < int(amount); i++ {
			ret = append(ret, glm.Vec2{this.p1[0] + dx*(float32(i)/amount), this.p1[1]})
		}
	}
	return ret
}

func (this *Line) LerpPointOnLine(t float32) glm.Vec2 {
	return LerpVec2(this.p1, this.p2, t)
}

// m,n
func CalculateLine(p1, p2 glm.Vec2) Line {
	if p1[1] == p2[1] {
		return Line{
			p1: p1, p2: p2,
			lineData: glm.Vec2{0, p1[1]},
			LineType: "horizontal",
		}
	}
	if p1[0] == p2[0] {
		return Line{
			p1: p1, p2: p2,
			lineData: glm.Vec2{math.Inf(1), p1[1]},
			LineType: "vertical",
		}
	}
	var dy = p2[1] - p1[1]
	var dx = p2[0] - p1[0]
	var m = dy / dx
	var n = p1[1] - m*p1[0]
	return Line{
		p1: p1, p2: p2,
		lineData: glm.Vec2{m, n},
		LineType: "normal",
	}
}

func (this *Line) ClosestPoint(p glm.Vec2) glm.Vec2 {
	return ClosestPointOnLine(this.p1, this.p2, p)
}

func ClosestPointOnLine(lp1, lp2, p glm.Vec2) glm.Vec2 {
	var y2 = lp2[1]
	var y1 = lp1[1]
	var x1 = lp1[0]
	var x2 = lp2[0]
	var a1 = y2 - y1
	var b1 = x1 - x2
	var c1 = (y2-y1)*x1 + (x1-x2)*y1
	var c2 = -b1*p[0] + a1*p[1]
	var det = a1*a1 - -b1*b1
	var cx float32 = 0
	var cy float32 = 0
	if det != 0 {
		cx = (a1*c1 - b1*c2) / det
		cy = (a1*c2 - -b1*c1) / det
	} else {
		cx = p[0]
		cy = p[1]
	}
	return glm.Vec2{cx, cy}
}
