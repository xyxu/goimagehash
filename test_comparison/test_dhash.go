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

	// Test dhash
	dhash1, _ := goimagehash.DifferenceHash(img1)
	fmt.Printf("DHash (sample1.jpg): %s\n", dhash1.ToString())

	// Also test with ExtDifferenceHash for comparison
	extdhash1, _ := goimagehash.ExtDifferenceHash(img1, 8, 8)
	fmt.Printf("ExtDHash 8x8 (sample1.jpg): %s\n", extdhash1.ToString())
}
