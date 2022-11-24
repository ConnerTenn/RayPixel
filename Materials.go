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

func (col Colour) Clamp() Colour {
	return Colour{
		R: math.Max(math.Min(col.R, 1.0), 0.0),
		G: math.Max(math.Min(col.G, 1.0), 0.0),
		B: math.Max(math.Min(col.B, 1.0), 0.0),
	}
}

type Material struct {
	SurfaceColour Colour
	Diffuse       float64
	Metallic      float64
	Emissive      float64
}

//https://www.scratchapixel.com/lessons/3d-basic-rendering/ray-tracing-overview/light-transport-ray-tracing-whitted

func DiffuseRay(normal Vec3, collide Vec3) Ray {
	randomDir := Vec3{FastRandF()*2 - 1, FastRandF()*2 - 1, FastRandF()*2 - 1}

	dir := normal.Add(randomDir)

	return Ray{
		Pos: collide,
		Dir: dir,
	}
}

func MetallicRay(incoming Ray, normal Vec3, collide Vec3) Ray {

	dir := incoming.Dir.Reflect(normal).Normalize()

	return Ray{
		Pos: collide,
		Dir: dir,
	}
}

func (mat *Material) RayCast(hitRay Ray, hitPos Vec3, normal Vec3, triangles *[]Triangle, bounces int, prevTri int) Colour {
	// var colour Colour

	var diffuse Colour
	var metallic Colour
	var emissive Colour

	var totalWeight float64 = 0.0

	if mat.Diffuse > Epsilon {
		ray := DiffuseRay(normal, hitPos)
		diffuse = ray.RayCast(triangles, bounces, prevTri)

		diffuse = mat.SurfaceColour.Mul(diffuse).MulScalar(mat.Diffuse)
		totalWeight += mat.Diffuse
	}

	if mat.Metallic > Epsilon {
		ray := MetallicRay(hitRay, normal, hitPos)
		metallic = ray.RayCast(triangles, bounces, prevTri)

		// metallic = mat.SurfaceColour.MulScalar(1-mat.Metallic).Add(metallic.MulScalar(mat.Metallic))
		metallic = metallic.MulScalar(mat.Metallic)
		totalWeight += mat.Metallic
	}

	if mat.Emissive > Epsilon {
		emissive = mat.SurfaceColour.MulScalar(mat.Emissive)
		totalWeight += mat.Emissive
	}

	return diffuse.Add(metallic).Add(emissive).DivScalar(totalWeight)
}

// //Calculate colour based on light that falls on it
// func (mat *Material) CalculateColour(diffuse Colour, metallic Colour) Colour {
// 	var colour Colour

// 	total := mat.Diffuse + mat.Metallic

// 	if mat.Diffuse > Epsilon {
// 		colour = mat.SurfaceColour.Mul(diffuse.MulScalar(mat.Diffuse))
// 	}

// 	if mat.Metallic > Epsilon {
// 		colour = colour.Add(metallic.MulScalar(mat.Metallic)).DivScalar(total)
// 	}

// 	colour = colour.Add(mat.SurfaceColour.MulScalar(mat.Emissive))
// 	if math.IsNaN(colour.R) {
// 		fmt.Println(colour)
// 	}

// 	return colour
// }
