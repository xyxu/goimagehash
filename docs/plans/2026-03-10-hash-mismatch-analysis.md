# Hash Mismatch Analysis Design
**Date**: 2026-03-10  
**Author**: Sisyphus (AI Agent)  
**Goal**: Identify root causes of mismatches between Go imagehash and Python imagehash implementations

## 1. Background

The `goimagehash` library is inspired by Python's `imagehash` library but shows mismatches in several hash algorithms:
- Average Hash (ahash)
- Difference Hash (dhash) 
- Crop-Resistant Hash

Other algorithms (Perception Hash, Wavelet Hash, Color Hash) appear to match.

## 2. Current Mismatch Status

From `test_comparison/compare_with_python.go` output:

### Sample1.jpg Results:
- **Average Hash**: `a:7cff3f030703c1f8` (Go) vs `7eff3f030703c1fc` (Python) - **MISMATCH**
- **Difference Hash**: `d:c89f6f63af4f8b69` (Go) vs `e89f6f63af4f8b69` (Python) - **MISMATCH**
- **Perception Hash**: `p:af95d2205c7b1f82` (Go) vs `af95d2205c7b1f82` (Python) - **MATCH** (prefix only difference)
- **Wavelet Hash**: `w:3cff3f010703c1f8` (Go) vs `3cff3f010703c1f8` (Python) - **MATCH** (prefix only difference)
- **Color Hash**: `c:07000000000` (Go) vs `07000000000` (Python) - **MATCH** (prefix only difference)
- **Crop-Resistant Hash**: Different segment counts (6 vs 4) and values - **MISMATCH**

## 3. Analysis Methodology

### Phase 1: Algorithm Decomposition
For each mismatching algorithm:
1. Document algorithm steps from both implementations
2. Identify key parameters: hash size, resize method, comparison thresholds
3. Map data flow: image → resize → grayscale → processing → bits → hex

### Phase 2: Step-by-Step Comparison
Create instrumentation to compare:
1. **Image loading**: Verify same pixel data
2. **Resize operation**: Compare output dimensions and pixel values
3. **Grayscale conversion**: Check RGB→grayscale formula matches
4. **Comparison logic**: For dhash, verify comparison direction
5. **Bit ordering**: Verify MSB/LSB and row-major vs column-major flattening
6. **Hex encoding**: Check endianness and string formatting

### Phase 3: Minimal Test Cases
Create simplified test cases:
- Small synthetic images (e.g., 4x4 gradients)
- Known edge cases (uniform colors, high contrast edges)
- Step-by-step value logging

### Phase 4: Fix Identification
For each discrepancy:
1. Determine if it's a **bug** (incorrect implementation) or **design difference** (intentional variation)
2. Propose fixes for bugs
3. Document design differences for user awareness

## 4. Specific Investigation Areas

### 4.1 Average Hash Mismatch
**Observation**: First nibble differs: `7c` vs `7e` (binary: `01111100` vs `01111110`)
**Possible causes**:
- Different resize algorithm (Lanczos vs bilinear vs nearest)
- Different grayscale conversion formula
- Different median calculation for threshold
- Bit ordering difference

### 4.2 Difference Hash Mismatch  
**Observation**: First nibble differs: `c8` vs `e8` (binary: `11001000` vs `11101000`)
**Possible causes**:
- Comparison direction: `pixels[i][j] < pixels[i][j+1]` vs `pixels[i][j] > pixels[i][j+1]`
- Different resize dimensions
- Bit ordering (MSB vs LSB)

### 4.3 Crop-Resistant Hash Mismatch
**Observation**: Different number of segments (6 vs 4) and different values
**Possible causes**:
- Different segmentation algorithm
- Different parameters (segment size, threshold)
- Different hash algorithm applied to segments

## 5. Tools and Infrastructure

### Existing Tools:
- `test_comparison/compare_with_python.go` - main comparison tool
- `test_comparison/python_reference_hashes.json` - Python reference hashes
- `test_comparison/debug_*.go` - debugging utilities
- `test_comparison/test_*.py` - Python test scripts

### New Tools Needed:
1. **Step-by-step logger** for both Go and Python implementations
2. **Pixel value comparator** to verify intermediate results match
3. **Bit visualization tool** to compare bit ordering
4. **Algorithm parameter extractor** to document default values

## 6. Success Criteria

1. **Root cause identification**: Understand why each hash mismatch occurs
2. **Fix recommendations**: Have actionable fixes for reproducible mismatches  
3. **Updated test suite**: Can verify hash compatibility after fixes
4. **Documentation**: Clear explanation of Go vs Python differences

## 7. Risks and Mitigations

- **Risk**: Python imagehash implementation details may be unclear
  - **Mitigation**: Study source code and create reference tests
- **Risk**: Some differences may be intentional design choices
  - **Mitigation**: Document as compatibility notes rather than bugs
- **Risk**: Fixing one algorithm may break others
  - **Mitigation**: Comprehensive testing across all hash types

## 8. Next Steps

1. Begin with Difference Hash analysis (simplest mismatch pattern)
2. Create minimal test case to isolate the comparison direction issue
3. Verify bit ordering matches Python's row-major flattening
4. Apply findings to Average Hash analysis
5. Finally analyze Crop-Resistant Hash (most complex)