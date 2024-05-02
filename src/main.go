package main

import (
	"fmt"
	_ "image/png"
	"runtime"
	"strconv"
	"time"
	"unsafe"

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
	var dirtTex = loadImage("assets/tileset1.png", gl.RGBA)
	var vao, vbo uint32 = 0, 0
	generateBuffers(&vao, &vbo, nil, cube, 0, nil, []int{3, 3, 2})
	//projection = glm.Ident4()
	var chunks = []*Chunk{}

	//24
	for y := 0; y < 24; y++ {
		var posXZ = []float32{
			0, 0,
			16, 0,
			16, 16,
			16, 32,
			0, 32,
			-16, 0,
			-16, 16,
			-16, 32,
		}
		for i := 0; i < len(posXZ); i += 2 {
			var chunk = newChunk(glm.Vec3{16, 16, 16}, glm.Vec3{posXZ[i], float32(y * 16), posXZ[i+1]}, dirtTex, &c, &factory.projection3D, factory.shadermap["cube"])
			chunks = append(chunks, &chunk)
		}
	}

	var singleCube = factory.newCube(glm.Vec3{0, 3, 17}, dirtTex)

	window.SetScrollCallback(c.scrollCb)
	window.SetCursorPosCallback(c.mouseCallback)
	window.SetInputMode(glfw.CursorMode, glfw.CursorDisabled)
	gl.PolygonMode(gl.FRONT_AND_BACK, gl.FILL)

	var isWireframeMode = false

	glfw.SwapInterval(0)

	var test = []byte{1, 2, 3, 4, 5}
	test = append(test, 0x05)
	println(test)
	var ptr = (*float32)(unsafe.Pointer(&test[0]))
	*ptr = 4.5
	println(*ptr)

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
		for _, x := range chunks {
			x.draw()
		}

		singleCube.draw()
		text.x = 0
		text.y = 0
		text.draw("FPS: " + strconv.FormatInt(int64(fpsCounter.fpsAverage), 10) + "!")

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
