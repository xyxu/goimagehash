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
			lum := pixel2Gray(r, g, b, 0)
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
// Uses Python PIL's exact formula: R*299/1000 + G*587/1000 + B*114/1000
// where R,G,B are in range 0-255
// Returns integer value 0-255 to match Python's PIL convert('L')
func pixel2Gray(r, g, b, a uint32) float64 {
	// Scale from 0-65535 to 0-255 (same as r/257)
	R := float64(r) / 257.0
	G := float64(g) / 257.0
	B := float64(b) / 257.0

	// Python PIL uses integer arithmetic: R*299/1000 + G*587/1000 + B*114/1000
	// First, scale to 0-255 integer like Python does
	R_int := int(R + 0.5) // Round to nearest integer
	G_int := int(G + 0.5)
	B_int := int(B + 0.5)

	// Apply Python's integer arithmetic formula
	// Note: Python does integer division: R*299/1000 means (R*299)//1000
	lum_int := (R_int*299 + G_int*587 + B_int*114) / 1000

	return float64(lum_int)
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
	_, hImg := bounds.Max.X-bounds.Min.X, bounds.Max.Y-bounds.Min.Y

	h = make([][]float64, hImg)
	s = make([][]float64, hImg)
	v = make([][]float64, hImg)

	switch c := colorImg.(type) {
	case *image.YCbCr:
		rgbToHSVYCbCr(c, h, s, v)
	case *image.RGBA:
		rgbToHSVRGBA(c, h, s, v)
	default:
		rgbToHSVDefault(colorImg, h, s, v)
	}

	return h, s, v
}

func rgbToHSVYCbCr(colorImg *image.YCbCr, h, s, v [][]float64) {
	hImg := len(h)
	w := len(h[0])

	for i := 0; i < hImg; i++ {
		for j := 0; j < w; j++ {
			ycbcr := colorImg.YCbCrAt(j, i)
			y, cb, cr := ycbcr.Y, ycbcr.Cb, ycbcr.Cr

			rf := float64(y)
			gf := float64(y)
			bf := float64(y)

			maxVal := rf
			minVal := rf

			v[i][j] = rf

			if cb != 128 || cr != 128 {
				rf = float64(y) + 1.402*(float64(cr)-128)
				gf = float64(y) - 0.344136*(float64(cb)-128) - 0.714136*(float64(cr)-128)
				bf = float64(y) + 1.772*(float64(cb)-128)

				maxVal = math.Max(math.Max(rf, gf), bf)
				minVal = math.Min(math.Min(rf, gf), bf)
			}

			delta := maxVal - minVal

			if maxVal > 0 {
				s[i][j] = delta / maxVal * 255
			} else {
				s[i][j] = 0
			}

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

			h[i][j] = h[i][j] / 360 * 255
		}
	}
}

func rgbToHSVRGBA(colorImg *image.RGBA, h, s, v [][]float64) {
	hImg := len(h)
	w := len(h[0])

	for i := 0; i < hImg; i++ {
		for j := 0; j < w; j++ {
			c := colorImg.RGBAAt(j, i)
			rf := float64(c.R)
			gf := float64(c.G)
			bf := float64(c.B)

			maxVal := math.Max(math.Max(rf, gf), bf)
			minVal := math.Min(math.Min(rf, gf), bf)
			delta := maxVal - minVal

			v[i][j] = maxVal

			if maxVal > 0 {
				s[i][j] = delta / maxVal * 255
			} else {
				s[i][j] = 0
			}

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

			h[i][j] = h[i][j] / 360 * 255
		}
	}
}

func rgbToHSVDefault(colorImg image.Image, h, s, v [][]float64) {
	hImg := len(h)
	w := len(h[0])

	for i := 0; i < hImg; i++ {
		for j := 0; j < w; j++ {
			r, g, b, _ := colorImg.At(j, i).RGBA()
			rf := float64(r / 257)
			gf := float64(g / 257)
			bf := float64(b / 257)

			maxVal := math.Max(math.Max(rf, gf), bf)
			minVal := math.Min(math.Min(rf, gf), bf)
			delta := maxVal - minVal

			v[i][j] = maxVal

			if maxVal > 0 {
				s[i][j] = delta / maxVal * 255
			} else {
				s[i][j] = 0
			}

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

			h[i][j] = h[i][j] / 360 * 255
		}
	}
}

// RGBToHSVAndIntensity converts an RGB image to HSV color space and intensity in a single pass.
// Returns: h, s, v, intensity
func RGBToHSVAndIntensity(colorImg image.Image) (h, s, v, intensity [][]float64) {
	bounds := colorImg.Bounds()
	w := bounds.Max.X - bounds.Min.X
	hImg := bounds.Max.Y - bounds.Min.Y

	h = make([][]float64, hImg)
	s = make([][]float64, hImg)
	v = make([][]float64, hImg)
	intensity = make([][]float64, hImg)

	for i := 0; i < hImg; i++ {
		h[i] = make([]float64, w)
		s[i] = make([]float64, w)
		v[i] = make([]float64, w)
		intensity[i] = make([]float64, w)
	}

	switch c := colorImg.(type) {
	case *image.YCbCr:
		rgbToHSVAndIntensityYCbCr(c, h, s, v, intensity)
	case *image.RGBA:
		rgbToHSVAndIntensityRGBA(c, h, s, v, intensity)
	default:
		rgbToHSVAndIntensityDefault(colorImg, h, s, v, intensity)
	}

	return h, s, v, intensity
}

func rgbToHSVAndIntensityYCbCr(colorImg *image.YCbCr, h, s, v, intensity [][]float64) {
	hImg := len(h)
	w := len(h[0])

	for i := 0; i < hImg; i++ {
		h[i] = make([]float64, w)
		s[i] = make([]float64, w)
		v[i] = make([]float64, w)
		intensity[i] = make([]float64, w)
		for j := 0; j < w; j++ {
			ycbcr := colorImg.YCbCrAt(j, i)
			y, cb, cr := ycbcr.Y, ycbcr.Cb, ycbcr.Cr

			rf := float64(y)
			gf := float64(y)
			bf := float64(y)

			maxVal := rf
			minVal := rf

			v[i][j] = rf
			intensity[i][j] = rf

			if cb != 128 || cr != 128 {
				rf = float64(y) + 1.402*(float64(cr)-128)
				gf = float64(y) - 0.344136*(float64(cb)-128) - 0.714136*(float64(cr)-128)
				bf = float64(y) + 1.772*(float64(cb)-128)

				maxVal = math.Max(math.Max(rf, gf), bf)
				minVal = math.Min(math.Min(rf, gf), bf)
			}

			delta := maxVal - minVal

			if maxVal > 0 {
				s[i][j] = delta / maxVal * 255
			} else {
				s[i][j] = 0
			}

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

			h[i][j] = h[i][j] / 360 * 255
		}
	}
}

func rgbToHSVAndIntensityRGBA(colorImg *image.RGBA, h, s, v, intensity [][]float64) {
	hImg := len(h)
	w := len(h[0])

	for i := 0; i < hImg; i++ {
		h[i] = make([]float64, w)
		s[i] = make([]float64, w)
		v[i] = make([]float64, w)
		intensity[i] = make([]float64, w)
		for j := 0; j < w; j++ {
			c := colorImg.RGBAAt(j, i)
			rf := float64(c.R)
			gf := float64(c.G)
			bf := float64(c.B)

			maxVal := math.Max(math.Max(rf, gf), bf)
			minVal := math.Min(math.Min(rf, gf), bf)
			delta := maxVal - minVal

			v[i][j] = maxVal
			intensity[i][j] = 0.299*rf + 0.587*gf + 0.114*bf

			if maxVal > 0 {
				s[i][j] = delta / maxVal * 255
			} else {
				s[i][j] = 0
			}

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

			h[i][j] = h[i][j] / 360 * 255
		}
	}
}

func rgbToHSVAndIntensityDefault(colorImg image.Image, h, s, v, intensity [][]float64) {
	hImg := len(h)
	w := len(h[0])

	for i := 0; i < hImg; i++ {
		h[i] = make([]float64, w)
		s[i] = make([]float64, w)
		v[i] = make([]float64, w)
		intensity[i] = make([]float64, w)
		for j := 0; j < w; j++ {
			r, g, b, _ := colorImg.At(j, i).RGBA()
			rf := float64(r / 257)
			gf := float64(g / 257)
			bf := float64(b / 257)

			maxVal := math.Max(math.Max(rf, gf), bf)
			minVal := math.Min(math.Min(rf, gf), bf)
			delta := maxVal - minVal

			v[i][j] = maxVal
			intensity[i][j] = 0.299*rf + 0.587*gf + 0.114*bf

			if maxVal > 0 {
				s[i][j] = delta / maxVal * 255
			} else {
				s[i][j] = 0
			}

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

			h[i][j] = h[i][j] / 360 * 255
		}
	}
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
			lum := 0.299*float64(r/257) + 0.587*float64(g/257) + 0.114*float64(b/257)
			pixels[i][j] = lum
		}
	}

	return pixels
}

// Point represents a 2D point
type Point struct {
	X int
	Y int
}

// Segment represents a region in the image
type Segment struct {
	Points []Point
}

// FindAllSegments finds all regions within an image pixel array.
// Uses watershed-like algorithm to segment bright and dark areas.
func FindAllSegments(pixels [][]float64, segmentThreshold float64, minSegmentSize int) []Segment {
	height := len(pixels)
	if height == 0 {
		return nil
	}
	width := len(pixels[0])

	thresholdPixels := make([][]bool, height)
	unassignedPixels := make([][]bool, height)

	for i := 0; i < height; i++ {
		thresholdPixels[i] = make([]bool, width)
		unassignedPixels[i] = make([]bool, width)
		for j := 0; j < width; j++ {
			thresholdPixels[i][j] = pixels[i][j] > segmentThreshold
			unassignedPixels[i][j] = true
		}
	}

	alreadySegmented := make(map[Point]bool)

	addBorderPixels(height, width, alreadySegmented)

	segments := findRegions(thresholdPixels, unassignedPixels, alreadySegmented, minSegmentSize, true)

	for i := 0; i < height; i++ {
		for j := 0; j < width; j++ {
			thresholdPixels[i][j] = !thresholdPixels[i][j]
		}
	}

	segments = append(segments, findRegions(thresholdPixels, unassignedPixels, alreadySegmented, minSegmentSize, false)...)

	return segments
}

func addBorderPixels(height int, width int, alreadySegmented map[Point]bool) {
	for i := 0; i < width; i++ {
		alreadySegmented[Point{X: i, Y: -1}] = true
		alreadySegmented[Point{X: i, Y: height}] = true
	}
	for i := 0; i < height; i++ {
		alreadySegmented[Point{X: -1, Y: i}] = true
		alreadySegmented[Point{X: width, Y: i}] = true
	}
}

func findRegions(thresholdPixels, unassignedPixels [][]bool, alreadySegmented map[Point]bool, minSegmentSize int, firstPass bool) []Segment {
	height := len(thresholdPixels)
	width := len(thresholdPixels[0])
	segments := []Segment{}

	for {
		found := false
		for i := 0; i < height; i++ {
			for j := 0; j < width; j++ {
				if thresholdPixels[i][j] && unassignedPixels[i][j] {
					segment := findRegion(i, j, thresholdPixels, unassignedPixels, alreadySegmented)
					if len(segment) >= minSegmentSize {
						points := make([]Point, len(segment))
						idx := 0
						for p := range segment {
							points[idx] = p
							idx++
						}
						segments = append(segments, Segment{Points: points})
					}
					found = true
					break
				}
			}
			if found {
				break
			}
		}
		if !found {
			break
		}
	}

	return segments
}

func findRegion(startX, startY int, thresholdPixels, unassignedPixels [][]bool, alreadySegmented map[Point]bool) map[Point]bool {
	height := len(thresholdPixels)
	width := len(thresholdPixels[0])

	segment := make(map[Point]bool)
	queue := make([]Point, 0)
	queue = append(queue, Point{X: startY, Y: startX})

	for len(queue) > 0 {
		point := queue[len(queue)-1]
		queue = queue[:len(queue)-1]

		if point.Y < 0 || point.Y >= height || point.X < 0 || point.X >= width {
			continue
		}

		if segment[point] {
			continue
		}

		if !thresholdPixels[point.Y][point.X] || !unassignedPixels[point.Y][point.X] {
			continue
		}

		if alreadySegmented[point] {
			continue
		}

		segment[point] = true
		unassignedPixels[point.Y][point.X] = false

		neighbors := []Point{
			{X: point.X - 1, Y: point.Y},
			{X: point.X + 1, Y: point.Y},
			{X: point.X, Y: point.Y - 1},
			{X: point.X, Y: point.Y + 1},
		}

		for _, n := range neighbors {
			if n.Y >= 0 && n.Y < height && n.X >= 0 && n.X < width {
				if thresholdPixels[n.Y][n.X] && unassignedPixels[n.Y][n.X] && !alreadySegmented[n] {
					queue = append(queue, n)
				}
			}
		}
	}

	return segment
}

// GetBoundingBox returns the bounding box of a segment
func GetBoundingBox(segment Segment) (minX, minY, maxX, maxY int) {
	if len(segment.Points) == 0 {
		return 0, 0, 0, 0
	}
	minX = segment.Points[0].X
	minY = segment.Points[0].Y
	maxX = segment.Points[0].X
	maxY = segment.Points[0].Y

	for _, p := range segment.Points {
		if p.X < minX {
			minX = p.X
		}
		if p.X > maxX {
			maxX = p.X
		}
		if p.Y < minY {
			minY = p.Y
		}
		if p.Y > maxY {
			maxY = p.Y
		}
	}
	return minX, minY, maxX, maxY
}
