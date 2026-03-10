package main

import (
	"fmt"
)

func main() {
	// Test bit ordering
	// Simulate setting bits at positions 63, 62, ..., 0
	var hash uint64 = 0
	
	// Set bit 63 (MSB)
	hash |= 1 << uint(63)
	fmt.Printf("Set bit 63: %064b = 0x%016x\n", hash, hash)
	
	// Set bit 62
	hash |= 1 << uint(62)
	fmt.Printf("Set bit 62: %064b = 0x%016x\n", hash, hash)
	
	// Set bit 1
	hash |= 1 << uint(1)
	fmt.Printf("Set bit 1:  %064b = 0x%016x\n", hash, hash)
	
	// Set bit 0 (LSB)
	hash |= 1 << uint(0)
	fmt.Printf("Set bit 0:  %064b = 0x%016x\n", hash, hash)
	
	// Now test what Python would produce for the same bits
	// Python flattens row-major and converts to hex
	// For a 2x2 array: [[1,0],[0,1]] -> [1,0,0,1] -> binary "1001" -> hex "9"
	// But is "1001" MSB first or LSB first?
	
	fmt.Println("\nPython-style (row-major flattening):")
	// Assume array is [[true, false], [false, true]]
	bits := []bool{true, false, false, true}
	
	// Convert to binary string MSB first
	binaryStr := ""
	for _, b := range bits {
		if b {
			binaryStr += "1"
		} else {
			binaryStr += "0"
		}
	}
	
	fmt.Printf("Binary string (MSB first): %s\n", binaryStr)
	fmt.Printf("As hex: 0x%x\n", uint64(0b1001))
	
	// What if it's LSB first?
	binaryStrLSB := ""
	for i := len(bits)-1; i >= 0; i-- {
		if bits[i] {
			binaryStrLSB += "1"
		} else {
			binaryStrLSB += "0"
		}
	}
	
	fmt.Printf("Binary string (LSB first): %s\n", binaryStrLSB)
	fmt.Printf("As hex: 0x%x\n", uint64(0b1001))
}
