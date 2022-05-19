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

	Edge1  Vec3
	Edge2  Vec3
	Cross  Vec3 //Normal
	Normal Vec3 //Normalized Normal
}

func NewTriangle(p1 Vec3, p2 Vec3, p3 Vec3, mat Material) Triangle {
	tri := Triangle{
		P1:  p1,
		P2:  p2,
		P3:  p3,
		Mat: mat,
	}

	tri.Edge1 = p2.Sub(p1)
	tri.Edge2 = p3.Sub(p1)
	tri.Cross = tri.Edge1.Cross(tri.Edge2)
	tri.Normal = tri.Cross.Normalize()
	return tri
}

//https://stackoverflow.com/questions/42740765/intersection-between-line-and-triangle-in-3d
//https://en.wikipedia.org/wiki/Möller–Trumbore_intersection_algorithm
//Intersect returns intersection and the barycentric coords
func (ray *Ray) Intersect(tri *Triangle) (bool, Vec3, Vec3) {
	det := -ray.Dir.Dot(tri.Cross)
	//Backface culling
	if det < 0 {
		return false, Vec3{}, Vec3{}
	}
	invDet := 1.0 / det
	a0 := ray.Pos.Sub(tri.P1)
	da0 := a0.Cross(ray.Dir)

	u := tri.Edge2.Dot(da0) * invDet
	v := -tri.Edge1.Dot(da0) * invDet
	t := a0.Dot(tri.Cross) * invDet

	intersect := t >= 0.0 && u >= 0.0 && v >= 0.0 && (u+v) <= 1.0
	if !intersect {
		return false, Vec3{}, Vec3{}
	}

	intersection := ray.Pos.Add(ray.Dir.MulScalar(t))
	barry := Vec3{X: u, Y: v, Z: 1.0 - u - v}
	return true, intersection, barry
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

func (ray Ray) RayCast(triangles *[]Triangle, bounces int, lastTri int) Colour {
	var nearest *Triangle = nil
	var hitPos Vec3
	hitI := -1
	depth := math.MaxFloat64

	//Find nearest triangle
	for i := range *triangles {
		//Bypass the triangle we've already intersected with
		if i == lastTri {
			continue
		}
		intersect, hitpos, _ := ray.Intersect(&(*triangles)[i])

		if intersect {
			dist := ray.Pos.Distance(hitpos)
			//Check if nearest
			if dist < depth {
				nearest = &(*triangles)[i]
				depth = dist
				hitPos = hitpos
				hitI = i
			}
		}
	}

	if nearest != nil {
		//At max bounces
		if bounces == MaxBounces-1 {
			//Just surface colour
			return nearest.Mat.SurfaceColour
		}

		//Recursive RayCast for each type of ray
		var diffuse Colour
		var metallic Colour
		//Diffuse
		if nearest.Mat.Diffuse > Epsilon {
			diffuse = DiffuseRay(ray, nearest.Normal, hitPos).RayCast(triangles, bounces+1, hitI)
		}
		//Metallic
		if nearest.Mat.Metallic > Epsilon {
			metallic = MetallicRay(ray, nearest.Normal, hitPos).RayCast(triangles, bounces+1, hitI)
		}

		//Calculate the colour based on the material
		return nearest.Mat.CalculateColour(diffuse, metallic)
	} else {
		//sky light calculation
		mag := (ray.Dir.Dot(Vec3{Z: 1}) + 0.5)
		if mag > 0 {
			return NewColour(mag, mag, mag)
		}
		return NewColour(0, 0, 0)
	}
}

var NumSamples int
var FrameBuf [WindowHeight][WindowWidth]Colour

func Render(tex *Texture, triangles []Triangle) {

	camPos := Vec3{X: 0, Y: -8, Z: 4}
	NumSamples++

	wait := sync.WaitGroup{}
	for y := 0; y < tex.Height; y++ {
		wait.Add(1)
		go func(y int) {
			for x := 0; x < tex.Width; x++ {
				//Generate Ray
				ray := Ray{
					Pos: camPos,
					Dir: RayDir(50, x, y, tex.Width, tex.Height),
				}

				//Start Raytracing
				colour := ray.RayCast(&triangles, 0, -1)

				//'Accumulate light' and average
				FrameBuf[y][x] = FrameBuf[y][x].Add(colour)
				avg := FrameBuf[y][x].DivScalar(float64(NumSamples))

				//Map to texture
				r := uint8(math.Max(math.Min(255*avg.R, 255), 0))
				g := uint8(math.Max(math.Min(255*avg.G, 255), 0))
				b := uint8(math.Max(math.Min(255*avg.B, 255), 0))
				tex.Set(x, y, Pixel{
					R: r,
					G: g,
					B: b,
				})
			}
			wait.Done()
		}(y)
	}
	wait.Wait()
}
