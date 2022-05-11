package main

import (
	"math"
	"sync"

	v2 "github.com/deeean/go-vector/vector2"
	v3 "github.com/deeean/go-vector/vector3"
)

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

type Ray struct {
	Pos *v3.Vector3
	Dir *v3.Vector3
}

func ToRadians(degrees float64) float64 {
	return degrees * math.Pi / 180.0
}

func RayDir(fov float64, x int, y int, width int, height int) *v3.Vector3 {
	size := v2.New(float64(width), float64(height))
	xy := v2.New(float64(x), float64(y)).Sub(size.DivScalar(2.0))

	halfFov := math.Tan(ToRadians(90.0 - fov*0.5))
	z := size.Y * 0.5 * halfFov

	return v3.New(xy.X, xy.Y, -z).Normalize()
}

func Render(tex *Texture) {
	wait := sync.WaitGroup{}
	for y := 0; y < tex.Height; y++ {
		wait.Add(1)
		go func(y int) {
			for x := 0; x < tex.Width; x++ {
				tex.Set(x, y, Pixel{
					Red:   byte(255 * (x + y) / (tex.Width + tex.Height)),
					Green: byte(255 * x / tex.Width),
					Blue:  byte(255 * y / tex.Height),
				})
			}
			wait.Done()
		}(y)
	}
	wait.Wait()
}
