package main

import (
	"github.com/EngoEngine/glm"
	"github.com/go-gl/gl/v4.1-core/gl"
)

type CharacterInfo struct {
	tex                                    *uint32
	texW, texH, charX, charY, charW, charH uint32
	asciicode                              byte
}

type Character struct {
	shader        *Shader
	x, y, w, h    float32
	tint          glm.Vec3
	charInfo      CharacterInfo
	vao, vbo, ebo uint32
	projection    *glm.Mat4
}

func newCharacter(charInfo CharacterInfo, shader *Shader, x, y, w, h float32, tint glm.Vec3, projection *glm.Mat4) Character {
	var char = Character{shader, x, y, w, h, tint, charInfo, 0, 0, 0, projection}
	generateBuffers(&char.vao, &char.vbo, &char.ebo, nil, 4*4*4, indicesQuad, []VertexInfo{{4, 0}})
	return char
}

func (this *Character) draw() {
	this.shader.use()
	this.shader.setUniformMatrix4("projection", this.projection)
	this.shader.setUniformVec3("colour", &this.tint)
	gl.ActiveTexture(gl.TEXTURE0)
	gl.BindTexture(gl.TEXTURE_2D, *this.charInfo.tex)
	this.shader.setUniform1i("text", 0)

	gl.BindVertexArray(this.vao)
	gl.BindBuffer(gl.ARRAY_BUFFER, this.vbo)
	gl.BindBuffer(gl.ELEMENT_ARRAY_BUFFER, this.ebo)
	var string = "aabb"
	var posX = this.x
	var posY = this.y
	var letterWidth float32 = 100
	var letterHeight float32 = 100
	var spacing = 10
	for _, x := range string {
		_ = x
		var startX float32 = 1 / (float32(this.charInfo.texW) / (float32(this.charInfo.charX)))
		var startY float32 = 1 / (float32(this.charInfo.texH) / (float32(this.charInfo.charY)))
		var endX float32 = 1 / (float32(this.charInfo.texW) / (float32(this.charInfo.charX) + float32(this.charInfo.charW)))
		var endY float32 = 1 / (float32(this.charInfo.texH) / (float32(this.charInfo.charY) + float32(this.charInfo.charH)))
		var vertices = []float32{
			posX + letterWidth, posY, endX, startY,
			posX + letterWidth, posY + letterHeight, endX, endY,
			posX, posY + letterHeight, startX, endY,
			posX, posY, startX, startY,
		}
		posX += letterWidth + float32(spacing)
		gl.BufferSubData(gl.ARRAY_BUFFER, 0, 4*len(vertices), gl.Ptr(vertices))
		gl.DrawElements(gl.TRIANGLES, 6, gl.UNSIGNED_INT, nil)
	}
}
