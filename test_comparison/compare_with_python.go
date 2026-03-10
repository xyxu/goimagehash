package main

import (
	"encoding/json"
	"fmt"
	"image/jpeg"
	"os"
	"strings"

	"github.com/corona10/goimagehash"
)

type PythonHashes struct {
	Image  string            `json:"image"`
	Hashes map[string]string `json:"hashes"`
}

func loadPythonHashes() (map[string]PythonHashes, error) {
	data, err := os.ReadFile("python_reference_hashes.json")
	if err != nil {
		return nil, err
	}

	var result map[string]PythonHashes
	err = json.Unmarshal(data, &result)
	return result, err
}

func compareHash(goHash, pyHash string, algo string) {
	if goHash == pyHash {
		fmt.Printf("  ✓ %s: MATCH\n", algo)
		fmt.Printf("    Go:  %s\n", goHash)
	} else {
		fmt.Printf("  ✗ %s: MISMATCH\n", algo)
		fmt.Printf("    Go:  %s\n", goHash)
		fmt.Printf("    Py:  %s\n", pyHash)

		// For color hash, also compare lengths
		if algo == "color" || algo == "color_binbits1" {
			fmt.Printf("    Go len: %d, Py len: %d\n", len(goHash), len(pyHash))
		}
	}
}

func main() {
	// Load Python reference hashes
	pyHashes, err := loadPythonHashes()
	if err != nil {
		fmt.Printf("Error loading Python hashes: %v\n", err)
		return
	}

	// Test with sample1.jpg
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

	pyRef := pyHashes["sample1.jpg"]

	fmt.Println("Comparing hashes for sample1.jpg:")
	fmt.Println("=================================")

	// Average hash
	goAhash, _ := goimagehash.AverageHash(img1)
	compareHash(goAhash.ToString(), pyRef.Hashes["average"], "average")

	// Difference hash
	goDhash, _ := goimagehash.DifferenceHash(img1)
	compareHash(goDhash.ToString(), pyRef.Hashes["difference"], "difference")

	// Perception hash
	goPhash, _ := goimagehash.PerceptionHash(img1)
	compareHash(goPhash.ToString(), pyRef.Hashes["perception"], "perception")

	// Wavelet hash
	goWhash, _ := goimagehash.WaveletHash(img1)
	compareHash(goWhash.ToString(), pyRef.Hashes["wavelet"], "wavelet")

	// Color hash (binbits=3)
	goColorhash, _ := goimagehash.ColorHash(img1, 3)
	// Python color hash output needs to be converted to match Go format
	// Python: "07000000000" (11 hex chars = 44 bits)
	// Go: "c:1c00000000000000" (prefix + 16 hex chars = 64 bits)
	// Extract just the hex part after "c:"
	goColorStr := goColorhash.ToString()
	if strings.HasPrefix(goColorStr, "c:") {
		goColorStr = goColorStr[2:]
	}
	// Python hash is shorter, pad or truncate for comparison
	pyColor := pyRef.Hashes["color"]
	// For now, just compare
	fmt.Printf("  ? color (binbits=3):\n")
	fmt.Printf("    Go:  %s (full: %s)\n", goColorStr, goColorhash.ToString())
	fmt.Printf("    Py:  %s\n", pyColor)

	// Color hash (binbits=1)
	goColorhash1, _ := goimagehash.ColorHash(img1, 1)
	goColorStr1 := goColorhash1.ToString()
	if strings.HasPrefix(goColorStr1, "c:") {
		goColorStr1 = goColorStr1[2:]
	}
	pyColor1 := pyRef.Hashes["color_binbits1"]
	fmt.Printf("  ? color (binbits=1):\n")
	fmt.Printf("    Go:  %s (full: %s)\n", goColorStr1, goColorhash1.ToString())
	fmt.Printf("    Py:  %s\n", pyColor1)

	// Crop-resistant hash
	goCrophash, _ := goimagehash.CropResistantHash(img1, nil, 128, 500, 300)
	// Build string similar to Python format
	var goCropStr string
	segments := goCrophash.GetSegmentHashes()
	for i, seg := range segments {
		if i > 0 {
			goCropStr += ","
		}
		segStr := seg.ToString()
		if strings.HasPrefix(segStr, "d:") {
			segStr = segStr[2:]
		}
		goCropStr += segStr
	}
	pyCrop := pyRef.Hashes["crop_resistant"]
	fmt.Printf("  ? crop_resistant:\n")
	fmt.Printf("    Go:  %s\n", goCropStr)
	fmt.Printf("    Py:  %s\n", pyCrop)
	fmt.Printf("    Go segments: %d, Py segments: %d\n", len(segments), len(strings.Split(pyCrop, ",")))
}
