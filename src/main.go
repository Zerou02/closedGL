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

type Vao = uint32
type Vbo = uint32
type Ebo = uint32
type Prog = uint32
type Texture = uint32

func main() {
	runtime.LockOSThread()
	var window = initGlfw()
	initOpenGL()

	var shader = initShader("./shader/base.vs", "./shader/base.fs")
	var projection = glm.Ortho(0, width, height, 0, -1, 1)
	var view = glm.Ident4()

	var textShader = initShader("./shader/text.vs", "./shader/text.fs")
	_ = textShader

	var pointShader = initShader("./shader/points.vs", "./shader/points.fs")
	_ = pointShader
	var textTex = loadImage("./assets/lower_letters.png", gl.RGBA)
	var ballTex = loadImage("./assets/ball.png", gl.RGBA)

	gl.PolygonMode(gl.FRONT_AND_BACK, gl.FILL)
	gl.UseProgram(shader.prog)
	shader.setUniformMatrix4("projection", &projection)
	shader.setUniformMatrix4("view", &view)

	var delta = 0.0
	_ = delta
	var lastFrame = 0.0
	var a = newCharacter(CharacterInfo{tex: textTex}, &textShader, 0, 0, width, height, glm.Vec4{1, 0, 0, 1}, &projection)
	var ball = newSprite2D(&shader, ballTex, glm.Vec2{100, 100}, glm.Vec2{40, 40}, glm.Vec4{1, 1, 1, 1})
	/* var p2 = newPoint(&pointShader, glm.Vec2{200, 200}, glm.Vec3{0, 0, 1}, &projection)
	var p = newPoint(&pointShader, glm.Vec2{250, 250}, glm.Vec3{1, 0, 0}, &projection)
	var p3 = newPoint(&pointShader, glm.Vec2{200, 250}, glm.Vec3{0, 1, 0}, &projection) */
	var rect = newRect(&pointShader, &projection, glm.Vec4{300, 300, 300, 300}, glm.Vec3{0, 0, 1})

	/* var line = newLine(&pointShader, &projection)
	line.addPoint(p)
	line.addPoint(p2)
	line.addPoint(p3)

	var start = glm.Vec2{20, 400}
	var end = glm.Vec2{500, 100}
	var control = glm.Vec2{400, 500}
	var points = []*Point{}
	var line2 = newLine(&pointShader, &projection)
	var res float32 = 30
	for i := 0; i < int(res); i++ {
		var t = float32(i+1) / res
		var pos = bezierLerp(start, control, end, t)
		fmt.Printf("t:%f, x:%f, y:%f\n", t, pos[0], pos[1])
		var p = newPoint(&pointShader, pos, glm.Vec3{t, t, t}, &projection)
		points = append(points, &p)
		line2.addPoint(p)
	} */

	gl.PolygonMode(gl.FRONT_AND_BACK, gl.FILL)
	for !window.ShouldClose() {

		var currFrame = glfw.GetTime()
		delta = currFrame - lastFrame
		lastFrame = currFrame

		gl.ClearColor(0.0, 0.0, 0.0, 1.0)
		gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)
		a.draw()
		ball.draw()
		line.draw()
		rect.draw()
		/* for _, x := range points {
			x.draw()
		} */
		line2.draw()
		process(window)
		glfw.PollEvents()
		window.SwapBuffers()
	}
	glfw.Terminate()
}

func aabbAabbCol(b1, b2 glm.Vec4) bool {
	var colX = b1.X()+b1.Z() >= b2[0] && b2[0]+b2[2] >= b1[0]
	var colY = b1[1]+b1[3] >= b2[1] && b2[1]+b2[3] >= b1[1]
	return colX && colY
}
func process(window *glfw.Window) {
	if window.GetKey(glfw.KeyEscape) == glfw.Press {
		window.SetShouldClose(true)
	}
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

func aabbCircleCol(circle glm.Vec3, aabb glm.Vec4) (bool, Direction, glm.Vec2) {
	var centre = glm.Vec2{circle[0] + circle[2], circle[1] + circle[2]}
	var aabbHalf = glm.Vec2{aabb[2] / 2, aabb[3] / 2}
	var aabbCentre = glm.Vec2{aabb[0] + aabb[2]/2, aabb[1] + aabb[3]/2}
	var diff = centre.Sub(&aabbCentre)
	var clamped = glm.Vec2{glm.Clamp(diff[0], -aabbHalf[0], aabbHalf[0]), glm.Clamp(diff[1], -aabbHalf[1], aabbHalf[1])}
	var closest = aabbCentre.Add(&clamped)
	diff = closest.Sub(&centre)
	if diff.Len() < circle[2] {
		return true, vectorDirection(diff), diff
	} else {
		return false, UP, glm.Vec2{0, 0}
	}

}
func clamp(val, min, max float32) float32 {
	return float32(math.Max(float64(min), math.Min(float64(max), float64(val))))
}

func ssVectorOrigionCol(ssVel, ssWall glm.Vec2) glm.Vec2 {

	var esVel = ssVel.ComponentProduct(&glm.Vec2{1, -1})
	var angle = glm.RadToDeg(float32(math.Acos(float64(ssWall.Dot(&esVel) / (esVel.Len() * ssWall.Len())))))
	var rotangle = 2 * angle
	if angle == 0 {
		rotangle = 180
	}
	var rotMat = glm.Rotate2D(glm.DegToRad(360 - rotangle))
	var newAngle = rotMat.Mul2x1(&esVel)
	newAngle.Normalize()
	return newAngle.ComponentProduct(&glm.Vec2{1, -1})
}
