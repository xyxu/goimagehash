#!/usr/bin/env python3
from PIL import Image
import numpy as np

# Load sample image
img = Image.open('../_examples/sample1.jpg')
print(f"Original size: {img.size}")
print(f"Original mode: {img.mode}")

# Python's approach: convert to grayscale FIRST, then resize
gray_first = img.convert('L')
gray_first_resized = gray_first.resize((8, 8), Image.Resampling.LANCZOS)
pixels_gray_first = np.array(gray_first_resized)
print(f"\nGray-first approach (Python):")
print(f"  Min pixel: {pixels_gray_first.min()}, Max: {pixels_gray_first.max()}, Mean: {pixels_gray_first.mean():.2f}")

# Go's approach: resize color image FIRST, then convert to grayscale
color_resized = img.resize((8, 8), Image.Resampling.LANCZOS)
color_resized_gray = color_resized.convert('L')
pixels_color_first = np.array(color_resized_gray)
print(f"\nColor-first approach (Go):")
print(f"  Min pixel: {pixels_color_first.min()}, Max: {pixels_color_first.max()}, Mean: {pixels_color_first.mean():.2f}")

# Compare pixel values
print(f"\nPixel differences (gray-first vs color-first):")
diff = pixels_gray_first - pixels_color_first
print(f"  Max difference: {abs(diff).max()}")
print(f"  Mean absolute difference: {np.mean(np.abs(diff)):.2f}")
print(f"  Different pixels: {np.sum(diff != 0)}/64")

# Show first few pixels
print(f"\nFirst 5 pixels comparison:")
for i in range(5):
    print(f"  Pixel {i}: Gray-first={pixels_gray_first.flatten()[i]}, Color-first={pixels_color_first.flatten()[i]}, diff={diff.flatten()[i]}")

# Check if this affects the hash
avg_gray = pixels_gray_first.mean()
avg_color = pixels_color_first.mean()
print(f"\nAverages:")
print(f"  Gray-first average: {avg_gray:.2f}")
print(f"  Color-first average: {avg_color:.2f}")
print(f"  Difference: {avg_gray - avg_color:.2f}")

# Compute what the hash would be
hash_gray = (pixels_gray_first > avg_gray).flatten()
hash_color = (pixels_color_first > avg_color).flatten()

print(f"\nHash differences:")
print(f"  Same bits: {np.sum(hash_gray == hash_color)}/64")
print(f"  Different bits: {np.sum(hash_gray != hash_color)}/64")

if np.any(hash_gray != hash_color):
    print(f"\nDifferent bit positions:")
    for i in range(64):
        if hash_gray[i] != hash_color[i]:
            row = i // 8
            col = i % 8
            print(f"  Bit {i} (row {row}, col {col}): Gray-first={hash_gray[i]}, Color-first={hash_color[i]}")