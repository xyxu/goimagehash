package main

import (
	"flag"
	"fmt"
	"image"
	"image/jpeg"
	"image/png"
	"os"
	"strings"

	"github.com/corona10/goimagehash"
)

var (
	hashType     string
	binbits      int
	hashSize     int
	allHashes    bool
	outputFormat string
)

func init() {
	flag.StringVar(&hashType, "type", "ahash", "Hash type: ahash, phash, dhash, whash, colorhash, cropresistant")
	flag.IntVar(&binbits, "binbits", 3, "Bin bits for colorhash (default: 3)")
	flag.IntVar(&hashSize, "size", 8, "Hash size for ahash, phash, dhash, whash (default: 8)")
	flag.BoolVar(&allHashes, "all", false, "Compute all hash types")
	flag.StringVar(&outputFormat, "format", "text", "Output format: text, json")
}

func main() {
	flag.Parse()

	// Validate format
	if outputFormat != "text" && outputFormat != "json" {
		fmt.Fprintf(os.Stderr, "Error: invalid format '%s'. Must be 'text' or 'json'\n", outputFormat)
		os.Exit(1)
	}

	if flag.NArg() < 1 {
		fmt.Fprintf(os.Stderr, "Usage: %s [options] <image_file>\n", os.Args[0])
		flag.PrintDefaults()
		os.Exit(1)
	}

	imagePath := flag.Arg(0)

	img, err := loadImage(imagePath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error loading image: %v\n", err)
		os.Exit(1)
	}

	if allHashes {
		printAllHashes(img)
	} else {
		printHash(img, hashType)
	}
}

func loadImage(path string) (image.Image, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	img, format, err := image.Decode(file)
	if err != nil {
		return nil, err
	}

	if format == "jpeg" {
		file.Seek(0, 0)
		return jpeg.Decode(file)
	} else if format == "png" {
		file.Seek(0, 0)
		return png.Decode(file)
	}

	return img, nil
}

func printHash(img image.Image, hashType string) {
	var hashStr string

	switch strings.ToLower(hashType) {
	case "ahash", "average":
		hash, err := goimagehash.AverageHash(img)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}
		hashStr = fmt.Sprintf("%s (bits: %d)", hash.ToString(), hash.Bits())

	case "phash", "perceptual":
		hash, err := goimagehash.PerceptionHash(img)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}
		hashStr = fmt.Sprintf("%s (bits: %d)", hash.ToString(), hash.Bits())

	case "dhash", "difference":
		hash, err := goimagehash.DifferenceHash(img)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}
		hashStr = fmt.Sprintf("%s (bits: %d)", hash.ToString(), hash.Bits())

	case "whash", "wavelet":
		hash, err := goimagehash.WaveletHash(img)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}
		hashStr = fmt.Sprintf("%s (bits: %d)", hash.ToString(), hash.Bits())

	case "colorhash", "color":
		hash, err := goimagehash.ColorHash(img, binbits)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}
		hashStr = fmt.Sprintf("%s (bits: %d)", hash.ToString(), hash.Bits())

	case "cropresistant", "crop":
		hash, err := goimagehash.CropResistantHash(img, nil, 0, 0, 0)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}
		segments := hash.GetSegmentHashes()
		hashStr = fmt.Sprintf("%d segments", len(segments))
		for i, seg := range segments {
			hashStr += fmt.Sprintf("\n  segment%d: %s (bits: %d)", i+1, seg.ToString(), seg.Bits())
		}

	case "extahash":
		hash, err := goimagehash.ExtAverageHash(img, hashSize, hashSize)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}
		hashStr = fmt.Sprintf("%s (bits: %d)", hash.ToString(), hash.Bits())

	case "extphash":
		hash, err := goimagehash.ExtPerceptionHash(img, hashSize, hashSize)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}
		hashStr = fmt.Sprintf("%s (bits: %d)", hash.ToString(), hash.Bits())

	case "extdhash":
		hash, err := goimagehash.ExtDifferenceHash(img, hashSize, hashSize)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}
		hashStr = fmt.Sprintf("%s (bits: %d)", hash.ToString(), hash.Bits())

	case "extwhash":
		hash, err := goimagehash.ExtWaveletHash(img, hashSize, hashSize)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}
		hashStr = fmt.Sprintf("%s (bits: %d)", hash.ToString(), hash.Bits())

	default:
		fmt.Fprintf(os.Stderr, "Unknown hash type: %s\n", hashType)
		os.Exit(1)
	}

	fmt.Printf("%s: %s\n", hashType, hashStr)
}

func printAllHashes(img image.Image) {
	fmt.Println("Computing all hashes:")
	fmt.Println()

	hashTypes := []string{"ahash", "phash", "dhash", "whash", "colorhash"}

	for _, ht := range hashTypes {
		fmt.Printf("%s: ", ht)
		printHash(img, ht)
	}
}
