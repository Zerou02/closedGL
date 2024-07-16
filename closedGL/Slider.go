package closedGL

/*

import (
	"math"
	"strconv"

	"github.com/EngoEngine/glm"
	"github.com/go-gl/glfw/v3.2/glfw"
)

type Slider struct {
	baseLine, marker Rectangle
	min, max, curr   float32
	step             float64
	label            string
	window           *glfw.Window
}

func newSlider(window *glfw.Window, lineDim, lineColour, markerDim, markerColour glm.Vec4, min, max, curr float32, step float64, label string) Slider {
	var slider = Slider{
		baseLine: factory.NewRect(lineDim, lineColour),
		marker:   factory.NewRect(markerDim, markerColour),
		min:      min,
		max:      max,
		curr:     curr,
		step:     step,
		window:   window,
		label:    label,
	}
	slider.alignMarker()
	return slider
}

func (this *Slider) alignMarker() {
	var perc = this.curr / this.max
	var posX = this.baseLine.dim[2]*perc + this.baseLine.dim[0]
	this.marker.dim[0] = posX - this.marker.dim[2]/2
}

func (this *Slider) draw() {
	this.baseLine.Draw()
	this.marker.Draw()
	var x = this.baseLine.dim[0] + this.baseLine.dim[2]/2 - 50
	var y = this.baseLine.dim[1] - 20
	text.createVertices(this.label+strconv.FormatFloat(float64(this.curr), 'f', -1, 32), x, y)

	x = this.baseLine.dim[0] - 50
	y = this.baseLine.dim[1]
	text.createVertices(strconv.FormatFloat(float64(this.min), 'f', -1, 32), x, y)

	x = this.baseLine.dim[0] + this.baseLine.dim[2] + 10
	y = this.baseLine.dim[1]
	text.createVertices(strconv.FormatFloat(float64(this.max), 'f', -1, 32), x, y)
}

func (this *Slider) process() {
	var mouseX, mouseY = this.window.GetCursorPos()
	var mousePosVec = glm.Vec2{float32(mouseX), float32(mouseY)}
	if this.window.GetMouseButton(glfw.MouseButton1) == glfw.Press {
		if isPointInRect(mousePosVec, this.baseLine.dim) {
			var mouseXNormal = mousePosVec[0] - this.baseLine.dim[0]
			var percentage = mouseXNormal / this.baseLine.dim[2]
			this.curr = lerp(float32(this.min), float32(this.max), percentage)
			var decimalPointOffset = float64(neededDecimalPlacesToNextInt(this.step))
			var adjustedDivisor = this.step * float64(decimalPointOffset) * 10
			var adjustedCurr = float64(this.curr) * float64(decimalPointOffset) * 10
			var ratio = math.Round(adjustedCurr / adjustedDivisor)
			var roundedCurr = ratio * adjustedDivisor
			this.curr = float32(roundedCurr / (decimalPointOffset * 10))

			this.alignMarker()
		}
	}

}
*/
