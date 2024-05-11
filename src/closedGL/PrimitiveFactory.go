package closed_gl

import (
	"github.com/EngoEngine/glm"
	"github.com/go-gl/gl/v4.1-core/gl"
)

type PrimitiveFactory struct {
	Shadermap                    map[string]*Shader
	projectionMatrix, viewMatrix glm.Mat4
	Projection3D                 glm.Mat4
	camera                       *Camera
}

func newPrimitiveFactory2D(width, height float32, camera *Camera) PrimitiveFactory {
	var factory = PrimitiveFactory{Shadermap: map[string]*Shader{}}
	var base = initShaderFromName("base")
	var light = initShaderFromName("light")
	var points = initShaderFromName("points")
	var text = initShaderFromName("text")
	var cube = initShaderFromName("cube")

	factory.camera = camera
	factory.Shadermap["base"] = &base
	factory.Shadermap["light"] = &light
	factory.Shadermap["points"] = &points
	factory.Shadermap["text"] = &text
	factory.Shadermap["cube"] = &cube

	factory.projectionMatrix = glm.Ortho2D(0, width, height, 0)
	factory.Projection3D = glm.Perspective(glm.DegToRad(45), width/height, 0.1, 2000)

	factory.viewMatrix = glm.Ident4()

	//TODO:Besser machen
	gl.UseProgram(factory.Shadermap["base"].prog)
	var shader = factory.Shadermap["base"]
	shader.setUniformMatrix4("projection", &factory.projectionMatrix)
	shader.setUniformMatrix4("view", &factory.viewMatrix)
	return factory
}

func (this *PrimitiveFactory) newRect(dim, colour glm.Vec4) Rectangle {
	return newRect(this.Shadermap["points"], &this.projectionMatrix, dim, colour)
}

func (this *PrimitiveFactory) newCube(pos glm.Vec3, tex *Texture) Cube {
	return newCube(this.Shadermap["base"], this.camera, &this.Projection3D, tex, pos)
}

func (this *PrimitiveFactory) NewChunk(dim, pos glm.Vec3, tex *Texture) *Chunk {
	var chunk = NewChunk(dim, pos, tex, this.camera, &this.Projection3D, this.Shadermap["cube"])
	return &chunk

}
