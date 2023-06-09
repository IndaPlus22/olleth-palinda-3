// Stefan Nilsson 2013-02-27

// This program creates pictures of Julia sets (en.wikipedia.org/wiki/Julia_set).

// Original Run-time: 10.509s
// Improved run-time: 1.628s
package main

import (
	"image"
	"image/color"
	"image/png"
	"log"
	"math/cmplx"
	"os"
	"runtime"
	"strconv"
	"sync"
)

type ComplexFunc func(complex128) complex128

var Funcs []ComplexFunc = []ComplexFunc{
	func(z complex128) complex128 { return z*z - 0.61803398875 },
	func(z complex128) complex128 { return z*z + complex(0, 1) },
	func(z complex128) complex128 { return z*z + complex(-0.835, -0.2321) },
	func(z complex128) complex128 { return z*z + complex(0.45, 0.1428) },
	func(z complex128) complex128 { return z*z*z + 0.400 },
	func(z complex128) complex128 { return cmplx.Exp(z*z*z) - 0.621 },
	func(z complex128) complex128 { return (z*z+z)/cmplx.Log(z) + complex(0.268, 0.060) },
	func(z complex128) complex128 { return cmplx.Sqrt(cmplx.Sinh(z*z)) + complex(0.065, 0.122) },
}

func main() {
	for n, fn := range Funcs {
		err := CreatePng("picture-"+strconv.Itoa(n)+".png", fn, 1024)
		if err != nil {
			log.Fatal(err)
		}
	}
}

// CreatePng creates a PNG picture file with a Julia image of size n x n.
func CreatePng(filename string, f ComplexFunc, n int) (err error) {
	file, err := os.Create(filename)
	if err != nil {
		return
	}
	defer file.Close()
	err = png.Encode(file, JuliaParallel(f, n))
	return
}

// Julia returns an image of size n x n of the Julia set for f.
func Julia(f ComplexFunc, n int) image.Image {
	bounds := image.Rect(-n/2, -n/2, n/2, n/2)
	img := image.NewRGBA(bounds)
	s := float64(n / 4)
	for i := bounds.Min.X; i < bounds.Max.X; i++ {
		for j := bounds.Min.Y; j < bounds.Max.Y; j++ {
			n := Iterate(f, complex(float64(i)/s, float64(j)/s), 256)
			r := uint8(0)
			g := uint8(0)
			b := uint8(n % 32 * 8)
			img.Set(i, j, color.RGBA{r, g, b, 255})
		}
	}
	return img
}

// Iterate sets z_0 = z, and repeatedly computes z_n = f(z_{n-1}), n â‰¥ 1,
// until |z_n| > 2  or n = max and returns this n.
func Iterate(f ComplexFunc, z complex128, max int) (n int) {
	for ; n < max; n++ {
		if real(z)*real(z)+imag(z)*imag(z) > 4 {
			break
		}
		z = f(z)
	}
	return
}

//Divides the images up into rectangles and computed parallel on all available CPUs.
/*
The JuliaParallel function divides the image into smaller rectangles of size
32x32 pixels, and sends each rectangle as a task to a channel. Then, it starts
a fixed number of workers (one per available CPU) that read tasks from the channel
and process each rectangle independently. Finally, it waits for all workers
to finish and returns the completed image.
*/
func JuliaParallel(f ComplexFunc, n int) image.Image {
	bounds := image.Rect(-n/2, -n/2, n/2, n/2)
	img := image.NewRGBA(bounds)
	s := float64(n / 4)

	// Divide the image into smaller rectangles.
	type task struct {
		bounds image.Rectangle
		done   chan struct{}
	}
	tasks := make(chan task)
	go func() {
		for i := bounds.Min.X; i < bounds.Max.X; i += 32 {
			for j := bounds.Min.Y; j < bounds.Max.Y; j += 32 {
				tasks <- task{image.Rect(i, j, i+32, j+32), make(chan struct{})}
			}
		}
		close(tasks)
	}()

	// Start workers to process the rectangles. Each
	var wg sync.WaitGroup
	for i := 0; i < runtime.NumCPU(); i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for t := range tasks {
				processRectangle(f, img, s, t.bounds)
				close(t.done)
			}
		}()
	}

	// Wait for all workers to finish.
	wg.Wait()

	return img
}

// Process a rectangle of the image.
func processRectangle(f ComplexFunc, img *image.RGBA, s float64, bounds image.Rectangle) {
	for i := bounds.Min.X; i < bounds.Max.X; i++ {
		for j := bounds.Min.Y; j < bounds.Max.Y; j++ {
			n := Iterate(f, complex(float64(i)/s, float64(j)/s), 256)
			r := uint8(0)
			g := uint8(0)
			b := uint8(n % 32 * 8)
			img.Set(i, j, color.RGBA{r, g, b, 255})
		}
	}
}
