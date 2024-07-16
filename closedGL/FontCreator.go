package closedGL

/*

import (
	"fmt"
	"math"
	"os"
	"strconv"
	"time"

	"github.com/EngoEngine/glm"
	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/go-gl/glfw/v3.2/glfw"
)

type FontCreator struct {
	grids               []Grid
	cellSize, gridWidth int
	keyBoardManager     *KeyBoardManager
	window              *glfw.Window
	currentIdx          int
	currColour          glm.Vec4
	previewRect         Rectangle
	slider              [4]*Slider
	autoUpdate          bool
}

func newFontCreator(cellSize, gridWidth int, gridShader *Shader, projection *glm.Mat4, keyboardManger *KeyBoardManager, window *glfw.Window) FontCreator {
	var grids = []Grid{}
	for i := 0; i < 128; i++ {
		grids = append(grids, newGrid(cellSize, gridWidth, gridShader, projection))
	}
	var fc = FontCreator{
		grids: grids, keyBoardManager: keyboardManger, window: window,
		cellSize: cellSize, gridWidth: gridWidth, currentIdx: 0, currColour: glm.Vec4{0, 1, 1, 1},
		previewRect: factory.NewRect(glm.Vec4{600, 200, 100, 100}, glm.Vec4{0, 1, 1, 1}),
		slider:      [4]*Slider{},
		autoUpdate:  false,
	}
	fc.loadFont("default")

	var slider1 = newSlider(window, glm.Vec4{550, 340, 200, 10}, glm.Vec4{1, 0, 0, 1}, glm.Vec4{650, 340, 10, 10}, glm.Vec4{0, 0, 1, 1}, 0, 1, 0, 0.01, "r: ")
	var slider2 = newSlider(window, glm.Vec4{550, 390, 200, 10}, glm.Vec4{1, 0, 0, 1}, glm.Vec4{650, 390, 10, 10}, glm.Vec4{0, 0, 1, 1}, 0, 1, 1, 0.01, "g: ")
	var slider3 = newSlider(window, glm.Vec4{550, 440, 200, 10}, glm.Vec4{1, 0, 0, 1}, glm.Vec4{650, 440, 10, 10}, glm.Vec4{0, 0, 1, 1}, 0, 1, 1, 0.01, "b: ")
	var slider4 = newSlider(window, glm.Vec4{550, 490, 200, 10}, glm.Vec4{1, 0, 0, 1}, glm.Vec4{650, 490, 10, 10}, glm.Vec4{0, 0, 1, 1}, 0, 1, 1, 0.01, "a: ")

	fc.slider[0] = &slider1
	fc.slider[1] = &slider2
	fc.slider[2] = &slider3
	fc.slider[3] = &slider4

	return fc
}

func (this *FontCreator) process() {
	for i, x := range this.slider {
		x.process()
		this.currColour[i] = x.curr
	}

	var mouseX, mouseY = this.window.GetCursorPos()
	if mouseX > 0 && mouseX < float64(this.cellSize*this.gridWidth) && mouseY > 0 && mouseY < float64(this.cellSize*this.gridWidth) {
		var gridX, gridY int = int(mouseX) / this.cellSize, int(mouseY) / this.cellSize
		var idx = gridY*this.gridWidth + gridX
		_ = idx
		if this.window.GetMouseButton(glfw.MouseButton1) == glfw.Press {
			this.grids[this.currentIdx].cells[idx].colour = this.currColour
			if this.autoUpdate {
				this.serializeIglbmf("default")
				text.deserializeIglbmt("default")
			}
		}
		if this.window.GetMouseButton(glfw.MouseButton2) == glfw.Press {
			this.grids[this.currentIdx].cells[idx].colour = glm.Vec4{0, 0, 0, 0}
			if this.autoUpdate {
				this.serializeIglbmf("default")
				text.deserializeIglbmt("default")
			}
		}
	}
	if this.keyBoardManager.IsPressed(glfw.KeyQ) {
		this.moveIdx(-1)
		this.printIdx()
	}
	if this.keyBoardManager.IsPressed(glfw.KeyW) {
		this.moveIdx(1)
		this.printIdx()
	}
	if this.keyBoardManager.IsPressed(glfw.KeyE) {
		this.moveIdx(-10)
		this.printIdx()
	}
	if this.keyBoardManager.IsPressed(glfw.KeyR) {
		this.moveIdx(10)
		this.printIdx()
	}
	var keys = []glfw.Key{
		glfw.Key7, glfw.KeyU, glfw.KeyJ, glfw.KeyM,
		glfw.Key8, glfw.KeyI, glfw.KeyK, glfw.KeyComma,
		glfw.Key9, glfw.KeyO, glfw.KeyL, glfw.KeyPeriod,
		glfw.Key0, glfw.KeyP, glfw.KeySemicolon, glfw.KeySlash,
	}
	var values = []float32{-0.1, -0.01, 0.01, 0.1}
	for i, x := range keys {
		if this.keyBoardManager.IsPressed(x) {
			this.changeColour(i/4, values[i%4])
			this.printColour()
		}
	}
	if this.keyBoardManager.IsPressed(glfw.KeyS) {
		this.serializeIglbmf("default")
	}
	this.previewRect.colour = this.currColour
}

func (this *FontCreator) printIdx() {
	println(this.currentIdx, string(rune(this.currentIdx)))
}

func (this *FontCreator) printColour() {
	fmt.Printf("currColour: r:%f ,g:%f ,b:%f ,a:%f\n", this.currColour[0], this.currColour[1], this.currColour[2], this.currColour[3])
}

func (this *FontCreator) changeColour(idx int, delta float32) {
	this.currColour[idx] += delta
	if this.currColour[idx] < 0 {
		this.currColour[idx] = 0
	}
	if this.currColour[idx] > 1 {
		this.currColour[idx] = 1
	}
}

func (this *FontCreator) moveIdx(offset int) {
	this.currentIdx += offset
	if this.currentIdx < 0 {
		this.currentIdx = 0
	}
	if this.currentIdx >= len(this.grids) {
		this.currentIdx = len(this.grids) - 1
	}
}

func (this *FontCreator) draw() {
	gl.Disable(gl.DEPTH_TEST)
	this.grids[this.currentIdx].draw()
	for _, x := range this.slider {
		x.draw()
	}
	this.previewRect.Draw()
	var x float32 = 500
	var y float32 = 100
	text.createVertices("preview: "+string(rune(this.currentIdx)), x, y)
	y = 50
	text.createVertices("current index: "+strconv.FormatInt(int64(this.currentIdx), 10), x, y)
	gl.Enable(gl.DEPTH_TEST)

}

func (this *FontCreator) loadFont(path string) {
	var file, _ = os.ReadFile("./font/" + path + ".iglbmf")
	file = RleDecode(file)
	var fileLen = len(file) / 128
	for i := 0; i < 128; i++ {

		var chunk = file[i*fileLen : (i+1)*fileLen]
		this.loadChunkInRect(&this.grids[i].cells, chunk)
	}
}

func (this *FontCreator) serializeIglbmf(path string) {
	var arr = []byte{}
	for i, x := range this.grids {
		var chunk = this.gridToChunk(x.cells, byte(i))
		arr = append(arr, chunk...)
	}
	var file, _ = os.Create("./font/" + path + ".iglbmf")
	this.iglbmfToIglbmt(arr, path)
	file.Write(RleEncode(arr))
	file.Close()
}

func (this *FontCreator) gridToChunk(grid []Rectangle, asciicode byte) []byte {
	var chunk = make([]byte, len(grid)*4+7)
	var topmostY, bottommostY, rightmostX, leftmostX int = 16, 0, 0, 16
	for i := 0; i < len(grid); i++ {
		if grid[i].colour[3] != 0 {
			var gridX, gridY = IdxToGridPos(i, 16, 16)
			if gridX < leftmostX {
				leftmostX = gridX
			}
			if gridX > rightmostX {
				rightmostX = gridX
			}
			if gridY < topmostY {
				topmostY = gridY
			}
			if gridY > bottommostY {
				bottommostY = gridY
			}
		}
	}
	//gridSize,[4]charDim,asciicode,dataOffset
	chunk[0] = byte(math.Sqrt(float64(len(grid))))
	chunk[1] = byte(leftmostX)
	chunk[2] = byte(topmostY)
	chunk[3] = byte(rightmostX) - byte(leftmostX) + 1
	chunk[4] = byte(bottommostY) - byte(topmostY) + 1
	chunk[5] = asciicode
	chunk[6] = 7

	for i, x := range grid {
		var r = byte(lerp(0, 255, x.colour[0]))
		var g = byte(lerp(0, 255, x.colour[1]))
		var b = byte(lerp(0, 255, x.colour[2]))
		var a = byte(lerp(0, 255, x.colour[3]))
		chunk[(i*4)+int(chunk[6])] = r
		chunk[(i*4)+int(chunk[6])+1] = g
		chunk[(i*4)+int(chunk[6])+2] = b
		chunk[(i*4)+int(chunk[6])+3] = a
	}
	return chunk
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
	_ = texPtr
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
				var posX, posY = IdxToGridPos(idx, texRowLen, texRowLen)
				var info = CharacterInfo{
					//tex:       &texPtr,
					//texW:      uint32(imgLenPx),
					//texH:      uint32(imgLenPx),
					asciicode: chunk[5], charX: byte(chunk[1]), charY: byte(chunk[2]),
					charW: byte(chunk[3]), charH: byte(chunk[4]),
					offsetX: uint32(posX) * uint32(chunkW), offsetY: uint32(posY) * uint32(chunkPxW),
				}
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
		//TODO: Change
		var texW = 192
		var bytes = []byte{
			byte(texW),
			byte(x.charX), byte(x.charY), byte(x.charW), byte(x.charH), byte(x.offsetX), byte(x.offsetY), x.asciicode}
		texData = append(texData, bytes...)
	}
	var outFile, _ = os.Create("font/" + path + ".iglbmt")
	outFile.Write(RleEncode(texData))
}

func (this *FontCreator) loadChunkInRect(grid *[]Rectangle, chunk []byte) {
	var dataOffset = int(chunk[6])
	if dataOffset == 0 {
		dataOffset = 7
	}
	for i := dataOffset; i < len(chunk); i += 4 {
		var rect = &(*grid)[(i-dataOffset)/4]
		rect.colour[0] = float32(chunk[i]) / 255.0
		rect.colour[1] = float32(chunk[i+1]) / 255.0
		rect.colour[2] = float32(chunk[i+2]) / 255.0
		rect.colour[3] = float32(chunk[i+3]) / 255.0
	}
}
*/
