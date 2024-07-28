package closedGL

import (
	"fmt"
	"os"

	"github.com/EngoEngine/glm"
	"github.com/go-gl/gl/v4.1-core/gl"
)

type Shader struct {
	prog uint32
}

func initShader(vertPath, fragPath string) Shader {
	var shader = Shader{0}
	shader.prog = gl.CreateProgram()
	var vert, err = os.ReadFile(vertPath)
	var frag, err2 = os.ReadFile(fragPath)
	if err != nil {
		println("Could not find vert shader")
	}
	if err2 != nil {
		println("Could not find frag shader")
	}
	var vertS = shader.compileShader(string(vert)+"\x00", true, false)
	var fragS = shader.compileShader(string(frag)+"\x00", false, false)

	gl.AttachShader(shader.prog, vertS)
	gl.AttachShader(shader.prog, fragS)
	//	gl.DeleteShader(vertS)
	//	gl.DeleteShader(fragS)
	gl.LinkProgram(shader.prog)
	var succ int32 = 0
	gl.GetProgramiv(shader.prog, gl.LINK_STATUS, &succ)
	var s [512]uint8
	if succ == 0 {
		println("ERROR")
		gl.GetProgramInfoLog(shader.prog, 512, nil, &s[0])
		fmt.Printf("%s", s)
	}
	return shader
}

func initCompShader(path string) Shader {
	var shader = Shader{0}
	shader.prog = gl.CreateProgram()
	var compute = shader.compileShader(path, false, true)
	gl.AttachShader(shader.prog, compute)
	gl.LinkProgram(shader.prog)
	var succ int32 = 0
	gl.GetProgramiv(shader.prog, gl.LINK_STATUS, &succ)
	var s [512]uint8
	if succ == 0 {
		println("ERROR")
		gl.GetProgramInfoLog(shader.prog, 512, nil, &s[0])
		fmt.Printf("%s", s)
	}
	return shader
}

func initShaderFromName(name string) Shader {
	return initShader("./assets/shader/"+name+".vs", "./assets/shader/"+name+".fs")
}

func (s *Shader) setUniformMatrix4(name string, value *glm.Mat4) {
	var location = gl.GetUniformLocation(s.prog, gl.Str(name+"\x00"))
	gl.UniformMatrix4fv(location, 1, false, &value[0])
}

func (s *Shader) setUniformVec3(name string, value *glm.Vec3) {
	var location = gl.GetUniformLocation(s.prog, gl.Str(name+"\x00"))
	gl.Uniform3f(location, value[0], value[1], value[2])
}

func (s *Shader) setUniformVec4(name string, value *glm.Vec4) {
	var location = gl.GetUniformLocation(s.prog, gl.Str(name+"\x00"))
	gl.Uniform4f(location, value[0], value[1], value[2], value[3])
}

func (s *Shader) setUniform1i(name string, value int32) {
	var location = gl.GetUniformLocation(s.prog, gl.Str(name+"\x00"))
	gl.Uniform1i(location, value)
}

func (s *Shader) setUniform1f(name string, value float32) {
	var location = gl.GetUniformLocation(s.prog, gl.Str(name+"\x00"))
	gl.Uniform1f(location, value)
}

func (s *Shader) setUniform2f(name string, value []float32) {
	var location = gl.GetUniformLocation(s.prog, gl.Str(name+"\x00"))
	gl.Uniform2f(location, value[0], value[1])
}

func (s *Shader) setUniform1fv(name string, values []float32) {
	var location = gl.GetUniformLocation(s.prog, gl.Str(name+"\x00"))
	gl.Uniform1fv(location, int32(len(values)), &values[0])
}

func (s *Shader) setUniform1iv(name string, values []int32) {
	var location = gl.GetUniformLocation(s.prog, gl.Str(name+"\x00"))
	gl.Uniform1iv(location, int32(len(values)), &values[0])
}
func (s *Shader) setUniform2fv(name string, value glm.Vec2) {
	var location = gl.GetUniformLocation(s.prog, gl.Str(name+"\x00"))
	gl.Uniform2f(location, value[0], value[1])
}

func (s *Shader) compileShader(shaderSrc string, vertex bool, compute bool) uint32 {
	var stype uint32 = gl.FRAGMENT_SHADER
	if vertex {
		stype = gl.VERTEX_SHADER
	}
	if compute {
		stype = gl.COMPUTE_SHADER
	}
	var shader = gl.CreateShader(stype)
	var vertSrc, _ = gl.Strs(shaderSrc)
	gl.ShaderSource(shader, 1, vertSrc, nil)
	gl.CompileShader(shader)
	var success int32 = 0
	gl.GetShaderiv(shader, gl.COMPILE_STATUS, &success)
	if success == 0 {
		var ptr uint8
		gl.GetShaderInfoLog(shader, 512, nil, &ptr)
		println("Error at shader Compiling", ptr)
	}
	return shader
}

func (s *Shader) use() {
	gl.UseProgram(s.prog)
}
