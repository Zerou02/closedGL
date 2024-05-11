package main

import (
	closedGL "closed_gl/src/closedGL"
	_ "image/png"
	"runtime"

	"github.com/EngoEngine/glm"
	"github.com/go-gl/gl/v2.1/gl"
	"github.com/go-gl/glfw/v3.2/glfw"
)

func main() {
	runtime.LockOSThread()

	var openGL = closedGL.InitClosedGL(800, 600)

	var dirtTex = closedGL.LoadImage("assets/tileset1alt.png", gl.RGBA)

	var chunks = []*closedGL.Chunk{}

	for x := 0; x < 10; x++ {
		for z := 0; z < 10; z++ {
			var chunk = openGL.Factory.NewChunk(glm.Vec3{32, 32, 32}, glm.Vec3{float32(x * 32), float32(0), float32(z * 32)}, dirtTex)
			chunks = append(chunks, chunk)
		}
	}
	openGL.Window.SetScrollCallback(openGL.Camera.ScrollCb)
	openGL.Window.SetCursorPosCallback(openGL.Camera.MouseCallback)
	openGL.Window.SetInputMode(glfw.CursorMode, glfw.CursorDisabled)
	var isWireframeMode = false

	glfw.SwapInterval(0)

	openGL.Camera.CameraPos = glm.Vec3{0, 0, 0}

	for !openGL.Window.ShouldClose() {
		closedGL.ClearBG()

		if openGL.KeyBoardManager.IsPressed(glfw.KeyF) {
			isWireframeMode = !isWireframeMode
			closedGL.SetWireFrameMode(isWireframeMode)
		}
		if openGL.KeyBoardManager.IsPressed(glfw.KeyP) {
			for i := 0; i < 10; i++ {
				var chunk = openGL.Factory.NewChunk(glm.Vec3{32, 32, 32}, glm.Vec3{0, 0, 0}, dirtTex)
				chunks = append(chunks, chunk)
			}
		}
		for _, x := range chunks {
			x.Draw()
		}
		openGL.DrawFPS(0, 0)
		openGL.Process()
	}
	openGL.Free()
}
