package main

import (
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"os"
	"strconv"

	"code.google.com/p/x-go-binding/ui"
	"code.google.com/p/x-go-binding/ui/x11"
)

func openDisplay() (window ui.Window) {
	window, err := x11.NewWindow()
	if err != nil {
		panic(fmt.Sprintf("%s", err))
	}
	return window
}

func displayImage(buffer image.Image, window ui.Window) {
	draw.Draw(window.Screen(), window.Screen().Bounds(), buffer, image.ZP, draw.Src)
	window.FlushImage()
}

func eventLoop(window ui.Window) bool {
	for e := range window.EventChan() {
		switch e := e.(type) {
		case ui.KeyEvent:
			if e.Key == ' ' {
				return true
			}
			if e.Key == 'R' {
				return false
			}
			continue
		case ui.ConfigEvent:
			return false
		case ui.MouseEvent:
			// TODO; zoom in on click locations
			continue
		}
		return false
	}
	return true
}

func iterate(x, y float64, iters int) int {
	var zx, zy float64 = 0, 0
	for iter := 1; iter <= iters; iter += 1 {
		zx, zy = zx*zx-zy*zy+x, 2*zx*zy+y
		if (zx*zx + zy*zy) > 1000 {
			return iter
		}
	}
	return 0
}

func calcMand(buffer *image.RGBA, x0, y0, x1, y1 float64) {
	bounds := buffer.Bounds()
	dx := (x1 - x0) / float64(bounds.Dx())
	dy := (y1 - y0) / float64(bounds.Dy())

	for y, my := bounds.Min.Y, y0; y < bounds.Max.Y; y, my = y+1, my+dy {
		for x, mx := bounds.Min.X, x0; x < bounds.Max.X; x, mx = x+1, mx+dx {
			maxiters := 1000
			iters := iterate(mx, my, maxiters)
			valR := uint8(255 - (230 / 7) * (iters % 7))
			valG := uint8(255 - (230 / 11) * (iters % 11))
			valB := uint8(255 - (230 / 13) * (iters % 13))
			c := color.RGBA{valR, valG, valB, 0xff}
			if iters == 0 {
				c = color.RGBA{0, 0, 0, 0}
			}
			buffer.Set(x, y, c)
		}
	}
}

func main() {
	if len(os.Args) != 4 {
		fmt.Printf("Usage: %s <x> <y> <zoom>\n", os.Args[0])
		return
	}

	center_x, err := strconv.ParseFloat(os.Args[1], 64)
	if err != nil {
		fmt.Println(os.Args[1], err)
		return
	}
	center_y, err := strconv.ParseFloat(os.Args[2], 64)
	if err != nil {
		fmt.Println(os.Args[2], err)
		return
	}
	zoom, err := strconv.ParseFloat(os.Args[3], 64)
	if err != nil {
		fmt.Println(os.Args[3], err)
		return
	}

	window := openDisplay()

	for ;; {
		bounds := window.Screen().Bounds()
		buffer := image.NewRGBA(image.Rect(0, 0, bounds.Dx(), bounds.Dy()))

		fmt.Println("Recalculating")
		calcMand(buffer, center_x-zoom, center_y-zoom, center_x+zoom, center_y+zoom)
		fmt.Println("Done")
		displayImage(buffer, window)
		exit := eventLoop(window)
		if exit {
			break
		}
	}
}
