package closedGL

import (
	"sort"
	"strconv"
	"unsafe"

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
	KeyBoardManager     *KeyBoardManager
	FPSCounter          *FPSCounter
	rectangleManager    *RectangleManager
	LineArr             *LineArr
	CircleManager       *CircleManager
	TriangleManager     *TriangleManager
	amountPrimitiveMans int
	primitiveManMap     map[depth]*[]unsafe.Pointer
	indexArr            []int
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
		LineArr: &lineArr, CircleManager: &cm, TriangleManager: &triMan,
		primitiveManMap: map[depth]*[]unsafe.Pointer{}, amountPrimitiveMans: 4, indexArr: []int{}}

	return con
}

func (this *ClosedGLContext) LimitFPS(val bool) {
	if val {
		glfw.SwapInterval(1)
	} else {
		glfw.SwapInterval(0)
	}
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

func (this *ClosedGLContext) ClearBG() {
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

func (this *ClosedGLContext) DrawFPS(posX, posY int, scale float32) {
	var average = this.FPSCounter.FpsAverage
	var nr = strconv.FormatInt(int64(average), 10)
	this.Text.DrawText(posX, posY, "FPS:"+nr+"!", scale)
}

func (this *ClosedGLContext) initEmptyMapAtDepth(depth int) {
	var newArr = []unsafe.Pointer{}
	for i := 0; i < this.amountPrimitiveMans; i++ {
		newArr = append(newArr, nil)
	}
	var rm = this.createRectMan()
	var cm = this.createCircleMan()
	var lm = this.createLineMan()
	var tm = this.createTriMan()

	this.primitiveManMap[depth] = &newArr

	this.setMapEntry(depth, 0, unsafe.Pointer(&rm))
	this.setMapEntry(depth, 1, unsafe.Pointer(&cm))
	this.setMapEntry(depth, 2, unsafe.Pointer(&lm))
	this.setMapEntry(depth, 3, unsafe.Pointer(&tm))

	this.indexArr = append(this.indexArr, depth)
	sort.Ints(this.indexArr)
}

func (this *ClosedGLContext) createRectMan() RectangleManager {
	return newRect(this.shaderCameraManager.Shadermap["rect"], &this.shaderCameraManager.projection2D)
}

func (this *ClosedGLContext) createCircleMan() CircleManager {
	return newCircleManger(this.shaderCameraManager.Shadermap["circle"], &this.shaderCameraManager.projection2D)
}

func (this *ClosedGLContext) createLineMan() LineArr {
	return NewLineArr(this.shaderCameraManager.Shadermap["points"], &this.shaderCameraManager.projection2D)
}

func (this *ClosedGLContext) createTriMan() TriangleManager {
	return newTriangleManager(this.shaderCameraManager.Shadermap["points"], &this.shaderCameraManager.projection2D)

}
func (this *ClosedGLContext) getMapEntry(depth int, idx int) unsafe.Pointer {
	return (*this.primitiveManMap[depth])[idx]
}

func (this *ClosedGLContext) setMapEntry(depth int, idx int, ptr unsafe.Pointer) {
	(*this.primitiveManMap[depth])[idx] = ptr
}
func (this *ClosedGLContext) DrawRect(dim, colour glm.Vec4, depth int) {
	if this.primitiveManMap[depth] == nil {
		this.initEmptyMapAtDepth(depth)
	}
	(*RectangleManager)(this.getMapEntry(depth, 0)).createVertices(dim, colour)
}

func (this *ClosedGLContext) DrawLine(dim1, dim2 glm.Vec2, colour1, colour2 glm.Vec4, depth int) {
	if this.primitiveManMap[depth] == nil {
		this.initEmptyMapAtDepth(depth)
	}
	(*LineArr)(this.getMapEntry(depth, 2)).addLine(dim1, dim2, colour1, colour2)
}

func (this *ClosedGLContext) DrawPath(pos []glm.Vec2, colours []glm.Vec4, depth int) {
	if this.primitiveManMap[depth] == nil {
		this.initEmptyMapAtDepth(depth)
	}
	(*LineArr)(this.getMapEntry(depth, 2)).AddPath(pos, colours)
}

// Basisfunktion für Kreise, andere rechnen in dieses Format um
func (this *ClosedGLContext) DrawCircleFaster(upperLeft glm.Vec2, colour, borderColour glm.Vec4, diameter, borderThickness float32, depth int) {
	if this.primitiveManMap[depth] == nil {
		this.initEmptyMapAtDepth(depth)
	}
	(*CircleManager)(this.getMapEntry(depth, 1)).createVertices(upperLeft, colour, borderColour, diameter, borderThickness)
}

func (this *ClosedGLContext) DrawCircle(centre glm.Vec2, colour, borderColour glm.Vec4, radius, borderThickness float32, depth int) {
	if this.primitiveManMap[depth] == nil {
		this.initEmptyMapAtDepth(depth)
	}
	(*CircleManager)(this.getMapEntry(depth, 1)).createVertices(glm.Vec2{centre[0] - radius, centre[1] - radius}, colour, borderColour, radius*2, borderThickness)
}
func (this *ClosedGLContext) DrawQuadraticBezier(p1, p2, controlPoint glm.Vec2, colour glm.Vec4, depth int) {
	if this.primitiveManMap[depth] == nil {
		this.initEmptyMapAtDepth(depth)
	}
	(*LineArr)(this.getMapEntry(depth, 2)).AddQuadraticBezier(p1, p2, controlPoint, colour)
}

func (this *ClosedGLContext) DrawQuadraticBezierLerp(p1, p2, controlPoint glm.Vec2, colour1, colour2 glm.Vec4, depth int) {
	if this.primitiveManMap[depth] == nil {
		this.initEmptyMapAtDepth(depth)
	}
	(*LineArr)(this.getMapEntry(depth, 2)).AddQuadraticBezierLerp(p1, p2, controlPoint, colour1, colour2)
}

func (this *ClosedGLContext) EndDrawing() {

	this.Text.draw()
	for _, x := range this.indexArr {
		var v = this.primitiveManMap[x]
		for i, x := range *v {
			if i == 0 {
				(*RectangleManager)(x).draw()
			} else if i == 1 {
				(*CircleManager)(x).draw()
			} else if i == 2 {
				(*LineArr)(x).draw()
			} else if i == 3 {
				(*TriangleManager)(x).draw()
			}
		}
	}
}

func (this *ClosedGLContext) BeginDrawing() {
	this.Text.clearBuffer()
	this.rectangleManager.beginDraw()
	this.CircleManager.beginDraw()
	this.LineArr.beginDraw()
	this.TriangleManager.beginDraw()

	for _, v := range this.primitiveManMap {
		for i, x := range *v {
			if i == 0 {
				(*RectangleManager)(x).beginDraw()
			} else if i == 1 {
				(*CircleManager)(x).beginDraw()
			} else if i == 2 {
				(*LineArr)(x).beginDraw()
			} else if i == 3 {
				(*TriangleManager)(x).beginDraw()
			}
		}
	}

}

func (this *ClosedGLContext) DrawTriangle(pos [3]glm.Vec2, colour glm.Vec4, depth int) {
	if this.primitiveManMap[depth] == nil {
		this.initEmptyMapAtDepth(depth)
	}
	(*TriangleManager)(this.getMapEntry(depth, 3)).createVertices(pos, colour)
}
