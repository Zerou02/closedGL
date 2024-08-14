package closedGL

import (
	"fmt"
	"math"

	"github.com/EngoEngine/glm"
)

type Direction int32

const (
	UP Direction = iota
	RIGHT
	DOWN
	LEFT
)

type Vec4 struct {
	x, y, z, w float32
}

func InitZeroVec4() Vec4 {
	return Vec4{0, 0, 0, 0}
}

func initVec4(x, y, z, w float32) Vec4 {
	return Vec4{x, y, z, w}
}

func (v Vec4) Print() {
	fmt.Printf("%f, %f, %f, %f \n", v.x, v.y, v.z, v.w)
}

type Mat4 struct {
	data [][]float32
}

func InitZeroMat4() Mat4 {
	return Mat4{[][]float32{{0, 0, 0, 0}, {0, 0, 0, 0}, {0, 0, 0, 0}, {0, 0, 0, 0}}}
}

func (m Mat4) print() {
	for _, y := range m.data {
		for _, x := range y {
			fmt.Printf("%f ,", x)
		}
		println()
	}
	println()

}

func initMatrixDiagonal(x, y, z, w float32) Mat4 {
	var retMat = InitZeroMat4()
	retMat.data[0][0] = x
	retMat.data[1][1] = y
	retMat.data[2][2] = z
	retMat.data[3][3] = w
	return retMat
}

func initMatrixTranslation(x, y, z float32) Mat4 {
	var retMat = initMatrixDiagonal(1, 1, 1, 1)
	retMat.data[0][3] = x
	retMat.data[1][3] = y
	retMat.data[2][3] = z
	return retMat
}

func (m Mat4) multiplyMat4(m2 Mat4) Mat4 {
	var retMat = InitZeroMat4()
	for i := 0; i < 4; i++ {
		for j := 0; j < 4; j++ {
			var result float32 = 0
			for k := 0; k < int(len(m.data)); k++ {
				result += m.data[i][k] * m2.data[k][j]
			}
			retMat.data[i][j] = result
		}
	}
	return retMat
}

func (m Mat4) multiplyVec4(vec Vec4) Vec4 {
	var retVec = InitZeroVec4()
	retVec.x = m.data[0][0]*vec.x + m.data[0][1]*vec.y + m.data[0][2]*vec.z + m.data[0][3]*vec.w
	retVec.y = m.data[1][0]*vec.x + m.data[1][1]*vec.y + m.data[1][2]*vec.z + m.data[1][3]*vec.w
	retVec.z = m.data[2][0]*vec.x + m.data[2][1]*vec.y + m.data[2][2]*vec.z + m.data[2][3]*vec.w
	retVec.w = m.data[2][0]*vec.x + m.data[3][1]*vec.y + m.data[3][2]*vec.z + m.data[3][3]*vec.w
	return retVec
}

func mulMat4Vec4(m glm.Mat4, vec glm.Vec4) glm.Vec4 {
	var retVec = glm.Vec4{0, 0, 0, 0}
	retVec[0] = m.At(0, 0)*vec[0] + m.At(0, 1)*vec[1] + m.At(0, 2)*vec[2] + m.At(0, 3)*vec[3]
	retVec[1] = m.At(1, 0)*vec[0] + m.At(1, 1)*vec[1] + m.At(1, 2)*vec[2] + m.At(1, 3)*vec[3]
	retVec[2] = m.At(2, 0)*vec[0] + m.At(2, 1)*vec[1] + m.At(2, 2)*vec[2] + m.At(2, 3)*vec[3]
	retVec[2] = m.At(3, 0)*vec[0] + m.At(3, 1)*vec[1] + m.At(3, 2)*vec[2] + m.At(3, 3)*vec[3]
	return retVec
}

func createTransformation(rot glm.Vec3, translation glm.Vec3, scale glm.Vec3) glm.Mat4 {
	var retMat = glm.Ident4()
	var rotX = glm.HomogRotate3DX(rot.X())
	var rotY = glm.HomogRotate3DY(rot.Y())
	var rotZ = glm.HomogRotate3DZ(rot.Z())
	var trans = glm.Translate3D(translation[0], translation[1], translation[2])
	var scaleMat = glm.Scale3D(scale[0], scale[1], scale[2])
	retMat.Mul4With(&rotX)
	retMat.Mul4With(&rotY)
	retMat.Mul4With(&rotZ)
	retMat.Mul4With(&trans)
	retMat.Mul4With(&scaleMat)
	return retMat
}

func vectorDirection(target glm.Vec2) Direction {
	var compass = []glm.Vec2{
		glm.Vec2{0, 1},
		glm.Vec2{1, 0},
		glm.Vec2{0, -1},
		glm.Vec2{-1, 0},
	}
	var max = 0.0
	var bestMatch = -1

	for i := 0; i < 4; i++ {
		var normalized = target.Normalized()
		var angle = normalized.Dot(&compass[i])
		if float64(angle) > max {
			max = float64(angle)
			bestMatch = i
		}
	}
	return Direction(bestMatch)
}

func Lerp(a, b, t float32) float32 {
	return a + (b-a)*t
}

func Clamp(a, b, val float32) float32 {
	if val < a {
		return a
	} else if val > b {
		return b
	} else {
		return val
	}
}

func lerpVec2(a, b glm.Vec2, t float32) glm.Vec2 {
	return glm.Vec2{Lerp(a[0], b[0], t), Lerp(a[1], b[1], t)}
}

func BezierLerp(a, b, controlPoint glm.Vec2, t float32) glm.Vec2 {
	var r = lerpVec2(a, controlPoint, t)
	var s = lerpVec2(controlPoint, b, t)
	return lerpVec2(r, s, t)
}

func IdxToGridPos(idx, w, h int) (int, int) {
	return idx % h, idx / w
}

func GridPosToIdx(posX, posY, w int) int {
	return posY*w + posX
}

func ssVectorOriginCol(ssVel, ssWall glm.Vec2) glm.Vec2 {

	var esVel = ssVel.ComponentProduct(&glm.Vec2{1, -1})
	var angle = glm.RadToDeg(float32(math.Acos(float64(ssWall.Dot(&esVel) / (esVel.Len() * ssWall.Len())))))
	var rotangle = 2 * angle
	if angle == 0 {
		rotangle = 180
	}
	var rotMat = glm.Rotate2D(glm.DegToRad(360 - rotangle))
	var newAngle = rotMat.Mul2x1(&esVel)
	newAngle.Normalize()
	return newAngle.ComponentProduct(&glm.Vec2{1, -1})
}

func aabbAabbCol(b1, b2 glm.Vec4) bool {
	var colX = b1.X()+b1.Z() >= b2[0] && b2[0]+b2[2] >= b1[0]
	var colY = b1[1]+b1[3] >= b2[1] && b2[1]+b2[3] >= b1[1]
	return colX && colY
}

func aabbCircleCol(circle glm.Vec3, aabb glm.Vec4) (bool, Direction, glm.Vec2) {
	var centre = glm.Vec2{circle[0] + circle[2], circle[1] + circle[2]}
	var aabbHalf = glm.Vec2{aabb[2] / 2, aabb[3] / 2}
	var aabbCentre = glm.Vec2{aabb[0] + aabb[2]/2, aabb[1] + aabb[3]/2}
	var diff = centre.Sub(&aabbCentre)
	var clamped = glm.Vec2{glm.Clamp(diff[0], -aabbHalf[0], aabbHalf[0]), glm.Clamp(diff[1], -aabbHalf[1], aabbHalf[1])}
	var closest = aabbCentre.Add(&clamped)
	diff = closest.Sub(&centre)
	if diff.Len() < circle[2] {
		return true, vectorDirection(diff), diff
	} else {
		return false, UP, glm.Vec2{0, 0}
	}
}

func IsPointInRect(p glm.Vec2, rect glm.Vec4) bool {
	return p[0] >= rect[0] && p[0] <= rect[0]+rect[2] && p[1] >= rect[1] && p[1] <= rect[1]+rect[3]
}

func neededDecimalPlacesToNextInt(x float64) int {
	var y = x
	var retVal = 0
	for y < 0.99999999 {
		y *= 10
		retVal += 1
	}
	return retVal
}

// x,y,z
func IdxToPos3(idx, x, y, z int) (int, int, int) {
	var yComp int = idx / int(x*z)
	var normalizedXIdx = idx - yComp*int(x*y)
	var xComp, zCom = IdxToGridPos(normalizedXIdx, int(y), int(z))
	return xComp, yComp, zCom
}

// pos = Vec{y,z,x}
func Pos3ToIdx(posX, posY, posZ int, dimX, dimY, dimZ int) int {
	var yLevel = posY * dimX * dimZ
	var idx = GridPosToIdx(posX, posZ, dimX)
	return yLevel + idx
}

func clampToIntegerMultipleOf(val float32, multiple float32) float32 {
	var new = int(val / multiple)
	return float32(new) * multiple
}

func multidimensionalNewton(startVec glm.Vec2) {
	var current = startVec
	for i := 0; i < 10; i++ {
		var mat = glm.Mat2{
			float32(math.Sin(float64(current[0]))),
			1, 1,
			float32(math.Sin(float64(current[1]))),
		}
		var inv = mat.Inverse()
		var fmat = glm.Vec2{
			current[1] - float32(math.Cos(float64(current[0]))),
			current[0] - float32(math.Cos(float64(current[1]))),
		}
		var tmpMat = inv.Mul2x1(&fmat)
		current = current.Sub(&tmpMat)
		PrintFloat(current[0])
		PrintFloat(current[1])
		println()
	}

}

func isPointInFrustum(c *Camera, worldPos glm.Vec3) bool {
	var view = c.lookAtMat
	var l = view.Mul4x1(&glm.Vec4{worldPos[0], worldPos[1], worldPos[2], 1})
	var p = glm.Vec3{l[0], l[1], l[2]}
	var v = p.Sub(&c.CameraPos)
	var vz = v.Dot(&glm.Vec3{0, 0, -1})
	var isVisible = !(vz >= -2000 && vz <= -0.1)
	var frustumH = 2.0 * float64(vz) * math.Tan(0.5*float64(glm.DegToRad(c.fov)))
	var isVis2y = !(-frustumH > float64(l[1]) || float64(l[1]) > frustumH)
	var w = frustumH * (float64(c.aspect))
	var isVis2x = !(l[0] > float32(w) || l[0] < -float32(w))

	return isVis2x && isVis2y && isVisible
}

/* func SspointsToCartesianLine(p1, p2 glm.Vec2) glm.Vec2 {
	var c1 = SsToCartesian(p1)
	var c2 = SsToCartesian(p2)
	var dy = c2[1] - c1[1]
	var dx = c2[0] - c1[0]
	var m = dy / dx
	var n = c1[1] - c1[0]*m
	return glm.Vec2{m, n}
} */

func subVec2(this *glm.Vec2, v glm.Vec2) glm.Vec2 {
	return glm.Vec2{this[0] - v[0], this[1] - v[1]}
}
func Dist(p1, p2 glm.Vec2) float32 {
	var a = subVec2(&p1, p2)
	return float32(math.Sqrt(float64(a[0]*a[0] + a[1]*a[1])))
}

// Ein-/Ausgabe in SS
/* func IntersectionOfLines(p1, p2, p3, p4 glm.Vec2) glm.Vec2 {
	//m,n
	var vec1 = SspointsToCartesianLine(p1, p2)
	var vec2 = SspointsToCartesianLine(p3, p4)

	//(d-b)/(a-c)
	var x = (vec2[1] - vec1[1]) / (vec1[0] - vec2[0])
	var y = vec1[0]*x + vec1[1]

	return CartesianToSS(glm.Vec2{x, y})
} */

// screen is upper right quadrant
func SsToCartesian(p glm.Vec2, vpHeight float32) glm.Vec2 {
	return glm.Vec2{p[0], vpHeight - p[1]}
}

func CartesianToSS(p glm.Vec2, vpHeight float32) glm.Vec2 {
	return glm.Vec2{p[0], vpHeight - p[1]}
}

func MiddlePoint(p1, p2 glm.Vec2) glm.Vec2 {
	return glm.Vec2{(p1[0] + p2[0]) / 2.0, (p1[1] + p2[1]) / 2}
}

func IsPointInCircle(p, circleCentre glm.Vec2, r float32) bool {
	var dist = circleCentre.Sub(&p)
	return dist.Len() <= r
}
