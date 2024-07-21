package closedGL

import (
	"bytes"
	"os"

	"github.com/ebitengine/oto/v3"
	"github.com/hajimehoshi/go-mp3"
)

type Audio struct {
	ctx    oto.Context
	sounds []*oto.Player
	music  []*Music
}

type Music struct {
	player *oto.Player
	file   *os.File
}

func newAudio() Audio {

	op := &oto.NewContextOptions{}
	op.SampleRate = 44100
	op.ChannelCount = 2
	op.Format = oto.FormatSignedInt16LE
	var otoCtx, readyChan, err = oto.NewContext(op)
	if err != nil {
		panic("oto.NewContext failed: " + err.Error())
	}
	<-readyChan
	return Audio{
		ctx:    *otoCtx,
		sounds: []*oto.Player{},
		music:  []*Music{},
	}
}

func (this *Audio) process() {
	//TODO: Optimize
	var newSounds = []*oto.Player{}
	for _, x := range this.sounds {
		if x.IsPlaying() {
			newSounds = append(newSounds, x)
		}
	}
	this.sounds = newSounds

	var newMusic = []*Music{}
	for _, x := range this.music {
		if x.player.IsPlaying() {
			newMusic = append(newMusic, x)
		} else {
			x.file.Close()
		}
	}
	this.music = newMusic
}
func (this *Audio) playSound(name string) {

	var fileBytes, readErr = os.ReadFile("./assets/audio/" + name + ".mp3")
	if readErr != nil {
		println("could not load sound", readErr.Error())
	}
	var fileBytesReader = bytes.NewReader(fileBytes)
	var player = this.ctx.NewPlayer(fileBytesReader)
	player.Play()
	this.sounds = append(this.sounds, player)
}

func (this *Audio) streamMusic(name string, volume float64) {
	var fileBytes, err = os.Open("./assets/audio/" + name + ".mp3")
	if err != nil {
		panic("reading my-file.mp3 failed: " + err.Error())
	}
	var decodedMp3, err2 = mp3.NewDecoder(fileBytes)
	if err2 != nil {
		println("could not decode", err2.Error())
	}
	var player = this.ctx.NewPlayer(decodedMp3)
	player.Play()
	player.SetVolume(volume)
	var music = Music{
		player: player,
		file:   fileBytes,
	}
	this.music = append(this.music, &music)
}
