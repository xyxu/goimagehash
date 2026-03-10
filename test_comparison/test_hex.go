package main

import (
	"fmt"
)

func main() {
	hash := uint64(2017612633061982208)
	fmt.Printf("Decimal: %d\n", hash)
	fmt.Printf("Hex: 0x%016x\n", hash)
	fmt.Printf("Binary: %064b\n", hash)
	
	// What are bits 63-22?
	fmt.Println("\nBits 63-22 (42 bits MSB):")
	for i := 63; i >= 22; i-- {
		if hash&(1<<uint(i)) != 0 {
			fmt.Printf("Bit %d is 1\n", i)
		}
	}
	
	// Extract first 3 bits (positions 63,62,61)
	first3 := (hash >> 61) & 0x7
	fmt.Printf("\nFirst 3 bits (63,62,61): %03b = %d\n", first3, first3)
}
