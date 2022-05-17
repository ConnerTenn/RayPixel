package main

type Material struct {
	Colour float64
}

func (mat Material) GetColour() float64 {
	return mat.Colour
}
