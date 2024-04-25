package main

import (
	"fmt"
	_ "image/png"
	"math"
	"os"
	"runtime"
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
	var projection = glm.Ortho(0, width, height, 0, -1, 1)
	var view = glm.Ident4()

	var textShader = initShader("./shader/text.vs", "./shader/text.fs")
	_ = textShader

	var pointShader = initShader("./shader/points.vs", "./shader/points.fs")
	_ = pointShader
	gl.PolygonMode(gl.FRONT_AND_BACK, gl.FILL)
	gl.UseProgram(shader.prog)
	shader.setUniformMatrix4("projection", &projection)
	shader.setUniformMatrix4("view", &view)

	var delta = 0.0
	_ = delta
	var lastFrame = 0.0

	var size, amount = 30, 16

	var rectHolder = [][]Rectangle{}
	for i := 0; i < 128; i++ {
		rectHolder = append(rectHolder, generateGrid(size, amount, &pointShader, &projection))
	}

	var _, info = deserializeIglbmf("default", &rectHolder)
	var text = newText(info, &textShader, 0, 500, 1, 1, glm.Vec3{1, 0, 1}, &projection)
	var currentIdx = 0
	var lines = []Line{}
	for y := 0; y < amount+1; y++ {
		var p1 = newPoint(&pointShader, glm.Vec2{1, float32(y * size)}, glm.Vec3{1, 0, 0}, &projection)
		var p2 = newPoint(&pointShader, glm.Vec2{1 + float32(size)*float32(amount), float32(y * size)}, glm.Vec3{0, 0, 1}, &projection)
		var line = newLine(&pointShader, &projection)
		line.addPoint(p1)
		line.addPoint(p2)
		lines = append(lines, line)
	}
	for x := 0; x < amount+1; x++ {
		var p1 = newPoint(&pointShader, glm.Vec2{1 + float32(x*size), 0}, glm.Vec3{1, 0, 0}, &projection)
		var p2 = newPoint(&pointShader, glm.Vec2{1 + float32(x*size), float32(size) * float32(amount)}, glm.Vec3{0, 0, 1}, &projection)
		var line = newLine(&pointShader, &projection)
		line.addPoint(p1)
		line.addPoint(p2)
		lines = append(lines, line)
	}

	gl.PolygonMode(gl.FRONT_AND_BACK, gl.FILL)
	var pPressed = false
	var ePressed = false
	var qPressed = false
	var rPressed = false
	var tPressed = false

	println(text.charInfo[0].asciicode)
	for !window.ShouldClose() {

		var currFrame = glfw.GetTime()
		delta = currFrame - lastFrame
		lastFrame = currFrame

		var mouseX, mouseY = window.GetCursorPos()
		if mouseX > 0 && mouseX < float64(size*amount) && mouseY > 0 && mouseY < float64(size*amount) {
			var gridX, gridY int = int(mouseX) / size, int(mouseY) / size
			var idx = gridY*amount + gridX
			if window.GetMouseButton(glfw.MouseButton1) == glfw.Press {
				rectHolder[currentIdx][idx].visible = true
			}
			if window.GetMouseButton(glfw.MouseButton2) == glfw.Press {
				rectHolder[currentIdx][idx].visible = false
			}
		}
		if window.GetKey(glfw.KeyP) == glfw.Press {
			if !pPressed {
				serializeIglbmf(rectHolder, "default")
				pPressed = true
			}
		}
		if window.GetKey(glfw.KeyP) == glfw.Release {
			pPressed = false
		}

		if window.GetKey(glfw.KeyE) == glfw.Press {
			if !ePressed {
				if currentIdx < len(rectHolder)-1 {
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
		for _, x := range rectHolder[currentIdx] {
			x.draw()
		}

		for _, x := range lines {
			x.draw()
		}
		text.draw("Hello,World!")
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

func deserializeIglbmf(path string, grid *[][]Rectangle) (*Texture, []CharacterInfo) {
	var start = time.Now()
	var file, _ = os.ReadFile("./font/" + path + "combined.iglbmf")
	var end = time.Now()

	var charInfo = []CharacterInfo{}
	var texData = []byte{}
	var texPtr uint32

	var chunkW = int(file[0])
	var dataOffset = int(file[6])
	var chunkSize = chunkW*chunkW + dataOffset
	var amountChunks = len(file) / chunkSize
	var texRowLen = int(math.Ceil(math.Sqrt(float64(amountChunks))))
	var texRowHeight = texRowLen
	var chunksPerRow = texRowLen
	var imgLenPx = texRowLen * chunkW
	var texRowHeightPx = chunkW
	var chunkPxW = imgLenPx / texRowLen

	for texLine := 0; texLine < texRowHeight; texLine++ {
		var chunks = [][]byte{}
		for i := 0; i < chunksPerRow; i++ {
			var chunk = []byte{}
			var idx = (i + texLine*chunksPerRow)
			if idx*chunkSize >= len(file) {
				chunk = make([]byte, chunkSize)
			} else {
				chunk = file[idx*chunkSize : (idx+1)*chunkSize]
			}
			if idx < len(*grid) {
				loadChunkInRect(&(*grid)[idx], chunk)
				var posX, posY = idxToGridPos(idx, texRowLen, texRowLen)
				var info = CharacterInfo{
					tex: &texPtr, texW: uint32(imgLenPx), texH: uint32(imgLenPx),
					asciicode: chunk[5], charX: uint32(chunk[1]), charY: uint32(chunk[2]),
					charW: uint32(chunk[3]), charH: uint32(chunk[4]),
					offsetX: uint32(posX) * uint32(chunkW), offsetY: uint32(posY) * uint32(chunkPxW),
				}
				charInfo = append(charInfo, info)
			}
			chunks = append(chunks, chunk)
		}
		for y := 0; y < texRowHeightPx; y++ {
			for i := 0; i < chunksPerRow; i++ {
				var currChunkData = chunks[i][y*chunkPxW+dataOffset : (y+1)*chunkPxW+dataOffset]
				for j := 0; j < chunkPxW; j++ {
					texData = append(texData, 0x00)
					if currChunkData[j] == 0x01 {
						texData = append(texData, 0xFF)
					} else {
						texData = append(texData, 0x00)
					}
					texData = append(texData, 0x00)
					texData = append(texData, 0xFF)
				}
			}
		}
	}

	start = time.Now()

	gl.GenTextures(1, &texPtr)
	gl.BindTexture(gl.TEXTURE_2D, texPtr)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_R, gl.REPEAT)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_T, gl.REPEAT)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MIN_FILTER, gl.NEAREST)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MAG_FILTER, gl.NEAREST)
	gl.TexImage2D(gl.TEXTURE_2D, 0, gl.RGBA, int32(imgLenPx), int32(imgLenPx), 0, gl.RGBA, gl.UNSIGNED_BYTE, gl.Ptr(texData))

	end = time.Now()
	fmt.Printf("other:%f", end.Sub(start).Seconds())
	return &texPtr, charInfo
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

func serializeIglbmf(grid [][]Rectangle, path string) {
	var arr = []byte{}
	for i, x := range grid {
		var chunk = gridToChunk(x, byte(i))
		arr = append(arr, chunk...)
	}
	var file, _ = os.Create("./font/" + path + "combined.iglbmf")
	file.Write(arr)
	file.Close()
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

func generateGrid(size, amount int, shader *Shader, projection *glm.Mat4) []Rectangle {
	var rects = []Rectangle{}
	for y := 0; y < amount; y++ {
		for x := 0; x < amount; x++ {
			rects = append(rects, newRect(shader, projection, glm.Vec4{float32(x * size), float32(y * size), float32(size), float32(size)}, glm.Vec3{0, 1, 0}))
		}
	}
	return rects
}
