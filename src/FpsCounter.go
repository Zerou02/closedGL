package main

import "github.com/go-gl/glfw/v3.2/glfw"

type FPSCounter struct {
	elapsed, fpsSum                   float64
	frameCount, fpsAmount, fpsAverage int
	delta                             float64
	lastFrame                         float64
}

func newFPSCounter() FPSCounter {
	var counter = FPSCounter{elapsed: 0, fpsSum: 0, fpsAmount: 0, fpsAverage: 0, delta: 0, frameCount: 0, lastFrame: 0}
	return counter
}

func (this *FPSCounter) process() {
	var currFrame = glfw.GetTime()
	this.elapsed += this.delta
	this.delta = currFrame - this.lastFrame
	this.lastFrame = currFrame
	this.frameCount += 1
	this.fpsAmount += 1
	this.fpsSum += 1 / this.delta
}

func (this *FPSCounter) calcAverage() {
	this.fpsAverage = int(this.fpsSum / float64(this.fpsAmount))
}

func (this *FPSCounter) clear() {
	this.elapsed = 0
	this.fpsSum = 0
	this.fpsAmount = 0
}
