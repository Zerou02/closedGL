package main

import (
	_ "image/png"
	"runtime"
	"time"

	"github.com/EngoEngine/glm"
	"github.com/Zerou02/closedGL/closedGL"
	"github.com/go-gl/glfw/v3.2/glfw"
)

func main() {
	StartClosedGL()
}

func StartClosedGL() {
	runtime.LockOSThread()

	var openGL = closedGL.InitClosedGL(800, 600)
	openGL.Window.SetScrollCallback(openGL.Camera.ScrollCb)
	openGL.Window.SetCursorPosCallback(openGL.Camera.MouseCallback)

	var isWireframeMode = false

	_ = isWireframeMode
	glfw.SwapInterval(0)
	openGL.Camera.CameraPos = glm.Vec3{0, 0, 0}

	var elapsed float64 = 0
	var last = glfw.GetTime()
	var frameCount = 0

	var movAnim = closedGL.NewAnimation(75, 365, 1, true, true)
	var movAnim2 = closedGL.NewAnimation(-200, 0, 1, true, false)
	var colourAnim = closedGL.NewAnimation(0, 1, 1, true, true)
	var colourAnim2 = closedGL.NewAnimation(1, 0, 1, true, true)
	var shiningAnim = closedGL.NewAnimation(0, 100, 1, true, true)

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
		openGL.ClearBG()

		movAnim.Process(float32(delta))
		movAnim2.Process(float32(delta))
		colourAnim.Process(float32(delta))
		colourAnim2.Process(float32(delta))
		shiningAnim.Process(float32(delta))

		var start = time.Now()
		openGL.DrawRect(glm.Vec4{145, 30, 10, 10}, glm.Vec4{0, 1, 1, 1}, 0)
		openGL.DrawRect(glm.Vec4{145 + 20, 30, 10, 40}, glm.Vec4{1, 0, 1, 1}, 0)
		openGL.DrawRect(glm.Vec4{145 + 40, 30, 10, 40}, glm.Vec4{1, 1, 0, 1}, 0)
		openGL.DrawRect(glm.Vec4{145 + 60, 30, 10, 40}, glm.Vec4{1, 0, 1, 1}, 0)
		openGL.DrawRect(glm.Vec4{145 + 80, 30, 10, 40}, glm.Vec4{1, 1, 0, 1}, 0)
		openGL.DrawLine(glm.Vec2{0, 0}, glm.Vec2{100, 100}, glm.Vec4{1, 0, 0, 1}, glm.Vec4{0, 1, 1, 1}, 0)

		openGL.DrawTriangle([3]glm.Vec2{{100, 100}, {0, 350}, {200, 350}}, glm.Vec4{1, 1, 0, 1}, 1)
		openGL.DrawCircle(glm.Vec2{150, 150}, glm.Vec4{1, 0, 0, 1}, glm.Vec4{1, 1, 0, 1}, 50, shiningAnim.GetValue(), 0)
		openGL.DrawCircle(glm.Vec2{150, 150}, glm.Vec4{1, 0, 1, 1}, glm.Vec4{1, 1, 0, 1}, 10, shiningAnim.GetValue(), 0)

		openGL.DrawRect(glm.Vec4{0, 50, 100, 100}, glm.Vec4{0, 1, 1, 1}, 2)
		openGL.DrawRect(glm.Vec4{0, 50, 50, 50}, glm.Vec4{1, 0, 0, 1}, 1)

		openGL.DrawFPS(500, 0)
		var end = time.Now()
		_, _ = end, start

		var start2 = time.Now()
		openGL.EndDrawing()
		var end2 = time.Now()
		//		fmt.Printf("draw:%f\n", end2.Sub(start2).Seconds())

		_, _ = end2, start2

		openGL.Process()
	}
	openGL.Free()
}
