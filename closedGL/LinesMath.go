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
			return glm.Vec2{this.p1[0], line.p1[1]}, true
		} else if line.LineType == "normal" {
			return line.EvalY(this.p1[1])
		} else {
			panic("unreachable")
		}
	} else if this.LineType == "vertical" {
		if line.LineType == "horizontal" {
			return glm.Vec2{line.p1[0], this.p1[1]}, true
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
