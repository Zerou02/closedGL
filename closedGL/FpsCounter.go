package closedGL

import "github.com/go-gl/glfw/v3.2/glfw"

type FPSCounter struct {
	Elapsed, fpsSum                   float64
	FrameCount, fpsAmount, FpsAverage int
	Delta                             float64 //in sec
	lastFrame                         float64
	Time                              float64
}

func NewFPSCounter() FPSCounter {
	var counter = FPSCounter{Elapsed: 0, fpsSum: 0, fpsAmount: 0, FpsAverage: 0, Delta: 0, FrameCount: 0, lastFrame: 0, Time: 0}
	return counter
}

func (this *FPSCounter) Process() {
	var currFrame = glfw.GetTime()
	this.Elapsed += this.Delta
	this.Delta = currFrame - this.lastFrame
	this.Time += this.Delta
	this.lastFrame = currFrame
	this.FrameCount += 1
	this.fpsAmount += 1
	this.fpsSum += 1 / this.Delta
}

func (this *FPSCounter) CalcAverage() {
	this.FpsAverage = int(this.fpsSum / float64(this.fpsAmount))
}

func (this *FPSCounter) Clear() {
	this.Elapsed = 0
	this.fpsSum = 0
	this.fpsAmount = 0
}
