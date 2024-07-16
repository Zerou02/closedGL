package closedGL

import (
	"strconv"

	"github.com/EngoEngine/glm"
	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/go-gl/glfw/v3.2/glfw"
	"neilpa.me/go-stbi"
)

var text Text

type Window struct {
	Window *glfw.Window
	Ww, Wh float32
}

func (this *Window) SetScrollCallback(fun glfw.ScrollCallback) {
	this.Window.SetScrollCallback(fun)
}

func (this *Window) SetCursorPosCallback(fun glfw.CursorPosCallback) {
	this.Window.SetCursorPosCallback(fun)
}

type ClosedGLContext struct {
	Window              *Window
	shaderCameraManager *ShaderCameraManager
	Camera              *Camera
	Text                *Text
	KeyBoardManager     *KeyBoardManager
	FPSCounter          *FPSCounter
	rectangleManager    *RectangleManager
	LineArr             *LineArr
	CircleManager       *CircleManager
	TriangleManager     *TriangleManager
}

func InitClosedGL(pWidth, pHeight float32) ClosedGLContext {
	var width = pWidth
	var height = pHeight
	var window = initGlfw(int(width), int(height))
	var fpsCounter = NewFPSCounter()
	var pWindow = Window{
		Window: window,
		Ww:     width,
		Wh:     height,
	}
	initOpenGL(width, height)

	var c = newCamera(width, height)
	var shaderManager = newShaderCameraManager(float32(width), float32(height), &c)
	text = NewText("default", shaderManager.Shadermap["text"], &shaderManager.projection2D)
	var key = newKeyBoardManager(window)
	var rm = newRect(shaderManager.Shadermap["rect"], &shaderManager.projection2D)
	var cm = newCircleManger(shaderManager.Shadermap["circle"], &shaderManager.projection2D)
	var lineArr = NewLineArr(shaderManager.Shadermap["points"], &shaderManager.projection2D)
	var triMan = newTriangleManager(shaderManager.Shadermap["points"], &shaderManager.projection2D)

	var con = ClosedGLContext{
		Window: &pWindow, shaderCameraManager: &shaderManager,
		Camera: &c, Text: &text, KeyBoardManager: &key,
		FPSCounter: &fpsCounter, rectangleManager: &rm,
		LineArr: &lineArr, CircleManager: &cm, TriangleManager: &triMan}
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

func initOpenGL(width, height float32) {
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

	gl.TexImage2D(gl.TEXTURE_2D, 0, gl.RGBA, int32(img.Bounds().Max.X), int32(img.Bounds().Max.Y), 0, format, gl.UNSIGNED_BYTE, gl.Ptr(img.Pix))
	return &texPtr
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
	if this.Window.Window.GetKey(glfw.KeyEscape) == glfw.Press {
		this.Window.Window.SetShouldClose(true)
	}
	this.KeyBoardManager.Process()
	this.FPSCounter.Process()
	this.Camera.Process(this.Window.Window, float32(0.16))
	glfw.PollEvents()
	this.Window.Window.SwapBuffers()
	if this.FPSCounter.Elapsed >= 0.5 {
		this.FPSCounter.CalcAverage()
		this.FPSCounter.Clear()
	}
}

func (this *ClosedGLContext) Free() {
	glfw.Terminate()
}

func (this *ClosedGLContext) DrawFPS(posX, posY int) {
	var average = this.FPSCounter.FpsAverage
	var nr = strconv.FormatInt(int64(average), 10)
	this.Text.DrawText(posX, posY, "FPS:"+nr+"!")
}

func (this *ClosedGLContext) DrawRect(dim, colour glm.Vec4) {
	this.rectangleManager.createVertices(dim, colour)
}

func (this *ClosedGLContext) DrawLine(dim1, dim2 glm.Vec2, colour1, colour2 glm.Vec4) {
	this.LineArr.addLine(dim1, dim2, colour1, colour2)
}

func (this *ClosedGLContext) DrawPath(pos []glm.Vec2, colours []glm.Vec4) {
	this.LineArr.AddPath(pos, colours)
}

// Basisfunktion für Kreise, andere rechnen in dieses Format um
func (this *ClosedGLContext) DrawCircleFaster(upperLeft glm.Vec2, colour, borderColour glm.Vec4, diameter, borderThickness float32) {
	this.CircleManager.createVertices(upperLeft, colour, borderColour, diameter, borderThickness)
}

func (this *ClosedGLContext) DrawCircle(centre glm.Vec2, colour, borderColour glm.Vec4, radius, borderThickness float32) {
	this.CircleManager.createVertices(glm.Vec2{centre[0] - radius, centre[1] - radius}, colour, borderColour, radius*2, borderThickness)
}
func (this *ClosedGLContext) DrawQuadraticBezier(p1, p2, controlPoint glm.Vec2, colour glm.Vec4) {
	this.LineArr.AddQuadraticBezier(p1, p2, controlPoint, colour)
}

func (this *ClosedGLContext) DrawQuadraticBezierLerp(p1, p2, controlPoint glm.Vec2, colour1, colour2 glm.Vec4) {
	this.LineArr.AddQuadraticBezierLerp(p1, p2, controlPoint, colour1, colour2)
}
func (this *ClosedGLContext) EndDrawing() {
	this.Text.draw()
	this.rectangleManager.Draw()
	this.CircleManager.Draw()
	this.LineArr.Draw()
	this.TriangleManager.Draw()
}

func (this *ClosedGLContext) BeginDrawing() {
	this.Text.clearBuffer()
	this.rectangleManager.beginDraw()
	this.CircleManager.beginDraw()
	this.LineArr.beginDraw()
	this.TriangleManager.beginDraw()

}

func (this *ClosedGLContext) DrawTriangle(pos [3]glm.Vec2, colour glm.Vec4) {
	this.TriangleManager.createVertices(pos, colour)
}
