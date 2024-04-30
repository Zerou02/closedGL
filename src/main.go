package main

import (
	"fmt"
	_ "image/png"
	"runtime"
	"strconv"
	"time"

	"github.com/EngoEngine/glm"
	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/go-gl/glfw/v3.2/glfw"
)

const width = 800
const height = 600

type Vao = uint32
type Vbo = uint32
type Ebo = uint32
type Prog = uint32
type Texture = uint32

var factory PrimitiveFactory
var text Text

func main() {
	startTime()
	runtime.LockOSThread()
	var window = initGlfw()
	initOpenGL()
	var c = CreateCamera()
	c.cameraPos = glm.Vec3{0, 3, 18}

	factory = newPrimitiveFactory2D(width, height, &c)
	var keyboardManger = newKeyBoardManager(window)
	var fpsCounter = newFPSCounter()
	text = newText("default", factory.shadermap["text"], 0, 500, 1, 1, glm.Vec3{1, 0, 1}, &factory.projectionMatrix)
	var dirtTex = loadImage("assets/dirt_side.jpg", gl.RGBA)
	var vao, vbo uint32 = 0, 0
	generateBuffers(&vao, &vbo, nil, cube, 0, nil, []int{3, 3, 2})
	//projection = glm.Ident4()
	var chunk = newChunk(glm.Vec3{16, 16, 16}, dirtTex)
	var singleCube = factory.newCube(glm.Vec3{0, 3, 17}, dirtTex)
	_, _ = singleCube, chunk

	window.SetScrollCallback(c.scrollCb)
	window.SetCursorPosCallback(c.mouseCallback)
	window.SetInputMode(glfw.CursorMode, glfw.CursorDisabled)
	gl.PolygonMode(gl.FRONT_AND_BACK, gl.FILL)

	for i := 0; i < 27; i++ {
		var vec = idxToPos3(i, glm.Vec3{3, 3, 3})
		var idx = pos3ToIdx(glm.Vec3{vec[1], vec[2], vec[0]}, glm.Vec3{3, 3, 3})
		fmt.Printf("%d: %f %f %f;; %d\n", i, vec[0], vec[1], vec[2], idx)
	}

	println()
	for i := 0; i < 9; i++ {
		var x, y = idxToGridPos(i, 3, 3)
		var idx = gridPosToIdx(x, y, 3)
		fmt.Printf("%d: x:%d, y:%d, idx:%d\n", i, x, y, idx)
	}

	var isWireframeMode = false
	for !window.ShouldClose() {
		gl.ClearColor(0.0, 0.0, 0.0, 1.0)
		gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)

		fpsCounter.process()
		c.process(window, float32(fpsCounter.delta))

		keyboardManger.process()
		if keyboardManger.isPressed(glfw.KeyF) {
			isWireframeMode = !isWireframeMode
			if isWireframeMode {
				gl.PolygonMode(gl.FRONT_AND_BACK, gl.FILL)
			} else {
				gl.PolygonMode(gl.FRONT_AND_BACK, gl.LINE)
			}
		}
		/* 		if keyboardManger.isPressed(glfw.KeyP) {
		   			cube.position[0] += 0.1
		   		}
		   		if keyboardManger.isPressed(glfw.KeyO) {
		   			cube.position[0] += -0.1
		   		} */
		chunk.draw()
		singleCube.draw()
		text.x = 0
		text.y = 0
		text.draw("FPS: " + strconv.FormatInt(int64(fpsCounter.fpsAverage), 10) + "!")

		glfw.SwapInterval(0)

		if fpsCounter.elapsed >= 0.5 {
			fpsCounter.calcAverage()
			fpsCounter.clear()
		}

		process(window)
		glfw.PollEvents()
		window.SwapBuffers()
	}
	glfw.Terminate()
}

func process(window *glfw.Window) {
	if window.GetKey(glfw.KeyEscape) == glfw.Press {
		window.SetShouldClose(true)
	}
}

var start = time.Now()
var end = time.Now()

func startTime() {
	start = time.Now()
}

func endTime(name string) {
	end = time.Now()
	var dur = end.Sub(start)
	fmt.Printf("%s:%f\n", name, dur.Seconds())
}

func stringToBool(s string) bool {
	if s == "true" {
		return true
	} else {
		return false
	}
}
