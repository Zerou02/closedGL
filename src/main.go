package main

import (
	closedGL "closed_gl/src/closedGL"
	"fmt"
	_ "image/png"
	"os"
	"runtime"

	"github.com/EngoEngine/glm"
	"github.com/ebitengine/oto/v3"
	"github.com/go-gl/glfw/v3.2/glfw"
)

func main() {
	startClosedGL()
}

func initAudio() *oto.Context {
	var op = &oto.NewContextOptions{}
	op.SampleRate = 44100
	op.ChannelCount = 2
	op.Format = oto.FormatSignedInt16LE
	var ctx, chann, _ = oto.NewContext(op)
	<-chann
	return ctx
}
func startClosedGL() {
	runtime.LockOSThread()

	var openGL = closedGL.InitClosedGL(800, 600)
	openGL.Window.SetScrollCallback(openGL.Camera.ScrollCb)
	openGL.Window.SetCursorPosCallback(openGL.Camera.MouseCallback)

	var isWireframeMode = false

	_ = isWireframeMode
	glfw.SwapInterval(0)
	openGL.Camera.CameraPos = glm.Vec3{0, 0, 0}

	for !openGL.Window.ShouldClose() {
		var delta = openGL.FPSCounter.Delta
		_ = delta

		if openGL.KeyBoardManager.IsPressed(glfw.KeyF) {
			isWireframeMode = !isWireframeMode
			closedGL.SetWireFrameMode(isWireframeMode)
		}

		openGL.BeginDrawing()

		closedGL.ClearBG()
		openGL.DrawFPS(0, 300)

		openGL.EndDrawing()

		openGL.Process()
	}
	openGL.Free()
}

func load() {
	var file, _ = os.ReadFile("./font/" + "default" + ".iglbmf")
	var file2, _ = os.ReadFile("./font/" + "default copy" + ".iglbmf")
	file = closedGL.RleDecode(file)
	for i, x := range file {
		if x != file2[i] {
			fmt.Printf("expected %d at %d, but got: %d", file2[i], i, x)
			break
		}
	}
	fmt.Printf("finished")
}

func save() {
	var file, _ = os.Create("./font/" + "default" + ".iglbmf")
	var file2, _ = os.ReadFile("./font/" + "default copy" + ".iglbmf")
	file2 = closedGL.RleEncode(file2)
	file.Write(file2)

	file.Close()
}
