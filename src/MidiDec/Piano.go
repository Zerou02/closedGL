package mididec

import (
	"bytes"
	closed_gl "closed_gl/src/closedGL"
	"os"
	"strings"

	"github.com/EngoEngine/glm"
	"github.com/ebitengine/oto/v3"
	"github.com/go-gl/glfw/v3.2/glfw"
	"github.com/hajimehoshi/go-mp3"
)

type MidiCode = int
type Piano struct {
	Pos, Size                glm.Vec2
	AudioCtx                 *oto.Context
	ClosedGLCtx              *closed_gl.ClosedGLContext
	notesPath, halfNotesPath []string
	notesKeys, halfNotesKeys []glfw.Key
	players                  []*oto.Player
	sprites                  []*closed_gl.Sprite2D
	overlaySprites           map[glfw.Key]*closed_gl.Sprite2D
	noteSprites              map[MidiCode]*closed_gl.Sprite2D
	colours                  []glm.Vec4
	pressedKeys              []glfw.Key
	keyNames                 []string
	octaves                  float32
}

func NewPiano(otoCtx *oto.Context, closedGLCtx *closed_gl.ClosedGLContext, ww, wh, pianoHeight float32) Piano {
	var notes = []string{"c", "d", "e", "f", "g", "a", "b"}
	var halfNotes = []string{"cs", "ds", "fs", "gs", "as"}
	var keys = []glfw.Key{glfw.KeyQ, glfw.KeyW, glfw.KeyE, glfw.KeyR, glfw.KeyU, glfw.KeyI, glfw.KeyO}
	var halfKeys = []glfw.Key{glfw.Key2, glfw.Key3, glfw.Key7, glfw.Key8, glfw.Key9}
	var sprites = []*closed_gl.Sprite2D{}
	var octaves float32 = 7.0
	var pianoWidth float32 = ww / octaves
	for i := 0; i < int(octaves); i++ {
		var sprite = closedGLCtx.Factory.NewSprite2D(glm.Vec2{pianoWidth * float32(i), wh - pianoHeight}, glm.Vec2{pianoWidth, pianoHeight}, glm.Vec4{1, 1, 1, 1}, "assets/piano.png")
		sprites = append(sprites, &sprite)
	}
	var dirs, _ = os.ReadDir("./assets/keys")
	var fileNames = []string{}
	for _, x := range dirs {
		fileNames = append(fileNames, x.Name())
	}
	return Piano{
		Pos:         glm.Vec2{0, 400},
		Size:        glm.Vec2{pianoWidth, pianoHeight},
		AudioCtx:    otoCtx,
		ClosedGLCtx: closedGLCtx,
		notesPath:   notes, halfNotesPath: halfNotes,
		notesKeys: keys, halfNotesKeys: halfKeys,
		players:        []*oto.Player{},
		sprites:        sprites,
		noteSprites:    map[int]*closed_gl.Sprite2D{},
		overlaySprites: map[glfw.Key]*closed_gl.Sprite2D{},
		colours:        []glm.Vec4{{0, 0, 1, 1}, {0, 1, 0, 1}, {0, 1, 1, 1}, {0, 0.5, 1, 1}},
		keyNames:       fileNames,
		octaves:        octaves,
	}
}

func (this *Piano) createPlayer(key glfw.Key, path string) {
	if this.ClosedGLCtx.KeyBoardManager.IsPressed(key) {
		var p = this.generatePlayer(path, this.AudioCtx)
		this.players = append(this.players, p)
		p.SetVolume(0.1)
		p.Play()

	}
}

func (this *Piano) PlaySound(midiNote int) {
	var p = this.generatePlayer("./assets/keys/"+this.keyNames[midiNote-21], this.AudioCtx)
	this.players = append(this.players, p)
	p.SetVolume(0.1)
	p.Play()
}

func (this *Piano) ShowKey(midiNote int) {
	var noteNames = []string{"a", "as", "b", "c", "cs", "d", "ds", "e", "f", "fs", "g", "gs"}
	var periodLength = 12
	var based = midiNote - 21
	var octave = int(based / periodLength)
	var noteIdx = based - octave*periodLength
	var noteName = noteNames[noteIdx]
	octave--
	var isHalfNote = strings.Contains(noteName, "s")
	var pos = glm.Vec2{0, 0}
	var size = glm.Vec2{10, 10}
	var spriteName = ""
	if isHalfNote {
		var i = 0
		var halfNotes = []string{"cs", "ds", "fs", "gs", "as"}
		for j, x := range halfNotes {
			if x == noteName {
				i = j
			}
		}
		var origPixelOffsets = []float32{32, 87, 169, 222, 276}
		var origImgWidth float32 = 329
		var origImgHeight float32 = 326
		var origNoteWidth float32 = 26
		var origNoteHeight float32 = 200
		_, _, _ = origImgHeight, origNoteHeight, origNoteWidth
		var offsetPos = glm.Vec2{float32((origPixelOffsets[i] / origImgWidth) * this.Size[0]), 0}
		pos = this.Pos.Add(&offsetPos)
		size = glm.Vec2{origNoteWidth / origImgWidth * this.Size[0], origNoteHeight / origImgHeight * this.Size[1]}
		spriteName = "./assets/halfNote.png"
	} else {
		var i = 0
		var notes = []string{"c", "d", "e", "f", "g", "a", "b"}
		var spriteNames = []string{"cis", "dis", "b", "cis", "dis", "dis", "b"}

		for j, x := range notes {
			if x == noteName {
				i = j
			}
		}
		var offsetPos = glm.Vec2{this.Size[0] / 7.0 * float32(i), 0}
		pos = this.Pos.Add(&offsetPos)
		size = glm.Vec2{this.Size[0] / 7.0, this.Size[1]}
		spriteName = "./assets/" + spriteNames[i] + ".png"
	}

	//Beginn ab Oktave1: Note 33
	var sprite = this.ClosedGLCtx.Factory.NewSprite2D(glm.Vec2{pos[0] + float32(octave+1)*this.Size[0], pos[1]}, size, this.colours[(len(this.noteSprites))%len(this.colours)], spriteName)
	this.noteSprites[midiNote] = &sprite
}

func (this *Piano) DeleteKey(midiNote int) {
	this.noteSprites[midiNote] = nil
}

func (this *Piano) Process() {
	for i, x := range this.notesKeys {
		this.createPlayer(x, "./lotusland/assets/keys/3-"+this.notesPath[i]+".mp3")
		if this.ClosedGLCtx.KeyBoardManager.IsPressed(x) {
			var contains = false
			for _, y := range this.pressedKeys {
				if y == x {
					contains = true
					break
				}
			}
			if !contains {
				var spriteNames = []string{"cis", "dis", "b", "cis", "dis", "dis", "b"}
				this.pressedKeys = append(this.pressedKeys, x)
				if this.ClosedGLCtx.KeyBoardManager.IsPressed(x) {
					var offsetPos = glm.Vec2{this.Size[0] / 7.0 * float32(i), 0}
					var finalPos = this.Pos.Add(&offsetPos)
					var newKeySprite = this.ClosedGLCtx.Factory.NewSprite2D(finalPos, glm.Vec2{this.Size[0] / 7.0, this.Size[1]}, this.colours[(len(this.pressedKeys)-1)%len(this.colours)], "./assets/"+spriteNames[i]+".png")
					this.overlaySprites[x] = &newKeySprite
				}
			}
		}
		for i, x := range this.halfNotesKeys {
			this.createPlayer(x, "./lotusland/assets/keys/3-"+this.halfNotesPath[i]+".mp3")
			if this.ClosedGLCtx.KeyBoardManager.IsPressed(x) {
				var contains = false
				for _, y := range this.pressedKeys {
					if y == x {
						contains = true
						break
					}
				}
				if !contains {
					this.pressedKeys = append(this.pressedKeys, x)
					var origPixelOffsets = []float32{32, 87, 169, 222, 276}
					var origImgWidth float32 = 329
					var origImgHeight float32 = 326
					var origNoteWidth float32 = 26
					var origNoteHeight float32 = 200
					var offsetPos = glm.Vec2{float32((origPixelOffsets[i] / origImgWidth) * this.Size[0]), 0}
					var finalPos = this.Pos.Add(&offsetPos)
					var newKeySprite = this.ClosedGLCtx.Factory.NewSprite2D(finalPos, glm.Vec2{origNoteWidth / origImgWidth * this.Size[0], origNoteHeight / origImgHeight * this.Size[1]}, this.colours[(len(this.pressedKeys)-1)%len(this.colours)], "./assets/halfNote.png")
					this.overlaySprites[x] = &newKeySprite
				}
			}
		}
		var newP = []*oto.Player{}
		for i := 0; i < len(this.players); i++ {
			if !this.players[i].IsPlaying() {
			} else {
				newP = append(newP, this.players[i])
			}
		}

		var newKeys = []glfw.Key{}
		for _, x := range this.pressedKeys {
			if this.ClosedGLCtx.KeyBoardManager.IsDown(x) {
				newKeys = append(newKeys, x)
			} else {
				var sprite = this.overlaySprites[x]
				sprite.Free()
				this.overlaySprites[x] = nil
			}
		}
		this.pressedKeys = newKeys

		this.players = newP
	}
}

func (this *Piano) Draw() {
	for _, x := range this.sprites {
		x.Draw()
	}
	for _, x := range this.noteSprites {
		if x != nil {
			x.Draw()
		}
	}
	for _, x := range this.overlaySprites {
		if x != nil {
			x.Draw()
		}
	}
}

func (this *Piano) generatePlayer(path string, otoCtx *oto.Context) *oto.Player {
	var file, _ = os.ReadFile(path)
	var reader = bytes.NewReader(file)
	var dec, _ = mp3.NewDecoder(reader)
	var p = otoCtx.NewPlayer(dec)
	return p
}
