package main

import (
	"fmt"
	"image/jpeg"
	"os"
	"strings"

	"github.com/corona10/goimagehash"
	"github.com/corona10/goimagehash/transforms"
)

func printBits(hash string) {
	// Remove prefix if present
	if strings.Contains(hash, ":") {
		parts := strings.Split(hash, ":")
		hash = parts[1]
	}

	// Convert hex to binary
	for i := 0; i < len(hash); i++ {
		hexChar := hash[i]
		var val byte
		if hexChar >= '0' && hexChar <= '9' {
			val = hexChar - '0'
		} else if hexChar >= 'a' && hexChar <= 'f' {
			val = hexChar - 'a' + 10
		} else if hexChar >= 'A' && hexChar <= 'F' {
			val = hexChar - 'A' + 10
		}

		// Print 4 bits
		for j := 3; j >= 0; j-- {
			if (val>>j)&1 == 1 {
				fmt.Print("1")
			} else {
				fmt.Print("0")
			}
		}
	}
	fmt.Println()
}

func main() {
	// Open image
	file1, err := os.Open("../_examples/sample1.jpg")
	if err != nil {
		fmt.Printf("Error opening image: %v\n", err)
		return
	}
	defer file1.Close()

	img1, err := jpeg.Decode(file1)
	if err != nil {
		fmt.Printf("Error decoding image: %v\n", err)
		return
	}

	fmt.Println("=== Go goimagehash debug ===")
	bounds := img1.Bounds()
	fmt.Printf("Image size: %dx%d\n", bounds.Dx(), bounds.Dy())

	// Debug grayscale conversion
	fmt.Println("\n--- Grayscale pixels (8x8) ---")
	// Manually compute 8x8 grayscale
	grayPixels := transforms.Rgb2Gray(img1)
	// Resize to 8x8 (simplified - actual implementation does more)
	fmt.Println("Note: Go's Rgb2Gray returns full size, not resized")

	// Compute average hash
	fmt.Println("\n--- Average Hash ---")
	ahash, err := goimagehash.AverageHash(img1)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}
	fmt.Printf("Hash: %s\n", ahash.ToString())
	fmt.Printf("Hash bits: ")
	printBits(ahash.ToString())

	// Compute difference hash
	fmt.Println("\n--- Difference Hash ---")
	dhash, err := goimagehash.DifferenceHash(img1)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}
	fmt.Printf("Hash: %s\n", dhash.ToString())
	fmt.Printf("Hash bits: ")
	printBits(dhash.ToString())

	// Compute perception hash
	fmt.Println("\n--- Perception Hash ---")
	phash, err := goimagehash.PerceptionHash(img1)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}
	fmt.Printf("Hash: %s\n", phash.ToString())
	fmt.Printf("Hash bits: ")
	printBits(phash.ToString())

	// Compute wavelet hash
	fmt.Println("\n--- Wavelet Hash ---")
	whash, err := goimagehash.WaveletHash(img1)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}
	fmt.Printf("Hash: %s\n", whash.ToString())
	fmt.Printf("Hash bits: ")
	printBits(whash.ToString())

	// Compute color hash
	fmt.Println("\n--- Color Hash (binbits=3) ---")
	colorhash, err := goimagehash.ColorHash(img1, 3)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}
	fmt.Printf("Hash: %s\n", colorhash.ToString())
	colorStr := colorhash.ToString()
	if strings.HasPrefix(colorStr, "c:") {
		colorStr = colorStr[2:]
	}
	fmt.Printf("Hash hex (clean): %s\n", colorStr)
	fmt.Printf("Length: %d\n", len(colorStr))
}
