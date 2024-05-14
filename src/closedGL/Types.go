package closed_gl

import "github.com/EngoEngine/glm"

type Plane struct {
	normal   glm.Vec3
	distance float32
}
type Frustum struct {
	//top,front,left,back,right,bottom
	planes []Plane
}
type Texture = uint32

type Buffer = uint32
type Vec2 = [2]float32
