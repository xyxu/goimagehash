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
	flag.StringVar(&hashType, "type", "ahash", "Hash type: ahash, phash, dhash, whash, colorhash, cropresistant, extahash, extphash, extdhash, extwhash")
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
		printAllHashes(img, imagePath)
	} else {
		printHash(img, hashType, imagePath)
	}
}

// Helper function to extract hex value from hash string
func extractHashValue(hashStr string) string {
	// Format is "a:ffff3f030703c1f0" or "c:1c00000000000000"
	// Extract part after ":"
	parts := strings.SplitN(hashStr, ":", 2)
	if len(parts) == 2 {
		return parts[1]
	}
	return hashStr
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

func computeHash(img image.Image, hashType string) (*HashResult, *CropResistantResult, error) {
	switch strings.ToLower(hashType) {
	case "ahash", "average":
		hash, err := goimagehash.AverageHash(img)
		if err != nil {
			return nil, nil, err
		}
		return &HashResult{
			Type:  "ahash",
			Value: extractHashValue(hash.ToString()),
			Bits:  hash.Bits(),
		}, nil, nil

	case "phash", "perceptual":
		hash, err := goimagehash.PerceptionHash(img)
		if err != nil {
			return nil, nil, err
		}
		return &HashResult{
			Type:  "phash",
			Value: extractHashValue(hash.ToString()),
			Bits:  hash.Bits(),
		}, nil, nil

	case "dhash", "difference":
		hash, err := goimagehash.DifferenceHash(img)
		if err != nil {
			return nil, nil, err
		}
		return &HashResult{
			Type:  "dhash",
			Value: extractHashValue(hash.ToString()),
			Bits:  hash.Bits(),
		}, nil, nil

	case "whash", "wavelet":
		hash, err := goimagehash.WaveletHash(img)
		if err != nil {
			return nil, nil, err
		}
		return &HashResult{
			Type:  "whash",
			Value: extractHashValue(hash.ToString()),
			Bits:  hash.Bits(),
		}, nil, nil

	case "colorhash", "color":
		hash, err := goimagehash.ColorHash(img, binbits)
		if err != nil {
			return nil, nil, err
		}
		return &HashResult{
			Type:    "colorhash",
			Value:   extractHashValue(hash.ToString()),
			Bits:    hash.Bits(),
			Binbits: binbits,
		}, nil, nil

	case "cropresistant", "crop":
		hash, err := goimagehash.CropResistantHash(img, nil, 0, 0, 0)
		if err != nil {
			return nil, nil, err
		}

		segments := hash.GetSegmentHashes()
		result := &CropResistantResult{
			Type:     "cropresistant",
			Segments: len(segments),
		}

		for _, seg := range segments {
			segmentHash := SegmentHash{
				Value: extractHashValue(seg.ToString()),
				Bits:  seg.Bits(),
			}
			result.SegmentHashes = append(result.SegmentHashes, segmentHash)
		}

		return nil, result, nil

	case "extahash":
		hash, err := goimagehash.ExtAverageHash(img, hashSize, hashSize)
		if err != nil {
			return nil, nil, err
		}
		return &HashResult{
			Type:  "extahash",
			Value: extractHashValue(hash.ToString()),
			Bits:  hash.Bits(),
		}, nil, nil

	case "extphash":
		hash, err := goimagehash.ExtPerceptionHash(img, hashSize, hashSize)
		if err != nil {
			return nil, nil, err
		}
		return &HashResult{
			Type:  "extphash",
			Value: extractHashValue(hash.ToString()),
			Bits:  hash.Bits(),
		}, nil, nil

	case "extdhash":
		hash, err := goimagehash.ExtDifferenceHash(img, hashSize, hashSize)
		if err != nil {
			return nil, nil, err
		}
		return &HashResult{
			Type:  "extdhash",
			Value: extractHashValue(hash.ToString()),
			Bits:  hash.Bits(),
		}, nil, nil

	case "extwhash":
		hash, err := goimagehash.ExtWaveletHash(img, hashSize, hashSize)
		if err != nil {
			return nil, nil, err
		}
		return &HashResult{
			Type:  "extwhash",
			Value: extractHashValue(hash.ToString()),
			Bits:  hash.Bits(),
		}, nil, nil

	default:
		return nil, nil, fmt.Errorf("unknown hash type: %s", hashType)
	}
}

func printHash(img image.Image, hashType string, imagePath string) {
	hashResult, cropResult, err := computeHash(img, hashType)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}

	if outputFormat == "json" {
		outputJSON(imagePath, hashResult, cropResult)
	} else {
		outputText(hashResult, cropResult)
	}
}

func printAllHashes(img image.Image, imagePath string) {
	hashTypes := []string{"ahash", "phash", "dhash", "whash", "colorhash", "cropresistant"}
	hashes := make(map[string]HashResult)
	var cropResult *CropResistantResult

	// Compute all hashes
	for _, ht := range hashTypes {
		hashResult, cropRes, err := computeHash(img, ht)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error computing %s: %v\n", ht, err)
			os.Exit(1)
		}

		if ht == "cropresistant" {
			cropResult = cropRes
		} else if hashResult != nil {
			hashes[ht] = *hashResult
		}
	}

	// Output based on format
	if outputFormat == "json" {
		outputAllJSON(imagePath, hashes, cropResult)
	} else {
		outputAllText(hashes, cropResult)
	}
}
