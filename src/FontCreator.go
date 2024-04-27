package main

import (
	"os"

	"github.com/EngoEngine/glm"
)

type FontCreator struct {
	grids []Grid
}

func newFontCreator(gridShader *Shader, projection *glm.Mat4) FontCreator {
	var grids = []Grid{}
	for i := 0; i < 128; i++ {
		grids = append(grids, newGrid(30, 16, gridShader, projection))
	}
	var fc = FontCreator{
		grids: grids,
	}
	fc.loadFont("default")
	return fc
}

func (this *FontCreator) draw(currIdx int) {
	this.grids[currIdx].draw()
}

func (this *FontCreator) loadFont(path string) {
	var file, _ = os.ReadFile("./font/" + path + ".iglbmf")

	var fileLen = len(file) / 128
	for i := 0; i < 128; i++ {

		var chunk = file[i*fileLen : (i+1)*fileLen]
		loadChunkInRect(&this.grids[i].cells, chunk)
	}
}
func (t *Text) serializeIglbmf(grid []Grid, path string) {
	var arr = []byte{}
	for i, x := range grid {
		var chunk = gridToChunk(x.cells, byte(i))
		arr = append(arr, chunk...)
	}
	var file, _ = os.Create("./font/" + path + ".iglbmf")
	t.iglbmfToIglbmt(arr, path)
	file.Write(arr)
	file.Close()
}
