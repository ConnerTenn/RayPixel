package main

import "math"

const Epsilon float64 = 0.000001

type Vec2 struct {
	X float64
	Y float64
}

func NewVec2(x float64, y float64) Vec2 {
	return Vec2{X: x, Y: y}
}

func (vec Vec2) Add(other Vec2) Vec2 {
	return Vec2{
		X: vec.X + other.X,
		Y: vec.Y + other.Y,
	}
}

func (vec Vec2) Sub(other Vec2) Vec2 {
	return Vec2{
		X: vec.X - other.X,
		Y: vec.Y - other.Y,
	}
}

func (vec Vec2) MulScalar(scalar float64) Vec2 {
	return Vec2{
		X: vec.X * scalar,
		Y: vec.Y * scalar,
	}
}

func (vec Vec2) DivScalar(scalar float64) Vec2 {
	return Vec2{
		X: vec.X / scalar,
		Y: vec.Y / scalar,
	}
}

type Vec3 struct {
	X float64
	Y float64
	Z float64
}

func NewVec3(x float64, y float64, z float64) Vec3 {
	return Vec3{X: x, Y: y, Z: z}
}

func (vec Vec3) Add(other Vec3) Vec3 {
	return Vec3{
		X: vec.X + other.X,
		Y: vec.Y + other.Y,
		Z: vec.Z + other.Z,
	}
}

func (vec Vec3) Sub(other Vec3) Vec3 {
	return Vec3{
		X: vec.X - other.X,
		Y: vec.Y - other.Y,
		Z: vec.Z - other.Z,
	}
}

func (vec Vec3) MulScalar(scalar float64) Vec3 {
	return Vec3{
		X: vec.X * scalar,
		Y: vec.Y * scalar,
		Z: vec.Z * scalar,
	}
}

func (vec Vec3) DivScalar(scalar float64) Vec3 {
	return Vec3{
		X: vec.X / scalar,
		Y: vec.Y / scalar,
		Z: vec.Z / scalar,
	}
}

func (vec Vec3) Magnitude() float64 {
	return math.Sqrt(vec.X*vec.X + vec.Y*vec.Y + vec.Z*vec.Z)
}

func (vec Vec3) Normalize() Vec3 {
	m := vec.Magnitude()

	if m > Epsilon {
		return vec.DivScalar(m)
	}
	return vec
}

func (vec Vec3) Distance(other Vec3) float64 {
	dx := vec.X - other.X
	dy := vec.Y - other.Y
	dz := vec.Z - other.Z
	return math.Sqrt(dx*dx + dy*dy + dz*dz)
}

func (vec Vec3) Dot(other Vec3) float64 {
	return vec.X*other.X + vec.Y*other.Y + vec.Z*other.Z
}

func (vec Vec3) Cross(other Vec3) Vec3 {
	return Vec3{
		X: vec.Y*other.Z - vec.Z*other.Y,
		Y: vec.Z*other.X - vec.X*other.Z,
		Z: vec.X*other.Y - vec.Y*other.X,
	}
}

func (vec Vec3) Lerp(other Vec3, t float64) Vec3 {
	return Vec3{
		X: vec.X + (other.X-vec.X)*t,
		Y: vec.Y + (other.Y-vec.Y)*t,
		Z: vec.Z + (other.Z-vec.Z)*t,
	}
}

//https://en.wikipedia.org/wiki/Rodrigues%27_rotation_formula
func (vec Vec3) Rotate(axis Vec3, angle float64) Vec3 {
	cos := math.Cos(angle)
	sin := math.Sin(angle)

	t1 := vec.MulScalar(cos)
	t2 := axis.Cross(vec).MulScalar(sin)
	t3 := axis.MulScalar(axis.Dot(vec) * (1.0 - cos))

	return t1.Add(t2).Add(t3)
}

func (vec Vec3) Perpendicular() Vec3 {
	return Vec3{X: vec.Y, Y: vec.Z, Z: vec.X}
}

//http://www.3dkingdoms.com/weekly/weekly.php?a=2
func (vec Vec3) Reflect(normal Vec3) Vec3 {
	return normal.MulScalar(-2.0 * vec.Dot(normal)).Add(vec)
}
