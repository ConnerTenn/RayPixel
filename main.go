package main

import (
	"flag"
	"fmt"
	"os"
	"runtime/pprof"
	"time"

	"github.com/veandco/go-sdl2/sdl"
)

func main() {
	flag.Parse()
	f, _ := os.Create("ray.prof")
	pprof.StartCPUProfile(f)
	defer pprof.StopCPUProfile()

	const (
		windowWidth  = 800
		windowHeight = 600
	)
	window, _ := sdl.CreateWindow("RayPixel", sdl.WINDOWPOS_UNDEFINED, sdl.WINDOWPOS_UNDEFINED, windowWidth, windowHeight, sdl.WINDOW_OPENGL)
	defer window.Destroy()

	renderer, _ := sdl.CreateRenderer(window, -1, sdl.RENDERER_ACCELERATED)
	defer renderer.Destroy()

	renderer.Clear()

	tex := NewTexture(renderer, windowWidth, windowHeight)
	defer tex.Destroy()

	ticker := time.NewTicker(100 * time.Millisecond)
	start := time.Now()
	last := start
	running := true
	for running {
		now := time.Now()
		dt := now.Sub(last)
		last = now

		for event := sdl.PollEvent(); event != nil; event = sdl.PollEvent() {
			switch event.(type) {
			case *sdl.QuitEvent:
				running = false
			}
		}

		// renderer.Clear()
		// renderer.SetDrawColor(0, 0, 0, 0x20)
		// renderer.FillRect(nil)

		Render(&tex)

		tex.Update()

		renderer.Present()

		select {
		case <-ticker.C:
			fmt.Print("\rFrameRate:", int(1.0/dt.Seconds()), "      ")
		default:
		}
	}
}
