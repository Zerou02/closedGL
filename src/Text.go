package main

import (
	"os"

	"github.com/EngoEngine/glm"
	"github.com/go-gl/gl/v4.1-core/gl"
)

type CharacterInfo struct {
	tex                          *uint32
	texW, texH, offsetX, offsetY uint32
	charX, charY, charW, charH   byte
	asciicode                    byte
}

type Text struct {
	shader     *Shader
	x, y, w, h float32
	tint       glm.Vec3
	charInfo   []CharacterInfo
	vao, vbo   uint32
	projection *glm.Mat4
	vertices   []float32
}

func newText(font string, shader *Shader, x, y, w, h float32, tint glm.Vec3, projection *glm.Mat4) Text {
	var text = Text{shader: shader, x: x, y: y, w: w, h: h, tint: tint, projection: projection}
	text.deserializeIglbmf(font)
	text.vertices = make([]float32, 4*6*64)
	generateBuffers(&text.vao, &text.vbo, nil, nil, len(text.vertices)*4, nil, []int{4})
	return text
}

func (this *Text) draw(text string) {
	this.shader.use()
	this.shader.setUniformMatrix4("projection", this.projection)
	this.shader.setUniformVec3("colour", &this.tint)
	gl.ActiveTexture(gl.TEXTURE0)
	gl.BindTexture(gl.TEXTURE_2D, *this.charInfo[0].tex)
	this.shader.setUniform1i("text", 0)

	gl.BindVertexArray(this.vao)
	gl.BindBuffer(gl.ARRAY_BUFFER, this.vbo)
	var posX = this.x
	var posY = this.y
	var letterWidth float32 = 10
	var letterHeight float32 = 10
	var spacing = 3
	var amountLetters = len(text)
	_ = amountLetters
	for i, x := range text {
		posX += letterWidth + float32(spacing)

		var info = this.charInfo[byte(x)]

		var startX = float32((uint32(info.charX) + info.offsetX)) / float32(info.texW)
		var startY = float32((uint32(info.charY) + info.offsetY)) / float32(info.texH)
		var endX = float32((uint32(info.charX) + info.offsetX + uint32(info.charW))) / float32(info.texW)
		var endY = float32((uint32(info.charY) + info.offsetY + uint32(info.charH))) / float32(info.texH)

		var newVertices = []float32{
			posX + letterWidth, posY, endX, startY, //tr
			posX, posY, startX, startY, //tl
			posX + letterWidth, posY + letterHeight, endX, endY, //br
			posX + letterWidth, posY + letterHeight, endX, endY, //br
			posX, posY, startX, startY, //tl
			posX, posY + letterHeight, startX, endY, //bl
		}
		for j := 0; j < len(newVertices); j++ {
			this.vertices[i*6*4+j] = newVertices[j]
		}
	}
	gl.BufferSubData(gl.ARRAY_BUFFER, 0, len(this.vertices)*4, gl.Ptr(this.vertices))
	gl.DrawArrays(gl.TRIANGLES, 0, int32(amountLetters)*6)

	gl.BindVertexArray(0)
	gl.BindBuffer(gl.ARRAY_BUFFER, 0)
	gl.BindTexture(gl.TEXTURE_2D, 0)
}

func (this *Text) deserializeIglbmf(path string) {
	var file, _ = os.ReadFile("font/" + path + ".iglbmt")
	var texLen = 147456
	var texData = file[0:texLen]
	var texPtr uint32
	var charDat = file[texLen:]
	var charInfo = []CharacterInfo{}
	for i := 0; i < len(charDat); i += 8 {
		var info = CharacterInfo{
			tex:       &texPtr,
			texW:      uint32(charDat[i]),
			texH:      uint32(charDat[i]),
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
	this.charInfo = charInfo
}
