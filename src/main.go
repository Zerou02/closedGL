package main

import (
	_ "image/png"
	"math"
	"os"
	"runtime"

	"github.com/EngoEngine/glm"
	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/go-gl/glfw/v3.2/glfw"
	"neilpa.me/go-stbi"
)

const width = 800
const height = 600

type Vao = uint32
type Vbo = uint32
type Ebo = uint32
type Prog = uint32
type Texture = uint32

var idxToCharMap = []byte{'a', 'b', 'c'}

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
	var ballTex = loadImage("./assets/ball.png", gl.RGBA)
	_ = ballTex
	gl.PolygonMode(gl.FRONT_AND_BACK, gl.FILL)
	gl.UseProgram(shader.prog)
	shader.setUniformMatrix4("projection", &projection)
	shader.setUniformMatrix4("view", &view)

	var delta = 0.0
	_ = delta
	var lastFrame = 0.0

	var size, amount = 30, 16
	var rects = generateGrid(size, amount, &pointShader, &projection)
	var rects2 = generateGrid(size, amount, &pointShader, &projection)
	var rects3 = generateGrid(size, amount, &pointShader, &projection)

	var rectHolder = [][]Rectangle{rects, rects2, rects3}
	var _, info = deserializeIglbmf("a", &rectHolder)
	var ball = newText(info, &textShader, 0, 0, 1, 1, glm.Vec3{1, 0, 1}, &projection)
	var letters = []rune{'a', 'b', 'c'}
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
				serializeIglbmf(rectHolder, string(letters[currentIdx]))
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
				}
				qPressed = true
			}
		}
		if window.GetKey(glfw.KeyQ) == glfw.Release {
			qPressed = false
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
		ball.draw()
		gl.Enable(gl.DEPTH_TEST)

		process(window)
		glfw.PollEvents()
		window.SwapBuffers()
	}
	glfw.Terminate()
}

func aabbAabbCol(b1, b2 glm.Vec4) bool {
	var colX = b1.X()+b1.Z() >= b2[0] && b2[0]+b2[2] >= b1[0]
	var colY = b1[1]+b1[3] >= b2[1] && b2[1]+b2[3] >= b1[1]
	return colX && colY
}
func process(window *glfw.Window) {
	if window.GetKey(glfw.KeyEscape) == glfw.Press {
		window.SetShouldClose(true)
	}
}

func initGlfw() *glfw.Window {
	glfw.Init()
	glfw.WindowHint(glfw.Resizable, glfw.False)
	glfw.WindowHint(glfw.ContextVersionMajor, 4)
	glfw.WindowHint(glfw.ContextVersionMinor, 1)
	glfw.WindowHint(glfw.OpenGLProfile, glfw.OpenGLCoreProfile)
	glfw.WindowHint(glfw.OpenGLForwardCompatible, glfw.True)
	var window, _ = glfw.CreateWindow(width, height, "light - i hate packages", nil, nil)
	window.MakeContextCurrent()
	return window
}

func initOpenGL() {
	gl.Init()
	gl.Viewport(0, 0, width, height)
	gl.Enable(gl.DEPTH_TEST)
	gl.Enable(gl.BLEND)
	gl.Enable(gl.PROGRAM_POINT_SIZE)
	gl.PointSize(1)
	gl.BlendFunc(gl.SRC_ALPHA, gl.ONE_MINUS_SRC_ALPHA)
}

func loadImage(path string, format uint32) *uint32 {
	var img, _ = stbi.Load(path)
	var texPtr uint32
	gl.GenTextures(1, &texPtr)
	gl.BindTexture(gl.TEXTURE_2D, texPtr)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_S, gl.REPEAT)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_T, gl.REPEAT)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MIN_FILTER, gl.NEAREST)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MAG_FILTER, gl.NEAREST)
	gl.TexImage2D(gl.TEXTURE_2D, 0, gl.RGBA, int32(img.Bounds().Max.X), int32(img.Bounds().Max.Y), 0, format, gl.UNSIGNED_BYTE, gl.Ptr(img.Pix))
	return &texPtr
}

func loadIglbmf(path string) (*Texture, CharacterInfo) {
	var file, _ = os.ReadFile("./font/" + path + "alt.iglbmf")
	var buffer = []byte{}
	var c = CharacterInfo{tex: nil, texW: uint32(file[0]), texH: uint32(file[0]), charX: uint32(file[1]), charY: uint32(file[2]), charW: uint32(file[3]), charH: uint32(file[4]), asciicode: file[5]}
	var dataOffset = int(file[6])
	for i := dataOffset; i < len(file); i++ {
		if file[i] == 0x01 {
			buffer = append(buffer, 0x00)
			buffer = append(buffer, 0xFF)
			buffer = append(buffer, 0x00)
			buffer = append(buffer, 0xFF)
		} else {
			buffer = append(buffer, 0x00)
			buffer = append(buffer, 0x00)
			buffer = append(buffer, 0x00)
			buffer = append(buffer, 0xFF)
		}
	}
	var texPtr uint32
	gl.GenTextures(1, &texPtr)
	gl.BindTexture(gl.TEXTURE_2D, texPtr)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_R, gl.REPEAT)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_T, gl.REPEAT)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MIN_FILTER, gl.NEAREST)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MAG_FILTER, gl.NEAREST)
	gl.TexImage2D(gl.TEXTURE_2D, 0, gl.RGBA, int32(16), int32(16), 0, gl.RGBA, gl.UNSIGNED_BYTE, gl.Ptr(buffer))
	return &texPtr, c
}

func deserializeIglbmf(path string, grid *[][]Rectangle) (*Texture, []CharacterInfo) {
	var file, _ = os.ReadFile("./font/" + path + "combined.iglbmf")
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
	//inklusiv, exklusiv

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
				var info = CharacterInfo{tex: &texPtr, texW: uint32(imgLenPx), texH: uint32(imgLenPx), asciicode: chunk[5], charX: uint32(chunk[1]), charY: uint32(chunk[2]), charW: uint32(chunk[3]), charH: uint32(chunk[4])}
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

	gl.GenTextures(1, &texPtr)
	gl.BindTexture(gl.TEXTURE_2D, texPtr)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_R, gl.REPEAT)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_T, gl.REPEAT)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MIN_FILTER, gl.NEAREST)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MAG_FILTER, gl.NEAREST)
	gl.TexImage2D(gl.TEXTURE_2D, 0, gl.RGBA, int32(imgLenPx), int32(imgLenPx), 0, gl.RGBA, gl.UNSIGNED_BYTE, gl.Ptr(texData))
	return &texPtr, charInfo
}

func printlnTexData(texData []byte) {
	for i := 0; i < 16*32*4*2; i += 4 {
		if texData[i+1] == 0xFF {
			print("1")
		} else {
			print("0")
		}
		if i%128 == 0 {
			println()
		}
	}
}
func generateVBO(vbo *uint32, vertices []float32) {
	gl.GenBuffers(1, vbo)
	gl.BindBuffer(gl.ARRAY_BUFFER, *vbo)
	gl.BufferData(gl.ARRAY_BUFFER, 4*len(vertices), gl.Ptr(vertices), gl.STATIC_DRAW)
}

func generateDynVBO(vbo *uint32, bytesLen int) {
	gl.GenBuffers(1, vbo)
	gl.BindBuffer(gl.ARRAY_BUFFER, *vbo)
	gl.BufferData(gl.ARRAY_BUFFER, bytesLen, gl.Ptr(nil), gl.DYNAMIC_DRAW)
}
func generateEBO(ebo *uint32, indices []uint32) {
	gl.GenBuffers(1, ebo)
	gl.BindBuffer(gl.ELEMENT_ARRAY_BUFFER, *ebo)
	gl.BufferData(gl.ELEMENT_ARRAY_BUFFER, 4*len(indices), gl.Ptr(indices), gl.STATIC_DRAW)
}

func generateBuffers(vao, vbo, ebo *uint32, vertices []float32, vboByteLen int, indices []uint32, vertexInfo []VertexInfo) {

	//vbo
	gl.GenBuffers(1, vbo)
	gl.GenVertexArrays(1, vao)
	gl.BindVertexArray(*vao)
	gl.BindBuffer(gl.ARRAY_BUFFER, *vbo)

	if vertices != nil {
		gl.BufferData(gl.ARRAY_BUFFER, 4*len(vertices), gl.Ptr(vertices), gl.STATIC_DRAW)
	} else {
		gl.BufferData(gl.ARRAY_BUFFER, vboByteLen, gl.Ptr(nil), gl.DYNAMIC_DRAW)
	}

	//vao
	gl.BindVertexArray(*vao)
	var stride = 0
	for i := 0; i < len(vertexInfo); i++ {
		stride += int(vertexInfo[i].amountBytes)
	}
	for i := 0; i < len(vertexInfo); i++ {
		var info = vertexInfo[i]
		gl.VertexAttribPointerWithOffset(uint32(i), int32(info.amountBytes), gl.FLOAT, false, int32(stride*4), info.offset)
		gl.EnableVertexAttribArray(uint32(i))
	}

	//ebo
	if ebo != nil {
		gl.GenBuffers(1, ebo)
		gl.BindBuffer(gl.ELEMENT_ARRAY_BUFFER, *ebo)
		gl.BufferData(gl.ELEMENT_ARRAY_BUFFER, 4*len(indices), gl.Ptr(indices), gl.STATIC_DRAW)
	}
	gl.BindVertexArray(0)
	gl.BindBuffer(gl.ARRAY_BUFFER, 0)
	gl.BindBuffer(gl.ELEMENT_ARRAY_BUFFER, 0)

}

type VertexInfo struct {
	amountBytes int
	offset      uintptr
}

func generateVAO(vao *uint32, vertexInfo []VertexInfo) {
	gl.GenVertexArrays(1, vao)
	gl.BindVertexArray(*vao)
	var stride = 0
	for i := 0; i < len(vertexInfo); i++ {
		stride += int(vertexInfo[i].amountBytes)
	}
	for i := 0; i < len(vertexInfo); i++ {
		var info = vertexInfo[i]
		gl.VertexAttribPointerWithOffset(uint32(i), int32(info.amountBytes), gl.FLOAT, false, int32(stride*4), info.offset)
		gl.EnableVertexAttribArray(uint32(i))
	}
}

func aabbCircleCol(circle glm.Vec3, aabb glm.Vec4) (bool, Direction, glm.Vec2) {
	var centre = glm.Vec2{circle[0] + circle[2], circle[1] + circle[2]}
	var aabbHalf = glm.Vec2{aabb[2] / 2, aabb[3] / 2}
	var aabbCentre = glm.Vec2{aabb[0] + aabb[2]/2, aabb[1] + aabb[3]/2}
	var diff = centre.Sub(&aabbCentre)
	var clamped = glm.Vec2{glm.Clamp(diff[0], -aabbHalf[0], aabbHalf[0]), glm.Clamp(diff[1], -aabbHalf[1], aabbHalf[1])}
	var closest = aabbCentre.Add(&clamped)
	diff = closest.Sub(&centre)
	if diff.Len() < circle[2] {
		return true, vectorDirection(diff), diff
	} else {
		return false, UP, glm.Vec2{0, 0}
	}

}
func clamp(val, min, max float32) float32 {
	return float32(math.Max(float64(min), math.Min(float64(max), float64(val))))
}

func ssVectorOrigionCol(ssVel, ssWall glm.Vec2) glm.Vec2 {

	var esVel = ssVel.ComponentProduct(&glm.Vec2{1, -1})
	var angle = glm.RadToDeg(float32(math.Acos(float64(ssWall.Dot(&esVel) / (esVel.Len() * ssWall.Len())))))
	var rotangle = 2 * angle
	if angle == 0 {
		rotangle = 180
	}
	var rotMat = glm.Rotate2D(glm.DegToRad(360 - rotangle))
	var newAngle = rotMat.Mul2x1(&esVel)
	newAngle.Normalize()
	return newAngle.ComponentProduct(&glm.Vec2{1, -1})
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
		var chunk = gridToChunk(x, idxToCharMap[i])
		arr = append(arr, chunk...)
	}
	var file, _ = os.Create("./font/" + path + "combined.iglbmf")
	file.Write(arr)
	file.Close()
}

func loadChunkInRect(grid *[]Rectangle, chunk []byte) {
	var dataOffset = int(chunk[6])
	for i := dataOffset; i < len(chunk); i++ {
		if chunk[i] == 0x01 {
			(*grid)[i-dataOffset].visible = true
		} else {
			(*grid)[i-dataOffset].visible = false
		}
	}
}
func loadData(grid *[]Rectangle, path string) {
	var content, _ = os.ReadFile("./font/" + path + "alt.iglbmf")
	var dataOffset = int(content[6])
	for i := dataOffset; i < len(content); i++ {
		if content[i] == 0x01 {
			(*grid)[i-dataOffset].visible = true
		} else {
			(*grid)[i-dataOffset].visible = false
		}
	}
	// var file, _ = os.Open("./font/" + path + ".iglbmf")
	// (*grid)[0]
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
