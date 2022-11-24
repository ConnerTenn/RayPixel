package main

import (
	"math"
	"sync"
)

const MaxBounces int = 8

type Ray struct {
	Pos Vec3
	Dir Vec3
}

type Line struct {
	P1 Vec3
	P2 Vec3
}

// Counter clockwise for normal facing camera
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

// https://stackoverflow.com/questions/42740765/intersection-between-line-and-triangle-in-3d
// https://en.wikipedia.org/wiki/Möller–Trumbore_intersection_algorithm
// Intersect returns intersection and the barycentric coords
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

func RayDir(fov float64, x float64, y float64) Vec3 {
	xz := NewVec2(x, y)

	halfFov := math.Tan(ToRadians(90.0 - fov*0.5))
	ypart := 0.5 * halfFov

	return NewVec3(xz.X, ypart, -xz.Y).Normalize()
}

func (ray Ray) RayCast(triangles *[]Triangle, bounces int, prevTri int) Colour {
	var nearest *Triangle = nil
	var hitPos Vec3
	hitI := -1
	depth := math.MaxFloat64

	//Find nearest triangle
	for i := range *triangles {
		//Bypass the triangle we've already intersected with
		if i == prevTri {
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

		//Raycast off next material
		distance := ray.Pos.Distance(hitPos)
		_ = distance
		return nearest.Mat.RayCast(ray, hitPos, nearest.Normal, triangles, bounces+1, hitI)
	} else {
		//sky light calculation
		sky := (ray.Dir.Dot(Vec3{Z: 1}) + 0.5)
		if sky > 0 {
			return NewColour(sky, sky, sky)
		}

		return NewColour(0, 0, 0)
	}
}
func (ray Ray) RayCastMany(triangles *[]Triangle, bounces int, prevTri int) Colour {
	var colour Colour
	for i := 0; i < 10; i++ {
		randRay := Vec3{FastRandF()*2 - 1, FastRandF()*2 - 1, FastRandF()*2 - 1}
		newray := ray
		newray.Pos = newray.Pos.Add(randRay.MulScalar(0.1))
		//Raycast
		colour = colour.Add(newray.RayCast(triangles, bounces, prevTri))
	}
	return colour.DivScalar(10)
}

type Camera struct {
	Position Vec3
	Rotation Vec3
}

func (cam *Camera) GetRay(x float64, y float64) Ray {
	ray := Ray{
		Pos: cam.Position,
		Dir: RayDir(50, x, y),
	}

	ray.Dir = ray.Dir.Rotate(
		Vec3{X: 0, Y: 1, Z: 0},
		cam.Rotation.Y,
	)
	ray.Dir = ray.Dir.Rotate(
		Vec3{X: 1, Y: 0, Z: 0},
		cam.Rotation.X,
	)
	ray.Dir = ray.Dir.Rotate(
		Vec3{X: 0, Y: 0, Z: 1},
		cam.Rotation.Z,
	)

	return ray
}

var NumSamples int

// Define the size of the pixel grid
const PixelSize int = 8

const AvgDepth int = 1

var FrameBuf [WindowHeight / PixelSize][WindowWidth / PixelSize][AvgDepth]Colour
var FrameIdx = 0

var timestep int

func Render(tex *Texture, triangles []Triangle) {
	timestep++
	cam := Camera{
		Position: Vec3{X: 0, Y: -18, Z: 8},
		Rotation: Vec3{X: -Tau * 0.04, Y: 0, Z: 0},
	}
	cam.Position = cam.Position.Rotate(Vec3{0, 0, 1}, float64(timestep)/100)
	cam.Rotation.Z = float64(timestep) / 100
	NumSamples++

	wait := sync.WaitGroup{}
	for y := 0; y < tex.Height; y++ {
		wait.Add(1)
		go func(y int) {
			for x := 0; x < tex.Width; x++ {
				xf := x / PixelSize
				yf := y / PixelSize
				xi := xf * PixelSize
				yi := yf * PixelSize
				if x%PixelSize == 0 && y%PixelSize == 0 {
					vx := float64(x) - float64(tex.Width)/2.0
					vy := float64(y) - float64(tex.Height)/2.0
					//Generate Ray
					ray := cam.GetRay(vx/float64(tex.Height), vy/float64(tex.Height))

					//Start Raytracing
					colour := ray.RayCastMany(&triangles, 0, -1)

					//'Accumulate light' and average
					FrameBuf[yf][xf][FrameIdx] = colour
					avg := Colour{}
					for i := 0; i < AvgDepth; i++ {
						avg = avg.Add(FrameBuf[yf][xf][(FrameIdx+i)%AvgDepth])
					}
					avg = avg.DivScalar(float64(AvgDepth))
					// avg := colour

					//Map to texture
					r := uint8(math.Max(math.Min(255*avg.R, 255), 0))
					g := uint8(math.Max(math.Min(255*avg.G, 255), 0))
					b := uint8(math.Max(math.Min(255*avg.B, 255), 0))
					tex.Set(x, y, Pixel{
						R: r,
						G: g,
						B: b,
					})
				} else {
					tex.Set(x, y, tex.Get(xi, yi))
				}
			}
			wait.Done()
		}(y)
	}
	wait.Wait()

	FrameIdx = (FrameIdx + 1) % AvgDepth
}
