package main

import (
	"math"

	v3 "github.com/deeean/go-vector/vector3"
)

type Material struct {
	Colour float64
}

func (mat Material) GetColour() float64 {
	return mat.Colour
}

var rseed uint

//https://stackoverflow.com/questions/26237419/faster-than-rand
func FastRandF() float64 {
	// return float64(time.Now().UnixNano()%1000) / 1000.0
	rseed = 214013*rseed + 2531011
	r := (rseed >> 16) & 0x7FFF
	return float64(r) / 32767
}

//https://www.scratchapixel.com/lessons/3d-basic-rendering/ray-tracing-overview/light-transport-ray-tracing-whitted
func (mat Material) NextRay(incoming Ray, normal *v3.Vector3, collide *v3.Vector3) Ray {

	diffuseDir := RotateV3(PerpendicularV3(normal).Lerp(normal, FastRandF()), normal, FastRandF()*math.Pi*2.0)
	// diffuseDir := v3.New(FastRandF()*2-1, FastRandF()*2-1, FastRandF()*2-1).Normalize()
	// if diffuseDir.Dot(normal) < 0 {
	// 	diffuseDir = diffuseDir.MulScalar(-1)
	// }
	// diffuseDir := ReflectV3(incoming.Dir, normal).Normalize()

	return Ray{
		Pos: collide,
		Dir: diffuseDir,
	}
}
