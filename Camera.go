package main

import (
	"math"

	"github.com/EngoEngine/glm"
	"github.com/go-gl/glfw/v3.2/glfw"
)

type Camera struct {
	pitch, yaw, roll     float64
	sensitivity          float64
	cameraPos, cameraDir glm.Vec3
	lookAtMat            glm.Mat4
}

type Test struct {
	x, y int
}

func (t1 Test) add(t2 Test) Test {
	return Test{t1.x + t2.x, t1.y + t2.y}
}

func (c Camera) process(w *glfw.Window) {
	c.cameraDir[0] = float32(math.Cos(float64((glm.DegToRad(yaw))))) * float32(math.Cos(float64(glm.DegToRad(pitch))))
	c.cameraDir[1] = float32(math.Sin(float64(glm.DegToRad(pitch))))
	c.cameraDir[2] = float32(math.Sin(float64((glm.DegToRad(yaw))))) * float32(math.Cos(float64(glm.DegToRad(pitch))))
	cameraFront = cameraDirection.Normalized()

	var speed float32 = 2
	if window.GetKey(glfw.KeyW) == glfw.Press {
		//var a = (cameraFront.Mul(speed))
		//a = a.Mul(float32(deltaTime))
		//cameraPos = cameraPos.Add(&a)
		cameraPos = cameraFront.MulNP(speed).MulNP(float32(deltaTime))
	}
	if window.GetKey(glfw.KeyS) == glfw.Press {
		var a = cameraFront.Mul(speed)
		a = a.Mul(float32(deltaTime))
		cameraPos = cameraPos.Sub(&a)
	}
	if window.GetKey(glfw.KeyA) == glfw.Press {
		var b = cameraFront.Cross(&cameraUp)
		b.Normalize()
		var a = b.Mul(speed)
		a = a.Mul(float32(deltaTime))
		cameraPos = cameraPos.Sub(&a)
	}
	if window.GetKey(glfw.KeyD) == glfw.Press {
		var b = cameraFront.Cross(&cameraUp)
		b.Normalize()
		var a = b.Mul(speed)
		a = a.Mul(float32(deltaTime))
		cameraPos = cameraPos.Add(&a)
	}
}
