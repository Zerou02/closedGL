package closedGL

import "github.com/EngoEngine/glm"

type Plane struct {
	normal   glm.Vec3
	distance float32
}
type Frustum struct {
	//top,front,left,back,right,bottom
	planes [6]Plane
}
type Texture = uint32
