package closedGL

/*

import (
	"fmt"

	"github.com/EngoEngine/glm"
	"github.com/go-gl/gl/v4.1-core/gl"
)

type Grid struct {
	cellSize, width int
	cells           []Rectangle
	shader          *Shader
	projection      *glm.Mat4
	vao, vbo        uint32
	vertices        []float32
	lineArr         *LineArr
}

func newGrid(cellSize, width int, cellShader *Shader, projection *glm.Mat4) Grid {

	var cells = generateGrid(cellSize, width, cellShader, projection)
	var line = newLineArr(cellShader, projection)
	var p1Colour = glm.Vec4{1, 0, 0, 1}
	var p2Colour = glm.Vec4{0, 0, 1, 1}
	for y := 0; y < width+1; y++ {
		var p1Pos = glm.Vec2{1, float32(y * cellSize)}
		var p2Pos = glm.Vec2{1 + float32(cellSize)*float32(width), float32(y * cellSize)}
		line.addPoint(p1Pos, p1Colour)
		line.addPoint(p2Pos, p2Colour)
	}
	for x := 0; x < width+1; x++ {
		var p1Pos = glm.Vec2{1 + float32(x*cellSize), 0}
		var p2Pos = glm.Vec2{1 + float32(x*cellSize), float32(cellSize) * float32(width)}
		line.addPoint(p1Pos, p1Colour)
		line.addPoint(p2Pos, p2Colour)
	}
	var grid = Grid{cellSize: cellSize, width: width, cells: cells, lineArr: &line}
	grid.optimize()
	return grid
}

func (this *Grid) draw() {
	this.shader.use()
	this.shader.setUniformMatrix4("projection", this.projection)
	gl.BindVertexArray(this.vao)
	gl.BindBuffer(gl.ARRAY_BUFFER, this.vbo)
	this.cells[0].createVertices()
	var offset = len(this.cells[0].vertices)
	for i, x := range this.cells {
		x.createVertices()
		for j, y := range x.vertices {
			this.vertices[i*offset+j] = y
		}
	}
	gl.BufferSubData(gl.ARRAY_BUFFER, 0, 4*len(this.vertices), gl.Ptr(this.vertices))
	gl.DrawArrays(gl.TRIANGLES, 0, int32((6 * len(this.cells))))
	if this.lineArr != nil {
		this.lineArr.draw()
	}
}

func generateGrid(size, amount int, shader *Shader, projection *glm.Mat4) []Rectangle {
	var rects = []Rectangle{}
	for y := 0; y < amount; y++ {
		for x := 0; x < amount; x++ {
			rects = append(rects, newRect(shader, projection, glm.Vec4{float32(x * size), float32(y * size), float32(size), float32(size)}, glm.Vec4{0, 0, 0, 0}))
		}
	}
	return rects
}

func (this *Grid) optimize() {
	this.shader = this.cells[0].shader
	this.projection = this.cells[0].projection
	var singleGridVerticesLen = len(this.cells[0].vertices)
	this.vertices = make([]float32, singleGridVerticesLen*len(this.cells))
	generateBuffers(&this.vao, &this.vbo, nil, nil, singleGridVerticesLen*4*len(this.cells), nil, []int{2, 4})
	for _, x := range this.cells {
		x.deleteBuffers()
	}
}

func (this *Grid) debugPrint() {
	for i, x := range this.cells {
		if i%this.width == 0 {
			println()
		}

		var c = x.colour
		var r = byte(lerp(0, 255, c[0]))
		var g = byte(lerp(0, 255, c[1]))
		var b = byte(lerp(0, 255, c[2]))
		var a = byte(lerp(0, 255, c[3]))
		fmt.Printf("%x%x%x%x ", r, g, b, a)

	}
}
*/
