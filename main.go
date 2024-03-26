package main

import (
	"image"
	_ "image/png"
	"os"
	"runtime"

	"github.com/EngoEngine/glm"
	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/go-gl/glfw/v3.2/glfw"
	"neilpa.me/go-stbi"
)

var vertices = []float32{
	-0.5, -0.5, 0.0,
	0.5, -0.5, 0.0,
	0.0, 0.5, 0.0,
}

var vertices2 = []float32{
	//pos ;; col;;tex
	-1.0, -1.0, 0.0, 1.0, 0.0, 0.0, 0.0, 0.0,
	-0.5, 0.0, 0.0, 0.0, 1.0, 0.0, 0.5, 0.5,
	-1.0, 1.0, 0.0, 0.0, 0.0, 1.0, 0.0, 1.0,
	//pos ;; col;;tex
	1.0, -1.0, 0.0, 1.0, 0.0, 0.0, 0.0, 0.0,
	0.5, 0.0, 0.0, 0.0, 1.0, 0.0, 0.5, 0.5,
	1.0, 1.0, 0.0, 0.0, 0.0, 1.0, 0.0, 1.0,
}

/* var vertices2 = []float32{
	//col ;; pso
	1.0, 0.0, 0.0, -1.0, -1.0, 0.0,
	0.0, 1.0, 0.0, -0.5, 0.0, 0.0,
	0.0, 0.0, 1.0, -1.0, 1.0, 0.0,
} */

var vertices3 = []float32{
	-0.5, -0.5, -0.5, 0.0, 0.0, 0.0, 0.0, 0.0,
	0.5, -0.5, -0.5, 1.0, 0.0, 0.0, 0.0, 0.0,
	0.5, 0.5, -0.5, 1.0, 1.0, 0.0, 0.0, 0.0,
	0.5, 0.5, -0.5, 1.0, 1.0, 0.0, 0.0, 0.0,
	-0.5, 0.5, -0.5, 0.0, 1.0, 0.0, 0.0, 0.0,
	-0.5, -0.5, -0.5, 0.0, 0.0, 0.0, 0.0, 0.0,

	-0.5, -0.5, 0.5, 0.0, 0.0, 0.0, 0.0, 0.0,
	0.5, -0.5, 0.5, 1.0, 0.0, 0.0, 0.0, 0.0,
	0.5, 0.5, 0.5, 1.0, 1.0, 0.0, 0.0, 0.0,
	0.5, 0.5, 0.5, 1.0, 1.0, 0.0, 0.0, 0.0,
	-0.5, 0.5, 0.5, 0.0, 1.0, 0.0, 0.0, 0.0,
	-0.5, -0.5, 0.5, 0.0, 0.0, 0.0, 0.0, 0.0,

	-0.5, 0.5, 0.5, 1.0, 0.0, 0.0, 0.0, 0.0,
	-0.5, 0.5, -0.5, 1.0, 1.0, 0.0, 0.0, 0.0,
	-0.5, -0.5, -0.5, 0.0, 1.0, 0.0, 0.0, 0.0,
	-0.5, -0.5, -0.5, 0.0, 1.0, 0.0, 0.0, 0.0,
	-0.5, -0.5, 0.5, 0.0, 0.0, 0.0, 0.0, 0.0,
	-0.5, 0.5, 0.5, 1.0, 0.0, 0.0, 0.0, 0.0,

	0.5, 0.5, 0.5, 1.0, 0.0, 0.0, 0.0, 0.0,
	0.5, 0.5, -0.5, 1.0, 1.0, 0.0, 0.0, 0.0,
	0.5, -0.5, -0.5, 0.0, 1.0, 0.0, 0.0, 0.0,
	0.5, -0.5, -0.5, 0.0, 1.0, 0.0, 0.0, 0.0,
	0.5, -0.5, 0.5, 0.0, 0.0, 0.0, 0.0, 0.0,
	0.5, 0.5, 0.5, 1.0, 0.0, 0.0, 0.0, 0.0,

	-0.5, -0.5, -0.5, 0.0, 1.0, 0.0, 0.0, 0.0,
	0.5, -0.5, -0.5, 1.0, 1.0, 0.0, 0.0, 0.0,
	0.5, -0.5, 0.5, 1.0, 0.0, 0.0, 0.0, 0.0,
	0.5, -0.5, 0.5, 1.0, 0.0, 0.0, 0.0, 0.0,
	-0.5, -0.5, 0.5, 0.0, 0.0, 0.0, 0.0, 0.0,
	-0.5, -0.5, -0.5, 0.0, 1.0, 0.0, 0.0, 0.0,

	-0.5, 0.5, -0.5, 0.0, 1.0, 0.0, 0.0, 0.0,
	0.5, 0.5, -0.5, 1.0, 1.0, 0.0, 0.0, 0.0,
	0.5, 0.5, 0.5, 1.0, 0.0, 0.0, 0.0, 0.0,
	0.5, 0.5, 0.5, 1.0, 0.0, 0.0, 0.0, 0.0,
	-0.5, 0.5, 0.5, 0.0, 0.0, 0.0, 0.0, 0.0,
	-0.5, 0.5, -0.5, 0.0, 1.0, 0.0, 0.0, 0.0,
}

/*
var vertices3 = []float32{
	//pos;;col;;tex
	0.5, 0.5, 0.0, 1.0, 0.0, 0.0, 1.0, 1.0, // top right
	0.5, -0.5, 0.0, 0.0, 1.0, 0.0, 1.0, 0.0, // bottom right
	-0.5, -0.5, 0.0, 0.0, 0.0, 1.0, 0.0, 0.0, // bottom left
	-0.5, 0.5, 0.0, 1.0, 1.0, 0.0, 0.0, 1.0, // top left
} */

var indices = []uint32{
	0, 1, 3, 1, 2, 3,
}

const width = 800
const height = 600

var vertexShaderSource = `
    #version 410 core
	layout (location = 0) in vec3 aPos;
	layout (location = 1) in vec3 aColor;
	layout (location = 2) in vec2 aTexCoord;
	uniform vec2 offset;
	uniform mat4 transform;
	uniform mat4 model;
	uniform mat4 view;
	uniform mat4 projection;

	out vec3 color;
	out vec2 texCoord;

    void main() {
		gl_Position = projection * view * model * vec4(aPos,1.0f);
        //gl_Position = vec4(aPos, 1.0f) + vec4(offset,0.0f,0.0f);
		color = aColor;
		texCoord = aTexCoord;
	}
` + "\x00"

var fragmentShaderSource = `
    #version 410 core
    out vec4 fragColour;
	in vec3 color;
	in vec2 texCoord;

	uniform sampler2D tex;
	uniform sampler2D tex2;	
    void main() {
		fragColour = mix(texture(tex,texCoord),texture(tex2,texCoord),0.2f) * vec4(color,1.0f);
		//fragColour = mix(texture(tex,texCoord),texture(tex2,texCoord),0.2f);
		//fragColour = vec4(color,1.0f);

    }
` + "\x00"

var cubePos = []glm.Vec3{
	{1, 1, 1},
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
type Shader = uint32
type Prog = uint32
type Texture = uint32

var vbo Vbo
var vertexShader Shader
var fragmentShader Shader
var shaderProg Prog
var vao Vao
var ebo Ebo
var texture Texture
var texture2 Texture

var cameraPos = glm.Vec3{0, 0, 3}
var cameraUp = glm.Vec3{0, 1, 0}
var cameraFront = glm.Vec3{0, 0, -1}

var lastFrame float64 = 0
var deltaTime float64 = 0

var yaw float32 = -90
var pitch float32 = 0
var cameraDirection = glm.Vec3{0, 0, 0}

var lastPosX float64 = 400
var lastPosY float64 = 300
var fov float64 = 45
var firstMouse = true

func mouseCallback(w *glfw.Window, xpos float64, ypos float64) {
	if firstMouse {
		lastPosX = xpos
		lastPosY = ypos
		firstMouse = false
	}
	var offsetX = xpos - lastPosX
	var offsetY = lastPosY - ypos
	lastPosX = xpos
	lastPosY = ypos
	var sensitivity = 0.1
	offsetX *= sensitivity
	offsetY *= sensitivity

	yaw += float32(offsetX)
	pitch += float32(offsetY)
	if pitch >= 89 {
		pitch = 89
	}
	if pitch <= -89 {
		pitch = -89
	}
	println("A")
}

func scrollCb(w *glfw.Window, xOffset float64, yOffset float64) {
	println(fov)
	fov -= yOffset
	if fov < 1 {
		fov = 1
	} else if fov > 45 {
		fov = 45
	}
}
func main() {
	runtime.LockOSThread()

	window = initGlfw()
	initOpenGL()
	window.SetInputMode(glfw.CursorMode, glfw.CursorDisabled)
	gl.Enable(gl.DEPTH_TEST)
	window.SetCursorPosCallback(mouseCallback)
	window.SetScrollCallback(scrollCb)

	vertexShader = compileShader(vertexShaderSource, true)
	fragmentShader = compileShader(fragmentShaderSource, false)
	shaderProg = createShaderProgram(vertexShader, fragmentShader)
	var succ int32 = 0
	gl.GetProgramiv(shaderProg, gl.LINK_STATUS, &succ)
	if succ == 0 {
		println("ERROR")
	} else {
		println("Shader successfully linked")
	}
	//var vertexColorLocation = gl.GetUniformLocation(shaderProg, gl.Str("color"+"\x00"))
	var offsetLoc = gl.GetUniformLocation(shaderProg, gl.Str("offset"+"\x00"))
	//var transformLoc = gl.GetUniformLocation(shaderProg, gl.Str("transform"+"\x00"))
	var modelLoc = gl.GetUniformLocation(shaderProg, gl.Str("model"+"\x00"))
	_ = modelLoc
	var viewLoc = gl.GetUniformLocation(shaderProg, gl.Str("view"+"\x00"))
	var projectionLoc = gl.GetUniformLocation(shaderProg, gl.Str("projection"+"\x00"))

	//gl.GenerateMipmap())
	var img, _ = stbi.Load("./container.jpg")
	gl.GenTextures(1, &texture)

	gl.BindTexture(gl.TEXTURE_2D, texture)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_S, gl.REPEAT)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_T, gl.REPEAT)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MIN_FILTER, gl.LINEAR)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MAG_FILTER, gl.LINEAR)
	gl.TexImage2D(gl.TEXTURE_2D, 0, gl.RGB, 512, 512, 0, gl.RGBA, gl.UNSIGNED_BYTE, gl.Ptr(img.Pix))
	gl.GenerateMipmap(gl.TEXTURE_2D)

	gl.GenTextures(1, &texture2)
	gl.BindTexture(gl.TEXTURE_2D, texture2)
	//	var face, faceTexel = decodePNG("./awesomeface.png")

	var img2, _ = stbi.Load("./awesomeface.png")
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_S, gl.REPEAT)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_T, gl.REPEAT)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MIN_FILTER, gl.LINEAR)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MAG_FILTER, gl.LINEAR)
	gl.TexImage2D(gl.TEXTURE_2D, 0, gl.RGBA, 512, 512, 0, gl.RGBA, gl.UNSIGNED_BYTE, gl.Ptr(img2.Pix))
	gl.GenerateMipmap(gl.TEXTURE_2D)

	gl.UseProgram(shaderProg)

	var vao2 Vao
	var vbo2 Vbo
	//amountArrays,ArrayPtr
	gl.GenVertexArrays(1, &vao2)
	gl.GenBuffers(1, &vbo2)
	gl.BindVertexArray(vao2)
	gl.BindBuffer(gl.ARRAY_BUFFER, vbo2)
	gl.BufferData(gl.ARRAY_BUFFER, 4*len(vertices2), gl.Ptr(vertices2), gl.STATIC_DRAW)
	//layout = 0: pos
	gl.VertexAttribPointer(0, 3, gl.FLOAT, false, 8*4, nil)
	gl.EnableVertexAttribArray(0)
	//layout = 1: colour;; erste Colour beginnt ab byte3
	gl.VertexAttribPointerWithOffset(1, 3, gl.FLOAT, false, 8*4, 12)
	gl.EnableVertexAttribArray(1)
	gl.VertexAttribPointerWithOffset(2, 2, gl.FLOAT, false, 8*4, 6*4)
	gl.EnableVertexAttribArray(2)

	gl.GenVertexArrays(1, &vao)
	gl.GenBuffers(1, &vbo)
	gl.GenBuffers(1, &ebo)
	gl.BindVertexArray(vao)
	gl.BindBuffer(gl.ARRAY_BUFFER, vbo)
	gl.BufferData(gl.ARRAY_BUFFER, 4*len(vertices3), gl.Ptr(vertices3), gl.STATIC_DRAW)
	gl.BindBuffer(gl.ELEMENT_ARRAY_BUFFER, ebo)
	gl.BufferData(gl.ELEMENT_ARRAY_BUFFER, 4*len(indices), gl.Ptr(indices), gl.STATIC_DRAW)
	//layout0,sizeof(vec3/pos),float,normalize,stride,offset
	gl.VertexAttribPointer(0, 3, gl.FLOAT, false, 8*4, nil)
	gl.EnableVertexAttribArray(0)
	gl.VertexAttribPointerWithOffset(1, 3, gl.FLOAT, false, 8*4, 12)
	gl.EnableVertexAttribArray(1)
	gl.VertexAttribPointerWithOffset(2, 2, gl.FLOAT, false, 8*4, 6*4)
	gl.EnableVertexAttribArray(2)
	//	gl.BindBuffer(gl.ARRAY_BUFFER, 0)

	//var nrAttri *int32
	//gl.GetIntegerv(gl.MAX_VERTEX_ATTRIBS, nrAttri)
	//println(nrAttri)
	glfw.SetTime(0)
	glfw.GetTimerFrequency()

	gl.Uniform1i(gl.GetUniformLocation(shaderProg, gl.Str("tex"+"\x00")), 0)
	gl.Uniform1i(gl.GetUniformLocation(shaderProg, gl.Str("tex2"+"\x00")), 1)

	var orthoProjection = glm.Ortho(0, 800, 0, 600, 0.1, 100)
	_ = orthoProjection
	_ = projectionLoc
	var viewMat = glm.Ident4()
	var viewTransMat = glm.Translate3D(0, 0, -3)
	viewMat = viewMat.Mul4(&viewTransMat)

	for !window.ShouldClose() {
		var currFrame = glfw.GetTime()
		deltaTime = currFrame - lastFrame
		lastFrame = currFrame
		//fmt.Printf("%f \n", glfw.GetTime())
		process(window)

		var perspectiveProjection = glm.Perspective(glm.DegToRad(float32(fov)), width/height, 0.1, 100)
		var a = cameraFront.Add(&cameraPos)
		var lookAt = glm.LookAtV(&cameraPos, &a, &(cameraUp))

		var modelMatrix = glm.Ident4()
		var rotMatModelX = glm.HomogRotate3DX(float32(glfw.GetTime()) * glm.DegToRad(50))
		var rotMatModelY = glm.HomogRotate3DY(float32(glfw.GetTime()) * glm.DegToRad(50))
		var rotMatModelZ = glm.HomogRotate3DZ(float32(glfw.GetTime()) * glm.DegToRad(50))

		//var rotMatModel = glm.HomogRotate3DX(glm.DegToRad(-55))
		var transModel = glm.Translate3D(1, 0, 0)
		var scaleMat = glm.Scale3D(0.5, 0.5, 0.5)
		_ = transModel
		_ = scaleMat
		//modelMatrix = modelMatrix.Mul4(&scaleMat)
		modelMatrix = modelMatrix.Mul4(&transModel)
		modelMatrix = modelMatrix.Mul4(&rotMatModelX)
		modelMatrix = modelMatrix.Mul4(&rotMatModelY)
		modelMatrix = modelMatrix.Mul4(&rotMatModelZ)

		//	gl.UniformMatrix4fv(transformLoc, 1, false, &destMat[0])
		gl.UniformMatrix4fv(viewLoc, 1, false, &lookAt[0])
		gl.UniformMatrix4fv(projectionLoc, 1, false, &perspectiveProjection[0])
		gl.UniformMatrix4fv(modelLoc, 1, false, &modelMatrix[0])

		gl.ClearColor(0.2, 0.3, 0.3, 1.0)
		gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)
		gl.DrawArrays(gl.TRIANGLES, 0, 36)

		/*
			var timeVal = glfw.GetTime()
			var greenVal = math.Sin(timeVal)/2.0 + 0.5
			gl.Uniform4f(vertexColorLocation, 0.0, float32(greenVal), 0.0, 1.0) */
		//gl.Uniform2f(offsetLoc, 1, -1.0)

		_ = offsetLoc
		gl.PolygonMode(gl.FRONT_AND_BACK, gl.FILL)
		gl.UseProgram(shaderProg)
		gl.BindVertexArray(vao2)
		gl.DrawArrays(gl.TRIANGLES, 0, int32(len(vertices2)/3))

		//	gl.PolygonMode(gl.FRONT_AND_BACK, gl.LINE)
		gl.ActiveTexture(gl.TEXTURE0)
		gl.BindTexture(gl.TEXTURE_2D, texture)
		gl.ActiveTexture(gl.TEXTURE1)
		gl.BindTexture(gl.TEXTURE_2D, texture2)
		gl.BindVertexArray(vao)
		for i := 0; i < 10; i++ {
			var modelMat = glm.Ident4()
			var transMat = glm.Translate3D(cubePos[i].X(), cubePos[i].Y(), cubePos[i].Z())
			var angle = 20.0 * (i + 1)
			var rotMatX = glm.HomogRotate3DX(float32(angle))
			var rotMatY = glm.HomogRotate3DY(float32(angle))
			var rotMatZ = glm.HomogRotate3DZ(float32(angle))

			modelMat = modelMat.Mul4(&transMat)
			modelMat = modelMat.Mul4(&rotMatX)
			modelMat = modelMat.Mul4(&rotMatY)
			modelMat = modelMat.Mul4(&rotMatZ)

			gl.UniformMatrix4fv(modelLoc, 1, false, &modelMat[0])
			gl.DrawArrays(gl.TRIANGLES, 0, 36)
		}
		//		gl.DrawElements(gl.TRIANGLES, 6, gl.UNSIGNED_INT, nil)

		glfw.PollEvents()
		window.SwapBuffers()
	}
	glfw.Terminate()
}

func process(window *glfw.Window) {
	if window.GetKey(glfw.KeyEscape) == glfw.Press {
		window.SetShouldClose(true)
	}
	var speed float32 = 0.2
	if window.GetKey(glfw.KeyW) == glfw.Press {
		//var a = (cameraFront.Mul(speed))
		//a = a.Mul(float32(deltaTime))
		//cameraPos = cameraPos.Add(&a)
		cameraPos = cameraFront.MulNP(speed).MulNP(float32(deltaTime))
	}
}

func createShaderProgram(shader ...uint32) uint32 {
	var shaderProg = gl.CreateProgram()
	for i := 0; i < len(shader); i++ {
		gl.AttachShader(shaderProg, shader[i])
	}
	gl.LinkProgram(shaderProg)
	for i := 0; i < len(shader); i++ {
		gl.DeleteShader(shader[i])
	}
	return shaderProg
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

func compileShader(shaderSrc string, vertex bool) uint32 {
	var stype uint32 = gl.FRAGMENT_SHADER
	if vertex {
		stype = gl.VERTEX_SHADER
	}
	var shader = gl.CreateShader(stype)
	var vertSrc, _ = gl.Strs(shaderSrc)
	gl.ShaderSource(shader, 1, vertSrc, nil)
	gl.CompileShader(shader)
	var success int32 = 0
	gl.GetShaderiv(shader, gl.COMPILE_STATUS, &success)
	if success == 0 {
		var ptr uint8
		gl.GetShaderInfoLog(shader, 512, nil, &ptr)
		println("Error at shader Compiling", ptr)
	} else {
		println("Shader successfully compiled")
	}
	return shader
}

func decodePNG(path string) (image.Image, []uint32) {
	var reader, _ = os.Open(path)
	var img, _, _ = image.Decode(reader)
	var texel = []uint32{}
	for y := 0; y < img.Bounds().Max.Y; y++ {
		for x := 0; x < img.Bounds().Max.X; x++ {
			var r, g, b, a = img.At(x, y).RGBA()
			texel = append(texel, r, g, b, a)
		}
	}
	return img, texel
}
