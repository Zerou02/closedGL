package main

import (
	closedGL "closed_gl/src/test"
	closed_gl "closed_gl/src/test"
	_ "image/png"
	"runtime"
	"strconv"

	"github.com/EngoEngine/glm"
	"github.com/go-gl/gl/v2.1/gl"
	"github.com/go-gl/glfw/v3.2/glfw"
)

const width = 800
const height = 600

func main() {
	runtime.LockOSThread()

	var openGL = closedGL.InitClosedGL(800, 600)

	var fpsCounter = closedGL.NewFPSCounter()
	var dirtTex = closedGL.LoadImage("assets/tileset1alt.png", gl.RGBA)
	var altTex = closedGL.LoadImage("assets/dirt_side.jpg", gl.RGBA)
	var testTex = closedGL.LoadImage("assets/tileset2.png", gl.RGBA)
	_ = testTex

	_ = fpsCounter
	_ = dirtTex
	_ = altTex

	var chunks = []*closedGL.Chunk{}

	for x := 0; x < 10; x++ {
		for z := 0; z < 10; z++ {
			_, _, _ = dirtTex, altTex, testTex
			var chunk = closedGL.NewChunk(glm.Vec3{32, 32, 32}, glm.Vec3{float32(x * 32), float32(0), float32(z * 32)}, dirtTex, openGL.Camera, &openGL.Factory.Projection3D, openGL.Factory.Shadermap["cube"])
			chunks = append(chunks, &chunk)
		}
	}
	openGL.Window.SetScrollCallback(openGL.Camera.ScrollCb)
	openGL.Window.SetCursorPosCallback(openGL.Camera.MouseCallback)
	openGL.Window.SetInputMode(glfw.CursorMode, glfw.CursorDisabled)
	var isWireframeMode = false

	glfw.SwapInterval(0)

	openGL.Camera.CameraPos = glm.Vec3{0, 0, 0}

	for !openGL.Window.ShouldClose() {
		closed_gl.ClearBG()
		fpsCounter.Process()

		openGL.Camera.Process(openGL.Window, float32(fpsCounter.Delta))
		openGL.KeyBoardManager.Process()
		if openGL.KeyBoardManager.IsPressed(glfw.KeyF) {
			isWireframeMode = !isWireframeMode
			closedGL.SetWireFrameMode(isWireframeMode)
		}
		if openGL.KeyBoardManager.IsPressed(glfw.KeyP) {
			for i := 0; i < 10; i++ {
				var chunk = closedGL.NewChunk(glm.Vec3{32, 32, 32}, glm.Vec3{float32(0), float32(0), float32(0)}, dirtTex, openGL.Camera, &openGL.Factory.Projection3D, openGL.Factory.Shadermap["cube"])
				chunks = append(chunks, &chunk)
			}
		}
		for _, x := range chunks {
			x.Draw()
		}
		openGL.Text.DrawText(0, 0, "FPS: "+strconv.FormatInt(int64(fpsCounter.FpsAverage), 10)+"!")

		if fpsCounter.Elapsed >= 0.5 {
			fpsCounter.CalcAverage()
			fpsCounter.Clear()
		}

		process(openGL.Window)
		glfw.PollEvents()
		openGL.Window.SwapBuffers()
	}
	glfw.Terminate()
}

func process(window *glfw.Window) {
	if window.GetKey(glfw.KeyEscape) == glfw.Press {
		window.SetShouldClose(true)
	}
}
