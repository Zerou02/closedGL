package closedGL

import (
	"github.com/EngoEngine/glm"
	"github.com/go-gl/gl/v4.1-core/gl"
)

type Cube struct {
	shader     *Shader
	camera     *Camera
	projection *glm.Mat4
	tex        *Texture
	position   glm.Vec3
	vao, vbo   uint32
	isInner    bool
}

func newCube(shader *Shader, camera *Camera, projection *glm.Mat4, tex *Texture, pos glm.Vec3) Cube {
	var retCube = Cube{shader: shader, camera: camera, projection: projection, tex: tex, position: pos, isInner: false}
	//TODO:Fix
	//	generateBuffers(&retCube.vao, &retCube.vbo, nil, cubeVertices, 0, nil, []int{3, 2})
	return retCube
}

func (this *Cube) draw() {
	this.shader.use()
	this.shader.setUniformMatrix4("view", &this.camera.lookAtMat)
	this.shader.setUniformMatrix4("projection", this.projection)
	var model = glm.Translate3D(this.position[0], this.position[1], this.position[2])
	this.shader.setUniformMatrix4("model", &model)
	this.shader.setUniform1i("tex", 0)

	gl.BindVertexArray(this.vao)
	gl.BindBuffer(gl.ARRAY_BUFFER, this.vbo)
	gl.ActiveTexture(gl.TEXTURE0)
	gl.BindTexture(gl.TEXTURE_2D, *this.tex)
	gl.DrawArrays(gl.TRIANGLES, 0, int32(len(cube)/(3+2)))
}
