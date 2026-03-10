# JSON Output Feature Implementation Plan

> **For Claude:** REQUIRED SUB-SKILL: Use superpowers:executing-plans to implement this plan task-by-task.

**Goal:** Add JSON output format support to goimagehash CLI tool with `-format json` flag.

**Architecture:** Add `-format` flag with `text` (default) and `json` options. Refactor output functions to support both formats while maintaining backward compatibility.

**Tech Stack:** Go, standard library (encoding/json, flag), goimagehash library

---

### Task 1: Add format flag and refactor main function

**Files:**
- Modify: `cmd/goimagehash/main.go:15-27` (flag definitions)
- Modify: `cmd/goimagehash/main.go:29-51` (main function)

**Step 1: Add format flag variable**

```go
// Add to flag variables
var (
    hashType    string
    binbits     int
    hashSize    int
    allHashes   bool
    outputFormat string  // New: "text" or "json"
)
```

**Step 2: Register format flag in init()**

```go
func init() {
    flag.StringVar(&hashType, "type", "ahash", "Hash type: ahash, phash, dhash, whash, colorhash, cropresistant")
    flag.IntVar(&binbits, "binbits", 3, "Bin bits for colorhash (default: 3)")
    flag.IntVar(&hashSize, "size", 8, "Hash size for ahash, phash, dhash, whash (default: 8)")
    flag.BoolVar(&allHashes, "all", false, "Compute all hash types")
    flag.StringVar(&outputFormat, "format", "text", "Output format: text, json")
}
```

**Step 3: Validate format flag in main()**

```go
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
    
    // Rest of main function...
}
```

**Step 4: Run basic test to verify flag works**

```bash
cd /home/deploy/Downloads/goimagehash
go run cmd/goimagehash/main.go -help 2>&1 | grep -A1 "format"
```
Expected: Shows format flag in help output

**Step 5: Commit**

```bash
git add cmd/goimagehash/main.go
git commit -m "feat: add -format flag for output format selection"
```

---

### Task 2: Create hash data structures

**Files:**
- Create: `cmd/goimagehash/types.go`

**Step 1: Create types.go with hash data structures**

```go
package main

// HashResult represents a single hash computation result
type HashResult struct {
    Type  string `json:"type"`
    Value string `json:"value"`
    Bits  int    `json:"bits"`
    Binbits int  `json:"binbits,omitempty"` // Only for colorhash
}

// SegmentHash represents a segment in crop-resistant hash
type SegmentHash struct {
    Value  string `json:"value"`
    Bits   int    `json:"bits"`
    Bounds struct {
        MinX int `json:"min_x"`
        MinY int `json:"min_y"`
        MaxX int `json:"max_x"`
        MaxY int `json:"max_y"`
    } `json:"bounds,omitempty"`
}

// CropResistantResult represents crop-resistant hash result
type CropResistantResult struct {
    Type          string        `json:"type"`
    Segments      int           `json:"segments"`
    SegmentHashes []SegmentHash `json:"segment_hashes,omitempty"`
}

// OutputData represents complete output data
type OutputData struct {
    Image  string                 `json:"image"`
    Hash   *HashResult            `json:"hash,omitempty"`           // For single hash
    Hashes map[string]HashResult  `json:"hashes,omitempty"`         // For -all flag
    Crop   *CropResistantResult   `json:"crop,omitempty"`           // For crop-resistant hash
    Error  string                 `json:"error,omitempty"`
}
```

**Step 2: Test compilation**

```bash
cd /home/deploy/Downloads/goimagehash
go build ./cmd/goimagehash
```
Expected: Build succeeds

**Step 3: Commit**

```bash
git add cmd/goimagehash/types.go
git commit -m "feat: add hash data structures for JSON output"
```

---

### Task 3: Refactor printHash to return data

**Files:**
- Modify: `cmd/goimagehash/main.go:76-170` (printHash function)

**Step 1: Change printHash signature and return HashResult**

```go
func computeHash(img image.Image, hashType string) (*HashResult, *CropResistantResult, error) {
    switch strings.ToLower(hashType) {
    case "ahash", "average":
        hash, err := goimagehash.AverageHash(img)
        if err != nil {
            return nil, nil, err
        }
        return &HashResult{
            Type:  "ahash",
            Value: hash.ToString(),
            Bits:  hash.Bits(),
        }, nil, nil
        
    case "phash", "perceptual":
        hash, err := goimagehash.PerceptionHash(img)
        if err != nil {
            return nil, nil, err
        }
        return &HashResult{
            Type:  "phash",
            Value: hash.ToString(),
            Bits:  hash.Bits(),
        }, nil, nil
        
    // Continue for other hash types...
    }
}
```

**Step 2: Update colorhash case to include binbits**

```go
case "colorhash", "color":
    hash, err := goimagehash.ColorHash(img, binbits)
    if err != nil {
        return nil, nil, err
    }
    return &HashResult{
        Type:    "colorhash",
        Value:   hash.ToString(),
        Bits:    hash.Bits(),
        Binbits: binbits,
    }, nil, nil
```

**Step 3: Handle crop-resistant hash separately**

```go
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
    
    // Note: bounds not available in current API
    for _, seg := range segments {
        segmentHash := SegmentHash{
            Value: seg.ToString(),
            Bits:  seg.Bits(),
        }
        result.SegmentHashes = append(result.SegmentHashes, segmentHash)
    }
    
    return nil, result, nil
```

**Step 4: Test compilation**

```bash
cd /home/deploy/Downloads/goimagehash
go build ./cmd/goimagehash
```
Expected: Build succeeds

**Step 5: Commit**

```bash
git add cmd/goimagehash/main.go
git commit -m "refactor: change printHash to computeHash returning data"
```

---

### Task 4: Create output formatters

**Files:**
- Create: `cmd/goimagehash/output.go`

**Step 1: Create outputText function**

```go
package main

import (
    "fmt"
    "os"
)

func outputText(data *OutputData) {
    if data.Error != "" {
        fmt.Fprintf(os.Stderr, "Error: %s\n", data.Error)
        return
    }
    
    if data.Crop != nil {
        fmt.Printf("cropresistant: %d segments\n", data.Crop.Segments)
        for i, seg := range data.Crop.SegmentHashes {
            fmt.Printf("  segment%d: %s (bits: %d)\n", i+1, seg.Value, seg.Bits)
        }
    } else if data.Hash != nil {
        fmt.Printf("%s: %s (bits: %d)\n", data.Hash.Type, data.Hash.Value, data.Hash.Bits)
    } else if data.Hashes != nil {
        fmt.Println("Computing all hashes:")
        fmt.Println()
        for _, hash := range data.Hashes {
            fmt.Printf("%s: %s (bits: %d)\n", hash.Type, hash.Value, hash.Bits)
        }
    }
}
```

**Step 2: Create outputJSON function**

```go
import (
    "encoding/json"
    "fmt"
    "os"
)

func outputJSON(data *OutputData) {
    jsonData, err := json.MarshalIndent(data, "", "  ")
    if err != nil {
        fmt.Fprintf(os.Stderr, "Error encoding JSON: %v\n", err)
        os.Exit(1)
    }
    fmt.Println(string(jsonData))
}
```

**Step 3: Create output function that routes based on format**

```go
func output(data *OutputData) {
    switch outputFormat {
    case "json":
        outputJSON(data)
    case "text":
        outputText(data)
    default:
        // Should not happen due to validation in main()
        fmt.Fprintf(os.Stderr, "Error: unknown format '%s'\n", outputFormat)
        os.Exit(1)
    }
}
```

**Step 4: Test compilation**

```bash
cd /home/deploy/Downloads/goimagehash
go build ./cmd/goimagehash
```
Expected: Build succeeds

**Step 5: Commit**

```bash
git add cmd/goimagehash/output.go
git commit -m "feat: add output formatters for text and JSON"
```

---

### Task 5: Update main function to use new system

**Files:**
- Modify: `cmd/goimagehash/main.go:29-51` (main function)
- Modify: `cmd/goimagehash/main.go:172-182` (printAllHashes function)

**Step 1: Update main function to use new system**

```go
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
        outputData := &OutputData{
            Image: imagePath,
            Error: fmt.Sprintf("Error loading image: %v", err),
        }
        output(outputData)
        os.Exit(1)
    }
    
    outputData := &OutputData{
        Image: imagePath,
    }
    
    if allHashes {
        // Compute all hashes
        hashes := make(map[string]HashResult)
        hashTypes := []string{"ahash", "phash", "dhash", "whash", "colorhash"}
        
        for _, ht := range hashTypes {
            hashResult, _, err := computeHash(img, ht)
            if err != nil {
                outputData.Error = fmt.Sprintf("Error computing %s: %v", ht, err)
                output(outputData)
                os.Exit(1)
            }
            hashes[ht] = *hashResult
        }
        
        outputData.Hashes = hashes
    } else {
        // Compute single hash
        hashResult, cropResult, err := computeHash(img, hashType)
        if err != nil {
            outputData.Error = fmt.Sprintf("Error computing hash: %v", err)
            output(outputData)
            os.Exit(1)
        }
        
        if cropResult != nil {
            outputData.Crop = cropResult
        } else {
            outputData.Hash = hashResult
        }
    }
    
    output(outputData)
}
```

**Step 2: Remove old printAllHashes function**

Delete the `printAllHashes` function (lines 172-182) as it's no longer needed.

**Step 3: Test basic functionality**

```bash
cd /home/deploy/Downloads/goimagehash
go run cmd/goimagehash/main.go -type ahash _examples/sample1.jpg
```
Expected: Text output similar to before

```bash
go run cmd/goimagehash/main.go -type ahash -format json _examples/sample1.jpg
```
Expected: JSON output

**Step 4: Test -all flag**

```bash
go run cmd/goimagehash/main.go -all _examples/sample1.jpg
```
Expected: Text output of all hashes

```bash
go run cmd/goimagehash/main.go -all -format json _examples/sample1.jpg
```
Expected: JSON output of all hashes

**Step 5: Commit**

```bash
git add cmd/goimagehash/main.go
git commit -m "feat: integrate new output system into main function"
```

---

### Task 6: Fix JSON output formatting issues

**Files:**
- Modify: `cmd/goimagehash/types.go` (HashResult.Value field)
- Modify: `cmd/goimagehash/main.go` (computeHash function)

**Step 1: Fix hash value extraction**

Current `hash.ToString()` returns format like `"a:ffff3f030703c1f0"`. For JSON, we should extract just the hex part.

Update `HashResult` creation in `computeHash`:

```go
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

// Update each case in computeHash:
return &HashResult{
    Type:  "ahash",
    Value: extractHashValue(hash.ToString()),
    Bits:  hash.Bits(),
}, nil, nil
```

**Step 2: Also update for ExtImageHash types**

For extended hashes (extahash, extphash, etc.), they also use `ToString()` method.

**Step 3: Test JSON output**

```bash
cd /home/deploy/Downloads/goimagehash
go run cmd/goimagehash/main.go -type ahash -format json _examples/sample1.jpg | jq .
```
Expected: Clean JSON with just hex value (no "a:" prefix)

**Step 4: Test all hash types**

```bash
for ht in ahash phash dhash whash colorhash; do
    echo "Testing $ht:"
    go run cmd/goimagehash/main.go -type $ht -format json _examples/sample1.jpg | jq -r '.hash.value'
done
```
Expected: Clean hex values for all hash types

**Step 5: Commit**

```bash
git add cmd/goimagehash/main.go cmd/goimagehash/types.go
git commit -m "fix: extract clean hex values for JSON output"
```

---

### Task 7: Add error handling and edge cases

**Files:**
- Modify: `cmd/goimagehash/main.go` (error handling)
- Modify: `cmd/goimagehash/output.go` (error output)

**Step 1: Improve error handling in outputText**

```go
func outputText(data *OutputData) {
    if data.Error != "" {
        fmt.Fprintf(os.Stderr, "Error: %s\n", data.Error)
        os.Exit(1)
    }
    
    // Rest of function...
}
```

**Step 2: Test error cases**

```bash
cd /home/deploy/Downloads/goimagehash
go run cmd/goimagehash/main.go nonexistent.jpg
```
Expected: Error message

```bash
go run cmd/goimagehash/main.go -type invalid _examples/sample1.jpg
```
Expected: Error message

```bash
go run cmd/goimagehash/main.go -type invalid -format json _examples/sample1.jpg
```
Expected: JSON with error field

**Step 3: Test with -format invalid**

```bash
go run cmd/goimagehash/main.go -format invalid _examples/sample1.jpg
```
Expected: Error about invalid format

**Step 4: Commit**

```bash
git add cmd/goimagehash/main.go cmd/goimagehash/output.go
git commit -m "feat: improve error handling for JSON output"
```

---

### Task 8: Final testing and cleanup

**Files:**
- All modified files

**Step 1: Run comprehensive tests**

```bash
cd /home/deploy/Downloads/goimagehash
# Test all combinations
echo "=== Testing text format ==="
go run cmd/goimagehash/main.go -type ahash _examples/sample1.jpg
go run cmd/goimagehash/main.go -all _examples/sample1.jpg

echo "=== Testing JSON format ==="
go run cmd/goimagehash/main.go -type ahash -format json _examples/sample1.jpg
go run cmd/goimagehash/main.go -all -format json _examples/sample1.jpg

echo "=== Testing extended hashes ==="
go run cmd/goimagehash/main.go -type extahash -size 16 _examples/sample1.jpg
go run cmd/goimagehash/main.go -type extahash -size 16 -format json _examples/sample1.jpg

echo "=== Testing crop-resistant hash ==="
go run cmd/goimagehash/main.go -type cropresistant _examples/sample1.jpg
go run cmd/goimagehash/main.go -type cropresistant -format json _examples/sample1.jpg
```

**Step 2: Build and verify**

```bash
go build ./cmd/goimagehash
./goimagehash -help
```
Expected: Shows -format flag in help

**Step 3: Run existing tests**

```bash
go test ./...
```
Expected: All tests pass

**Step 4: Final commit**

```bash
git add .
git commit -m "feat: complete JSON output feature for CLI tool"
```

---

**Plan complete and saved to `docs/plans/2026-03-09-json-output-feature.md`. Two execution options:**

**1. Subagent-Driven (this session)** - I dispatch fresh subagent per task, review between tasks, fast iteration

**2. Parallel Session (separate)** - Open new session with executing-plans, batch execution with checkpoints

**Which approach?**