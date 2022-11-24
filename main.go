package main

import (
	"fmt"
	"os"
	"runtime/pprof"
	"time"

	"github.com/hschendel/stl"
	"github.com/veandco/go-sdl2/sdl"
)

const (
	WindowWidth  = 800
	WindowHeight = 600
)

func main() {
	f, _ := os.Create("ray.prof")
	pprof.StartCPUProfile(f)
	defer pprof.StopCPUProfile()

	model, _ := stl.ReadFile("LowPolyAnimal.stl")
	model.ScaleLinearDowntoSizeBox(stl.Vec3{10, 10, 10})
	model.Scale(3)
	model.Rotate(stl.Vec3{0, 0, 0}, stl.Vec3{1, 0, 0}, Tau/4)
	model.Rotate(stl.Vec3{0, 0, 0}, stl.Vec3{0, 0, 1}, Tau/8)
	min := model.Measure().Min
	model.Translate(stl.Vec3{-1.5, 0, -min[2]})

	triangles := make([]Triangle, 1)
	for _, tri := range model.Triangles {
		triangles = append(triangles,
			NewTriangle(
				NewVec3(float64(tri.Vertices[0][0]), float64(tri.Vertices[0][1]), float64(tri.Vertices[0][2])),
				NewVec3(float64(tri.Vertices[1][0]), float64(tri.Vertices[1][1]), float64(tri.Vertices[1][2])),
				NewVec3(float64(tri.Vertices[2][0]), float64(tri.Vertices[2][1]), float64(tri.Vertices[2][2])),
				Material{
					SurfaceColour: NewColour(1, 0.8, 0.5),
					Diffuse:       1.0,
					Metallic:      0.2,
					Emissive:      0,
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
				Metallic:      0.1,
				Emissive:      0,
			},
		),
		NewTriangle(
			NewVec3(-100, 100, 0),
			NewVec3(-100, -100, 0),
			NewVec3(100, -100, 0),
			Material{
				SurfaceColour: NewColour(0.3, 0.3, 0.3),
				Diffuse:       1.0,
				Metallic:      0.1,
				Emissive:      0,
			},
		),
		NewTriangle(
			NewVec3(8, 10, 0),
			NewVec3(8, -10, 0),
			NewVec3(8, -10, 10),
			Material{
				SurfaceColour: NewColour(1, 0.2, 1),
				Diffuse:       0,
				Metallic:      0,
				Emissive:      0.5,
			},
		),
		NewTriangle(
			NewVec3(-8, -10, 0),
			NewVec3(-8, 10, 0),
			NewVec3(-8, -10, 10),
			Material{
				SurfaceColour: NewColour(0.2, 0.2, 1),
				Diffuse:       0,
				Metallic:      0,
				Emissive:      0.5,
			},
		),
	}...)

	//Create Window
	window, _ := sdl.CreateWindow("RayPixel", sdl.WINDOWPOS_UNDEFINED, sdl.WINDOWPOS_UNDEFINED, WindowWidth, WindowHeight, sdl.WINDOW_OPENGL)
	defer window.Destroy()

	renderer, _ := sdl.CreateRenderer(window, -1, sdl.RENDERER_ACCELERATED)
	defer renderer.Destroy()

	renderer.Clear()

	tex := NewTexture(renderer, WindowWidth, WindowHeight)
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
			fmt.Print("\rFrameRate:", int(1.0/dt.Seconds()), "  FrameTime:", int(dt.Milliseconds()), "ms", "  Samples:", NumSamples, "      ")
		default:
		}
	}

	fmt.Println()
	fmt.Println()
}
