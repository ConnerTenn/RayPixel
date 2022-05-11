package main

import (
	"unsafe"

	"github.com/veandco/go-sdl2/sdl"
)

type Pixel struct {
	Blue  uint8
	Green uint8
	Red   uint8
	_     uint8 //Padding
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
