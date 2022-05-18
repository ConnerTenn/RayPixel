package main

import (
	"math"
	"sync"

	v2 "github.com/deeean/go-vector/vector2"
	v3 "github.com/deeean/go-vector/vector3"
)

const Epsilon float64 = 0.000001
const MaxBounces int = 3

type Ray struct {
	Pos *v3.Vector3
	Dir *v3.Vector3

	LastHit *Triangle
	Colour  *v3.Vector3
	Bounces int
}

type Line struct {
	P1 *v3.Vector3
	P2 *v3.Vector3
}

type Triangle struct {
	P1 *v3.Vector3
	P2 *v3.Vector3
	P3 *v3.Vector3

	Mat       Material
	TotalDist float64
}

//https://stackoverflow.com/questions/42740765/intersection-between-line-and-triangle-in-3d
//https://en.wikipedia.org/wiki/Möller–Trumbore_intersection_algorithm
//Intersect returns intersection and the barycentric coords
func (ray *Ray) Intersect(tri *Triangle) (bool, v3.Vector3, v3.Vector3) {
	edge1 := tri.P2.Sub(tri.P1)
	edge2 := tri.P3.Sub(tri.P1)
	n := edge1.Cross(edge2)
	det := -ray.Dir.Dot(n)
	invDet := 1.0 / det
	a0 := ray.Pos.Sub(tri.P1)
	da0 := a0.Cross(ray.Dir)

	u := edge2.Dot(da0) * invDet
	v := -edge1.Dot(da0) * invDet
	t := a0.Dot(n) * invDet

	intersect := math.Abs(det) >= Epsilon && t >= 0.0 && u >= 0.0 && v >= 0.0 && (u+v) <= 1.0

	intersection := *ray.Pos.Add(ray.Dir.MulScalar(t))
	barry := v3.Vector3{X: u, Y: v, Z: 1.0 - u - v}
	return intersect, intersection, barry
}

func (tri *Triangle) Normal() *v3.Vector3 {
	edge1 := tri.P2.Sub(tri.P1)
	edge2 := tri.P3.Sub(tri.P1)
	return edge1.Cross(edge2).Normalize()
}

// func rotationXY(vec *v3.Vector3, x float64, y float64) *v3.Vector3 {
// 	cosX := math.Cos(x)
// 	cosY := math.Cos(y)
// 	sinX := math.Sin(x)
// 	sinY := math.Sin(y)

// 	row1 := v3.New(cosY, 0.0, -sinY)
// 	row2 := v3.New(sinY*sinX, cosX, cosY*sinX)
// 	row3 := v3.New(sinY*cosX, -sinX, cosY*cosX)

// 	return vec.Mul(row1).Mul(row2).Mul(row3)
// }

//https://en.wikipedia.org/wiki/Rodrigues%27_rotation_formula
func RotateV3(vec *v3.Vector3, axis *v3.Vector3, angle float64) *v3.Vector3 {
	cos := math.Cos(angle)
	sin := math.Sin(angle)

	t1 := vec.MulScalar(cos)
	t2 := axis.Cross(vec).MulScalar(sin)
	t3 := axis.MulScalar(axis.Dot(vec) * (1.0 - sin))

	return t1.Add(t2).Add(t3)
}

func PerpendicularV3(vec *v3.Vector3) *v3.Vector3 {
	return &v3.Vector3{X: vec.Y, Y: vec.Z, Z: vec.X}
}

//http://www.3dkingdoms.com/weekly/weekly.php?a=2
func ReflectV3(vec *v3.Vector3, normal *v3.Vector3) *v3.Vector3 {
	return normal.MulScalar(-2.0 * vec.Dot(normal)).Add(vec)
}

func ToRadians(degrees float64) float64 {
	return degrees * math.Pi / 180.0
}

func RayDir(fov float64, x int, y int, width int, height int) *v3.Vector3 {
	size := v2.New(float64(width), float64(height))
	xz := v2.New(float64(x), float64(y)).Sub(size.DivScalar(2.0))

	halfFov := math.Tan(ToRadians(90.0 - fov*0.5))
	ypart := size.Y * 0.5 * halfFov

	return v3.New(xz.X, ypart, -xz.Y).Normalize()
}

func RayCast(ray Ray, triangles *[]Triangle) float64 {
	var nearest *Triangle = nil
	var hitPos *v3.Vector3 = nil
	depth := math.MaxFloat64

	for i, tri := range *triangles {
		if &(*triangles)[i] == ray.LastHit {
			continue
		}
		intersect, hitpos, _ := ray.Intersect(&tri)

		if intersect {
			dist := ray.Pos.Distance(&hitpos)
			if dist < depth {
				nearest = &(*triangles)[i]
				depth = dist
				hitPos = hitpos.Copy()
			}
		}
	}

	if nearest != nil {
		colour := nearest.Mat.GetColour()
		if ray.Bounces == MaxBounces-1 {
			return colour
		}

		nextray := nearest.Mat.NextRay(ray, nearest.Normal(), hitPos)
		nextray.LastHit = nearest
		nextcolour := RayCast(nextray, triangles)

		return colour * nextcolour
	} else {
		mag := (ray.Dir.Dot(&v3.Vector3{Z: 1}) + 0.5)
		if mag > 0 {
			return mag
		}
		return 0
	}
}

var t float64

func Render(tex *Texture) {

	triangles := []Triangle{
		{
			P1:  v3.New(-100, 100, 0),
			P2:  v3.New(100, -100, 0),
			P3:  v3.New(100, 100, 0),
			Mat: Material{Colour: 0.5},
		},
		{
			P1:  v3.New(-100, 100, 0),
			P2:  v3.New(100, -100, 0),
			P3:  v3.New(-100, -100, 0),
			Mat: Material{Colour: 0.5},
		},
		{
			P1:  v3.New(-1, 3, 0),
			P2:  v3.New(1, 3, 0),
			P3:  v3.New(0, 2, 1),
			Mat: Material{Colour: 0.8},
		},
		{
			P1:  v3.New(-0.5, 4, 0.5),
			P2:  v3.New(1.5, 4, 0.5),
			P3:  v3.New(0.5, 3.8, 1.5),
			Mat: Material{Colour: 0.5},
		},
		{
			P1:  v3.New(-2, 4, 0),
			P2:  v3.New(-2, 7, 5),
			P3:  v3.New(-5, 7, 5),
			Mat: Material{Colour: 1},
		},
	}
	t += 3.1415 / 60.0

	wait := sync.WaitGroup{}
	for y := 0; y < tex.Height; y++ {
		wait.Add(1)
		go func(y int) {
			for x := 0; x < tex.Width; x++ {
				ray := Ray{
					Pos: &v3.Vector3{X: 2 * math.Sin(t), Y: 2*math.Cos(t) - 4, Z: 2},
					Dir: RayDir(50, x, y, tex.Width, tex.Height),
				}

				colour := RayCast(ray, &triangles)

				rgb := uint8(math.Max(math.Min(255*colour, 255), 0))
				tex.Set(x, y, Pixel{
					Red:   rgb,
					Green: rgb,
					Blue:  rgb,
				})
			}
			wait.Done()
		}(y)
	}
	wait.Wait()
}
