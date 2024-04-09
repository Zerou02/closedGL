package main

import (
	_ "image/png"
	"math"
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

var cubeVbo Vbo
var cubeVao Vao
var cubeEbo Ebo

var lastFrame float64 = 0
var deltaTime float64 = 0

var camera = CreateCamera()

func main() {
	runtime.LockOSThread()
	camera.cameraPos = glm.Vec3{1, 0, 5}
	window = initGlfw()
	initOpenGL()
	gl.Enable(gl.DEPTH_TEST)
	var sunShader = initShader("./base.vs", "./sun.fs")
	var objectShader = initShader("./light.vs", "./light.fs")

	window.SetInputMode(glfw.CursorMode, glfw.CursorDisabled)
	window.SetCursorPosCallback(camera.mouseCallback)
	window.SetScrollCallback(camera.scrollCb)

	var texture Texture
	var texture2 Texture
	var texture3 Texture

	loadImage(&texture, "./3d_artists_call_it_diffuse_map.png", gl.RGBA)
	loadImage(&texture2, "./specular_map.png", gl.RGBA)
	loadImage(&texture3, "./emission_map.jpg", gl.RGBA)

	var doubleTriangleVao Vao
	var doubleTriangleVbo Vbo
	generateVBO(&doubleTriangleVbo, doubleTriangle)
	generateVAO(&doubleTriangleVao, []VertexInfo{{3, 0}, {3, 12}, {2, 24}})

	generateVBO(&cubeVbo, cube)
	generateEBO(&cubeEbo, indicesQuad)
	generateVAO(&cubeVao, []VertexInfo{{3, 0}, {3, 12}, {2, 24}})

	var quadVao Vao
	var quadVBO Vbo
	var quadEbo Ebo

	generateVBO(&quadVBO, cube24)
	generateEBO(&quadEbo, indicesCube24)
	generateVAO(&quadVao, []VertexInfo{{3, 0}, {3, 12}, {2, 24}})

	var sunVao Vao
	var sunVbo Vbo
	generateVBO(&sunVbo, cube)
	generateVAO(&sunVao, []VertexInfo{{3, 0}, {3, 12}, {2, 24}})

	var containerVao Vao
	//var containerVbo Vbo
	//generateVBO(&containerVbo, diffuseCube)
	generateVAO(&containerVao, []VertexInfo{{3, 0}, {3, 12}, {2, 24}})

	var light = glm.Vec3{1.0, 1.0, 1.0}

	for !window.ShouldClose() {
		var currFrame = glfw.GetTime()
		deltaTime = currFrame - lastFrame
		lastFrame = currFrame
		process(window)

		//Draw
		gl.ClearColor(0.0, 0.0, 0.0, 1.0)
		gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)
		gl.PolygonMode(gl.FRONT_AND_BACK, gl.FILL)
		//	gl.PolygonMode(gl.FRONT_AND_BACK, gl.LINE)

		gl.ActiveTexture(gl.TEXTURE0)
		gl.BindTexture(gl.TEXTURE_2D, texture)
		gl.ActiveTexture(gl.TEXTURE1)
		gl.BindTexture(gl.TEXTURE_2D, texture2)
		gl.ActiveTexture(gl.TEXTURE2)
		gl.BindTexture(gl.TEXTURE_2D, texture3)

		gl.UseProgram(sunShader.prog)
		sunShader.setUniformMatrix4("view", &camera.lookAtMat)
		sunShader.setUniformMatrix4("projection", &camera.perspective)
		var sunPos = glm.Vec3{0, 0, 0}
		var modelMat = createTransformation(glm.Vec3{0, 0, 0}, sunPos, glm.Vec3{0.2, 0.2, 0.2})
		sunShader.setUniformMatrix4("model", &modelMat)

		gl.BindVertexArray(sunVao)
		gl.DrawArrays(gl.TRIANGLES, 0, int32(len(cube)/3))
		var secondPos = sunPos.Add(&glm.Vec3{4, 0, 0})
		modelMat = createTransformation(glm.Vec3{0, 0, 0}, secondPos, glm.Vec3{0.2, 0.2, 0.2})

		sunShader.setUniformMatrix4("model", &modelMat)
		gl.DrawArrays(gl.TRIANGLES, 0, int32(len(cube)/3))

		gl.UseProgram(objectShader.prog)
		modelMat = glm.Translate3D(2, 0, 0)
		objectShader.setUniformMatrix4("view", &camera.lookAtMat)
		objectShader.setUniformMatrix4("projection", &camera.perspective)
		modelMat = createTransformation(glm.Vec3{float32(glfw.GetTime()), 0, 0}, glm.Vec3{2, 0, 0}, glm.Vec3{1, 1, 1})
		objectShader.setUniformMatrix4("model", &modelMat)
		objectShader.setUniformVec3("lightColour", &light)
		objectShader.setUniformVec3("lightPos", &sunPos)
		objectShader.setUniformVec3("viewPos", &camera.cameraPos)

		objectShader.setUniformVec3("dirLight.direction", &glm.Vec3{0.0, -1.0, 0.0})
		objectShader.setUniformVec3("dirLight.ambient", &glm.Vec3{0.05, 0.05, 0.05})
		objectShader.setUniformVec3("dirLight.diffuse", &glm.Vec3{0.4, 0.4, 0.4})
		objectShader.setUniformVec3("dirLight.specular", &glm.Vec3{0.5, 0.5, 0.5})

		objectShader.setUniformVec3("pointLights[0].position", &sunPos)
		objectShader.setUniformVec3("pointLights[0].ambient", &glm.Vec3{0.05, 0.05, 0.05})
		objectShader.setUniformVec3("pointLights[0].diffuse", &glm.Vec3{0.8, 0.8, 0.8})
		objectShader.setUniformVec3("pointLights[0].specular", &glm.Vec3{1.0, 1.0, 1.0})
		objectShader.setUniform1f("pointLights[0].constant", 1.0)
		objectShader.setUniform1f("pointLights[0].linear", 0.09)
		objectShader.setUniform1f("pointLights[0].quadratic", 0.032)

		/* 	objectShader.setUniformVec3("pointLights[1].position", &secondPos)
		objectShader.setUniformVec3("pointLights[1].ambient", &glm.Vec3{0.05, 0.05, 0.05})
		objectShader.setUniformVec3("pointLights[1].diffuse", &glm.Vec3{0.8, 0.8, 0.8})
		objectShader.setUniformVec3("pointLights[1].specular", &glm.Vec3{1.0, 1.0, 1.0})
		objectShader.setUniform1f("pointLights[1].constant", 1.0)
		objectShader.setUniform1f("pointLights[1].linear", 0.09)
		objectShader.setUniform1f("pointLights[1].quadratic", 0.032) */

		objectShader.setUniformVec3("material.ambient", &glm.Vec3{1, 0.5, 0.31})
		objectShader.setUniform1i("material.diffuse", 0)
		objectShader.setUniform1i("material.specular", 1)
		objectShader.setUniform1i("material.emission", 2)
		objectShader.setUniform1i("material.shininess", 32)

		objectShader.setUniformVec3("spotLight.position", &camera.cameraPos)
		println(camera.cameraFront[2])
		var test = camera.cameraFront.Add(&glm.Vec3{-0.25, 0, 0})
		objectShader.setUniformVec3("spotLight.direction", &test)
		objectShader.setUniformVec3("spotLight.ambient", &glm.Vec3{0.0, 0.0, 0.0})
		objectShader.setUniformVec3("spotLight.diffuse", &glm.Vec3{1.0, 1.0, 1.0})
		objectShader.setUniformVec3("spotLight.specular", &glm.Vec3{1.0, 1.0, 1.0})
		objectShader.setUniform1f("spotLight.constant", 1.0)
		objectShader.setUniform1f("spotLight.linear", 0.09)
		objectShader.setUniform1f("spotLight.quadratic", 0.032)
		objectShader.setUniform1f("spotLight.cutOff", float32(math.Cos(float64(glm.DegToRad(5)))))
		objectShader.setUniform1f("spotLight.outerCutOff", float32(math.Cos(float64(glm.DegToRad(7.5)))))

		gl.BindVertexArray(containerVao)
		gl.DrawArrays(gl.TRIANGLES, 0, int32(len(cube)/3))

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
	var window, _ = glfw.CreateWindow(width, height, "light - i hate packages", nil, nil)
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
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MIN_FILTER, gl.NEAREST)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MAG_FILTER, gl.NEAREST)
	gl.TexImage2D(gl.TEXTURE_2D, 0, gl.RGBA, int32(img.Bounds().Max.X), int32(img.Bounds().Max.Y), 0, format, gl.UNSIGNED_BYTE, gl.Ptr(img.Pix))
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
