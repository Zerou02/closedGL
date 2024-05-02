package main

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

func lerp(a, b, t float32) float32 {
	return a + (b-a)*t
}

func lerpVec2(a, b glm.Vec2, t float32) glm.Vec2 {
	return glm.Vec2{lerp(a[0], b[0], t), lerp(a[1], b[1], t)}
}

func bezierLerp(a, b, c glm.Vec2, t float32) glm.Vec2 {
	var r = lerpVec2(a, b, t)
	var s = lerpVec2(b, c, t)
	return lerpVec2(r, s, t)
}

func idxToGridPos(idx, w, h int) (int, int) {
	return idx % h, idx / w
}

func gridPosToIdx(posX, posY, w int) int {
	return posY*w + posX
}

func clamp(val, min, max float32) float32 {
	return float32(math.Max(float64(min), math.Min(float64(max), float64(val))))
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

func isPointInRect(p glm.Vec2, rect glm.Vec4) bool {
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
func idxToPos3(idx, x, y, z int) (int, int, int) {
	var yComp int = idx / int(x*y)
	var normalizedXIdx = idx - yComp*int(x*y)
	var xComp, zCom = idxToGridPos(normalizedXIdx, int(y), int(z))
	return xComp, yComp, zCom
}

// pos = Vec{y,z,x}
func pos3ToIdx(posX, posY, posZ int, dimX, dimY, dimZ int) int {
	var yLevel = posY * dimX * dimZ
	var idx = gridPosToIdx(posX, posZ, dimX)
	return yLevel + idx
}
