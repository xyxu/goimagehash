package main

import (
	"fmt"
	"github.com/corona10/goimagehash"
	"image/jpeg"
	"os"
)

func main() {
	file1, _ := os.Open("../_examples/sample1.jpg")
	defer file1.Close()

	img1, _ := jpeg.Decode(file1)

	// Test whash
	whash1, _ := goimagehash.WaveletHash(img1)
	fmt.Printf("WHash (sample1.jpg): %s\n", whash1.ToString())

	// Also test with ExtWaveletHash for comparison
	extwhash1, _ := goimagehash.ExtWaveletHash(img1, 8, 8)
	fmt.Printf("ExtWHash 8x8 (sample1.jpg): %s\n", extwhash1.ToString())
}
