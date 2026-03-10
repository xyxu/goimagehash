#!/usr/bin/env python3
from PIL import Image
import imagehash
import numpy as np

# Load the sample image
img = Image.open('../_examples/sample1.jpg')

print("=== Python imagehash debug ===")
print(f"Image size: {img.size}")
print(f"Image mode: {img.mode}")

# Compute average hash
print("\n--- Average Hash ---")
ahash = imagehash.average_hash(img)
print(f"Hash: {ahash}")
print(f"Hash hex: {ahash.hash.flatten().tobytes().hex()}")

# Let's manually compute to understand
# First, convert to grayscale and resize to 8x8
gray = img.convert('L')
small = gray.resize((8, 8), Image.Resampling.LANCZOS)
print(f"\n8x8 grayscale pixels (0-255):")
pixels = np.array(small)
for i in range(8):
    row = pixels[i]
    print(f"Row {i}: {row}")

# Compute average
avg = np.mean(pixels)
print(f"\nAverage pixel value: {avg}")

# Create hash: 1 if pixel > avg, else 0
hash_matrix = pixels > avg
print(f"\nHash matrix (True = 1, False = 0):")
for i in range(8):
    row = hash_matrix[i]
    row_str = ''.join(['1' if x else '0' for x in row])
    print(f"Row {i}: {row_str}")

# Convert to hex
# Imagehash packs bits row-major, MSB first
bits = hash_matrix.flatten()
print(f"\nFlattened bits (64 bits): {''.join(['1' if x else '0' for x in bits])}")

# Compute difference hash
print("\n--- Difference Hash ---")
dhash = imagehash.dhash(img)
print(f"Hash: {dhash}")
print(f"Hash hex: {dhash.hash.flatten().tobytes().hex()}")

# Difference hash uses 9x8 image, compares horizontally
gray = img.convert('L')
small = gray.resize((9, 8), Image.Resampling.LANCZOS)
pixels = np.array(small)
print(f"\n9x8 grayscale pixels for dhash:")
for i in range(8):
    row = pixels[i]
    print(f"Row {i}: {row}")

# Compute differences: pixel[i] > pixel[i+1]
diff_matrix = pixels[:, :-1] > pixels[:, 1:]
print(f"\nDifference matrix (9 columns -> 8 comparisons per row):")
for i in range(8):
    row = diff_matrix[i]
    row_str = ''.join(['1' if x else '0' for x in row])
    print(f"Row {i}: {row_str}")

# Compute perception hash
print("\n--- Perception Hash ---")
phash = imagehash.phash(img)
print(f"Hash: {phash}")
print(f"Hash hex: {phash.hash.flatten().tobytes().hex()}")

# Compute wavelet hash
print("\n--- Wavelet Hash ---")
whash = imagehash.whash(img)
print(f"Hash: {whash}")
print(f"Hash hex: {whash.hash.flatten().tobytes().hex()}")

# Compute color hash
print("\n--- Color Hash (binbits=3) ---")
colorhash = imagehash.colorhash(img, binbits=3)
print(f"Hash: {colorhash}")
print(f"Hash hex (raw): {colorhash}")

# Let's see what the color hash actually contains
# Color hash returns a string like "07000000000"
print(f"Color hash string: {colorhash}")
print(f"Length: {len(str(colorhash))}")
print(f"Interpreted as hex: {int(str(colorhash), 16):#x}")