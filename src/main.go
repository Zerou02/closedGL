package main

import (
	mididec "closed_gl/src/MidiDec"
	closedGL "closed_gl/src/closedGL"
	"fmt"
	_ "image/png"
	"os"
	"runtime"

	"github.com/EngoEngine/glm"
	"github.com/ebitengine/oto/v3"
	"github.com/go-gl/gl/v4.1-core/gl"
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
	var otoCtx = initAudio()
	openGL.Window.SetScrollCallback(openGL.Camera.ScrollCb)
	openGL.Window.SetCursorPosCallback(openGL.Camera.MouseCallback)

	var isWireframeMode = false

	var ww = float32(800)
	var wh = float32(600)
	var sizeY float32 = 200
	_ = isWireframeMode
	glfw.SwapInterval(0)
	var piano = mididec.NewPiano(otoCtx, &openGL, ww, wh, sizeY)
	openGL.Camera.CameraPos = glm.Vec3{0, 0, 0}

	var dec = mididec.NewMidiDecoder()
	_ = dec
	dec.Deserialize("./assets/midisongs/milky_way_explore.midi")

	var texArr uint32 = 0
	gl.CreateTextures(gl.TEXTURE_2D_ARRAY, 1, &texArr)

	for !openGL.Window.ShouldClose() {
		var delta = openGL.FPSCounter.Delta
		_ = delta

		piano.Process()
		piano.ShowKey(36)
		dec.Process(float32(delta))
		var notes = dec.CurrNoteEvents
		for _, x := range notes {
			if x.EvType == mididec.NoteOn {
				var ev = x.Event.GetBytes()
				piano.PlaySound(int(ev[1]))
				piano.ShowKey(int(ev[1]))
			} else if x.EvType == mididec.NoteOff {
				var ev = x.Event.GetBytes()
				piano.DeleteKey(int(ev[1]))

			} else {
				panic("invalid event")
			}
		}

		closedGL.ClearBG()
		if openGL.KeyBoardManager.IsPressed(glfw.KeyF) {
			isWireframeMode = !isWireframeMode
			closedGL.SetWireFrameMode(isWireframeMode)
		}
		gl.Disable(gl.DEPTH_TEST)
		piano.Draw()
		gl.Enable(gl.DEPTH_TEST)

		openGL.DrawFPS(0, 0)

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
