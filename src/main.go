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

func main() {

	runtime.LockOSThread()
	var window = initGlfw()
	initOpenGL()

	var shader = initShader("./shader/base.vs", "./shader/base.fs")
	var view = glm.Ident4()

	var textShader = initShader("./shader/text.vs", "./shader/text.fs")
	_ = textShader

	var pointShader = initShader("./shader/points.vs", "./shader/points.fs")
	_ = pointShader

	var projection = glm.Ortho(0, width, height, 0, -1, 1)

	gl.PolygonMode(gl.FRONT_AND_BACK, gl.FILL)
	gl.UseProgram(shader.prog)
	shader.setUniformMatrix4("projection", &projection)
	shader.setUniformMatrix4("view", &view)
	var keyboardManger = newKeyBoardManager(window)

	_ = keyboardManger

	var fpsCounter = newFPSCounter()
	var fc = newFontCreator(30, 16, &pointShader, &projection, &keyboardManger, window)

	startTime()

	var text = newText("default", &textShader, 0, 500, 1, 1, glm.Vec3{1, 0, 1}, &projection)
	endTime("startup")

	gl.PolygonMode(gl.FRONT_AND_BACK, gl.FILL)

	glfw.SwapInterval(0)
	for !window.ShouldClose() {
		fpsCounter.process()
		keyboardManger.process()
		fc.process()

		gl.ClearColor(0.0, 0.0, 0.0, 1.0)
		gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)

		fc.draw()
		text.x = 0
		text.y = 500
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
