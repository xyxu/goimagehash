package main

import (
	"fmt"
	"image/jpeg"
	"os"
	"strings"

	"github.com/corona10/goimagehash"
)

func main() {
	file1, err := os.Open("../_examples/sample1.jpg")
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}
	defer file1.Close()

	img1, err := jpeg.Decode(file1)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}

	// Average hash
	ahash, _ := goimagehash.AverageHash(img1)
	hashStr := ahash.ToString()
	if strings.Contains(hashStr, ":") {
		parts := strings.Split(hashStr, ":")
		hashStr = parts[1]
	}

	fmt.Println("Go average hash hex:", hashStr)
	fmt.Print("Go average hash bits: ")

	for i := 0; i < len(hashStr); i++ {
		hexChar := hashStr[i]
		var val byte
		if hexChar >= '0' && hexChar <= '9' {
			val = hexChar - '0'
		} else if hexChar >= 'a' && hexChar <= 'f' {
			val = hexChar - 'a' + 10
		} else if hexChar >= 'A' && hexChar <= 'F' {
			val = hexChar - 'A' + 10
		}

		for j := 3; j >= 0; j-- {
			if (val>>j)&1 == 1 {
				fmt.Print("1")
			} else {
				fmt.Print("0")
			}
		}
	}
	fmt.Println()

	// Python bits for comparison
	fmt.Println("\nPython average hash bits: 0111111011111111001111110000001100000111000000111100000111111100")

	// Let's compare bit by bit
	fmt.Println("\nBit comparison (Go vs Python):")
	goBits := ""
	for i := 0; i < len(hashStr); i++ {
		hexChar := hashStr[i]
		var val byte
		if hexChar >= '0' && hexChar <= '9' {
			val = hexChar - '0'
		} else if hexChar >= 'a' && hexChar <= 'f' {
			val = hexChar - 'a' + 10
		} else if hexChar >= 'A' && hexChar <= 'F' {
			val = hexChar - 'A' + 10
		}

		for j := 3; j >= 0; j-- {
			if (val>>j)&1 == 1 {
				goBits += "1"
			} else {
				goBits += "0"
			}
		}
	}

	pyBits := "0111111011111111001111110000001100000111000000111100000111111100"

	differences := 0
	for i := 0; i < len(goBits) && i < len(pyBits); i++ {
		if goBits[i] != pyBits[i] {
			differences++
			row := i / 8
			col := i % 8
			fmt.Printf("  Bit %d (row %d, col %d): Go=%c, Py=%c\n", i, row, col, goBits[i], pyBits[i])
		}
	}
	fmt.Printf("\nTotal differences: %d/64\n", differences)
}
