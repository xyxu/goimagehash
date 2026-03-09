// Copyright 2017 The goimagehash Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package transforms

import (
	"image"
	"math"
)

// Rgb2Gray function converts RGB to a gray scale array.
func Rgb2Gray(colorImg image.Image) [][]float64 {
	bounds := colorImg.Bounds()
	w, h := bounds.Max.X-bounds.Min.X, bounds.Max.Y-bounds.Min.Y
	pixels := make([][]float64, h)

	for i := range pixels {
		pixels[i] = make([]float64, w)
		for j := range pixels[i] {
			color := colorImg.At(j, i)
			r, g, b, _ := color.RGBA()
			lum := 0.299*float64(r/257) + 0.587*float64(g/257) + 0.114*float64(b/256)
			pixels[i][j] = lum
		}
	}

	return pixels
}

// Rgb2GrayFast function converts RGB to a gray scale array.
func Rgb2GrayFast(colorImg image.Image, pixels *[]float64) {
	bounds := colorImg.Bounds()
	w, h := bounds.Max.X-bounds.Min.X, bounds.Max.Y-bounds.Min.Y
	if w != h {
		return
	}
	switch c := colorImg.(type) {
	case *image.YCbCr:
		rgb2GrayYCbCR(c, *pixels, w)
	case *image.RGBA:
		rgb2GrayRGBA(c, *pixels, w)
	default:
		rgb2GrayDefault(c, *pixels, w)
	}
}

// pixel2Gray converts a pixel to grayscale value base on luminosity
func pixel2Gray(r, g, b, a uint32) float64 {
	return 0.299*float64(r/257) + 0.587*float64(g/257) + 0.114*float64(b/256)
}

// rgb2GrayDefault uses the image.Image interface
func rgb2GrayDefault(colorImg image.Image, pixels []float64, s int) {
	for i := 0; i < s; i++ {
		for j := 0; j < s; j++ {
			pixels[j+(i*s)] = pixel2Gray(colorImg.At(j, i).RGBA())
		}
	}
}

// rgb2GrayYCbCR uses *image.YCbCr which is significantly faster than the image.Image interface.
func rgb2GrayYCbCR(colorImg *image.YCbCr, pixels []float64, s int) {
	for i := 0; i < s; i++ {
		for j := 0; j < s; j++ {
			pixels[j+(i*s)] = pixel2Gray(colorImg.YCbCrAt(j, i).RGBA())
		}
	}
}

// rgb2GrayRGBA uses *image.RGBA which is significantly faster than the image.Image interface.
func rgb2GrayRGBA(colorImg *image.RGBA, pixels []float64, s int) {
	for i := 0; i < s; i++ {
		for j := 0; j < s; j++ {
			pixels[(i*s)+j] = pixel2Gray(colorImg.At(j, i).RGBA())
		}
	}
}

// FlattenPixels function flattens 2d array into 1d array.
func FlattenPixels(pixels [][]float64, x int, y int) []float64 {
	flattens := make([]float64, x*y)
	for i := 0; i < y; i++ {
		for j := 0; j < x; j++ {
			flattens[y*i+j] = pixels[i][j]
		}
	}
	return flattens
}

// FlattenPixelsFast64 function flattens 2d array into 1d array.
func FlattenPixelsFast64(pixels []float64, x int, y int) []float64 {
	flattens := [64]float64{}
	for i := 0; i < y; i++ {
		for j := 0; j < x; j++ {
			flattens[y*i+j] = pixels[(i*64)+j]
		}
	}
	return flattens[:]
}

// RGBToHSV converts an RGB image to HSV color space.
// Returns three 2D arrays: H (hue 0-360), S (saturation 0-255), V (value 0-255).
func RGBToHSV(colorImg image.Image) (h, s, v [][]float64) {
	bounds := colorImg.Bounds()
	w, hImg := bounds.Max.X-bounds.Min.X, bounds.Max.Y-bounds.Min.Y

	h = make([][]float64, hImg)
	s = make([][]float64, hImg)
	v = make([][]float64, hImg)

	for i := 0; i < hImg; i++ {
		h[i] = make([]float64, w)
		s[i] = make([]float64, w)
		v[i] = make([]float64, w)
		for j := 0; j < w; j++ {
			r, g, b, _ := colorImg.At(j, i).RGBA()
			rf := float64(r / 257)
			gf := float64(g / 257)
			bf := float64(b / 256)

			maxVal := math.Max(math.Max(rf, gf), bf)
			minVal := math.Min(math.Min(rf, gf), bf)
			delta := maxVal - minVal

			// Value (brightness)
			v[i][j] = maxVal

			// Saturation
			if maxVal > 0 {
				s[i][j] = delta / maxVal * 255
			} else {
				s[i][j] = 0
			}

			// Hue
			if delta == 0 {
				h[i][j] = 0
			} else if maxVal == rf {
				h[i][j] = 60 * (math.Mod((gf-bf)/delta, 6))
			} else if maxVal == gf {
				h[i][j] = 60 * ((bf-rf)/delta + 2)
			} else {
				h[i][j] = 60 * ((rf-gf)/delta + 4)
			}

			if h[i][j] < 0 {
				h[i][j] += 360
			}

			// Scale to 0-255
			h[i][j] = h[i][j] / 360 * 255
		}
	}

	return h, s, v
}

// GetIntensity returns grayscale intensity of an image (0-255).
func GetIntensity(colorImg image.Image) [][]float64 {
	bounds := colorImg.Bounds()
	w, h := bounds.Max.X-bounds.Min.X, bounds.Max.Y-bounds.Min.Y
	pixels := make([][]float64, h)

	for i := range pixels {
		pixels[i] = make([]float64, w)
		for j := range pixels[i] {
			color := colorImg.At(j, i)
			r, g, b, _ := color.RGBA()
			lum := 0.299*float64(r/257) + 0.587*float64(g/257) + 0.114*float64(b/256)
			pixels[i][j] = lum
		}
	}

	return pixels
}
