package closedGL

/*

import (
	"github.com/EngoEngine/glm"
	"github.com/go-gl/glfw/v3.2/glfw"
)

type Player struct {
	chunk           *Chunk
	pos             glm.Vec3
	vel             glm.Vec3
	model           Cube
	keyboardManager *KeyBoardManager
	camera          *Camera
}

func newPlayer(pos glm.Vec3, chunk *Chunk, tex *Texture, keyboardManager *KeyBoardManager, camera *Camera) Player {
	var p = Player{
		chunk:           chunk,
		pos:             pos,
		model:           factory.NewCube(pos, tex),
		keyboardManager: keyboardManager,
		camera:          camera,
	}
	return p
}

func (this *Player) draw() {
	this.model.position = this.pos
	//this.model.draw()
}

func (this *Player) process(delta float32) {
	this.vel[0] = 0
	this.vel[2] = 0
	if this.keyboardManager.IsDown(glfw.KeyW) {
		this.vel = this.camera.cameraFront.Add(&this.vel)
	}
	if this.keyboardManager.IsDown(glfw.KeyS) {
		var new = this.camera.cameraFront.Mul(-1)
		this.vel = this.vel.Add(&new)
	}
	if this.keyboardManager.IsDown(glfw.KeyD) {
		var rotRight = glm.Rotate3DY(glm.DegToRad(-90))
		var right = rotRight.Mul3x1(&this.camera.cameraFront)
		this.vel = this.vel.Add(&right)
	}
	if this.keyboardManager.IsDown(glfw.KeyA) {
		var rotLeft = glm.Rotate3DY(glm.DegToRad(90))
		var left = rotLeft.Mul3x1(&this.camera.cameraFront)
		this.vel = this.vel.Add(&left)
	}
	this.vel = this.vel.Scale(10)

	this.vel = glm.Vec3{this.vel[0], -9.8, this.vel[2]}
	var dt = this.vel.Mul(delta)

	var oldPos = this.pos
	_ = oldPos
	this.pos = this.pos.Add(&dt)
	var isCollided = false
	for y := this.pos[1]; y <= oldPos[1]; y++ {
		var idx = pos3ToIdx(int(this.pos[0]), int(y), int(this.pos[2]), 32, 32, 32)
		if idx >= len(this.chunk.cubes) || idx < 0 {
			continue
		}
		//var exists = !this.chunk.cubes[idx].isInner
		/* if exists {
			isCollided = true
			break
		}
	}
	if isCollided {
		this.pos[1] = oldPos[1]
	}

	var isInChunk = true
	_ = isInChunk
	if this.pos[0] < this.chunk.pos[0] || this.pos[0] > this.chunk.pos[0]+this.chunk.dim[0] {
		isInChunk = false
	}
	if this.pos[1] < this.chunk.pos[1] || this.pos[1] > this.chunk.pos[1]+this.chunk.dim[1] {
		isInChunk = false
	}
	if this.pos[2] < this.chunk.pos[2] || this.pos[2] > this.chunk.pos[2]+this.chunk.dim[2] {
		isInChunk = false
	}
}
*/
