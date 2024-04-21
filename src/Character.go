package main

import (
	"github.com/EngoEngine/glm"
	"github.com/go-gl/gl/v4.1-core/gl"
)

type CharacterInfo struct {
	tex       *uint32
	pos, size glm.Vec2
	advance   uint
}

type Character struct {
	shader        *Shader
	x, y, w, h    float32
	tint          glm.Vec4
	charInfo      CharacterInfo
	vao, vbo, ebo uint32
	projection    *glm.Mat4
}

func newCharacter(charInfo CharacterInfo, shader *Shader, x, y, w, h float32, tint glm.Vec4, projection *glm.Mat4) Character {
	var char = Character{shader, x, y, w, h, tint, charInfo, 0, 0, 0, projection}
	generateBuffers(&char.vao, &char.vbo, &char.ebo, nil, 4*4*4, indicesQuad, []VertexInfo{{4, 0}})
	return char
}

func (this *Character) draw() {
	this.shader.use()
	this.shader.setUniformMatrix4("projection", this.projection)
	this.shader.setUniformVec4("colour", &this.tint)
	gl.ActiveTexture(gl.TEXTURE0)
	gl.BindTexture(gl.TEXTURE_2D, *this.charInfo.tex)
	this.shader.setUniform1i("text", 0)

	var xpos float32 = 0
	var ypos float32 = 0
	var sizeX, sizeY float32 = 100, 40
	var h float32 = 1 * sizeX
	var w float32 = 1 * sizeY

	gl.BindVertexArray(this.vao)
	gl.BindBuffer(gl.ARRAY_BUFFER, this.vbo)
	gl.BindBuffer(gl.ELEMENT_ARRAY_BUFFER, this.ebo)
	var imgWidth float32 = 260
	_ = imgWidth
	//316x31
	var string = "helloworld"
	for _, x := range string {
		var idx float32 = float32(x - 'a')
		var basePercentage float32 = 0.0385
		var baseSpacing float32 = 40
		var glyphwidth float32 = 1.0 / 27.0
		_ = basePercentage
		var vertices = []float32{
			xpos + w, ypos, (idx)*basePercentage + glyphwidth, 0.0,
			xpos + w, ypos + h, (idx)*basePercentage + glyphwidth, 1.0,
			xpos, ypos + h, (idx) * basePercentage, 1.0,
			xpos, ypos, (idx) * basePercentage, 0.0,
		}
		xpos += baseSpacing
		gl.BufferSubData(gl.ARRAY_BUFFER, 0, 4*len(vertices), gl.Ptr(vertices))
		gl.DrawElements(gl.TRIANGLES, 6, gl.UNSIGNED_INT, nil)
	}
}
