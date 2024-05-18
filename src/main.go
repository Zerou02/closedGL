package main

import (
	closedGL "closed_gl/src/closedGL"
	_ "image/png"
	"runtime"

	"github.com/EngoEngine/glm"
	"github.com/go-gl/glfw/v3.2/glfw"
)

func main() {
	runtime.LockOSThread()

	var openGL = closedGL.InitClosedGL(800, 600)

	openGL.Window.SetScrollCallback(openGL.Camera.ScrollCb)
	openGL.Window.SetCursorPosCallback(openGL.Camera.MouseCallback)
	//openGL.Window.SetInputMode(glfw.CursorMode, glfw.CursorDisabled)
	var isWireframeMode = false

	glfw.SwapInterval(0)

	openGL.Camera.CameraPos = glm.Vec3{0, 0, 0}

	for !openGL.Window.ShouldClose() {
		var delta = openGL.FPSCounter.Delta
		_ = delta

		closedGL.ClearBG()
		if openGL.KeyBoardManager.IsPressed(glfw.KeyF) {
			isWireframeMode = !isWireframeMode
			closedGL.SetWireFrameMode(isWireframeMode)
		}
		/* for _, x := range chunks {
		x.Draw()
		} */
		openGL.DrawFPS(0, 0)

		openGL.Process()
	}
	openGL.Free()
}
