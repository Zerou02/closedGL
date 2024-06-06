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
	var circle = initShaderFromName("circle")
	var tileMap = initShaderFromName("tilemap")

	factory.camera = camera
	factory.Shadermap["base"] = &base
	factory.Shadermap["light"] = &light
	factory.Shadermap["points"] = &points
	factory.Shadermap["text"] = &text
	factory.Shadermap["cube"] = &cube
	factory.Shadermap["circle"] = &circle
	factory.Shadermap["tilemap"] = &tileMap

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

func (this *PrimitiveFactory) newRectManager() RectangleManager {
	return newRect(this.Shadermap["points"], &this.projectionMatrix)
}

func (this *PrimitiveFactory) NewLine(p1, p2 Vec2, colour1, colour2 glm.Vec3) Line {
	var l = newLine(this.Shadermap["points"], &this.projectionMatrix)
	l.addPoint(p1, colour1)
	l.addPoint(p2, colour2)
	return l
}

func (this *PrimitiveFactory) NewCircle(centreColour, borderColour glm.Vec4, radius float32, centre Vec2, borderThickness float32) Circle {
	return newCircle(this.Shadermap["circle"], &this.projectionMatrix, radius, centre, centreColour, borderColour, borderThickness)
}

func (this *PrimitiveFactory) NewCube(pos glm.Vec3, tex *Texture) Cube {
	return newCube(this.Shadermap["base"], this.camera, &this.Projection3D, tex, pos)
}

func (this *PrimitiveFactory) NewChunk(dim, pos glm.Vec3, tex *Texture) *Chunk {
	var chunk = NewChunk(dim, pos, tex, this.camera, &this.Projection3D, this.Shadermap["cube"])
	return &chunk
}

func (this *PrimitiveFactory) NewTriangle(points []Vec2, colour glm.Vec4) Triangle {
	return newTriangle(this.Shadermap["points"], &this.projectionMatrix, colour, points)
}

func (this *PrimitiveFactory) NewSprite2D(pos, size glm.Vec2, tint glm.Vec4, texPath string) Sprite2D {
	return newSprite2D(this.Shadermap["base"], pos, size, tint, texPath)
}

func (this *PrimitiveFactory) NewTileMap() TileMap {
	return NewTileMap(this.Shadermap["tilemap"], &this.projectionMatrix)
}
