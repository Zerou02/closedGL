package closed_gl

import (
	"math"

	"github.com/EngoEngine/glm"
	"github.com/go-gl/glfw/v3.2/glfw"
)

type Camera struct {
	pitch, yaw, roll                            float32
	sensitivity                                 float64
	fov                                         float32
	CameraPos, cameraFront, cameraUp, cameraDir glm.Vec3
	perspective                                 glm.Mat4
	lookAtMat                                   glm.Mat4
	alreadyEntered                              bool
	lastPosX, lastPosY                          float64
	isOrtho                                     bool
	aspect                                      float32
}

func newCamera(width, height float32) Camera {
	return Camera{
		pitch: 0, yaw: 32, roll: 0,
		sensitivity: 0.1,
		fov:         45,
		CameraPos:   glm.Vec3{0, 0, 0}, cameraFront: glm.Vec3{0, 0, -1}, cameraUp: glm.Vec3{0, 1, 0},
		cameraDir: glm.Vec3{0.0, 0, 0}, perspective: glm.Ident4(), lookAtMat: glm.Ident4(), alreadyEntered: false, lastPosX: 0, lastPosY: 0, isOrtho: false, aspect: width / height,
	}
}

func (c *Camera) Process(w *glfw.Window, deltaTime float32) {
	c.cameraDir[0] = float32(math.Cos(float64((glm.DegToRad(c.yaw))))) * float32(math.Cos(float64(glm.DegToRad(c.pitch))))
	c.cameraDir[1] = float32(math.Sin(float64(glm.DegToRad(c.pitch))))
	c.cameraDir[2] = float32(math.Sin(float64((glm.DegToRad(c.yaw))))) * float32(math.Cos(float64(glm.DegToRad(c.pitch))))
	c.cameraFront = c.cameraDir.Normalized()

	var speed float32 = 0.2
	if w.GetKey(glfw.KeyW) == glfw.Press {
		c.CameraPos = c.CameraPos.AddNP(c.cameraFront.Scale(speed).Scale(float32(deltaTime)))
	}
	if w.GetKey(glfw.KeyS) == glfw.Press {
		c.CameraPos = c.CameraPos.SubNP(c.cameraFront.Scale(speed).Scale(float32(deltaTime)))
	}
	if w.GetKey(glfw.KeyA) == glfw.Press {
		var b = c.cameraFront.CrossNP(c.cameraUp).NormalizedNP().Scale(speed).Scale(float32(deltaTime))
		c.CameraPos = c.CameraPos.Sub(&b)
	}
	if w.GetKey(glfw.KeyD) == glfw.Press {
		var b = c.cameraFront.CrossNP(c.cameraUp).NormalizedNP().Scale(speed).Scale(float32(deltaTime))
		c.CameraPos = c.CameraPos.Add(&b)
	}
	if w.GetKey(glfw.KeyQ) == glfw.Press {
		c.CameraPos.AddWith(&glm.Vec3{0, speed * deltaTime, 0})
	}
	if w.GetKey(glfw.KeyE) == glfw.Press {
		c.CameraPos.AddWith(&glm.Vec3{0, -speed * deltaTime, 0})

	}
	var a = c.cameraFront.AddNP(c.CameraPos)
	c.lookAtMat = glm.LookAtV(&c.CameraPos, &a, &(c.cameraUp))

	if c.isOrtho {
		c.perspective = glm.Ortho2D(0, 800, 0, 600)
		c.lookAtMat = glm.Ident4()

	} else {
		c.perspective = glm.Perspective(glm.DegToRad(float32(c.fov)), c.aspect, 0.1, 100)
	}
}

func (c *Camera) ScrollCb(w *glfw.Window, xOffset float64, yOffset float64) {
	c.fov -= float32(yOffset)
	if c.fov < 1 {
		c.fov = 1
	} else if c.fov > 45 {
		c.fov = 45
	}
}

func (c *Camera) MouseCallback(w *glfw.Window, xpos float64, ypos float64) {
	if !c.alreadyEntered {
		c.lastPosX = xpos
		c.lastPosY = ypos
		c.alreadyEntered = true
	}
	var offsetX = xpos - c.lastPosX
	var offsetY = c.lastPosY - ypos
	c.lastPosX = xpos
	c.lastPosY = ypos
	var sensitivity = 0.1
	offsetX *= sensitivity
	offsetY *= sensitivity

	c.yaw += float32(offsetX)
	c.pitch += float32(offsetY)
	if c.pitch >= 89 {
		c.pitch = 89
	}
	if c.pitch <= -89 {
		c.pitch = -89
	}
}
