package main

import (
	v3 "github.com/deeean/go-vector/vector3"
)

type Material struct {
	Colour float64
}

func (mat Material) GetColour() float64 {
	return mat.Colour
}

//https://www.scratchapixel.com/lessons/3d-basic-rendering/ray-tracing-overview/light-transport-ray-tracing-whitted
func (mat Material) NextRay(incoming Ray, normal *v3.Vector3, collide *v3.Vector3) Ray {

	// diffuseDir := RotateV3(PerpendicularV3(normal).Lerp(normal, rand.Float64()), normal, rand.Float64()*math.Pi*2.0)
	diffuseDir := ReflectV3(incoming.Dir, normal).Normalize()

	return Ray{
		Pos: collide,
		Dir: diffuseDir,

		Bounces: incoming.Bounces + 1,
	}
}
