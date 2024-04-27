package main

import (
	"math"
	"os"
	"time"

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

func (this *FontCreator) serializeIglbmf(grid []Grid, path string) {
	var arr = []byte{}
	for i, x := range grid {
		var chunk = gridToChunk(x.cells, byte(i))
		arr = append(arr, chunk...)
	}
	var file, _ = os.Create("./font/" + path + ".iglbmf")
	this.iglbmfToIglbmt(arr, path)
	file.Write(arr)
	file.Close()
}

func (this *FontCreator) iglbmfToIglbmt(iglbmf []byte, path string) {
	var start = time.Now()
	var file = iglbmf
	var end = time.Now()
	var readTime = end.Sub(start).Seconds()
	_ = readTime
	var charInfo = []CharacterInfo{}
	var texData = []byte{}
	var texPtr uint32

	var chunkW = int(file[0])
	var dataOffset = int(file[6])
	var chunkSize = chunkW*chunkW*4 + dataOffset
	var amountChunks = len(file) / chunkSize
	var texRowLen = int(math.Ceil(math.Sqrt(float64(amountChunks))))
	var texRowHeight = texRowLen
	var chunksPerRow = texRowLen
	var imgLenPx = texRowLen * chunkW
	var texRowHeightPx = chunkW
	var chunkPxW = imgLenPx / texRowLen

	for texLine := 0; texLine < texRowHeight; texLine++ {
		var chunks = [][]byte{}
		for i := 0; i < chunksPerRow; i++ {
			var chunk = []byte{}
			var idx = (i + texLine*chunksPerRow)
			if idx*chunkSize >= len(file) {
				chunk = make([]byte, chunkSize)
			} else {
				chunk = file[idx*chunkSize : (idx+1)*chunkSize]
			}
			if idx < 128 {
				var posX, posY = idxToGridPos(idx, texRowLen, texRowLen)
				var info = CharacterInfo{
					tex: &texPtr, texW: uint32(imgLenPx), texH: uint32(imgLenPx),
					asciicode: chunk[5], charX: byte(chunk[1]), charY: byte(chunk[2]),
					charW: byte(chunk[3]), charH: byte(chunk[4]),
					offsetX: uint32(posX) * uint32(chunkW), offsetY: uint32(posY) * uint32(chunkPxW),
				}
				//TODO: Besser machen
				info.charX = 0
				info.charY = 0
				info.charW = 16
				info.charH = 16
				//TODO: End
				charInfo = append(charInfo, info)
			}
			chunks = append(chunks, chunk)
		}
		for y := 0; y < texRowHeightPx; y++ {
			for i := 0; i < chunksPerRow; i++ {
				var currChunkData = chunks[i][y*chunkPxW*4+dataOffset : (y+1)*chunkPxW*4+dataOffset]
				for j := 0; j < chunkPxW*4; j += 4 {
					texData = append(texData, currChunkData[j])
					texData = append(texData, currChunkData[j+1])
					texData = append(texData, currChunkData[j+2])
					texData = append(texData, currChunkData[j+3])
				}
			}
		}
	}
	for _, x := range charInfo {
		var bytes = []byte{byte(x.texW), byte(x.charX), byte(x.charY), byte(x.charW), byte(x.charH), byte(x.offsetX), byte(x.offsetY), x.asciicode}
		texData = append(texData, bytes...)
	}
	var outFile, _ = os.Create("font/" + path + ".iglbmt")
	outFile.Write(texData)
}
