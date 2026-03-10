# Hash Mismatch Root Cause Analysis Implementation Plan

> **For Claude:** REQUIRED SUB-SKILL: Use superpowers:executing-plans to implement this plan task-by-task.

**Goal:** Identify and fix root causes of hash mismatches between Go imagehash and Python imagehash implementations

**Architecture:** Analyze each mismatching algorithm (average hash, difference hash, crop-resistant hash) independently using step-by-step comparison, create minimal test cases, and implement fixes.

**Tech Stack:** Go, Python, image processing, hash algorithms

---

## Task 1: Analyze Difference Hash Mismatch

**Files:**
- Modify: `test_comparison/debug_dhash.go` - Add detailed logging
- Create: `test_comparison/analyze_dhash.py` - Python comparison script
- Modify: `test_comparison/test_python_dhash_pattern.py` - Extend analysis

**Step 1: Create Python analysis script to log intermediate values**

```python
#!/usr/bin/env python3
import imagehash
from PIL import Image
import numpy as np

# Load sample image
img = Image.open('_examples/sample1.jpg')

# Get dhash with logging
print("Python dhash analysis for sample1.jpg")
print("=" * 50)

# Resize to 9x8 as per dhash default
resized = img.resize((9, 8), Image.Resampling.LANCZOS)
print(f"Resized dimensions: {resized.size}")

# Convert to grayscale  
gray = resized.convert('L')
print("Grayscale pixel values (first row):")
pixels = list(gray.getdata())
width, height = gray.size
for y in range(height):
    row = pixels[y*width:(y+1)*width]
    print(f"Row {y}: {row[:10]}...")

# Compute differences
print("\nDifference comparisons (pixels[i][j] > pixels[i][j+1]):")
bits = []
for y in range(height):
    row = pixels[y*width:(y+1)*width]
    for x in range(width-1):
        bit = 1 if row[x] > row[x+1] else 0
        bits.append(bit)
        if y < 2 and x < 5:  # Print first few
            print(f"  pixels[{y}][{x}]={row[x]:3d} > pixels[{y}][{x+1}]={row[x+1]:3d} ? {bit}")

print(f"\nTotal bits: {len(bits)}")
print(f"Bits (first 16): {bits[:16]}")

# Convert to hex
hash_obj = imagehash.dhash(img)
print(f"\nPython dhash result: {hash_obj}")
print(f"Hash hex: {hash_obj.hash.flatten()}")
```

**Step 2: Run Python analysis to get reference values**

Run: `cd test_comparison && python3 analyze_dhash.py`
Expected: Detailed output showing pixel values, comparisons, and final hash

**Step 3: Create Go equivalent with matching logging**

Modify `debug_dhash.go` to:
1. Load same image
2. Log resize dimensions and pixel values
3. Log comparison logic
4. Show bit ordering

**Step 4: Compare Go and Python intermediate values**

Run both scripts and compare:
- Resize dimensions
- Grayscale pixel values (first few rows)
- Comparison results
- Bit ordering

**Step 5: Identify the mismatch cause**

Based on comparison, determine if issue is:
- Comparison direction (`<` vs `>`)
- Bit ordering (MSB vs LSB)
- Resize algorithm difference
- Grayscale conversion formula

**Step 6: Commit analysis findings**

```bash
git add test_comparison/analyze_dhash.py test_comparison/debug_dhash.go
git commit -m "analysis: detailed dhash comparison between Go and Python"
```

---

## Task 2: Fix Difference Hash Implementation

**Files:**
- Modify: `hashcompute.go` - Difference hash implementation
- Test: `hashcompute_test.go` - Add test for Python compatibility
- Modify: `test_comparison/compare_with_python.go` - Update to verify fix

**Step 1: Write test for Python-compatible dhash**

Add to `hashcompute_test.go`:

```go
func TestDifferenceHashPythonCompatible(t *testing.T) {
    // Load test image
    file, err := os.Open("_examples/sample1.jpg")
    require.NoError(t, err)
    defer file.Close()
    
    img, err := jpeg.Decode(file)
    require.NoError(t, err)
    
    // Compute hash
    hash, err := DifferenceHash(img)
    require.NoError(t, err)
    
    // Python reference hash for sample1.jpg: e89f6f63af4f8b69
    expected := "e89f6f63af4f8b69"
    
    // Remove 'd:' prefix for comparison
    hashStr := hash.ToString()
    if strings.HasPrefix(hashStr, "d:") {
        hashStr = hashStr[2:]
    }
    
    assert.Equal(t, expected, hashStr, "Difference hash should match Python imagehash")
}
```

**Step 2: Run test to verify it fails**

Run: `go test -v -run TestDifferenceHashPythonCompatible ./...`
Expected: FAIL with hash mismatch

**Step 3: Analyze hashcompute.go DifferenceHash implementation**

Examine the `DifferenceHash` function in `hashcompute.go`:
1. Check resize dimensions (should be 9x8 for horizontal differences)
2. Check comparison logic (`pixels[i][j] < pixels[i][j+1]` vs `>`)
3. Check bit ordering (MSB vs LSB)

**Step 4: Implement fix based on Task 1 findings**

If comparison direction is wrong:
```go
// Change from:
if pixels[y][x] < pixels[y][x+1] {
    hash |= 1 << uint(idx)
}

// To:
if pixels[y][x] > pixels[y][x+1] {
    hash |= 1 << uint(idx)
}
```

If bit ordering is wrong, adjust the `idx` calculation.

**Step 5: Run test to verify fix works**

Run: `go test -v -run TestDifferenceHashPythonCompatible ./...`
Expected: PASS

**Step 6: Update comparison tool and verify**

Run: `cd test_comparison && go run compare_with_python.go`
Check that difference hash now shows "✓ difference: MATCH"

**Step 7: Commit the fix**

```bash
git add hashcompute.go hashcompute_test.go
git commit -m "fix: align difference hash with Python imagehash implementation"
```

---

## Task 3: Analyze Average Hash Mismatch

**Files:**
- Create: `test_comparison/analyze_ahash.py` - Python analysis
- Modify: `test_comparison/debug_go_hashes2.go` - Extend for ahash
- Examine: `hashcompute.go` - AverageHash function

**Step 1: Create Python ahash analysis**

```python
#!/usr/bin/env python3
import imagehash
from PIL import Image
import numpy as np

img = Image.open('_examples/sample1.jpg')

print("Python average hash analysis")
print("=" * 50)

# Get ahash with default size 8x8
hash_obj = imagehash.average_hash(img)
print(f"Python ahash: {hash_obj}")

# Manual calculation to verify
from PIL import Image
import numpy as np

# Resize to 8x8
resized = img.resize((8, 8), Image.Resampling.LANCZOS)
print(f"Resized to: {resized.size}")

# Convert to grayscale
gray = resized.convert('L')
pixels = np.array(gray.getdata()).reshape((8, 8))
print("Pixel matrix:")
for y in range(8):
    print(f"Row {y}: {pixels[y].tolist()}")

# Calculate mean
mean = np.mean(pixels)
print(f"\nMean pixel value: {mean:.2f}")

# Create hash
bits = []
for y in range(8):
    for x in range(8):
        bits.append(1 if pixels[y, x] > mean else 0)

print(f"\nBits (row-major): {bits[:16]}...")
print(f"Total bits: {len(bits)}")
```

**Step 2: Run Python analysis**

Run: `cd test_comparison && python3 analyze_ahash.py`
Expected: Detailed output showing resize, pixels, mean, bits

**Step 3: Create Go equivalent analysis**

Extend `debug_go_hashes2.go` to log:
1. Resize dimensions and method
2. Grayscale pixel values
3. Mean calculation
4. Bit comparisons

**Step 4: Compare Go vs Python**

Identify differences in:
- Resize algorithm (Lanczos vs bilinear vs nearest)
- Grayscale conversion formula
- Mean calculation (integer vs float)
- Bit threshold (pixel > mean vs pixel >= mean)

**Step 5: Commit analysis**

```bash
git add test_comparison/analyze_ahash.py test_comparison/debug_go_hashes2.go
git commit -m "analysis: average hash comparison between Go and Python"
```

---

## Task 4: Fix Average Hash Implementation

**Files:**
- Modify: `hashcompute.go` - AverageHash function
- Test: `hashcompute_test.go` - Add Python compatibility test
- Modify: `transforms/pixels.go` - Check resize implementation

**Step 1: Write test for Python-compatible ahash**

Add to `hashcompute_test.go`:

```go
func TestAverageHashPythonCompatible(t *testing.T) {
    file, err := os.Open("_examples/sample1.jpg")
    require.NoError(t, err)
    defer file.Close()
    
    img, err := jpeg.Decode(file)
    require.NoError(t, err)
    
    hash, err := AverageHash(img)
    require.NoError(t, err)
    
    // Python reference: 7eff3f030703c1fc
    expected := "7eff3f030703c1fc"
    
    hashStr := hash.ToString()
    if strings.HasPrefix(hashStr, "a:") {
        hashStr = hashStr[2:]
    }
    
    assert.Equal(t, expected, hashStr, "Average hash should match Python imagehash")
}
```

**Step 2: Run test to verify failure**

Run: `go test -v -run TestAverageHashPythonCompatible ./...`
Expected: FAIL

**Step 3: Implement fix based on analysis**

Possible fixes:
1. Change resize algorithm to match Python's LANCZOS
2. Adjust grayscale formula (check RGB weights)
3. Fix mean calculation (use float64)
4. Adjust bit threshold

**Step 4: Test the fix**

Run: `go test -v -run TestAverageHashPythonCompatible ./...`
Expected: PASS

**Step 5: Update comparison tool**

Run: `cd test_comparison && go run compare_with_python.go`
Check average hash shows "✓ average: MATCH"

**Step 6: Commit fix**

```bash
git add hashcompute.go hashcompute_test.go
git commit -m "fix: align average hash with Python imagehash"
```

---

## Task 5: Analyze Crop-Resistant Hash Mismatch

**Files:**
- Create: `test_comparison/analyze_crophash.py` - Python analysis
- Examine: `hashcompute.go` - CropResistantHash function
- Check: Python imagehash source for crop_resistant_hash algorithm

**Step 1: Research Python crop_resistant_hash algorithm**

Check Python imagehash source or documentation for:
- Segment size calculation
- Hash algorithm used per segment
- Threshold for segment inclusion
- Segment overlap or grid pattern

**Step 2: Create Python analysis script**

```python
#!/usr/bin/env python3
import imagehash
from PIL import Image

img = Image.open('_examples/sample1.jpg')

print("Python crop-resistant hash analysis")
print("=" * 50)

# Get crop resistant hash
hash_obj = imagehash.crop_resistant_hash(img)
print(f"Python crop-resistant hash: {hash_obj}")

# Examine segments
print(f"\nNumber of segments: {len(hash_obj.segment_hashes)}")
for i, seg_hash in enumerate(hash_obj.segment_hashes):
    print(f"Segment {i}: {seg_hash}")
    
# Check parameters
print("\nDefault parameters in Python imagehash:")
print(f"  segment_threshold: {imagehash.crop_resistant_hash.__defaults__}")
```

**Step 3: Compare with Go implementation**

Examine `CropResistantHash` in `hashcompute.go`:
1. Segment size calculation
2. Hash algorithm (appears to use dhash)
3. Segment count logic
4. Threshold values

**Step 4: Identify differences**

Key areas:
- Different default segment size
- Different hash algorithm per segment
- Different threshold for including segments
- Different segmentation pattern (grid vs other)

**Step 5: Commit analysis**

```bash
git add test_comparison/analyze_crophash.py
git commit -m "analysis: crop-resistant hash comparison"
```

---

## Task 6: Document Findings and Compatibility Notes

**Files:**
- Create: `COMPATIBILITY.md` - Document Go vs Python differences
- Update: `README.md` - Add compatibility section
- Create: `test_comparison/verification_suite.go` - Comprehensive test

**Step 1: Create COMPATIBILITY.md**

```markdown
# Go imagehash vs Python imagehash Compatibility

## Algorithms with Full Compatibility
- Perception Hash (phash) - ✓ Matches exactly
- Wavelet Hash (whash) - ✓ Matches exactly  
- Color Hash - ✓ Matches exactly (with binbits parameter)

## Algorithms with Differences

### Difference Hash (dhash)
**Status**: Fixed in version X.X.X
**Previous mismatch**: Comparison direction (`<` vs `>`)
**Fix**: Changed to `pixels[y][x] > pixels[y][x+1]` to match Python

### Average Hash (ahash)
**Status**: Fixed in version X.X.X  
**Previous mismatch**: Resize algorithm and bit threshold
**Fix**: [Describe fix applied]

### Crop-Resistant Hash
**Status**: Intentional design difference
**Difference**: 
- Go uses 6 segments with dhash algorithm
- Python uses 4 segments with dhash algorithm
- Different segmentation logic and thresholds

**Recommendation**: Not directly comparable between implementations

## Usage Notes
1. For cross-language compatibility, use phash or whash
2. dhash and ahash now match Python after fixes
3. Crop-resistant hashes are implementation-specific
```

**Step 2: Update README.md**

Add compatibility section referencing COMPATIBILITY.md

**Step 3: Create verification test suite**

Create `test_comparison/verification_suite.go` that runs all comparisons and reports status.

**Step 4: Run final verification**

Run: `cd test_comparison && go run verification_suite.go`
Expected: All compatible algorithms show "MATCH"

**Step 5: Commit documentation**

```bash
git add COMPATIBILITY.md README.md test_comparison/verification_suite.go
git commit -m "docs: add compatibility documentation and verification suite"
```

---

## Task 7: Final Integration and Testing

**Files:**
- Run: All existing tests
- Check: `go test ./...` - Ensure no regressions
- Verify: `test_comparison/compare_with_python.go` - All matches

**Step 1: Run full test suite**

```bash
go test ./...
```
Expected: All tests pass

**Step 2: Run comparison with all sample images**

```bash
cd test_comparison && go run compare_with_python.go
```
Check output for all 4 sample images.

**Step 3: Create summary report**

Generate final compatibility report showing:
- Algorithms that now match
- Remaining differences (if any)
- Recommendations for users

**Step 4: Final commit**

```bash
git add .
git commit -m "chore: complete hash mismatch analysis and fixes"
```

**Step 5: Optional - Create PR or tag release**

If appropriate, create GitHub PR or tag new release version.