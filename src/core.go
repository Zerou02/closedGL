package main

import (
	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/go-gl/glfw/v3.2/glfw"
	"neilpa.me/go-stbi"
)

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
