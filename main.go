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
	var chunks = []ynnebcraft.Chunk{}

	for i := 0; i < 1; i++ {
		chunks = append(chunks, ynnebcraft.NewChunk(glm.Vec3{float32(i) * 32, 0, 0}, glm.Vec3{32, 32, 32}, &openGL))
	}
	openGL.Logger.Enabled = false
	for !openGL.WindowShouldClose() {
		if openGL.KeyBoardManager.IsPressed(glfw.KeyF) {
			val = !val
			closedGL.SetWireFrameMode(val)
		}
		openGL.Logger.Start("all")
		openGL.BeginDrawing()
		openGL.ClearBG(glm.Vec4{0, 0, 0, 0})
		openGL.DrawFPS(500, 0, 1)
		openGL.DrawSprite(glm.Vec4{0, 0, 20, 20}, "./assets/sprites/fence_small.png", 1)

		openGL.Logger.Start("cpu")
		for i := 0; i < len(chunks); i++ {
			chunks[i].Draw()
		}
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
