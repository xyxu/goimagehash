// Copyright 2017 The goimagehash Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package goimagehash

import (
	"errors"
	"image"
	"math"
	"math/bits"
	"sync"

	"github.com/corona10/goimagehash/etcs"
	"github.com/corona10/goimagehash/transforms"
	"github.com/nfnt/resize"
)

// AverageHash function returns a hash computation of average hash.
// Implementation follows
// http://www.hackerfactor.com/blog/index.php?/archives/432-Looks-Like-It.html
func AverageHash(img image.Image) (*ImageHash, error) {
	if img == nil {
		return nil, errors.New("image object can not be nil")
	}

	// Create 64bits hash.
	ahash := NewImageHash(0, AHash)
	resized := resize.Resize(8, 8, img, resize.Bilinear)
	pixels := transforms.Rgb2Gray(resized)
	flattens := transforms.FlattenPixels(pixels, 8, 8)
	avg := etcs.MeanOfPixels(flattens)

	for idx, p := range flattens {
		if p > avg {
			ahash.leftShiftSet(len(flattens) - idx - 1)
		}
	}

	return ahash, nil
}

// DifferenceHash function returns a hash computation of difference hash.
// Implementation follows
// http://www.hackerfactor.com/blog/?/archives/529-Kind-of-Like-That.html
func DifferenceHash(img image.Image) (*ImageHash, error) {
	if img == nil {
		return nil, errors.New("image object can not be nil")
	}

	dhash := NewImageHash(0, DHash)
	resized := resize.Resize(9, 8, img, resize.Bilinear)
	pixels := transforms.Rgb2Gray(resized)
	idx := 0
	for i := 0; i < len(pixels); i++ {
		for j := 0; j < len(pixels[i])-1; j++ {
			if pixels[i][j] < pixels[i][j+1] {
				dhash.leftShiftSet(64 - idx - 1)
			}
			idx++
		}
	}

	return dhash, nil
}

// PerceptionHash function returns a hash computation of phash.
// Implementation follows
// http://www.hackerfactor.com/blog/index.php?/archives/432-Looks-Like-It.html
func PerceptionHash(img image.Image) (*ImageHash, error) {
	if img == nil {
		return nil, errors.New("image object can not be nil")
	}

	phash := NewImageHash(0, PHash)
	resized := resize.Resize(64, 64, img, resize.Bilinear)

	pixels := pixelPool64.Get().(*[]float64)

	transforms.Rgb2GrayFast(resized, pixels)
	flattens := transforms.DCT2DFast64(pixels)

	pixelPool64.Put(pixels)

	median := etcs.MedianOfPixelsFast64(flattens[:])

	for idx, p := range flattens {
		if p > median {
			phash.leftShiftSet(64 - idx - 1) // leftShiftSet
		}
	}

	return phash, nil
}

func WaveletHash(img image.Image) (*ImageHash, error) {
	if img == nil {
		return nil, errors.New("image object can not be nil")
	}
	bounds := img.Bounds()
	if transforms.Min(bounds.Max.X, bounds.Max.Y) < 8 {
		return nil, errors.New("image width and height should be more than 8")
	}
	whash := NewImageHash(0, WHash)

	imgScale := transforms.Floorp2(transforms.Min(bounds.Max.X, bounds.Max.Y))
	resized := resize.Resize(imgScale, imgScale, img, resize.Bilinear)
	pixels := transforms.Rgb2Gray(resized)
	maxlevel := bits.Len(imgScale) - 1
	transforms.DWT2D(pixels, maxlevel)
	pixels[0][0] = 0.0
	transforms.IDWT2D(pixels, maxlevel)
	transforms.DWT2D(pixels, maxlevel-3)
	flattens := transforms.FlattenPixels(pixels, 8, 8)
	median := etcs.MedianOfPixelsFast64(flattens[:])
	for idx, p := range flattens {
		if p > median {
			whash.leftShiftSet(64 - idx - 1)
		}
	}
	return whash, nil
}

var pixelPool64 = sync.Pool{
	New: func() interface{} {
		p := make([]float64, 4096)
		return &p
	},
}

// ExtPerceptionHash function returns phash of which the size can be set larger than uint64
// Some variable name refer to https://github.com/JohannesBuchner/imagehash/blob/master/imagehash/__init__.py
// Support 64bits phash (width=8, height=8) and 256bits phash (width=16, height=16)
// Important: width * height should be the power of 2
func ExtPerceptionHash(img image.Image, width, height int) (*ExtImageHash, error) {
	imgSize := width * height
	if img == nil {
		return nil, errors.New("image object can not be nil")
	}
	if imgSize <= 0 || imgSize&(imgSize-1) != 0 {
		return nil, errors.New("width * height should be power of 2")
	}
	var phash []uint64
	resized := resize.Resize(uint(imgSize), uint(imgSize), img, resize.Bilinear)
	pixels := transforms.Rgb2Gray(resized)
	dct := transforms.DCT2D(pixels, imgSize, imgSize)
	flattens := transforms.FlattenPixels(dct, width, height)
	median := etcs.MedianOfPixels(flattens)

	lenOfUnit := 64
	if imgSize%lenOfUnit == 0 {
		phash = make([]uint64, imgSize/lenOfUnit)
	} else {
		phash = make([]uint64, imgSize/lenOfUnit+1)
	}
	for idx, p := range flattens {
		indexOfArray := idx / lenOfUnit
		indexOfBit := lenOfUnit - idx%lenOfUnit - 1
		if p > median {
			phash[indexOfArray] |= 1 << uint(indexOfBit)
		}
	}
	return NewExtImageHash(phash, PHash, imgSize), nil
}

// ExtAverageHash function returns ahash of which the size can be set larger than uint64
// Support 64bits ahash (width=8, height=8) and 256bits ahash (width=16, height=16)
func ExtAverageHash(img image.Image, width, height int) (*ExtImageHash, error) {
	if img == nil {
		return nil, errors.New("image object can not be nil")
	}
	var ahash []uint64
	imgSize := width * height

	resized := resize.Resize(uint(width), uint(height), img, resize.Bilinear)
	pixels := transforms.Rgb2Gray(resized)
	flattens := transforms.FlattenPixels(pixels, width, height)
	avg := etcs.MeanOfPixels(flattens)

	lenOfUnit := 64
	if imgSize%lenOfUnit == 0 {
		ahash = make([]uint64, imgSize/lenOfUnit)
	} else {
		ahash = make([]uint64, imgSize/lenOfUnit+1)
	}
	for idx, p := range flattens {
		indexOfArray := idx / lenOfUnit
		indexOfBit := lenOfUnit - idx%lenOfUnit - 1
		if p > avg {
			ahash[indexOfArray] |= 1 << uint(indexOfBit)
		}
	}
	return NewExtImageHash(ahash, AHash, imgSize), nil
}

// ExtDifferenceHash function returns dhash of which the size can be set larger than uint64
// Support 64bits dhash (width=8, height=8) and 256bits dhash (width=16, height=16)
func ExtDifferenceHash(img image.Image, width, height int) (*ExtImageHash, error) {
	if img == nil {
		return nil, errors.New("image object can not be nil")
	}

	var dhash []uint64
	imgSize := width * height

	resized := resize.Resize(uint(width)+1, uint(height), img, resize.Bilinear)
	pixels := transforms.Rgb2Gray(resized)

	lenOfUnit := 64
	if imgSize%lenOfUnit == 0 {
		dhash = make([]uint64, imgSize/lenOfUnit)
	} else {
		dhash = make([]uint64, imgSize/lenOfUnit+1)
	}
	idx := 0
	for i := 0; i < len(pixels); i++ {
		for j := 0; j < len(pixels[i])-1; j++ {
			indexOfArray := idx / lenOfUnit
			indexOfBit := lenOfUnit - idx%lenOfUnit - 1
			if pixels[i][j] < pixels[i][j+1] {
				dhash[indexOfArray] |= 1 << uint(indexOfBit)
			}
			idx++
		}
	}
	return NewExtImageHash(dhash, DHash, imgSize), nil
}

// ExtWaveletHash function returns whash of which the size can be set larger than uint64
// Support 64bits whash (width=8, height=8) and 256bits whash (width=16, height=16)
// Important: width * height should be the power of 2
func ExtWaveletHash(img image.Image, width, height int) (*ExtImageHash, error) {
	if img == nil {
		return nil, errors.New("image object can not be nil")
	}
	imgSize := width * height
	if imgSize <= 0 || imgSize&(imgSize-1) != 0 {
		return nil, errors.New("width * height should be power of 2")
	}
	var whash []uint64
	bounds := img.Bounds()
	imgScale := transforms.Floorp2(transforms.Min(bounds.Max.X, bounds.Max.Y))
	resized := resize.Resize(imgScale, imgScale, img, resize.Bilinear)
	pixels := transforms.Rgb2Gray(resized)
	maxlevel := bits.Len(imgScale) - 1
	hashlevel := bits.Len(uint(math.Sqrt(float64(imgSize)))) - 1
	transforms.DWT2D(pixels, maxlevel)
	pixels[0][0] = 0.0
	transforms.IDWT2D(pixels, maxlevel)
	transforms.DWT2D(pixels, maxlevel-hashlevel)
	flattens := transforms.FlattenPixels(pixels, width, height)
	median := etcs.MedianOfPixelsFast64(flattens[:])
	lenOfUnit := 64
	if imgSize%lenOfUnit == 0 {
		whash = make([]uint64, imgSize/lenOfUnit)
	} else {
		whash = make([]uint64, imgSize/lenOfUnit+1)
	}
	for idx, p := range flattens {
		indexOfArray := idx / lenOfUnit
		indexOfBit := lenOfUnit - idx%lenOfUnit - 1
		if p > median {
			whash[indexOfArray] |= 1 << uint(indexOfBit)
		}
	}
	return NewExtImageHash(whash, WHash, imgSize), nil
}

// ColorHash function returns a hash computation of color hash.
// Implementation follows Python imagehash library.
// It computes fractions of image in intensity, hue and saturation bins:
// - the first binbits encode the black fraction of the image
// - the next binbits encode the gray fraction of the remaining image (low saturation)
// - the next 6*binbits encode the fraction in 6 bins of saturation, for highly saturated parts
// - the next 6*binbits encode the fraction in 6 bins of saturation, for mildly saturated parts
func ColorHash(img image.Image, binbits int) (*ExtImageHash, error) {
	if img == nil {
		return nil, errors.New("image object can not be nil")
	}
	if binbits <= 0 {
		return nil, errors.New("binbits must be greater than 0")
	}

	h, s, _ := transforms.RGBToHSV(img)
	height := len(h)
	width := len(h[0])
	totalPixels := float64(height * width)

	intensity := transforms.GetIntensity(img)

	maskBlack := make([][]bool, height)
	maskGray := make([][]bool, height)
	maskColors := make([][]bool, height)
	maskFaintColors := make([][]bool, height)
	maskBrightColors := make([][]bool, height)

	countBlack := 0
	countGray := 0
	countColors := 0
	countFaintColors := 0
	countBrightColors := 0

	blackThreshold := 256 / 8     // 32
	grayThreshold := 256 / 3      // 85
	colorThreshold := 256 * 2 / 3 // 170

	for i := 0; i < height; i++ {
		maskBlack[i] = make([]bool, width)
		maskGray[i] = make([]bool, width)
		maskColors[i] = make([]bool, width)
		maskFaintColors[i] = make([]bool, width)
		maskBrightColors[i] = make([]bool, width)
		for j := 0; j < width; j++ {
			intVal := intensity[i][j]
			satVal := s[i][j]

			isBlack := intVal < float64(blackThreshold)
			isGray := satVal < float64(grayThreshold)

			maskBlack[i][j] = isBlack
			if isBlack {
				countBlack++
			} else if isGray {
				maskGray[i][j] = true
				countGray++
			} else {
				maskColors[i][j] = true
				countColors++

				if satVal < float64(colorThreshold) {
					maskFaintColors[i][j] = true
					countFaintColors++
				} else {
					maskBrightColors[i][j] = true
					countBrightColors++
				}
			}
		}
	}

	fracBlack := float64(countBlack) / totalPixels
	fracGray := float64(countGray) / totalPixels

	hueBins := 6
	hueBinEdges := make([]float64, hueBins+1)
	for i := 0; i <= hueBins; i++ {
		hueBinEdges[i] = float64(i) * 256.0 / float64(hueBins)
	}

	hFaintCounts := make([]int, hueBins)
	hBrightCounts := make([]int, hueBins)

	for i := 0; i < height; i++ {
		for j := 0; j < width; j++ {
			hueVal := h[i][j]
			if maskFaintColors[i][j] {
				bin := 0
				for k := 0; k < hueBins; k++ {
					if hueVal >= hueBinEdges[k] && hueVal < hueBinEdges[k+1] {
						bin = k
						break
					}
				}
				hFaintCounts[bin]++
			}
			if maskBrightColors[i][j] {
				bin := 0
				for k := 0; k < hueBins; k++ {
					if hueVal >= hueBinEdges[k] && hueVal < hueBinEdges[k+1] {
						bin = k
						break
					}
				}
				hBrightCounts[bin]++
			}
		}
	}

	maxValue := 1 << binbits
	numBins := 2 + hueBins*2

	values := make([]int, numBins)
	values[0] = min(maxValue-1, int(fracBlack*float64(maxValue)))
	values[1] = min(maxValue-1, int(fracGray*float64(maxValue)))

	c := float64(countColors)
	if c == 0 {
		c = 1
	}

	idx := 2
	for _, count := range hFaintCounts {
		values[idx] = min(maxValue-1, int(float64(count)/c*float64(maxValue)))
		idx++
	}
	for _, count := range hBrightCounts {
		values[idx] = min(maxValue-1, int(float64(count)/c*float64(maxValue)))
		idx++
	}

	bitLen := numBins * binbits
	hash := make([]uint64, (bitLen+63)/64)

	for binIdx, val := range values {
		for bit := 0; bit < binbits; bit++ {
			overallBit := binIdx*binbits + bit
			if val&(1<<(binbits-bit-1)) != 0 {
				hashIndex := overallBit / 64
				bitIndex := 63 - (overallBit % 64)
				hash[hashIndex] |= 1 << uint(bitIndex)
			}
		}
	}

	return NewExtImageHash(hash, CHash, bitLen), nil
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
