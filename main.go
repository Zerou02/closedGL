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

	openGL.LimitFPS(false)
	var timer = closedGL.NewTimer(1, false)
	var anim = closedGL.NewAnimation(100, 500, 3, false, true)
	for !openGL.Window.Window.ShouldClose() {

		var delta = openGL.FPSCounter.Delta
		_ = delta
		anim.Process(float32(delta))
		timer.Process(float32(delta))

		if openGL.KeyBoardManager.IsPressed(glfw.KeyF) {
			isWireframeMode = !isWireframeMode
			closedGL.SetWireFrameMode(isWireframeMode)
		}
		if timer.IsTick() {
			println(openGL.FPSCounter.FpsAverage)
		}

		openGL.BeginDrawing()
		openGL.DrawFPS(500, 0, 1)
		openGL.DrawSprite(glm.Vec4{0, 0, 200, 200}, "./assets/sprites/fence.png", 1)
		openGL.ClearBG(glm.Vec4{0, 0, 0, 0})

		openGL.EndDrawing()
		openGL.Process()
	}
	openGL.Free()
}
