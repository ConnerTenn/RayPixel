package main

import (
	"math"
	"unsafe"

	"github.com/veandco/go-sdl2/sdl"
)

const (
	Pi  = math.Pi
	Tau = Pi * 2.0
)

func ToRadians(degrees float64) float64 {
	return degrees * Tau / 360.0
}

type Pixel struct {
	B uint8
	G uint8
	R uint8
	_ uint8 //Padding
}

type Texture struct {
	Width    int
	Height   int
	Data     []Pixel
	Renderer *sdl.Renderer
	SdlTex   *sdl.Texture
}

func NewTexture(renderer *sdl.Renderer, width int, height int) Texture {
	sdltex, _ := renderer.CreateTexture(uint32(sdl.PIXELFORMAT_RGB888), sdl.TEXTUREACCESS_STREAMING, int32(width), int32(height))
	return Texture{
		Width:    width,
		Height:   height,
		Data:     make([]Pixel, width*height),
		Renderer: renderer,
		SdlTex:   sdltex,
	}
}

func (tex *Texture) ToByteArray() []byte {
	return (*(*[1]byte)(unsafe.Pointer(&tex)))[:]
}

func (tex *Texture) Set(x int, y int, pixel Pixel) {
	tex.Data[y*tex.Width+x] = pixel
}

func (tex *Texture) Update() {
	tex.SdlTex.Update(nil, (*(*[1]byte)(unsafe.Pointer(&tex.Data[0])))[:], tex.Width*4)
	tex.Renderer.Copy(tex.SdlTex, nil, nil)
}

func (tex *Texture) Destroy() {
	tex.SdlTex.Destroy()
}

func IntMax(a int, b int) int {
	if a >= b {
		return a
	}
	return b
}
func IntMin(a int, b int) int {
	if a <= b {
		return a
	}
	return b
}

// https://stackoverflow.com/questions/26237419/faster-than-rand
// Range [0,1)
var rseed uint32

func FastRandF() float64 {
	rseed = 214013*rseed + 2531011
	// return float64((rseed>>16)&0x7FFF) / 32767.0
	return float64(rseed>>16) / 65536.0
}
