package main

import (
	"sync"
	"time"

	"github.com/veandco/go-sdl2/sdl"
)

const (
	WindowWidth  = 800
	WindowHeight = 600
	RectWidth    = 20
	RectHeight   = 20
)

func main() {
	window, _ := sdl.CreateWindow("RayPixel", sdl.WINDOWPOS_UNDEFINED, sdl.WINDOWPOS_UNDEFINED, WindowWidth, WindowHeight, sdl.WINDOW_OPENGL)
	defer window.Destroy()

	renderer, _ := sdl.CreateRenderer(window, -1, sdl.RENDERER_ACCELERATED)
	defer renderer.Destroy()

	renderer.Clear()

	tex := NewTexture(renderer, WindowWidth, WindowHeight)
	defer tex.Destroy()

	start := time.Now()
	running := true
	for running {
		now := time.Now()
		dt := now.Sub(start)

		for event := sdl.PollEvent(); event != nil; event = sdl.PollEvent() {
			switch event.(type) {
			case *sdl.QuitEvent:
				running = false
			}
		}

		// renderer.Clear()
		// renderer.SetDrawColor(0, 0, 0, 0x20)
		// renderer.FillRect(nil)

		wg := sync.WaitGroup{}
		for y := 0; y < WindowHeight; y++ {
			wg.Add(1)
			go func(y int) {
				for x := 0; x < WindowWidth; x++ {
					tex.Set(x, y, Pixel{
						Red:   byte(int(dt.Milliseconds()/10) + 255*(x-y)/(WindowWidth+WindowHeight)),
						Green: byte(int(dt.Milliseconds()/20) + 255*x/WindowWidth),
						Blue:  byte(int(dt.Milliseconds()/40) + 255*y/WindowHeight),
					})
				}
				wg.Done()
			}(y)
		}

		tex.Update()

		renderer.Present()
	}
}
