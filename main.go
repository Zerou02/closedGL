package main

import (
	_ "image/png"

	"github.com/EngoEngine/glm"
	"github.com/Zerou02/closedGL/closedGL"
	"github.com/Zerou02/closedGL/ynnebcraft"
	"github.com/go-gl/glfw/v3.2/glfw"
)

func main() {
	StartClosedGL()
}

func StartClosedGL() {

	var openGL = closedGL.InitClosedGL(800, 600, "demo")
	openGL.LimitFPS(false)
	var val = true
	var chunk = ynnebcraft.NewChunk(glm.Vec3{0, -16, 0}, glm.Vec3{32, 32, 32}, &openGL)
	var chunk2 = ynnebcraft.NewChunk(glm.Vec3{32, 0, 0}, glm.Vec3{32, 32, 32}, &openGL)
	_, _ = chunk2, chunk
	for !openGL.WindowShouldClose() {
		if openGL.KeyBoardManager.IsPressed(glfw.KeyF) {
			val = !val
			closedGL.SetWireFrameMode(val)
		}
		openGL.Logger.Start("all")
		openGL.BeginDrawing()
		openGL.ClearBG(glm.Vec4{0, 0, 0, 0})
		openGL.DrawFPS(500, 0, 1)
		openGL.DrawSprite(glm.Vec4{0, 0, 20, 20}, "./assets/sprites/fence.png", 1)
		openGL.Logger.Start("idx")
		/* 		for i := 0; i < 10_000; i++ {
		   			var a, b, c = closedGL.IdxToPos3(i, 100, 100, 100)
		   			_, _, _ = a, b, c
		   		}
		   		openGL.Logger.End("idx")
		   		openGL.Logger.Start("cubes")
		   		for i := 0; i < 10_000; i++ {
		   			var a, b, c = closedGL.IdxToPos3(i, 100, 100, 100)
		   			openGL.DrawCube(glm.Vec3{float32(a), float32(b), float32(c)}, "./assets/sprites/fence.png", 1)
		   		}
		   		openGL.Logger.End("cubes") */
		openGL.Logger.Start("cpu")
		chunk.Draw()
		chunk2.Draw()
		openGL.Logger.End("cpu")
		openGL.Logger.Start("gpu")

		openGL.EndDrawing()
		openGL.Logger.End("gpu")
		openGL.Logger.End("all")
		openGL.Logger.Print()
		openGL.Process()
	}
	openGL.Free()
}
