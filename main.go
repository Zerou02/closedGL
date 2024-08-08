package main

import (
	_ "image/png"
	"os"
	"strconv"

	"github.com/EngoEngine/glm"
	"github.com/Zerou02/closedGL/closedGL"
	"github.com/Zerou02/closedGL/ynnebcraft"
	"github.com/go-gl/glfw/v3.2/glfw"
)

func convertCubeVertices() {
	var cube = closedGL.CubeVertices
	var retArr = []byte{}
	for i := 0; i < len(cube); i += 5 {
		var newEntry byte = 0
		for j := 0; j < 5; j++ {
			var val = cube[i+j]
			var bit byte = 0
			if val >= 0.5 {
				bit = 1
			}
			bit <<= (5 - 1 - j)
			newEntry |= bit
		}
		retArr = append(retArr, newEntry)
	}
	var f, _ = os.Create("./vert.txt")
	for _, x := range retArr {
		println(strconv.FormatInt(int64(x), 10))
		f.WriteString(strconv.FormatInt(int64(x), 10) + ",")
	}
}

func main() {
	StartClosedGL()
}

func decodeTest() {
	for i, x := range closedGL.CompressedCubeVertices {
		var data float32 = float32(x)
		var v float32 = float32(int(data) & 1)
		var u float32 = float32((int(data) & 2) >> 1)
		var z float32 = float32((int(data) & 4) >> 2)
		var y float32 = float32((int(data) & 8) >> 3)
		var x float32 = float32((int(data) & 16) >> 4)
		if x == 1 {
			x = 0.5
		} else {
			x = -0.5
		}
		if y == 1 {
			y = 0.5
		} else {
			y = -0.5
		}
		if z == 1 {
			z = 0.5
		} else {
			z = -0.5
		}
		var val = closedGL.CubeVertices[i*5] == x &&
			(closedGL.CubeVertices[i*5+1] == y) &&
			(closedGL.CubeVertices[i*5+2] == z) &&
			(closedGL.CubeVertices[i*5+3] == u) &&
			(closedGL.CubeVertices[i*5+4] == v)
		if !val {
			print(i, ",")
			print(x, y, z, u, v, ",")
			closedGL.PrintlnFloat(data)
		}
	}
}

func StartClosedGL() {
	//convertCubeVertices()
	var openGL = closedGL.InitClosedGL(800, 600, "demo")
	openGL.LimitFPS(false)
	glfw.GetCurrentContext().SetInputMode(glfw.CursorMode, glfw.CursorDisabled)
	openGL.Window.SetScrollCallback(openGL.Camera.ScrollCb)
	openGL.Window.SetCursorPosCallback(openGL.Camera.MouseCallback)

	var val = true
	var chunks = []ynnebcraft.Chunk{}
	for i := 0; i < 1; i++ {
		for j := 0; j < 1; j++ {
			chunks = append(chunks, ynnebcraft.NewChunk(glm.Vec3{float32(i) * 32, 0, float32(j) * 32}, glm.Vec3{32, 32, 32}, &openGL))
		}
	}
	openGL.Logger.Enabled = true
	for !openGL.WindowShouldClose() {
		if openGL.KeyBoardManager.IsPressed(glfw.KeyF) {
			val = !val
			closedGL.SetWireFrameMode(val)
		}
		openGL.Logger.Start("all")
		openGL.BeginDrawing()
		openGL.ClearBG(glm.Vec4{0, 0, 0, 0})
		openGL.DrawFPS(500, 0, 1)
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

		openGL.EndDrawing()
		openGL.Logger.End("gpu")
		openGL.Logger.End("all")
		openGL.Logger.Print()
		openGL.Process()

	}
	openGL.Free()
}
