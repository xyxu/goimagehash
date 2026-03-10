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
	fmt.Printf("ColorHash (sample1.jpg, binbits=3): %s\n", colorhash1.ToString())

	// Also test with binbits=1 for simpler output
	colorhash2, _ := goimagehash.ColorHash(img1, 1)
	fmt.Printf("ColorHash (sample1.jpg, binbits=1): %s\n", colorhash2.ToString())
}
