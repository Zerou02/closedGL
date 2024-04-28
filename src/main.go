package main

import (
	"bufio"
	"fmt"
	_ "image/png"
	"os"
	"runtime"
	"strconv"
	"strings"
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
	factory = newPrimitiveFactory2D(width, height)
	var keyboardManger = newKeyBoardManager(window)
	var fpsCounter = newFPSCounter()
	var fc = newFontCreator(30, 16, factory.shadermap["points"], &factory.projectionMatrix, &keyboardManger, window)
	text = newText("default", factory.shadermap["text"], 0, 500, 1, 1, glm.Vec3{1, 0, 1}, &factory.projectionMatrix)

	var file, _ = os.Open("assets/config.ini")
	var sc = bufio.NewScanner(file)
	for sc.Scan() {
		var line = sc.Text()
		if !strings.HasPrefix(line, "[") {
			var splitted = strings.Split(line, "=")
			if splitted[0] == "free_fps" {
				if splitted[1] == "false" {
					glfw.SwapInterval(1)
				} else {
					glfw.SwapInterval(0)
				}
			} else if splitted[0] == "potato-friendliness" {
				fc.autoUpdate = !stringToBool(splitted[1])
				println(fc.autoUpdate)
			} else if splitted[0] == "default_font" {
				fc.loadFont(splitted[1])
				text.deserializeIglbmf(splitted[1])
			}
		}
	}

	endTime("startup")

	fmt.Printf("%0x", byte(lerp(0, 255, 0.85)))
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

		text.x = 0
		text.y = 550
		text.draw("the quick brown fox jumps over the lazy dog. 0123456789")
		text.y = 570
		text.draw("THE QUICK BROWN FOX JUMPS OVER THE LAZY DOG")

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
