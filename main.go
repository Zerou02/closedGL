package main

import (
	_ "image/png"
	"strconv"

	"github.com/EngoEngine/glm"
	"github.com/Zerou02/closedGL/closedGL"
	turingfontparser "github.com/Zerou02/closedGL/turing_font_parser"
	"github.com/go-gl/glfw/v3.2/glfw"
)

func main() {
	StartTTF()
	//StartClosedGL()
	//StartTuwuing()
}
func StartTTF() {
	var opengl = closedGL.InitClosedGL(800, 800, "comic")
	opengl.LimitFPS(false)
	var p = turingfontparser.NewTuringFont("./assets/font/jetbrains_mono_medium.ttf", &opengl)
	var points = []turingfontparser.GlyfPoints{}
	var offset float32 = 0
	for _, x := range "A" {
		var newPoints = p.ParseGlyf(uint32(x), 1).AddXOffset(offset)
		var biggestX float32 = 0
		for i, y := range newPoints {
			closedGL.PrintlnVec2(y.Pos)
			println(y.OnCurve)
			println("---")
			var offset = glm.Vec2{0, -50}
			newPoints[i].Pos = newPoints[i].Pos.Add(&offset)
			if y.Pos[0] > biggestX {
				biggestX = y.Pos[0]
			}
		}
		offset += (biggestX - offset)
		points = append(points, newPoints...)
	}
	println("points", points)
	var lines = opengl.CreateLineMesh()
	//var startOfLine = points[0].Pos
	/* 	for i, x := range points {
		points[i].Pos[1] -= 20
		println("nr", i, "----")
		closedGL.PrintlnVec2(x.Pos)
		println("on", x.OnCurve)
		println("end", x.EndPoint)
	} */

	println("II")
	for i := 0; i < len(points)-6; i += 2 {
		closedGL.PrintlnVec2(points[i].Pos)
		println("--")
		lines.AddQuadraticBezier(points[i].Pos, points[i+2].Pos, points[i+1].Pos, glm.Vec4{1, 1, 1, 1})
		if points[i].EndPoint {
			i += 2
			//	i--
			//i += 2
		}
	}
	lines.Copy()
	for !opengl.WindowShouldClose() {
		opengl.BeginDrawing()
		opengl.ClearBG(glm.Vec4{0, 0, 0, 0})
		lines.Draw()
		opengl.DrawFPS(600, 0, 1)
		opengl.EndDrawing()
	}
	opengl.Free()
}

/* func StartTuwuing() {
	var openGL = closedGL.InitClosedGL(1400, 800, "demo")

	openGL.LimitFPS(true)
	var complete = tuwuing_complete.NewTuwuingComplete(&openGL)

	for !openGL.WindowShouldClose() {
		complete.Process()

		openGL.ClearBG(glm.Vec4{0, 0, 0, 0})
		openGL.BeginDrawing()
		complete.Draw()

		openGL.DrawFPS(500, 0, 1)

		openGL.EndDrawing()
	}
	openGL.Free()
} */

func StartClosedGL() {
	var openGL = closedGL.InitClosedGL(800, 600, "demo")
	openGL.LimitFPS(false)
	glfw.GetCurrentContext().SetInputMode(glfw.CursorMode, glfw.CursorDisabled)

	openGL.Window.SetScrollCallback(openGL.Camera.ScrollCb)
	openGL.Window.SetCursorPosCallback(openGL.Camera.MouseCallback)

	var val = true

	//	var chunks = []ynnebcraft.Chunk{}
	//	var mesher = ynnebcraft.NewGreedyMesher()

	/* 	for i := 0; i < 20; i++ {
		for j := 0; j < 10; j++ {
			chunks = append(chunks, ynnebcraft.NewChunk(glm.Vec3{float32(i) * 32, 0, float32(j) * 32}, glm.Vec3{32, 32, 32}, &openGL, &mesher))
		}
	} */
	openGL.Logger.Enabled = true

	var pxTest = openGL.CreatePixelMesh()
	var lines = openGL.CreateLineMesh()

	for i := 0; i < 800; i++ {
		for j := 0; j < 600; j++ {
			pxTest.AddPixel(glm.Vec2{float32(i), float32(j)}, glm.Vec4{float32(j) / 600, float32(i) / 800, 1, 1})
		}
	}
	pxTest.Copy()
	lines.AddPath([]glm.Vec2{{0, 0}, {200, 200}}, []glm.Vec4{{1, 1, 1, 1}, glm.Vec4{1, 1, 1, 1}})
	lines.AddQuadraticBezier(glm.Vec2{0, 0}, glm.Vec2{200, 200}, glm.Vec2{0, 200}, glm.Vec4{0, 0.54, 0.57, 1})
	lines.Copy()

	var cam = openGL.NewCam2D()
	openGL.LimitFPS(false)
	for !openGL.WindowShouldClose() {
		cam.Process(float32(openGL.FPSCounter.Delta))
		pxTest.View = cam.ViewMat
		if openGL.IsKeyPressed(glfw.Key(glfw.KeyF)) {
			val = !val
			closedGL.SetWireFrameMode(val)
		}
		openGL.BeginDrawing()
		openGL.ClearBG(glm.Vec4{0, 0, 0, 0})
		openGL.DrawFPS(500, 0, 1)

		//openGL.DrawSprite(glm.Vec4{0, 0, 20, 20}, "./assets/sprites/fence_small.png", glm.Vec4{0, 0, 1, 1}, glm.Vec2{32, 1024}, 1)

		openGL.Text.DrawText(500, 50, "x:"+strconv.FormatFloat(float64(openGL.Camera.CameraFront[0]), 'f', 2, 64), 1)
		openGL.Text.DrawText(500, 75, "y:"+strconv.FormatFloat(float64(openGL.Camera.CameraFront[1]), 'f', 2, 64), 1)
		openGL.Text.DrawText(500, 100, "z:"+strconv.FormatFloat(float64(openGL.Camera.CameraFront[2]), 'f', 2, 64), 1)

		openGL.Text.DrawText(600, 50, "x:"+strconv.FormatFloat(float64(openGL.Camera.CameraPos[0]), 'f', 2, 64), 1)
		openGL.Text.DrawText(600, 75, "y:"+strconv.FormatFloat(float64(openGL.Camera.CameraPos[1]), 'f', 2, 64), 1)
		openGL.Text.DrawText(600, 100, "z:"+strconv.FormatFloat(float64(openGL.Camera.CameraPos[2]), 'f', 2, 64), 1)

		/* for i := 0; i < len(chunks); i++ {
			chunks[i].Draw()
		} */
		//	lines.Draw()
		pxTest.Draw()
		//openGL.DrawRect(glm.Vec4{0, 0, 100, 100}, glm.Vec4{0, 1, 0, 1}, 1)
		openGL.EndDrawing()
	}
	openGL.Free()
}
