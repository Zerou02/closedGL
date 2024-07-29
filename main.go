package main

import (
	_ "image/png"

	"github.com/EngoEngine/glm"
	"github.com/Zerou02/closedGL/closedGL"
)

func main() {
	StartClosedGL()
}

func StartClosedGL() {

	var openGL = closedGL.InitClosedGL(800, 600, "demo")
	openGL.LimitFPS(false)
	for !openGL.WindowShouldClose() {
		openGL.BeginDrawing()
		openGL.ClearBG(glm.Vec4{0, 0, 0, 0})
		openGL.DrawFPS(500, 0, 1)
		openGL.DrawSprite(glm.Vec4{0, 0, 200, 200}, "./assets/sprites/fence.png", 1)
		openGL.DrawSprite(glm.Vec4{1, 1, 200, 200}, "./assets/sprites/fence2.png", 1)
		openGL.DrawSprite(glm.Vec4{2, 2, 200, 200}, "./assets/sprites/fence3.png", 1)

		openGL.EndDrawing()
		openGL.Process()
	}
	openGL.Free()
}
