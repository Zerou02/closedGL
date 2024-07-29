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
	vao        uint32
	isInner    bool
	buffer     BufferFloat
}

func NewCube(shader *Shader, camera *Camera, projection *glm.Mat4, tex *Texture, pos glm.Vec3) Cube {
	var retCube = Cube{shader: shader, camera: camera, projection: projection, tex: tex, position: pos, isInner: false}
	//TODO:Fix
	retCube.vao = genVAO()
	retCube.buffer = generateInterleavedVBOFloat2(retCube.vao, 0, []int{3, 2})
	retCube.buffer.cpuArr = cube
	retCube.buffer.copyToGPU()
	return retCube
}

func (this *Cube) Draw() {
	this.shader.use()
	gl.Disable(gl.CULL_FACE)
	this.shader.setUniformMatrix4("view", &this.camera.lookAtMat)
	//var proj = glm.Ident4()
	this.shader.setUniformMatrix4("projection", &this.camera.perspective)
	//	var model = glm.Translate3D(this.position[0], this.position[1], this.position[2])
	//	this.shader.setUniformMatrix4("model", &model)
	this.shader.setUniform1i("tex", 0)

	gl.BindVertexArray(this.vao)
	this.buffer.copyToGPU()
	gl.ActiveTexture(gl.TEXTURE0)
	gl.BindTexture(gl.TEXTURE_2D, *this.tex)
	gl.DrawArrays(gl.TRIANGLES, 0, int32(len(cube)/(3+2)))
	gl.Enable(gl.CULL_FACE)

}
