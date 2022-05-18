package main

import (
	"math"
)

type Colour struct {
	R float64
	G float64
	B float64
}

func NewColour(r float64, g float64, b float64) Colour {
	return Colour{R: r, G: g, B: b}
}

func (col Colour) Add(other Colour) Colour {
	return Colour{
		R: col.R + other.R,
		G: col.G + other.G,
		B: col.B + other.B,
	}
}

func (col Colour) Mul(other Colour) Colour {
	return Colour{
		R: col.R * other.R,
		G: col.G * other.G,
		B: col.B * other.B,
	}
}

func (col Colour) MulScalar(scalar float64) Colour {
	return Colour{
		R: col.R * scalar,
		G: col.G * scalar,
		B: col.B * scalar,
	}
}

func (col Colour) DivScalar(scalar float64) Colour {
	return Colour{
		R: col.R / scalar,
		G: col.G / scalar,
		B: col.B / scalar,
	}
}

type Material struct {
	SurfaceColour Colour
}

func (mat Material) GetColour() Colour {
	return mat.SurfaceColour
}

var rseed uint

//https://stackoverflow.com/questions/26237419/faster-than-rand
//Range [0,1)
func FastRandF() float64 {
	// return float64(time.Now().UnixNano()%1000) / 1000.0
	rseed = 214013*rseed + 2531011
	return float64((rseed>>16)&0x7FFF) / 32767.0
}

//https://www.scratchapixel.com/lessons/3d-basic-rendering/ray-tracing-overview/light-transport-ray-tracing-whitted
func (mat Material) NextRay(incoming Ray, normal Vec3, collide Vec3) Ray {

	diffuseDir := normal.Perpendicular().Lerp(normal, FastRandF()).Rotate(normal, FastRandF()*math.Pi*2.0)
	// diffuseDir := v3.New(FastRandF()*2-1, FastRandF()*2-1, FastRandF()*2-1).Normalize()
	// if diffuseDir.Dot(normal) < 0 {
	// 	diffuseDir = diffuseDir.MulScalar(-1)
	// }
	// diffuseDir := incoming.Dir.Reflect(normal).Normalize()

	return Ray{
		Pos: collide,
		Dir: diffuseDir,
	}
}
