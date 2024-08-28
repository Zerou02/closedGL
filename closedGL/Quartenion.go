package closedGL

import (
	"fmt"
	"math"
	"time"
)

type Vector struct {
	x, y, z float64
}

type Matrix struct {
	m11, m12, m13, m21, m22, m23, m31, m32, m33 float64
}

func (v Vector) Add(v2 Vector) Vector {
	return Vector{v.x + v2.x, v.y + v2.y, v.z + v2.z}
}

func (v Vector) Sub(v2 Vector) Vector {
	return v.Add(v2.Scale(-1))
}

func (v Vector) Scale(s float64) Vector {
	return Vector{v.x * s, v.y * s, v.z * s}
}

func (v Vector) Dot(v2 Vector) float64 {
	return v.x*v2.x + v.y*v2.y + v.z*v2.z
}

func (v Vector) Cross(v2 Vector) Vector {
	return Vector{v.y*v2.z - v.z*v2.y, v.z*v2.x - v.x*v2.z, v.x*v2.y - v.y*v2.x}
}

func (v Vector) Apply(m Matrix) Vector {
	return Vector{
		v.x*m.m11 + v.y*m.m12 + v.z*m.m13,
		v.x*m.m21 + v.y*m.m22 + v.z*m.m23,
		v.x*m.m31 + v.y*m.m32 + v.z*m.m33,
	}
}

func (m Matrix) Multiply(m2 Matrix) Matrix {
	return Matrix{
		m.m11*m2.m11 + m.m12*m2.m21 + m.m13*m2.m31,
		m.m11*m2.m12 + m.m12*m2.m22 + m.m13*m2.m32,
		m.m11*m2.m13 + m.m12*m2.m23 + m.m13*m2.m33,
		m.m21*m2.m11 + m.m22*m2.m21 + m.m23*m2.m31,
		m.m21*m2.m12 + m.m22*m2.m22 + m.m23*m2.m32,
		m.m21*m2.m13 + m.m22*m2.m23 + m.m23*m2.m33,
		m.m31*m2.m11 + m.m32*m2.m21 + m.m33*m2.m31,
		m.m31*m2.m12 + m.m32*m2.m22 + m.m33*m2.m32,
		m.m31*m2.m13 + m.m32*m2.m23 + m.m33*m2.m33,
	}
}

func (m Matrix) Transpose() Matrix {
	return Matrix{
		m.m11, m.m21, m.m31,
		m.m12, m.m22, m.m32,
		m.m13, m.m23, m.m33,
	}
}

func (m Matrix) Determinant() float64 {
	return m.m11*(m.m22*m.m33-m.m23*m.m32) - m.m12*(m.m21*m.m33-m.m23*m.m31) + m.m13*(m.m21*m.m32-m.m22*m.m31)
}

func (m Matrix) Inverse() (Matrix, error) {
	det := m.Determinant()
	if det == 0 {
		return Matrix{}, fmt.Errorf("Matrix is not invertible")
	}
	return Matrix{
		(m.m22*m.m33 - m.m23*m.m32) / det,
		(m.m13*m.m32 - m.m12*m.m33) / det,
		(m.m12*m.m23 - m.m13*m.m22) / det,
		(m.m23*m.m31 - m.m21*m.m33) / det,
		(m.m11*m.m33 - m.m13*m.m31) / det,
		(m.m13*m.m21 - m.m11*m.m23) / det,
		(m.m21*m.m32 - m.m22*m.m31) / det,
		(m.m12*m.m31 - m.m11*m.m32) / det,
		(m.m11*m.m22 - m.m12*m.m21) / det,
	}, nil
}

func (m Matrix) Add(m2 Matrix) Matrix {
	return Matrix{
		m.m11 + m2.m11, m.m12 + m2.m12, m.m13 + m2.m13,
		m.m21 + m2.m21, m.m22 + m2.m22, m.m23 + m2.m23,
		m.m31 + m2.m31, m.m32 + m2.m32, m.m33 + m2.m33,
	}
}

func (m Matrix) Scale(s float64) Matrix {
	return Matrix{
		m.m11 * s, m.m12 * s, m.m13 * s,
		m.m21 * s, m.m22 * s, m.m23 * s,
		m.m31 * s, m.m32 * s, m.m33 * s,
	}
}

func (m Matrix) Sub(m2 Matrix) Matrix {
	return m.Add(m2.Scale(-1))
}

type Quaternion struct {
	scale float64
	axis  Vector
}

func (q Quaternion) Add(q2 Quaternion) Quaternion {
	return Quaternion{q.scale + q2.scale, q.axis.Add(q2.axis)}
}

func (q Quaternion) Scale(s float64) Quaternion {
	return Quaternion{q.scale * s, q.axis.Scale(s)}
}

func (q Quaternion) Sub(q2 Quaternion) Quaternion {
	return q.Add(q2.Scale(-1))
}

func (q Quaternion) Multiply(q2 Quaternion) Quaternion {
	return Quaternion{q.scale*q2.scale - q.axis.Dot(q2.axis), q.axis.Cross(q2.axis).Add(q2.axis.Scale(q.scale)).Add(q.axis.Scale(q2.scale))}
}

func (q Quaternion) Norm() float64 {
	return math.Sqrt(q.scale*q.scale + q.axis.Dot(q.axis))
}

func (q Quaternion) Normalize() Quaternion {
	norm := q.Norm()
	return Quaternion{q.scale / norm, q.axis.Scale(1 / norm)}
}

func (q Quaternion) Conjugate() Quaternion {
	return Quaternion{q.scale, q.axis.Scale(-1)}
}

func (q Quaternion) Inverse() Quaternion {
	return q.Conjugate().Scale(1 / q.scale)
}

func (q Quaternion) Rotate(v Vector) Vector {
	qv := Quaternion{0, v}
	return q.Multiply(qv).Multiply(q.Conjugate()).axis
}

type Axis int

const (
	X Axis = iota
	Y
	Z
)

type Rotation struct {
	axis  Axis
	angle float64
}

func (a Axis) ToVector() Vector {
	switch a {
	case X:
		return Vector{1, 0, 0}
	case Y:
		return Vector{0, 1, 0}
	case Z:
		return Vector{0, 0, 1}
	}
	panic("Invalid axis")
}

func (r Rotation) ToQuaternion() Quaternion {
	angle := r.angle / 2
	return Quaternion{math.Cos(angle), r.axis.ToVector().Scale(math.Sin(angle))}
}

func (r Rotation) ToMatrix() Matrix {
	switch r.axis {
	case X:
		return Matrix{1, 0, 0, 0, math.Cos(r.angle), -math.Sin(r.angle), 0, math.Sin(r.angle), math.Cos(r.angle)}
	case Y:
		return Matrix{math.Cos(r.angle), 0, math.Sin(r.angle), 0, 1, 0, -math.Sin(r.angle), 0, math.Cos(r.angle)}
	case Z:
		return Matrix{math.Cos(r.angle), -math.Sin(r.angle), 0, math.Sin(r.angle), math.Cos(r.angle), 0, 0, 0, 1}
	}
	panic("Invalid axis")
}

func (v Vector) RotateWithQuaternion(rotations ...Rotation) Vector {
	for _, r := range rotations {
		v = r.ToQuaternion().Rotate(v)
	}
	return v
}

func (v Vector) RotateWithMatrix(rotations ...Rotation) Vector {
	for _, r := range rotations {
		v = v.Apply(r.ToMatrix())
	}
	return v
}

func main2op() {
	benchmark()
	// v := X.ToVector()
	// fmt.Println(Quaternion{scale: 0.7071068, axis: Vector{0, 0, 0.7071068}}.Rotate(v))
	// fmt.Println(v.RotateWithQuaternion(Rotation{Z, math.Pi / 2}))
}

func benchmark() {
	v := X.ToVector()
	start := time.Now()
	for i := 0; i < 1000000; i++ {
		v.RotateWithQuaternion(Rotation{Z, math.Pi / 2}, Rotation{X, math.Pi / 2}, Rotation{Y, math.Pi / 2})
	}
	end := time.Now()
	fmt.Println("Took ", end.Sub(start), "to rotate with quaternion. Got", v.RotateWithQuaternion(Rotation{Z, math.Pi / 2}, Rotation{X, math.Pi / 2}, Rotation{Y, math.Pi / 2}))
	start = time.Now()
	for i := 0; i < 1000000; i++ {
		v.RotateWithMatrix(Rotation{Z, math.Pi / 2}, Rotation{X, math.Pi / 2}, Rotation{Y, math.Pi / 2})
	}
	end = time.Now()
	fmt.Println("Took ", end.Sub(start), "to rotate with matrix. Got", v.RotateWithMatrix(Rotation{Z, math.Pi / 2}, Rotation{X, math.Pi / 2}, Rotation{Y, math.Pi / 2}))
}
