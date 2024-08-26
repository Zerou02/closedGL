package closedGL

import (
	"os"

	"github.com/EngoEngine/glm"
	"github.com/go-gl/gl/v4.6-core/gl"
)

type CharacterInfo struct {
	offsetX, offsetY           uint32
	charX, charY, charW, charH byte
	asciicode                  byte
}

type Text struct {
	tex         *uint32
	texW, texH  uint32
	shader      *Shader
	charInfo    []CharacterInfo
	vao         uint32
	dataBuffer  BufferFloat
	projection  *glm.Mat4
	amountChars int
}

func NewText(font string, shader *Shader, projection *glm.Mat4) Text {
	var text = Text{shader: shader, projection: projection}
	text.deserializeIglbmt(font)
	text.vao = genVAO()
	text.dataBuffer = genSingularBufferFloat(text.vao, 0, 4, gl.FLOAT, false, 0)
	return text
}

func (this *Text) DrawText(posX, posY int, text string, scale float32) {
	this.createVertices(text, float32(posX), float32(posY), scale)
}

func (this *Text) clearBuffer() {
	this.dataBuffer.clear()
	this.amountChars = 0
}

func (this *Text) createVertices(text string, posX, posY float32, scale float32) {
	var letterWidth float32 = 10 * scale
	var letterHeight float32 = 10 * scale
	var spacing = 3
	var stride = 4 * 6
	this.dataBuffer.resizeCPUData((this.amountChars + len(text)) * stride)
	for _, x := range text {

		posX += letterWidth + float32(spacing)

		var info = this.charInfo[byte(x)]

		var startX = float32((uint32(info.charX) + info.offsetX)) / float32(this.texW)
		var startY = float32((uint32(info.charY) + info.offsetY)) / float32(this.texH)
		var endX = float32((uint32(info.charX) + info.offsetX + uint32(info.charW))) / float32(this.texW)
		var endY = float32((uint32(info.charY) + info.offsetY + uint32(info.charH))) / float32(this.texH)

		var newVertices = []float32{
			posX + letterWidth, posY, endX, startY, //tr
			posX, posY, startX, startY, //tl
			posX + letterWidth, posY + letterHeight, endX, endY, //br
			posX + letterWidth, posY + letterHeight, endX, endY, //br
			posX, posY, startX, startY, //tl
			posX, posY + letterHeight, startX, endY, //bl
		}
		for i := 0; i < len(newVertices); i++ {
			this.dataBuffer.cpuArr[this.amountChars*stride+i] = newVertices[i]
		}

		this.amountChars++
	}
}
func (this *Text) draw() {
	this.shader.use()
	gl.Disable(gl.DEPTH_TEST)
	this.shader.setUniformMatrix4("projection", this.projection)
	gl.ActiveTexture(gl.TEXTURE0)
	gl.BindTexture(gl.TEXTURE_2D, *this.tex)
	this.shader.setUniform1i("text", 0)

	gl.BindVertexArray(this.vao)
	this.dataBuffer.copyToGPU()

	gl.DrawArrays(gl.TRIANGLES, 0, int32(this.amountChars*6))
	gl.Enable(gl.DEPTH_TEST)

}

func (this *Text) deserializeIglbmt(path string) {
	var file, _ = os.ReadFile("./assets/font/" + path + ".iglbmt")
	file = RleDecode(file)
	var texLen = 147456
	var texData = file[0:texLen]
	var texPtr uint32
	var charDat = file[texLen:]
	var charInfo = []CharacterInfo{}
	for i := 0; i < len(charDat); i += 8 {
		this.texW = uint32(charDat[i])
		this.texH = uint32(charDat[i])
		var info = CharacterInfo{
			charX:     (charDat[i+1]),
			charY:     (charDat[i+2]),
			charW:     (charDat[i+3]),
			charH:     (charDat[i+4]),
			offsetX:   uint32(charDat[i+5]),
			offsetY:   uint32(charDat[i+6]),
			asciicode: charDat[i+7],
		}
		charInfo = append(charInfo, info)
	}
	gl.DeleteTextures(1, &texPtr)
	gl.GenTextures(1, &texPtr)
	gl.BindTexture(gl.TEXTURE_2D, texPtr)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_R, gl.REPEAT)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_T, gl.REPEAT)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MIN_FILTER, gl.NEAREST)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MAG_FILTER, gl.NEAREST)
	gl.TexImage2D(gl.TEXTURE_2D, 0, gl.RGBA, int32(192), int32(192), 0, gl.RGBA, gl.UNSIGNED_BYTE, gl.Ptr(texData))
	this.tex = &texPtr
	this.charInfo = charInfo
}
