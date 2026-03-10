package main

import (
	"fmt"
)

func main() {
	pixels := [][]float64{
		{10, 20, 30, 40},
		{50, 60, 70, 80},
		{90, 100, 110, 120},
	}
	
	fmt.Println("Go dhash comparison: pixels[i][j] < pixels[i][j+1]")
	bits := []int{}
	for i := 0; i < len(pixels); i++ {
		for j := 0; j < len(pixels[i])-1; j++ {
			if pixels[i][j] < pixels[i][j+1] {
				bits = append(bits, 1)
			} else {
				bits = append(bits, 0)
			}
		}
	}
	
	fmt.Println("Result bits:", bits)
	
	var hexVal uint64 = 0
	for i, b := range bits {
		if b == 1 {
			hexVal |= 1 << uint(len(bits)-i-1)
		}
	}
	
	fmt.Printf("As hex: 0x%x\n", hexVal)
	
	fmt.Println("\nOpposite comparison: pixels[i][j] > pixels[i][j+1]")
	bits2 := []int{}
	for i := 0; i < len(pixels); i++ {
		for j := 0; j < len(pixels[i])-1; j++ {
			if pixels[i][j] > pixels[i][j+1] {
				bits2 = append(bits2, 1)
			} else {
				bits2 = append(bits2, 0)
			}
		}
	}
	
	fmt.Println("Result bits:", bits2)
	
	var hexVal2 uint64 = 0
	for i, b := range bits2 {
		if b == 1 {
			hexVal2 |= 1 << uint(len(bits2)-i-1)
		}
	}
	
	fmt.Printf("As hex: 0x%x\n", hexVal2)
}
