package main

import "github.com/EngoEngine/glm"

type Chunk struct {
	dim   glm.Vec3
	cubes []Cube
}

func newChunk(dim glm.Vec3, tex *Texture) Chunk {
	var chunk = Chunk{dim: dim}
	var cubeArr = make([]Cube, int(dim[0]*dim[1]*dim[2]))
	var count = 0
	for y := 0; y < int(dim[1]); y++ {
		for z := 0; z < int(dim[2]); z++ {
			for x := 0; x < int(dim[0]); x++ {
				cubeArr[count] = factory.newCube(glm.Vec3{float32(x), float32(y), float32(z)}, tex)
				count += 1
			}
		}
	}
	chunk.cubes = cubeArr
	return chunk
}

func (this *Chunk) draw() {
	for i := 0; i < 16*16*16; i++ {
		this.cubes[i].draw()
	}
}
