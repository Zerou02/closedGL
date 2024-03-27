package main

import (
	"fmt"
	"os"

	"github.com/EngoEngine/glm"
	"github.com/go-gl/gl/v4.1-core/gl"
)

type Shader struct {
	prog     uint32
	uniforms map[string]int32
}

func initShader(vertPath, fragPath string) Shader {
	var shader = Shader{0, map[string]int32{}}
	shader.prog = gl.CreateProgram()
	var vert, _ = os.ReadFile(vertPath)
	var frag, _ = os.ReadFile(fragPath)
	var vertS = shader.compileShader(string(vert)+"\x00", true)
	var fragS = shader.compileShader(string(frag)+"\x00", false)

	shader.prog = gl.CreateProgram()
	gl.AttachShader(shader.prog, vertS)
	gl.AttachShader(shader.prog, fragS)
	gl.DeleteShader(vertS)
	gl.DeleteShader(fragS)
	gl.LinkProgram(shader.prog)
	var succ int32 = 0
	gl.GetProgramiv(shader.prog, gl.LINK_STATUS, &succ)
	var s [512]uint8
	if succ == 0 {
		println("ERROR")
		gl.GetProgramInfoLog(shader.prog, 512, nil, &s[0])
		fmt.Printf("%s", s)
	} else {
		println("Shader successfully linked")
	}

	shader.parseUniforms(string(vert))
	shader.parseUniforms(string(frag))
	return shader
}

func (s *Shader) parseUniforms(shader string) {
	var ip = 0
	for ip < len(shader) {
		var c = shader[ip]
		if c == 'u' {
			if checkIsText(shader, ip, "uniform") {
				ip += len("uniform")
				parseUniformType(shader, &ip)
				var name = parseLiteral(shader, &ip)
				s.uniforms[name] = (gl.GetUniformLocation(s.prog, gl.Str(name+"\x00")))
			} else {
				ip++
			}
		} else {
			ip++
		}
	}
}

func (s *Shader) setUniformUInt(name string, value uint32) {
	gl.Uniform1ui(shader.uniforms[name], value)

}
func (s *Shader) setUniformMatrix4(name string, value *glm.Mat4) {
	gl.UniformMatrix4fv(shader.uniforms[name], 1, false, &value[0])
}

func parseLiteral(shader string, ip *int) string {
	var found = false
	var retStr = ""
	for !found {
		var c = shader[*ip]
		if c == ' ' {
			*ip += 1
		} else {
			for !found {
				c = shader[*ip]
				if c == ' ' || c == ';' {
					found = true
				} else {
					retStr += string(c)
					*ip += 1
				}
			}
		}
	}
	return retStr
}
func parseUniformType(shader string, ip *int) {
	var found = false
	for !found {
		var c = shader[*ip]
		if c == ' ' {
			*ip += 1
		} else {
			for !found {
				if c == ' ' {
					found = true
				} else {
					c = shader[*ip]
					*ip += 1
				}
			}
		}
	}
}

func checkIsText(text string, offset int, word string) bool {
	var isCorrect = true
	for i := 0; i < len(word); i++ {
		if word[i] != text[offset+i] {
			isCorrect = false
			break
		}
	}
	return isCorrect
}

func (s *Shader) compileShader(shaderSrc string, vertex bool) uint32 {
	var stype uint32 = gl.FRAGMENT_SHADER
	if vertex {
		stype = gl.VERTEX_SHADER
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
	} else {
		println("Shader successfully compiled")
	}
	return shader
}
