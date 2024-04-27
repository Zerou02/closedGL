package main

import (
	"fmt"
	_ "image/png"
	"math"
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

	var delta = 0.0
	_ = delta
	var lastFrame = 0.0

	var size, amount = 30, 16

	var fc = newFontCreator(&pointShader, &projection)

	var _, info = deserializeIglbmf("default")
	startTime()

	var text = newText(info, &textShader, 0, 500, 1, 1, glm.Vec3{1, 0, 1}, &projection)
	endTime("startup")

	var currentIdx = 0

	gl.PolygonMode(gl.FRONT_AND_BACK, gl.FILL)
	var pPressed = false
	var ePressed = false
	var qPressed = false
	var rPressed = false
	var tPressed = false

	var frameCount int64 = 0
	glfw.SwapInterval(0)
	var elapsed = 0.0
	var fpsSum = 0.0
	var fpsAmount = 0
	var fpsAverage = 0
	for !window.ShouldClose() {
		frameCount += 1
		fpsAmount += 1
		fpsSum += 1 / delta
		var currFrame = glfw.GetTime()
		delta = currFrame - lastFrame
		elapsed += delta
		lastFrame = currFrame

		var mouseX, mouseY = window.GetCursorPos()
		if mouseX > 0 && mouseX < float64(size*amount) && mouseY > 0 && mouseY < float64(size*amount) {
			var gridX, gridY int = int(mouseX) / size, int(mouseY) / size
			var idx = gridY*amount + gridX
			_ = idx
			if window.GetMouseButton(glfw.MouseButton1) == glfw.Press {
				fc.grids[currentIdx].cells[idx].visible = true
			}
			if window.GetMouseButton(glfw.MouseButton2) == glfw.Press {
				fc.grids[currentIdx].cells[idx].visible = false
			}
		}
		if window.GetKey(glfw.KeyP) == glfw.Press {
			if !pPressed {
				text.serializeIglbmf(fc.grids, "default")
				pPressed = true
			}
		}
		if window.GetKey(glfw.KeyP) == glfw.Release {
			pPressed = false
		}

		if window.GetKey(glfw.KeyE) == glfw.Press {
			if !ePressed {
				if currentIdx < len(fc.grids)-1 {
					currentIdx += 1
					println(currentIdx, string(rune(currentIdx)))
				}
				ePressed = true
			}
		}
		if window.GetKey(glfw.KeyE) == glfw.Release {
			ePressed = false
		}

		if window.GetKey(glfw.KeyQ) == glfw.Press {
			if !qPressed {
				if currentIdx > 0 {
					currentIdx -= 1
					println(currentIdx, string(rune(currentIdx)))

				}
				qPressed = true
			}
		}
		if window.GetKey(glfw.KeyQ) == glfw.Release {
			qPressed = false
		}

		if window.GetKey(glfw.KeyR) == glfw.Press {
			if !rPressed {
				currentIdx -= 10
				if currentIdx < 0 {
					currentIdx = 0
				}
				println(currentIdx, string(rune(currentIdx)))
				rPressed = true
			}
		}
		if window.GetKey(glfw.KeyR) == glfw.Release {
			rPressed = false
		}

		if window.GetKey(glfw.KeyT) == glfw.Press {
			if !tPressed {
				currentIdx += 10
				if currentIdx > 127 {
					currentIdx = 127
				}
				println(currentIdx, string(rune(currentIdx)))
				tPressed = true
			}
		}
		if window.GetKey(glfw.KeyT) == glfw.Release {
			tPressed = false
		}
		gl.ClearColor(0.0, 0.0, 0.0, 1.0)
		gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)
		gl.Disable(gl.DEPTH_TEST)

		startTime()
		fc.draw(currentIdx)
		//	endTime("grindRender: ")

		startTime()
		//	endTime("lineRender: ")
		text.x = 0
		text.y = 450
		text.draw("1234567890")
		text.y = 500

		if elapsed >= 0.5 {
			elapsed = 0
			fpsAverage = int(fpsSum / float64(fpsAmount))
			fpsSum = 0
			fpsAmount = 0
		}
		text.draw("FPS: " + strconv.FormatInt(int64(fpsAverage), 10))

		gl.Enable(gl.DEPTH_TEST)

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

func gridToChunk(grid []Rectangle, asciicode byte) []byte {
	var chunk = make([]byte, len(grid)+7)
	var topmostY, bottommostY, rightmostX, leftmostX int = 16, 0, 0, 16
	for i := 0; i < len(grid); i++ {
		if grid[i].visible {
			var gridX, gridY = idxToGridPos(i, 16, 16)
			if gridX < leftmostX {
				leftmostX = gridX
			}
			if gridX > rightmostX {
				rightmostX = gridX
			}
			if gridY < topmostY {
				topmostY = gridY
			}
			if gridY > bottommostY {
				bottommostY = gridY
			}
		}
	}
	//gridSize,[4]charDim,asciicode,dataOffset
	chunk[0] = byte(math.Sqrt(float64(len(grid))))
	chunk[1] = byte(leftmostX)
	chunk[2] = byte(topmostY)
	chunk[3] = byte(rightmostX) - byte(leftmostX) + 1
	chunk[4] = byte(bottommostY) - byte(topmostY) + 1
	chunk[5] = asciicode
	chunk[6] = 7

	for i, x := range grid {
		if x.visible {
			chunk[i+int(chunk[6])] = 1
		} else {
			chunk[i+int(chunk[6])] = 0
		}
	}
	return chunk
}

func loadChunkInRect(grid *[]Rectangle, chunk []byte) {
	var dataOffset = int(chunk[6])
	if dataOffset == 0 {
		dataOffset = 7
	}
	for i := dataOffset; i < len(chunk); i++ {
		if chunk[i] == 0x01 {
			(*grid)[i-dataOffset].visible = true
		} else {
			(*grid)[i-dataOffset].visible = false
		}
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
