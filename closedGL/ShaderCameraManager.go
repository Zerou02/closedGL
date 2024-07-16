package closedGL

import (
	"os"
	"strings"

	"github.com/EngoEngine/glm"
)

type ShaderCameraManager struct {
	Shadermap                  map[string]*Shader
	viewMatrix                 glm.Mat4
	projection2D, Projection3D glm.Mat4
	camera                     *Camera
}

func newShaderCameraManager(width, height float32, camera *Camera) ShaderCameraManager {
	var mane = ShaderCameraManager{Shadermap: map[string]*Shader{}}
	var dir, _ = os.ReadDir("./assets/shader")
	for _, x := range dir {
		var shaderName = strings.Split(x.Name(), ".")[0]
		mane.Shadermap[shaderName] = nil
	}

	for k := range mane.Shadermap {
		var shader = initShaderFromName(k)
		mane.Shadermap[k] = &shader
	}

	mane.camera = camera
	mane.projection2D = glm.Ortho2D(0, width, height, 0)
	mane.Projection3D = glm.Perspective(glm.DegToRad(45), width/height, 0.1, 2000)
	mane.viewMatrix = glm.Ident4()

	return mane
}
