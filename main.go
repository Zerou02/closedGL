package main

import (
	"fmt"
	"runtime"
	"strings"

	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/go-gl/glfw/v3.2/glfw"
)

var vertices = []float32{
	0, 0.5, 0, // top
	-0.5, -0.5, 0, // left
	0.5, -0.5, 0, // right
}

const width = 800
const height = 600
const vertexShaderSource = `
    #version 410 core
    layout (location = 0) in vec3 aPos;
    void main() {
        gl_Position = vec4(aPos, 1.0f);
    }
` + "\x00"

const fragmentShaderSource = `
    #version 410
    out vec4 frag_colour;
    void main() {
        frag_colour = vec4(1.0f, 0.5f, 0.2f, 1.0f);
    }
` + "\x00"

var window *glfw.Window

var vbo uint32 = 0
var vertexShader uint32 = 0
var fragmentShader uint32 = 0
var shaderProg uint32 = 0
var vao uint32 = 0

func main() {
	runtime.LockOSThread()

	window = initGlfw()
	initOpenGL()
	gl.GenBuffers(1, &vbo)
	gl.BindBuffer(gl.ARRAY_BUFFER, vbo)
	gl.BufferData(gl.ARRAY_BUFFER, len(vertices), gl.Ptr(vertices), gl.STATIC_DRAW)

	vertexShader = gl.CreateShader(gl.VERTEX_SHADER)
	var src, _ = gl.Strs(vertexShaderSource)
	gl.ShaderSource(vertexShader, 1, src, nil)
	gl.CompileShader(vertexShader)

	var succ int32 = 0
	gl.GetShaderiv(vertexShader, gl.COMPILE_STATUS, &succ)
	if succ != 1 {
		println("ERRRO")
		return
	}

	fragmentShader = gl.CreateShader(gl.FRAGMENT_SHADER)
	var fragSrc, _ = gl.Strs(fragmentShaderSource)
	gl.ShaderSource(fragmentShader, 1, fragSrc, nil)
	gl.CompileShader(fragmentShader)

	shaderProg = gl.CreateProgram()
	gl.AttachShader(shaderProg, vertexShader)
	gl.AttachShader(shaderProg, fragmentShader)
	//gl.DeleteShader(vertexShader)
	//gl.DeleteShader(fragmentShader)

	gl.VertexAttribPointer(0, 3, gl.FLOAT, false, 3*4, nil)

	gl.GenVertexArrays(1, &vao)
	gl.BindBuffer(gl.ARRAY_BUFFER, vbo)
	gl.BufferData(gl.ARRAY_BUFFER, len(vertices), gl.Ptr(vertices), gl.STATIC_DRAW)
	//var mao = createMao(triangle)

	gl.LinkProgram(shaderProg)
	gl.UseProgram(shaderProg)
	gl.BindVertexArray(vao)
	gl.EnableVertexAttribArray(0)

	for !window.ShouldClose() {

		process(window)
		gl.DrawArrays(gl.TRIANGLES, 0, 3)
		//	draw(mao, window, prog)

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

// mertex Array Object
func createMao(points []float32) uint32 {
	var vbo uint32 = 0
	gl.GenBuffers(1, &vbo)
	gl.BindBuffer(gl.ARRAY_BUFFER, vbo)
	gl.BufferData(gl.ARRAY_BUFFER, 4*len(points), gl.Ptr(points), gl.STATIC_DRAW)

	var mao uint32 = 0
	gl.GenVertexArrays(1, &mao)
	gl.BindVertexArray(mao)
	gl.EnableVertexAttribArray(0)
	gl.BindBuffer(gl.ARRAY_BUFFER, vbo)
	gl.VertexAttribPointer(0, 3, gl.FLOAT, false, 0, nil)
	return mao
}

func draw(mao uint32, window *glfw.Window, prog uint32) {
	gl.ClearColor(0.2, 0.3, 0.3, 1.0)
	gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)
	gl.UseProgram(prog)
	gl.BindVertexArray(mao)
	gl.DrawArrays(gl.TRIANGLES, 0, int32(len(vertices)/3))

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

func compileShader(source string, shaderType uint32) (uint32, error) {
	shader := gl.CreateShader(shaderType)

	csources, free := gl.Strs(source)
	gl.ShaderSource(shader, 1, csources, nil)
	free()
	gl.CompileShader(shader)

	var status int32
	gl.GetShaderiv(shader, gl.COMPILE_STATUS, &status)
	if status == gl.FALSE {
		var logLength int32
		gl.GetShaderiv(shader, gl.INFO_LOG_LENGTH, &logLength)

		log := strings.Repeat("\x00", int(logLength+1))
		gl.GetShaderInfoLog(shader, logLength, nil, gl.Str(log))

		return 0, fmt.Errorf("failed to compile %v: %v", source, log)
	}

	return shader, nil
}
