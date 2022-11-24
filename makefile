
run: build
	./RayPixel

build:
	go build

prof:
	go tool pprof RayPixel ray.prof

clean:
	rm -f RayPixel ray.prof