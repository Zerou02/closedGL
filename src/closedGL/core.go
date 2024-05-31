package closed_gl

import (
	"strconv"
	"unsafe"

	"github.com/EngoEngine/glm"
	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/go-gl/glfw/v3.2/glfw"
	"neilpa.me/go-stbi"
)

var factory PrimitiveFactory
var width float32
var height float32
var text Text

type ClosedGLContext struct {
	Window          *glfw.Window
	Factory         *PrimitiveFactory
	Camera          *Camera
	Text            *Text
	KeyBoardManager *KeyBoardManager
	FPSCounter      *FPSCounter
}

func InitClosedGL(pWidth, pHeight float32) ClosedGLContext {
	width = pWidth
	height = pHeight
	var window = initGlfw(int(width), int(height))
	var fpsCounter = NewFPSCounter()
	initOpenGL()

	var c = newCamera(width, height)
	factory = newPrimitiveFactory2D(float32(width), float32(height), &c)
	text = NewText("default", factory.Shadermap["text"], 0, 500, 1, 1, glm.Vec3{1, 0, 1}, &factory.projectionMatrix)
	var key = newKeyBoardManager(window)
	var con = ClosedGLContext{Window: window, Factory: &factory, Camera: &c, Text: &text, KeyBoardManager: &key, FPSCounter: &fpsCounter}
	return con
}

func initGlfw(width, height int) *glfw.Window {
	glfw.Init()
	glfw.WindowHint(glfw.Resizable, glfw.False)
	glfw.WindowHint(glfw.ContextVersionMajor, 4)
	glfw.WindowHint(glfw.ContextVersionMinor, 1)
	glfw.WindowHint(glfw.OpenGLProfile, glfw.OpenGLCoreProfile)
	glfw.WindowHint(glfw.OpenGLForwardCompatible, glfw.True)
	var window, _ = glfw.CreateWindow(width, height, "fast blazinglycraft", nil, nil)
	window.MakeContextCurrent()
	return window
}

func initOpenGL() {
	gl.Init()
	gl.Viewport(0, 0, int32(width), int32(height))
	gl.Enable(gl.DEPTH_TEST)
	gl.Enable(gl.BLEND)
	gl.Enable(gl.CULL_FACE)
	gl.Enable(gl.PROGRAM_POINT_SIZE)
	gl.PointSize(1)
	gl.BlendFunc(gl.SRC_ALPHA, gl.ONE_MINUS_SRC_ALPHA)
}

func LoadImage(path string, format uint32) *uint32 {
	var img, _ = stbi.Load(path)
	var texPtr uint32
	gl.GenTextures(1, &texPtr)
	gl.BindTexture(gl.TEXTURE_2D, texPtr)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_S, gl.CLAMP_TO_EDGE)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_T, gl.CLAMP_TO_EDGE)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MIN_FILTER, gl.NEAREST)

	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MAG_FILTER, gl.NEAREST)
	//gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MAX_LEVEL, 4)

	gl.TexImage2D(gl.TEXTURE_2D, 0, gl.RGBA, int32(img.Bounds().Max.X), int32(img.Bounds().Max.Y), 0, format, gl.UNSIGNED_BYTE, gl.Ptr(img.Pix))
	//gl.GenerateMipmap(gl.TEXTURE_2D)
	return &texPtr
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

func generateBuffers(vao, vbo, ebo *uint32, vertices []float32, vboByteLen int, indices []uint32, vertexAttribBytes []int) {
	generateBuffersSuper(vao, vbo, ebo, vertices, vboByteLen, indices, vertexAttribBytes, false)
}

func GenerateBuffers(vao, vbo, ebo *uint32, vertices []float32, vboByteLen int, indices []uint32, vertexAttribBytes []int) {
	generateBuffersSuper(vao, vbo, ebo, vertices, vboByteLen, indices, vertexAttribBytes, false)
}

func generateBuffersCopy2(vao, vbo, ebo *uint32, vertices []uint32, vboByteLen int, indices []uint32, vertexAttribBytes []int) {
	generateBuffersSuper2(vao, vbo, ebo, vertices, vboByteLen, indices, vertexAttribBytes, true)
}

func generateBuffersCopy(vao, vbo, ebo *uint32, vertices []float32, vboByteLen int, indices []uint32, vertexAttribBytes []int) {
	generateBuffersSuper(vao, vbo, ebo, vertices, vboByteLen, indices, vertexAttribBytes, true)
}
func generateBuffersSuper2(vao, vbo, ebo *uint32, vertices []uint32, vboByteLen int, indices []uint32, vertexAttribBytes []int, copyVertices bool) {
	//vbo
	gl.GenBuffers(1, vbo)
	gl.GenVertexArrays(1, vao)
	gl.BindVertexArray(*vao)
	gl.BindBuffer(gl.ARRAY_BUFFER, *vbo)

	if vertices != nil {
		var length = vboByteLen
		if vboByteLen <= 0 {
			length = 4 * len(vertices)
		}
		if copyVertices {
			gl.BufferData(gl.ARRAY_BUFFER, length, nil, gl.STATIC_DRAW)
			var baseAddr = (*uint32)(gl.MapBuffer(gl.ARRAY_BUFFER, gl.WRITE_ONLY))
			var gpuSlice = unsafe.Slice(baseAddr, length)
			copy(gpuSlice, vertices)
			gl.UnmapBuffer(gl.ARRAY_BUFFER)

		} else {
			gl.BufferData(gl.ARRAY_BUFFER, length, gl.Ptr(vertices), gl.STATIC_DRAW)
		}
	} else {
		gl.BufferData(gl.ARRAY_BUFFER, vboByteLen, gl.Ptr(nil), gl.DYNAMIC_DRAW)
	}
	//vao
	gl.BindVertexArray(*vao)
	var stride = 0
	for i := 0; i < len(vertexAttribBytes); i++ {
		stride += int(vertexAttribBytes[i])
	}
	var currOffset = 0
	for i := 0; i < len(vertexAttribBytes); i++ {
		gl.VertexAttribIPointerWithOffset(uint32(i), int32(vertexAttribBytes[i]), gl.UNSIGNED_INT, int32(stride*4), uintptr(currOffset)*4)
		gl.EnableVertexAttribArray(uint32(i))
		currOffset += vertexAttribBytes[i]
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

func generateBuffersSuper(vao, vbo, ebo *uint32, vertices []float32, vboByteLen int, indices []uint32, vertexAttribBytes []int, copyVertices bool) {
	//vbo
	gl.GenBuffers(1, vbo)
	gl.GenVertexArrays(1, vao)
	gl.BindVertexArray(*vao)
	gl.BindBuffer(gl.ARRAY_BUFFER, *vbo)

	if vertices != nil {
		var length = vboByteLen
		if vboByteLen <= 0 {
			length = 4 * len(vertices)
		}
		if copyVertices {
			gl.BufferData(gl.ARRAY_BUFFER, length, nil, gl.STATIC_DRAW)
			var baseAddr = (*float32)(gl.MapBuffer(gl.ARRAY_BUFFER, gl.WRITE_ONLY))
			for i := 0; i < len(vertices); i++ {
				var b = (*float32)(unsafe.Add(unsafe.Pointer(baseAddr), i*4))
				*b = vertices[i]
			}
			gl.UnmapBuffer(gl.ARRAY_BUFFER)
		} else {
			gl.BufferData(gl.ARRAY_BUFFER, length, gl.Ptr(vertices), gl.STATIC_DRAW)
		}
	} else {
		gl.BufferData(gl.ARRAY_BUFFER, vboByteLen, gl.Ptr(nil), gl.DYNAMIC_DRAW)
	}

	//vao
	gl.BindVertexArray(*vao)
	var stride = 0
	for i := 0; i < len(vertexAttribBytes); i++ {
		stride += int(vertexAttribBytes[i])
	}
	var currOffset = 0
	for i := 0; i < len(vertexAttribBytes); i++ {
		gl.VertexAttribPointerWithOffset(uint32(i), int32(vertexAttribBytes[i]), gl.FLOAT, false, int32(stride*4), uintptr(currOffset)*4)
		gl.EnableVertexAttribArray(uint32(i))
		currOffset += vertexAttribBytes[i]
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

func ClearBG() {
	gl.ClearColor(0.0, 0.0, 0.0, 1.0)
	gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)
}

func SetWireFrameMode(val bool) {
	if val {
		gl.PolygonMode(gl.FRONT_AND_BACK, gl.FILL)
	} else {
		gl.PolygonMode(gl.FRONT_AND_BACK, gl.LINE)

	}
}

func (this *ClosedGLContext) Process() {
	if this.Window.GetKey(glfw.KeyEscape) == glfw.Press {
		this.Window.SetShouldClose(true)
	}
	this.KeyBoardManager.Process()
	this.FPSCounter.Process()
	this.Camera.Process(this.Window, float32(0.16))
	glfw.PollEvents()
	this.Window.SwapBuffers()
	if this.FPSCounter.Elapsed >= 0.5 {
		this.FPSCounter.CalcAverage()
		this.FPSCounter.Clear()
	}
}

func (this *ClosedGLContext) Free() {
	glfw.Terminate()
}

func (this *ClosedGLContext) DrawFPS(posX, posY int) {
	this.Text.DrawText(posX, posY, "FPS: "+strconv.FormatInt(int64(this.FPSCounter.FpsAverage), 10)+"!")

}
