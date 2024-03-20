package main

import (
	"fmt"
	"runtime"
	"strings"

	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/go-gl/glfw/v3.2/glfw"
)

var vertices = []float32{
	-0.5, -0.5, 0.0, // left
	0.5, -0.5, 0.0, // right
	0.0, 0.5, 0.0, // top
}

const width = 800
const height = 600

var vertexShaderSource = `
    #version 410
    in vec3 vp;
    void main() {
        gl_Position = vec4(vp, 1.0);
    }
` + "\x00"

var fragmentShaderSource = `
    #version 410
    out vec4 frag_colour;
    void main() {
        frag_colour = vec4(1, 1, 1, 1);
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

	//vertexShader = gl.CreateShader(gl.VERTEX_SHADER)
	//var vertSrc, _ = gl.Strs(vertexShaderSource)
	//gl.ShaderSource(vertexShader, 1, vertSrc, nil)
	//gl.CompileShader(vertexShader)
	//var success int32 = 0
	//gl.GetShaderiv(vertexShader, gl.COMPILE_STATUS, &success)
	//println("success:Vertex", success)
	//
	//fragmentShader = gl.CreateShader(gl.FRAGMENT_SHADER)
	//var fragSrc, _ = gl.Strs(fragmentShaderSource)
	//gl.ShaderSource(fragmentShader, 1, fragSrc, nil)
	//gl.CompileShader(fragmentShader)
	//var success_frag int32 = 0
	//gl.GetShaderiv(vertexShader, gl.COMPILE_STATUS, &success_frag)
	//println("success:Frag:", success_frag)

	var err error
	vertexShader, err = compileShader(vertexShaderSource, gl.VERTEX_SHADER)
	fragmentShader, _ = compileShader(fragmentShaderSource, gl.FRAGMENT_SHADER)
	if err != nil {
		print("ERRR!!!")
	}
	shaderProg = gl.CreateProgram()
	gl.AttachShader(shaderProg, vertexShader)
	gl.AttachShader(shaderProg, fragmentShader)
	gl.LinkProgram(shaderProg)
	gl.DeleteShader(vertexShader)
	gl.DeleteShader(fragmentShader)

	gl.GenVertexArrays(1, &vao)
	gl.GenBuffers(1, &vbo)
	gl.BindVertexArray(vao)
	gl.BindBuffer(gl.ARRAY_BUFFER, vbo)
	gl.BufferData(gl.ARRAY_BUFFER, len(vertices), gl.Ptr(vertices), gl.STATIC_DRAW)
	gl.VertexAttribPointer(0, 3, gl.FLOAT, false, 12, nil)
	gl.EnableVertexAttribArray(0)

	gl.BindBuffer(gl.ARRAY_BUFFER, 0)
	gl.BindVertexArray(0)

	//gl.PolygonMode(gl.FRONT_AND_BACK,gl.LINE)

	for !window.ShouldClose() {

		process(window)
		gl.ClearColor(0.2, 0.3, 0.3, 1.0)
		gl.Clear(gl.COLOR_BUFFER_BIT)

		gl.UseProgram(shaderProg)
		gl.BindVertexArray(vao)
		gl.DrawArrays(gl.TRIANGLES, 0, 3)

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
