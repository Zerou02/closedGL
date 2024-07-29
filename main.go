package main

import (
	_ "image/png"

	"github.com/EngoEngine/glm"
	"github.com/Zerou02/closedGL/closedGL"
	"github.com/go-gl/glfw/v3.2/glfw"
)

func main() {
	StartClosedGL()
}

func StartClosedGL() {

	var openGL = closedGL.InitClosedGL(800, 600, "demo")
	openGL.LimitFPS(false)
	var val = true
	for !openGL.WindowShouldClose() {
		if openGL.KeyBoardManager.IsPressed(glfw.KeyF) {
			val = !val
			closedGL.SetWireFrameMode(val)
		}
		openGL.BeginDrawing()
		openGL.ClearBG(glm.Vec4{0, 0, 0, 0})
		openGL.DrawFPS(500, 0, 1)
		openGL.DrawSprite(glm.Vec4{0, 0, 20, 20}, "./assets/sprites/fence.png", 1)
		var amount = 4_000
		for i := 0; i < amount; i++ {
			openGL.DrawCube(glm.Vec3{float32(i % (amount / 10)), float32(i / (amount / 10)), 0}, glm.Vec4{1, 1, 1, 1}, "./assets/sprites/fence.png", 1)
		}

		openGL.EndDrawing()
		openGL.Process()
	}
	openGL.Free()
}
