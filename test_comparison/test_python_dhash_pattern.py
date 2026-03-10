#!/usr/bin/env python3
import imagehash
import numpy as np
from PIL import Image

# Create same test pattern as Go test
# 9x8 image with increasing values horizontally
img = Image.new('L', (9, 8))
pixels = np.array(img)

# Set pixel values: row i has values i*10, i*10+1, i*10+2, ...
for i in range(8):
    for j in range(9):
        pixels[i, j] = i*10 + j

img = Image.fromarray(pixels)

# Compute hash
dhash = imagehash.dhash(img)
print(f"Python dhash: {dhash}")
print(f"Python dhash hex: {dhash.hash.flatten().tobytes().hex()}")

# Get bit matrix
print(f"\nBit matrix (8x8):")
for i in range(8):
    row = dhash.hash[i]
    row_str = ''.join(['1' if x else '0' for x in row])
    print(f"Row {i}: {row_str}")

# For increasing values: pixel[j] < pixel[j+1] always True
# So all bits should be 1
expected = 'ffffffffffffffff'
print(f"\nExpected (all 1s): {expected}")

if str(dhash) == expected:
    print("✓ MATCH! All bits are 1 as expected")
else:
    print("✗ MISMATCH")
    
    # Convert to binary
    hash_int = int(str(dhash), 16)
    hash_bin = format(hash_int, '064b')
    print(f"Hash binary: {hash_bin}")
    
    # Check which bits are 0
    print("Bits that are 0 (should be 1):")
    for i in range(64):
        if hash_bin[i] == '0':
            row = i // 8
            col = i % 8
            print(f"  Bit {i} (row {row}, col {col})")