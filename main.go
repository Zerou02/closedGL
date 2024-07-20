package main

import (
	_ "image/png"

	"github.com/EngoEngine/glm"
	"github.com/Zerou02/closedGL/closedGL"
	"github.com/go-gl/glfw/v3.2/glfw"
)

func main() {
	StartClosedGL()
}

func StartClosedGL() {

	var openGL = closedGL.InitClosedGL(800, 600, "demo")
	openGL.Window.SetScrollCallback(openGL.Camera.ScrollCb)
	openGL.Window.SetCursorPosCallback(openGL.Camera.MouseCallback)
	openGL.Window.Window.SetMouseButtonCallback(closedGL.StandardMouseClickCB)
	var isWireframeMode = false

	_ = isWireframeMode
	openGL.Camera.CameraPos = glm.Vec3{0, 0, 0}

	var elapsed float64 = 0
	var last = glfw.GetTime()
	var frameCount = 0

	openGL.PlayMusic("bgm")
	var anim = closedGL.NewStaggeredAnimation([]closedGL.Animation{
		closedGL.NewAnimation(100, 500, 1, false, false),
		closedGL.NewAnimation(100, 500, 1, false, false),
	})

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
		anim.Process(float32(delta))

		if openGL.KeyBoardManager.IsPressed(glfw.KeyF) {
			isWireframeMode = !isWireframeMode
			closedGL.SetWireFrameMode(isWireframeMode)
		}

		openGL.BeginDrawing()
		openGL.ClearBG(glm.Vec4{0, 0, 0, 0})
		var valueArr = anim.GetValueArr()
		openGL.DrawRect(glm.Vec4{0, 0, valueArr[0], valueArr[1]}, glm.Vec4{1, 1, 1, 1}, 1)
		openGL.EndDrawing()
		openGL.Process()
	}
	openGL.Free()
}
