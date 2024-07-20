package main

import (
	_ "image/png"
	"os"
	"strings"

	"github.com/EngoEngine/glm"
	"github.com/Zerou02/closedGL/closedGL"
	"github.com/go-gl/glfw/v3.2/glfw"
)

func main() {
	StartClosedGL()
}

func StartClosedGL() {

	var openGL = closedGL.InitClosedGL(800, 600)
	openGL.Window.SetScrollCallback(openGL.Camera.ScrollCb)
	openGL.Window.SetCursorPosCallback(openGL.Camera.MouseCallback)
	openGL.Window.Window.SetMouseButtonCallback(closedGL.StandardMouseClickCB)
	var isWireframeMode = false

	_ = isWireframeMode
	glfw.SwapInterval(0)
	openGL.Camera.CameraPos = glm.Vec3{0, 0, 0}

	var elapsed float64 = 0
	var last = glfw.GetTime()
	var frameCount = 0

	var bytees, _ = os.ReadFile("./assets/intro.txt")
	var contents = string(bytees)
	var lines []string = []string{}
	var amountChars = 0
	var secondPerChars float32 = 0.03
	for _, x := range strings.Split(contents, "\n") {
		amountChars += len(x)
		lines = append(lines, x)
	}
	var time = float32(amountChars) * secondPerChars
	var anim = closedGL.NewAnimation(0, float32(amountChars), time, false, false)

	var samples = []string{}
	for i := 0; i < 60; i++ {
		var sample = ""
		for j := 0; j < i+1; j++ {
			sample += "a"
		}
		samples = append(samples, sample)
	}

	for !openGL.Window.Window.ShouldClose() {
		var curr = glfw.GetTime()
		elapsed += curr - last
		last = curr
		frameCount++
		if elapsed > 1 {
			closedGL.PrintlnFloat(float32(openGL.FPSCounter.FpsAverage))

			elapsed = 0
			frameCount = 0
		}
		var delta = openGL.FPSCounter.Delta
		_ = delta
		if openGL.KeyBoardManager.IsPressed(glfw.KeyF) {
			isWireframeMode = !isWireframeMode
			closedGL.SetWireFrameMode(isWireframeMode)
		}

		openGL.BeginDrawing()
		openGL.ClearBG(glm.Vec4{0, 0, 0, 0})

		anim.Process(float32(delta))
		openGL.LimitFPS(true)

		if !anim.IsFinished() {
			var currTextLen = anim.GetValue()
			println("len", int(currTextLen))
			_ = currTextLen

			var lineToDraw = ""
			for i := -1; i < int(currTextLen); i++ {
				lineToDraw += "a"
			}
			/* 			for i := 0; i < len(lines); i++ {
			if alreadyDrawn >= int(currTextLen) {
				break
			}
			var line = lines[i]
			var lineToDraw = ""
			if len(line)+alreadyDrawn < int(currTextLen) {
				alreadyDrawn += len(line)
				lineToDraw = line
			} else {
				var copy = alreadyDrawn
				for j := 0; j < int(currTextLen)-copy; j++ {
					alreadyDrawn++
					lineToDraw += string(line[j])
				}
			} */
			//openGL.DrawRect(glm.Vec4{100 + currTextLen, 100, 200, 200}, glm.Vec4{1, 1, 1, 1}, int(currTextLen))
			openGL.Text.DrawText(0, 100+50, samples[openGL.FPSCounter.FrameCount%len(samples)], 1)
			openGL.DrawRect(glm.Vec4{100 + currTextLen, 100, 100, 100}, glm.Vec4{0, currTextLen / 128, 1, 1}, int(currTextLen))
			openGL.DrawRect(glm.Vec4{100 + currTextLen, 150, 100, 100}, glm.Vec4{0, currTextLen / 128, 1, 1}, int(currTextLen))
			openGL.DrawRect(glm.Vec4{100 + currTextLen, 250, 100, 100}, glm.Vec4{0, currTextLen / 128, 1, 1}, int(currTextLen))

			//	openGL.DrawFPS(0, 0, 1)

		}
		openGL.EndDrawing()
		openGL.Process()
	}
	openGL.Free()
}
