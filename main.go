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
	//convertCubeVertices()
	var openGL = closedGL.InitClosedGL(800, 600, "demo")
	openGL.LimitFPS(false)
//	glfw.GetCurrentContext().SetInputMode(glfw.CursorMode, glfw.CursorDisabled)
//
//	openGL.Window.SetScrollCallback(openGL.Camera.ScrollCb)
//	openGL.Window.SetCursorPosCallback(openGL.Camera.MouseCallback)

	//var val = true
	/*

	var chunks = []ynnebcraft.Chunk{}
	var mesher = ynnebcraft.NewGreedyMesher()

	for i := 0; i < 10; i++ {
		for j := 0; j < 10; j++ {
			chunks = append(chunks, ynnebcraft.NewChunk(glm.Vec3{float32(i) * 32, 0, float32(j) * 32}, glm.Vec3{32, 32, 32}, &openGL, &mesher))
		}
	} */
	openGL.Logger.Enabled = true
	for !openGL.WindowShouldClose() {
	/* 	if openGL.KeyBoardManager.IsPressed(glfw.Key(glfw.KeyF)) {
			val = !val
			closedGL.SetWireFrameMode(val)
			} */
			openGL.Logger.Start("all")
			openGL.BeginDrawing()
			openGL.ClearBG(glm.Vec4{0, 0, 0, 0})
			openGL.DrawFPS(500, 0, 1)
			openGL.EndDrawing()
			
			/* 	
		openGL.Text.DrawText(500, 50, "x:"+strconv.FormatFloat(float64(openGL.Camera.CameraFront[0]), 'f', 2, 64), 1)
		openGL.Text.DrawText(500, 75, "y:"+strconv.FormatFloat(float64(openGL.Camera.CameraFront[1]), 'f', 2, 64), 1)
		openGL.Text.DrawText(500, 100, "z:"+strconv.FormatFloat(float64(openGL.Camera.CameraFront[2]), 'f', 2, 64), 1)

		openGL.Text.DrawText(600, 50, "x:"+strconv.FormatFloat(float64(openGL.Camera.CameraPos[0]), 'f', 2, 64), 1)
		openGL.Text.DrawText(600, 75, "y:"+strconv.FormatFloat(float64(openGL.Camera.CameraPos[1]), 'f', 2, 64), 1)
		openGL.Text.DrawText(600, 100, "z:"+strconv.FormatFloat(float64(openGL.Camera.CameraPos[2]), 'f', 2, 64), 1)
		openGL.DrawSprite(glm.Vec4{0, 0, 20, 20}, "./assets/sprites/fence_small.png", 1)

		openGL.Logger.Start("cpu")
		for i := 0; i < len(chunks); i++ {
			chunks[i].Draw()
		}
		openGL.Logger.End("cpu")
		openGL.Logger.Start("gpu")

		openGL.Logger.End("gpu")
		openGL.Logger.End("all")
		openGL.Logger.Print() */
		openGL.Process()

	}
	openGL.Free()
}
