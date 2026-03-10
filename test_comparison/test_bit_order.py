#!/usr/bin/env python3
import imagehash
import numpy as np
from PIL import Image

# Create a simple test pattern
# 8x8 image with predictable pattern
img = Image.new('L', (8, 8), 0)
pixels = np.array(img)

# Set a specific pattern: diagonal line
for i in range(8):
    pixels[i, i] = 255

img = Image.fromarray(pixels)

# Compute hash
ahash = imagehash.average_hash(img)
print(f"Hash: {ahash}")
print(f"Hash hex: {ahash.hash.flatten().tobytes().hex()}")

# Get the bit matrix
bits = ahash.hash.flatten()
print(f"\nBit array (flattened): {bits}")

# Check if it's row-major
print("\nBit matrix (8x8):")
for i in range(8):
    row = ahash.hash[i]
    row_str = ''.join(['1' if x else '0' for x in row])
    print(f"Row {i}: {row_str}")

# Now let's manually compute what we expect
# Average of pixels: (8*255 + 56*0) / 64 = 31.875
# So diagonal pixels > average, others < average
print("\nExpected pattern:")
print("Diagonal should be 1, others 0")
for i in range(8):
    expected = ['1' if j == i else '0' for j in range(8)]
    print(f"Row {i}: {''.join(expected)}")

# Check if Python uses row-major, MSB-first
print("\nChecking bit order...")
# The hash should have bits along the diagonal
hash_hex = str(ahash)
print(f"Hash string: {hash_hex}")

# Convert to binary
hash_int = int(hash_hex, 16)
hash_bin = format(hash_int, '064b')
print(f"Hash binary (64 bits): {hash_bin}")

# Check if diagonal bits are set
print("\nChecking diagonal bits in binary string:")
for i in range(8):
    bit_pos = i * 8 + i  # row-major, MSB first?
    bit_value = hash_bin[bit_pos]
    print(f"  Position {bit_pos} (row {i}, col {i}): {bit_value}")