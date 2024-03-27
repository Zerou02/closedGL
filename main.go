package main

import (
	_ "image/png"
	"runtime"

	"github.com/EngoEngine/glm"
	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/go-gl/glfw/v3.2/glfw"
	"neilpa.me/go-stbi"
)

const width = 800
const height = 600

var cubePos = []glm.Vec3{
	{0, 0, -5},
	{2.0, 5.0, -15.0},
	{-1.5, -2.2, -2.5},
	{-3.8, -2.0, -12.3},
	{2.4, -0.4, -3.5},
	{-1.7, 3.0, -7.5},
	{1.3, -2.0, -2.5},
	{1.5, 2.0, -2.5},
	{1.5, 0.2, -1.5},
	{-1.3, 1.0, -1.5},
}

var texCoords = []float32{
	0.0, 0.0, // lower-left corner
	1.0, 0.0, // lower-right corner
	0.5, 1.0, // top-center corner
}

var window *glfw.Window

type Vao = uint32
type Vbo = uint32
type Ebo = uint32
type Prog = uint32
type Texture = uint32

var vbo Vbo

var vao Vao
var ebo Ebo
var texture Texture
var texture2 Texture

var lastFrame float64 = 0
var deltaTime float64 = 0

var camera = CreateCamera()

var shader Shader

func main() {
	runtime.LockOSThread()

	window = initGlfw()
	initOpenGL()
	gl.Enable(gl.DEPTH_TEST)
	shader = initShader("./vertexShader.glsl", "./fragShader.glsl")
	window.SetInputMode(glfw.CursorMode, glfw.CursorDisabled)
	window.SetCursorPosCallback(camera.mouseCallback)
	window.SetScrollCallback(camera.scrollCb)

	loadImage(&texture, "./container.jpg", gl.RGBA)
	loadImage(&texture2, "./awesomeface.png", gl.RGBA)

	var vao2 Vao
	var vbo2 Vbo
	generateVBO(&vbo2, doubleTriangle)
	generateVAO(&vao2, []VertexInfo{{3, 0}, {3, 12}, {2, 24}})

	generateVBO(&vbo, cube)
	generateEBO(&ebo, indicesQuad)
	generateVAO(&vao, []VertexInfo{{3, 0}, {3, 12}, {2, 24}})

	shader.setUniformUInt("tex", texture)
	shader.setUniformUInt("tex2", texture2)

	for !window.ShouldClose() {
		var currFrame = glfw.GetTime()
		deltaTime = currFrame - lastFrame
		lastFrame = currFrame
		process(window)

		var rotVec = glm.Vec3{
			float32(glfw.GetTime()) * glm.DegToRad(50),
			float32(glfw.GetTime()) * glm.DegToRad(50),
			float32(glfw.GetTime()) * glm.DegToRad(50),
		}
		var modelMat = createTransformation(rotVec, glm.Vec3{1, 0, 0}, glm.Vec3{1, 1, 1})
		shader.setUniformMatrix4("view", &camera.lookAtMat)
		shader.setUniformMatrix4("projection", &camera.perspective)
		shader.setUniformMatrix4("model", &modelMat)

		//Draw
		gl.ClearColor(0.2, 0.3, 0.3, 1.0)
		gl.UseProgram(shader.prog)
		gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)
		gl.DrawArrays(gl.TRIANGLES, 0, 36)

		gl.PolygonMode(gl.FRONT_AND_BACK, gl.FILL)
		gl.UseProgram(shader.prog)
		gl.BindVertexArray(vao2)
		gl.DrawArrays(gl.TRIANGLES, 0, int32(len(doubleTriangle)/3))

		//	gl.PolygonMode(gl.FRONT_AND_BACK, gl.LINE)
		gl.ActiveTexture(gl.TEXTURE0)
		gl.BindTexture(gl.TEXTURE_2D, texture)
		gl.ActiveTexture(gl.TEXTURE1)
		gl.BindTexture(gl.TEXTURE_2D, texture2)
		gl.BindVertexArray(vao)
		for i := 0; i < 10; i++ {
			var angle float32 = 20.0 * (float32(i) + 1)
			var modelMat = createTransformation(glm.Vec3{angle, angle, angle}, glm.Vec3{cubePos[i].X(), cubePos[i].Y(), cubePos[i].Z()}, glm.Vec3{1, 1, 1})
			shader.setUniformMatrix4("model", &modelMat)
			gl.DrawArrays(gl.TRIANGLES, 0, 36)
		}
		glfw.PollEvents()
		window.SwapBuffers()
	}
	glfw.Terminate()
}

func process(window *glfw.Window) {
	if window.GetKey(glfw.KeyEscape) == glfw.Press {
		window.SetShouldClose(true)
	}
	camera.process(window)
}

func initGlfw() *glfw.Window {
	glfw.Init()
	glfw.WindowHint(glfw.Resizable, glfw.False)
	glfw.WindowHint(glfw.ContextVersionMajor, 4)
	glfw.WindowHint(glfw.ContextVersionMinor, 1)
	glfw.WindowHint(glfw.OpenGLProfile, glfw.OpenGLCoreProfile)
	glfw.WindowHint(glfw.OpenGLForwardCompatible, glfw.True)
	var window, _ = glfw.CreateWindow(width, height, "Nggght", nil, nil)
	window.MakeContextCurrent()
	return window
}

func initOpenGL() {
	gl.Init()
	gl.Viewport(0, 0, width, height)
}

func loadImage(texture *uint32, path string, format uint32) {
	var img, _ = stbi.Load(path)
	gl.GenTextures(1, texture)
	gl.BindTexture(gl.TEXTURE_2D, *texture)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_S, gl.REPEAT)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_T, gl.REPEAT)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MIN_FILTER, gl.LINEAR)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MAG_FILTER, gl.LINEAR)
	gl.TexImage2D(gl.TEXTURE_2D, 0, gl.RGBA, 512, 512, 0, format, gl.UNSIGNED_BYTE, gl.Ptr(img.Pix))
	gl.GenerateMipmap(gl.TEXTURE_2D)
}

func generateVBO(vbo *uint32, vertices []float32) {
	gl.GenBuffers(1, vbo)
	gl.BindBuffer(gl.ARRAY_BUFFER, *vbo)
	gl.BufferData(gl.ARRAY_BUFFER, 4*len(vertices), gl.Ptr(vertices), gl.STATIC_DRAW)
}

func generateEBO(ebo *uint32, indices []uint32) {
	gl.GenBuffers(1, ebo)
	gl.BindBuffer(gl.ELEMENT_ARRAY_BUFFER, *ebo)
	gl.BufferData(gl.ELEMENT_ARRAY_BUFFER, 4*len(indices), gl.Ptr(indices), gl.STATIC_DRAW)
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
	println(stride)
	for i := 0; i < len(vertexInfo); i++ {
		var info = vertexInfo[i]
		gl.VertexAttribPointerWithOffset(uint32(i), int32(info.amountBytes), gl.FLOAT, false, int32(stride*4), info.offset)
		gl.EnableVertexAttribArray(uint32(i))
	}
}
