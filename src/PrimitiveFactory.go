package main

import (
	"github.com/EngoEngine/glm"
	"github.com/go-gl/gl/v4.1-core/gl"
)

type PrimitiveFactory struct {
	shadermap                    map[string]*Shader
	projectionMatrix, viewMatrix glm.Mat4
}

func newPrimitiveFactory2D(width, height float32) PrimitiveFactory {
	var factory = PrimitiveFactory{shadermap: map[string]*Shader{}}
	var base = initShaderFromName("base")
	var light = initShaderFromName("light")
	var points = initShaderFromName("points")
	var text = initShaderFromName("text")

	factory.shadermap["base"] = &base
	factory.shadermap["light"] = &light
	factory.shadermap["points"] = &points
	factory.shadermap["text"] = &text
	factory.projectionMatrix = glm.Ortho2D(0, width, height, 0)
	factory.viewMatrix = glm.Ident4()

	//TODO:Besser machen
	gl.UseProgram(factory.shadermap["base"].prog)
	var shader = factory.shadermap["base"]
	shader.setUniformMatrix4("projection", &factory.projectionMatrix)
	shader.setUniformMatrix4("view", &factory.viewMatrix)
	return factory
}

func (this *PrimitiveFactory) newRect(dim, colour glm.Vec4) Rectangle {
	return newRect(this.shadermap["points"], &this.projectionMatrix, dim, colour)
}
