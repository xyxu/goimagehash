package main

import (
	"encoding/json"
	"fmt"
	"os"
)

// outputText prints hash results in text format
func outputText(hashResult *HashResult, cropResult *CropResistantResult) {
	if cropResult != nil {
		fmt.Printf("%s: %d segments\n", cropResult.Type, cropResult.Segments)
		for i, seg := range cropResult.SegmentHashes {
			fmt.Printf("  segment%d: %s (bits: %d)\n", i+1, seg.Value, seg.Bits)
		}
	} else if hashResult != nil {
		if hashResult.Type == "colorhash" {
			fmt.Printf("%s: %s (bits: %d, binbits: %d)\n",
				hashResult.Type, hashResult.Value, hashResult.Bits, hashResult.Binbits)
		} else {
			fmt.Printf("%s: %s (bits: %d)\n", hashResult.Type, hashResult.Value, hashResult.Bits)
		}
	}
}

// outputAllText prints all hashes in text format
func outputAllText(hashes map[string]HashResult, cropResult *CropResistantResult) {
	fmt.Println("Computing all hashes:")
	fmt.Println()

	// Define consistent order for output
	hashOrder := []string{"ahash", "phash", "dhash", "whash", "colorhash"}

	for _, hashType := range hashOrder {
		if hash, exists := hashes[hashType]; exists {
			if hash.Type == "colorhash" {
				fmt.Printf("%s: %s (bits: %d, binbits: %d)\n",
					hash.Type, hash.Value, hash.Bits, hash.Binbits)
			} else {
				fmt.Printf("%s: %s (bits: %d)\n", hash.Type, hash.Value, hash.Bits)
			}
		}
	}

	// Output cropresistant hash if available
	if cropResult != nil {
		fmt.Printf("\n%s: %d segments\n", cropResult.Type, cropResult.Segments)
		for i, seg := range cropResult.SegmentHashes {
			fmt.Printf("  segment%d: %s (bits: %d)\n", i+1, seg.Value, seg.Bits)
		}
	}
}

// outputJSON prints hash results in JSON format
func outputJSON(imagePath string, hashResult *HashResult, cropResult *CropResistantResult) {
	outputData := OutputData{
		Image: imagePath,
	}

	if cropResult != nil {
		outputData.Crop = cropResult
	} else if hashResult != nil {
		outputData.Hash = hashResult
	}

	jsonBytes, err := json.MarshalIndent(outputData, "", "  ")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error encoding JSON: %v\n", err)
		os.Exit(1)
	}

	fmt.Println(string(jsonBytes))
}

// outputAllJSON prints all hashes in JSON format
func outputAllJSON(imagePath string, hashes map[string]HashResult, cropResult *CropResistantResult) {
	outputData := OutputData{
		Image:  imagePath,
		Hashes: hashes,
		Crop:   cropResult,
	}

	jsonBytes, err := json.MarshalIndent(outputData, "", "  ")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error encoding JSON: %v\n", err)
		os.Exit(1)
	}

	fmt.Println(string(jsonBytes))
}
