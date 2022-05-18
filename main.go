package main

import (
	"fmt"
	"math"
	"os"
	"runtime/pprof"
	"time"

	"github.com/hschendel/stl"
	"github.com/veandco/go-sdl2/sdl"
)

func main() {
	f, _ := os.Create("ray.prof")
	pprof.StartCPUProfile(f)
	defer pprof.StopCPUProfile()

	model, _ := stl.ReadFile("LowPolyAnimal.stl")
	model.ScaleLinearDowntoSizeBox(stl.Vec3{10, 10, 10})
	model.Scale(3)
	model.Rotate(stl.Vec3{0, 0, 0}, stl.Vec3{1, 0, 0}, math.Pi/2)
	model.Rotate(stl.Vec3{0, 0, 0}, stl.Vec3{0, 0, 1}, math.Pi/4)
	min := model.Measure().Min
	model.Translate(stl.Vec3{-2, 10, -min[2]})

	triangles := make([]Triangle, 1)
	for _, tri := range model.Triangles {
		triangles = append(triangles,
			NewTriangle(
				NewVec3(float64(tri.Vertices[0][0]), float64(tri.Vertices[0][1]), float64(tri.Vertices[0][2])),
				NewVec3(float64(tri.Vertices[1][0]), float64(tri.Vertices[1][1]), float64(tri.Vertices[1][2])),
				NewVec3(float64(tri.Vertices[2][0]), float64(tri.Vertices[2][1]), float64(tri.Vertices[2][2])),
				Material{
					SurfaceColour: NewColour(0.1, 0.5, 0.5),
					Diffuse:       1.0,
					Metallic:      0.2,
				},
			),
		)
	}

	//Ground plane
	triangles = append(triangles, []Triangle{
		NewTriangle(
			NewVec3(-100, 100, 0),
			NewVec3(100, -100, 0),
			NewVec3(100, 100, 0),
			Material{
				SurfaceColour: NewColour(0.3, 0.3, 0.3),
				Diffuse:       1.0,
				Metallic:      1.0,
			},
		),
		NewTriangle(
			NewVec3(-100, 100, 0),
			NewVec3(-100, -100, 0),
			NewVec3(100, -100, 0),
			Material{
				SurfaceColour: NewColour(0.3, 0.3, 0.3),
				Diffuse:       1.0,
				Metallic:      1.0,
			},
		),
	}...)

	//Create Window
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

		Render(&tex, triangles)

		tex.Update()

		renderer.Present()

		select {
		case <-ticker.C:
			fmt.Print("\rFrameRate:", int(1.0/dt.Seconds()), "      ")
		default:
		}
	}
}
