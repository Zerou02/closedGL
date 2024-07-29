package main

import (
	_ "image/png"

	"github.com/EngoEngine/glm"
	"github.com/Zerou02/closedGL/closedGL"
	"github.com/go-gl/gl/v4.1-core/gl"
)

func main() {
	StartClosedGL()
}

func StartClosedGL() {

	var openGL = closedGL.InitClosedGL(800, 600, "demo")
	openGL.LimitFPS(false)
	var tex = closedGL.LoadImage("./assets/sprites/fence.png", gl.RGBA)
	var cube = openGL.CreateCube(tex)
	for !openGL.WindowShouldClose() {
		openGL.BeginDrawing()
		openGL.ClearBG(glm.Vec4{0, 0, 0, 0})
		openGL.DrawFPS(500, 0, 1)
		//openGL.DrawSprite(glm.Vec4{0, 0, 200, 200}, "./assets/sprites/fence.png", 1)
		cube.Draw()
		/* openGL.DrawSprite(glm.Vec4{1, 1, 200, 200}, "./assets/sprites/fence2.png", 1)
		openGL.DrawSprite(glm.Vec4{2, 2, 200, 200}, "./assets/sprites/fence3.png", 1) */

		openGL.EndDrawing()
		openGL.Process()
	}
	openGL.Free()
}
