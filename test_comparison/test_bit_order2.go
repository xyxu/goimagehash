package main

import (
	"fmt"
)

func main() {
	// Test with pattern 1010
	bits := []bool{true, false, true, false}
	
	// MSB first: 1010 = 0xa
	binaryStrMSB := ""
	for _, b := range bits {
		if b {
			binaryStrMSB += "1"
		} else {
			binaryStrMSB += "0"
		}
	}
	
	// LSB first: 0101 = 0x5  
	binaryStrLSB := ""
	for i := len(bits)-1; i >= 0; i-- {
		if bits[i] {
			binaryStrLSB += "1"
		} else {
			binaryStrLSB += "0"
		}
	}
	
	fmt.Printf("Bits: %v\n", bits)
	fmt.Printf("MSB first: %s = 0x%x\n", binaryStrMSB, uint64(0b1010))
	fmt.Printf("LSB first: %s = 0x%x\n", binaryStrLSB, uint64(0b0101))
	
	// Now test Go's leftShiftSet
	fmt.Println("\nGo leftShiftSet simulation:")
	var hash uint64 = 0
	// Set bits at positions 63, 61 (simulating idx=0, idx=2 with 64-idx-1)
	hash |= 1 << uint(63)  // idx=0 -> 64-0-1=63
	hash |= 1 << uint(61)  // idx=2 -> 64-2-1=61
	fmt.Printf("Hash: %064b = 0x%016x\n", hash, hash)
	
	// What if we set bits in opposite order?
	var hash2 uint64 = 0
	hash2 |= 1 << uint(61)  // bit 61
	hash2 |= 1 << uint(63)  // bit 63
	fmt.Printf("Hash2: %064b = 0x%016x\n", hash2, hash2)
}
