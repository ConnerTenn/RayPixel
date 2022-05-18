package main

import (
	"math"
	"sync"
)

const MaxBounces int = 3

type Ray struct {
	Pos Vec3
	Dir Vec3
}

type Line struct {
	P1 Vec3
	P2 Vec3
}

//Counter clockwise for normal facing camera
type Triangle struct {
	P1 Vec3
	P2 Vec3
	P3 Vec3

	Mat Material
}

//https://stackoverflow.com/questions/42740765/intersection-between-line-and-triangle-in-3d
//https://en.wikipedia.org/wiki/Möller–Trumbore_intersection_algorithm
//Intersect returns intersection and the barycentric coords
func (ray *Ray) Intersect(tri *Triangle) (bool, Vec3, Vec3) {
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

	intersection := ray.Pos.Add(ray.Dir.MulScalar(t))
	barry := Vec3{X: u, Y: v, Z: 1.0 - u - v}
	return intersect, intersection, barry
}

func (tri *Triangle) Normal() Vec3 {
	edge1 := tri.P2.Sub(tri.P1)
	edge2 := tri.P3.Sub(tri.P1)
	return edge1.Cross(edge2).Normalize()
}

func ToRadians(degrees float64) float64 {
	return degrees * math.Pi / 180.0
}

func RayDir(fov float64, x int, y int, width int, height int) Vec3 {
	size := NewVec2(float64(width), float64(height))
	xz := NewVec2(float64(x), float64(y)).Sub(size.DivScalar(2.0))

	halfFov := math.Tan(ToRadians(90.0 - fov*0.5))
	ypart := size.Y * 0.5 * halfFov

	return NewVec3(xz.X, ypart, -xz.Y).Normalize()
}

func RayCast(ray Ray, triangles *[]Triangle, bounces int, lastTri int) float64 {
	var nearest *Triangle = nil
	var hitPos Vec3
	hitI := -1
	depth := math.MaxFloat64

	//Find nearest triangle
	for i, tri := range *triangles {
		if i == lastTri {
			continue
		}
		intersect, hitpos, _ := ray.Intersect(&tri)

		if intersect {
			dist := ray.Pos.Distance(hitpos)
			if dist < depth {
				nearest = &(*triangles)[i]
				depth = dist
				hitPos = hitpos
				hitI = i
			}
		}
	}

	if nearest != nil {
		colour := nearest.Mat.GetColour()
		if bounces == MaxBounces-1 {
			return colour
		}

		nextray := nearest.Mat.NextRay(ray, nearest.Normal(), hitPos)
		nextcolour := RayCast(nextray, triangles, bounces+1, hitI)

		return colour * nextcolour
	} else {
		mag := (ray.Dir.Dot(Vec3{Z: 1}) + 0.5)
		if mag > 0 {
			return mag
		}
		return 0
	}
}

var t float64

var NumSamples int
var FrameBuf [600][800]float64

func Render(tex *Texture) {

	triangles := []Triangle{
		{
			P1:  NewVec3(-100, 100, 0),
			P2:  NewVec3(100, -100, 0),
			P3:  NewVec3(100, 100, 0),
			Mat: Material{Colour: 0.5},
		},
		{
			P1:  NewVec3(-100, 100, 0),
			P2:  NewVec3(-100, -100, 0),
			P3:  NewVec3(100, -100, 0),
			Mat: Material{Colour: 0.5},
		},
		{
			P1:  NewVec3(-1, 3, 0),
			P2:  NewVec3(1, 3, 0),
			P3:  NewVec3(0, 2, 1),
			Mat: Material{Colour: 0.8},
		},
		{
			P1:  NewVec3(-0.5, 4, 0.5),
			P2:  NewVec3(1.5, 4, 0.5),
			P3:  NewVec3(0.5, 3.8, 1.5),
			Mat: Material{Colour: 0.5},
		},
		{
			P1:  NewVec3(-2, 4, 0),
			P2:  NewVec3(-2, 7, 5),
			P3:  NewVec3(-5, 7, 5),
			Mat: Material{Colour: 1},
		},
	}
	// t += 3.1415 / 60.0

	camPos := Vec3{X: 0, Y: -4, Z: 2}
	NumSamples++

	wait := sync.WaitGroup{}
	for y := 0; y < tex.Height; y++ {
		wait.Add(1)
		go func(y int) {
			for x := 0; x < tex.Width; x++ {
				ray := Ray{
					Pos: camPos,
					Dir: RayDir(50, x, y, tex.Width, tex.Height),
				}

				colour := RayCast(ray, &triangles, 0, -1)

				FrameBuf[y][x] += colour

				rgb := uint8(math.Max(math.Min(255*FrameBuf[y][x]/float64(NumSamples), 255), 0))
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
