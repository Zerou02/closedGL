package closedGL

import (
	"runtime"
	"strconv"
	"unsafe"

	"github.com/EngoEngine/glm"
	"github.com/go-gl/gl/v4.6-core/gl"

	"github.com/go-gl/glfw/v3.2/glfw"
	"neilpa.me/go-stbi"
)

var text Text

type TextMesh = TriangleMesh
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

func (this *Window) SetMouseButtonCB(fun glfw.MouseButtonCallback) {
	this.Window.SetMouseButtonCallback(fun)
}

type PrimitiveMan interface {
	beginDraw()
	endDraw()
}

type depth = int
type primitiveManPtr = unsafe.Pointer

type ClosedGLContext struct {
	Window              *Window
	shaderCameraManager *ShaderCameraManager
	Camera              *Camera
	Text                *Text
	audio               Audio
	KeyBoardManager     *KeyBoardManager
	FPSCounter          *FPSCounter
	Logger              PerfLogger
	TextManager         *TextManager
	Config              map[string]string

	mouseThisFramePressed      bool
	mouseLastFramePressed      bool
	mouseRightThisFramePressed bool
	mouseRightLastFramePressed bool
	drawWireframe              bool
}

func InitClosedGL(pWidth, pHeight float32, name string) ClosedGLContext {
	runtime.LockOSThread()

	var width = pWidth
	var height = pHeight
	var window = initGlfw(int(width), int(height), name)
	var fpsCounter = NewFPSCounter()
	var pWindow = Window{
		Window: window,
		Ww:     width,
		Wh:     height,
	}
	initOpenGL(width, height)

	var c = newCamera(width, height)
	//glfw.GetCurrentContext().SetScrollCallback(c.MouseCallback)
	var shaderManager = newShaderCameraManager(float32(width), float32(height), &c)
	var config = parseConfig("./assets/config.ini")
	if config["default_font"] != "" {
		text = NewText(config["default_font"], shaderManager.Shadermap["text"], &shaderManager.projection2D)
	} else {
		text = NewText("default", shaderManager.Shadermap["text"], &shaderManager.projection2D)
	}
	var key = newKeyBoardManager(window)
	var con = ClosedGLContext{
		Window: &pWindow, shaderCameraManager: &shaderManager,
		Camera: &c, Text: &text, KeyBoardManager: key,
		FPSCounter: &fpsCounter,
		Config:     config,
		audio:      newAudio(),
		Logger:     NewLogger(),
	}
	if config["potato-friendliness"] != "" {
		con.LimitFPS(strToBool(config["potato-friendliness"]))
	}
	con.LoadFont("./assets/font/jetbrains_mono_medium.ttf")

	return con
}

func (this *ClosedGLContext) LimitFPS(val bool) {
	if val {
		glfw.SwapInterval(1)
	} else {
		glfw.SwapInterval(0)
	}
}

func initGlfw(width, height int, name string) *glfw.Window {
	glfw.Init()
	glfw.WindowHint(glfw.Resizable, glfw.False)
	glfw.WindowHint(glfw.ContextVersionMajor, 4)
	glfw.WindowHint(glfw.ContextVersionMinor, 1)
	glfw.WindowHint(glfw.OpenGLProfile, glfw.OpenGLCoreProfile)
	glfw.WindowHint(glfw.OpenGLForwardCompatible, glfw.True)
	var window, err = glfw.CreateWindow(width, height, name, nil, nil)
	if err != nil {
		panic(err.Error())
	}
	window.MakeContextCurrent()
	return window
}

func initOpenGL(width, height float32) {
	var err = gl.Init()
	if err != nil {
		panic(err.Error())
	}
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

func (this *ClosedGLContext) ClearBG(clearColour glm.Vec4) {
	gl.ClearColor(clearColour[0], clearColour[1], clearColour[2], clearColour[3])
	gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)
}

func (this *ClosedGLContext) SetWireFrameMode(val bool) {
	this.drawWireframe = val
	if val {
		gl.PolygonMode(gl.FRONT_AND_BACK, gl.LINE)
	} else {
		gl.PolygonMode(gl.FRONT_AND_BACK, gl.FILL)
	}
}

func (this *ClosedGLContext) GetWireFrameMode() bool {
	return this.drawWireframe
}

func (this *ClosedGLContext) Process() {
	if this.Window.Window.GetKey(glfw.KeyEscape) == glfw.Press {
		this.Window.Window.SetShouldClose(true)
	}
	this.KeyBoardManager.Process(this.FPSCounter.FrameCount)
	this.FPSCounter.Process()
	this.Camera.Process(this.Window.Window, float32(this.FPSCounter.Delta))
	this.audio.process()
	glfw.PollEvents()
	this.Window.Window.SwapBuffers()
	if this.FPSCounter.Elapsed >= 0.5 {
		this.FPSCounter.CalcAverage()
		this.FPSCounter.Clear()
	}

	this.mouseLastFramePressed = this.mouseThisFramePressed
	this.mouseThisFramePressed = this.IsMouseDown()
	this.mouseRightLastFramePressed = this.mouseRightThisFramePressed
	this.mouseRightThisFramePressed = this.IsMouseRightDown()
}

func (this *ClosedGLContext) Free() {
	glfw.Terminate()
}

func (this *ClosedGLContext) DrawFPS(posX, posY int, scale float32) {
	var average = this.FPSCounter.FpsAverage
	var nr = strconv.FormatInt(int64(average), 10)
	this.Text.DrawText(posX, posY, "FPS:"+nr+"!", scale)
}

func (this *ClosedGLContext) EndDrawing() {
	this.Text.draw()
	this.Process()
}

func (this *ClosedGLContext) BeginDrawing() {
	this.Text.clearBuffer()
}

func (this *ClosedGLContext) PlaySound(name string) {
	this.audio.playSound(name)
}

func (this *ClosedGLContext) PlayMusic(name string, volume float64) {
	this.audio.streamMusic(name, volume, true)
}

func (this *ClosedGLContext) EndMusic(name string) {
	this.audio.closeMusic(name)
}

func (this *ClosedGLContext) WindowShouldClose() bool {
	return this.Window.Window.ShouldClose()
}

func (this *ClosedGLContext) IsMouseInRect(rect glm.Vec4) bool {
	return IsPointInRect(this.GetMousePos(), rect)
}

func (this *ClosedGLContext) IsMouseInCircle(centre glm.Vec2, r float32) bool {
	return IsPointInCircle(this.GetMousePos(), centre, r)
}

func (this *ClosedGLContext) IsMouseDown() bool {
	return this.Window.Window.GetMouseButton(glfw.MouseButtonLeft) == glfw.Press
}

func (this *ClosedGLContext) IsMouseRightDown() bool {
	return this.Window.Window.GetMouseButton(glfw.MouseButtonRight) == glfw.Press
}

func (this *ClosedGLContext) GetMousePos() glm.Vec2 {
	var x, y = this.Window.Window.GetCursorPos()
	return glm.Vec2{float32(x), float32(y)}
}

// true if and only if mouse first pressed this frame
func (this *ClosedGLContext) MouseClicked() bool {
	return this.mouseThisFramePressed && !this.mouseLastFramePressed
}

func (this *ClosedGLContext) MouseRightClicked() bool {
	return this.mouseRightThisFramePressed && !this.mouseRightLastFramePressed
}

func (this *ClosedGLContext) GetThisFramePressedKeys() []glfw.Key {
	return this.KeyBoardManager.thisFramePressed
}

func (this *ClosedGLContext) GetThisFramePressedKey() *glfw.Key {

	var keys = this.KeyBoardManager.thisFramePressed
	if len(keys) == 0 {
		return nil
	} else {
		var key = keys[0]
		return &key
	}
}

func (this *ClosedGLContext) IsKeyDown(key glfw.Key) bool {
	return this.KeyBoardManager.IsDown(key)
}

func (this *ClosedGLContext) IsKeyPressed(key glfw.Key) bool {
	return Contains(&this.KeyBoardManager.thisFramePressed, key)
}

func (this *ClosedGLContext) CreateRectMesh() RectangleMesh {
	return newRectMesh(this.shaderCameraManager.Shadermap["rect"], this.shaderCameraManager.projection2D)
}

func (this *ClosedGLContext) CreatePixelMesh() PixelMesh {
	return newPixelMesh(this.shaderCameraManager.Shadermap["pixel"], this.shaderCameraManager.projection2D)
}

func (this *ClosedGLContext) CreateLineMesh() LineMesh {
	return newLineMesh(this.shaderCameraManager.Shadermap["pixel"], this.shaderCameraManager.projection2D)
}

func (this *ClosedGLContext) CreateTextMesh() TextMesh {
	return newTriMesh(this.shaderCameraManager.Shadermap["ttf"], this.shaderCameraManager.projection2D)
}

func (this *ClosedGLContext) NewCam2D() Camera2D {
	return newCamera2D(this.Window.Ww, this.Window.Wh, this)
}

func (this *ClosedGLContext) CartesianToSS(vec glm.Vec2) glm.Vec2 {
	return CartesianToSS(vec, this.Window.Wh)
}

func (this *ClosedGLContext) SSToCartesian(vec glm.Vec2) glm.Vec2 {
	return SsToCartesian(vec, this.Window.Wh)
}

func (this *ClosedGLContext) LoadFont(path string) {
	var t = newTextManager(path, this)
	this.TextManager = &t
}

func (this *ClosedGLContext) DrawText(x, y float32, size float32, text string, textMesh *TextMesh) {
	this.TextManager.drawText(x, y, size, text, textMesh)
}
