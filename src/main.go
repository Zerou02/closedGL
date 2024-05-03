package main

import (
	"fmt"
	_ "image/png"
	"runtime"
	"strconv"

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
var profiler Profiler

func main() {
	profiler = newProfiler()
	profiler.startTime("123")
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
	var chunks = []*Chunk{}

	profiler.startTime("chunks")
	for y := 0; y < 16; y++ {
		/* 		var posXZ = []float32{
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
		   		} */
	}
	profiler.endTime("chunks")
	var singleCube = factory.newCube(glm.Vec3{0, 3, 17}, dirtTex)

	window.SetScrollCallback(c.scrollCb)
	window.SetCursorPosCallback(c.mouseCallback)
	window.SetInputMode(glfw.CursorMode, glfw.CursorDisabled)
	gl.PolygonMode(gl.FRONT_AND_BACK, gl.FILL)

	var isWireframeMode = false

	profiler.startTime("1")
	profiler.endTime("1")

	glfw.SwapInterval(0)
	profiler.endTime("123")

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
		if keyboardManger.isPressed(glfw.KeyP) {
			for i := 0; i < 1; i++ {
				profiler.startTime("generating")
				var chunk = newChunk(glm.Vec3{16, 16, 16}, glm.Vec3{0, 0, 0}, dirtTex, &c, &factory.projection3D, factory.shadermap["cube"])
				profiler.endTime("generating")

				chunks = append(chunks, &chunk)
			}
		}
		if keyboardManger.isPressed(glfw.KeyL) {
			for _, x := range chunks {
				x.delete()
			}
			chunks = []*Chunk{}
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

func stringToBool(s string) bool {
	if s == "true" {
		return true
	} else {
		return false
	}
}

func decodeVertex(vertex uint32) {

	var val = uint32(vertex)
	println("val", val)
	var modelZ = val & 31
	val >>= 5
	var modelY = val & 31
	val >>= 5
	var modelX = val & 31
	val >>= 5
	var texY = val & 31
	val >>= 5
	var texX = val & 31
	val >>= 5
	var ndcIdx = val & 7

	var x = (ndcIdx >> 2) & 1
	var y = (ndcIdx >> 1) & 1
	var z = (ndcIdx >> 0) & 1

	println(modelX, ",", modelY, ";", modelZ)
	println(texX, ",", texY, ";")
	println(x, ",", y, ",", z)
}

func encodeVertex(ndc glm.Vec3, tex glm.Vec3, modelX, modelY, modelZ int) uint {
	//copy pos-3bit
	var vertex uint32 = 0
	if ndc[0] == 1.0 {
		vertex |= 0b100
	}
	if ndc[1] == 1.0 {
		vertex |= 0b010
	}
	if ndc[2] == 1.0 {
		vertex |= 0b001
	}
	vertex <<= 5
	//cumulatedVertices[vboSize*cubeStride+j*blockVerticesBytes] = float32(ndcBitMask)
	//copy tex-10bit
	var texX = byte(tex[0])
	var texY = byte(tex[1])
	vertex |= uint32(texX)
	vertex <<= 5
	vertex |= uint32(texY)
	vertex <<= 5

	//copy model-15bit
	vertex |= uint32(modelX)
	vertex <<= 5
	vertex |= uint32(modelY)
	vertex <<= 5
	vertex |= uint32(modelZ)
	return uint(vertex)
}

func printBitPattern(val uint32) {
	var str = strconv.FormatInt(int64(val), 2)
	for len(str) < 32 {
		str = "0" + str
	}
	fmt.Println(str, uint64(val)) //
}

func printBitPatternF(val float32) {
	var str = strconv.FormatInt(int64(val), 2)
	for len(str) < 32 {
		str = "0" + str
	}
	fmt.Println(str, uint64(val)) //
}

func printBitPatternF64(val float64) {
	var str = strconv.FormatInt(int64(val), 2)
	for len(str) < 32 {
		str = "0" + str
	}
	fmt.Println(str, uint64(val)) //
}
