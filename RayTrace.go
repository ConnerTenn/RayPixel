package main

import (
	"math"
	"sync"

	v2 "github.com/deeean/go-vector/vector2"
	v3 "github.com/deeean/go-vector/vector3"
)

const Epsilon float64 = 0.000001
const MaxBounces int = 20

type Ray struct {
	Pos *v3.Vector3
	Dir *v3.Vector3

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

	intersect := det >= Epsilon && t >= 0.0 && u >= 0.0 && v >= 0.0 && (u+v) <= 1.0

	intersection := *ray.Pos.Add(ray.Dir.MulScalar(t))
	barry := v3.Vector3{X: u, Y: v, Z: 1.0 - u - v}
	return intersect, intersection, barry
}

func rotationXY(vec *v3.Vector3, x float64, y float64) *v3.Vector3 {
	cosX := math.Cos(x)
	cosY := math.Cos(y)
	sinX := math.Sin(x)
	sinY := math.Sin(y)

	row1 := v3.New(cosY, 0.0, -sinY)
	row2 := v3.New(sinY*sinX, cosX, cosY*sinX)
	row3 := v3.New(sinY*cosX, -sinX, cosY*cosX)

	return vec.Mul(row1).Mul(row2).Mul(row3)
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

func RayTrace(ray Ray) {
}

var t float64

func Render(tex *Texture) {

	triangles := []Triangle{
		{
			P1:  &v3.Vector3{X: -1, Y: 3, Z: 0},
			P2:  &v3.Vector3{X: 1, Y: 3, Z: 0},
			P3:  &v3.Vector3{X: 0, Y: 3, Z: 1},
			Mat: Material{Colour: 1},
		},
		{
			P1:  &v3.Vector3{X: 0, Y: 5, Z: 1},
			P2:  &v3.Vector3{X: 2, Y: 5, Z: 1},
			P3:  &v3.Vector3{X: 1, Y: 5, Z: 2},
			Mat: Material{Colour: 0.5},
		},
	}
	t += 3.1415 / 60.0

	wait := sync.WaitGroup{}
	for y := 0; y < tex.Height; y++ {
		wait.Add(1)
		go func(y int) {
			for x := 0; x < tex.Width; x++ {
				ray := Ray{
					Pos: &v3.Vector3{X: math.Sin(t), Y: 0, Z: 0},
					Dir: RayDir(50, x, y, tex.Width, tex.Height),
				}

				var intersected *Triangle = nil
				depth := math.MaxFloat64
				for i, tri := range triangles {
					intersect, iPos, _ := ray.Intersect(&tri)

					if intersect {
						dist := ray.Pos.Distance(&iPos)
						if dist < depth {
							intersected = &triangles[i]
							depth = dist
						}
					}
				}

				if intersected != nil {
					// tex.Set(x, y, Pixel{
					// 	Red:   byte(255 * (x + y) / (tex.Width + tex.Height)),
					// 	Green: byte(255 * x / tex.Width),
					// 	Blue:  byte(255 * y / tex.Height),
					// })
					tex.Set(x, y, Pixel{
						Red:   uint8(255 * intersected.Mat.GetColour()),
						Green: uint8(255 * intersected.Mat.GetColour()),
						Blue:  uint8(255 * intersected.Mat.GetColour()),
					})
				} else {
					tex.Set(x, y, Pixel{
						Red:   0,
						Green: 0,
						Blue:  0,
					})
				}
			}
			wait.Done()
		}(y)
	}
	wait.Wait()
}
