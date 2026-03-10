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

	colorhash1, _ := goimagehash.ColorHash(img1, 3)
	fmt.Printf("ColorHash (binbits=3): %s\n", colorhash1.ToString())
	fmt.Printf("Hash bits: %d\n", colorhash1.Bits())

	// Let's also check the raw hash values
	fmt.Printf("Raw hash: %v\n", colorhash1.GetHash())

	// Test with binbits=1
	colorhash2, _ := goimagehash.ColorHash(img1, 1)
	fmt.Printf("\nColorHash (binbits=1): %s\n", colorhash2.ToString())
	fmt.Printf("Hash bits: %d\n", colorhash2.Bits())
	fmt.Printf("Raw hash: %v\n", colorhash2.GetHash())

	// Test with binbits=2
	colorhash3, _ := goimagehash.ColorHash(img1, 2)
	fmt.Printf("\nColorHash (binbits=2): %s\n", colorhash3.ToString())
	fmt.Printf("Hash bits: %d\n", colorhash3.Bits())
}
