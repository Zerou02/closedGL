package closedGL

/* import (
	"github.com/EngoEngine/glm"
	"github.com/go-gl/gl/v4.1-core/gl"
)

type Point struct {
	shader     *Shader
	projection *glm.Mat4
	pos        glm.Vec2
	colour     glm.Vec3
	vao, vbo   uint32
}

func newPoint(shader *Shader, pos glm.Vec2, colour glm.Vec3, projection *glm.Mat4) Point {
	var p = Point{pos: pos, colour: colour, shader: shader, projection: projection}
	var vertex = []float32{
		pos[0], pos[1], colour[0], colour[1], colour[2],
	}
	generateBuffers(&p.vao, &p.vbo, nil, vertex, 0, nil, []int{2, 3})
	return p
}

func (this *Point) draw() {
	this.shader.use()
	this.shader.setUniformMatrix4("projection", this.projection)
	gl.BindVertexArray(this.vao)
	//gl.BindBuffer(gl.ARRAY_BUFFER, this.vbo)
	gl.DrawArrays(gl.POINTS, 0, 3)
} */
