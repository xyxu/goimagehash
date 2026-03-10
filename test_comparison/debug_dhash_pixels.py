#!/usr/bin/env python3
from PIL import Image
import numpy as np

# Load sample image
img = Image.open('../_examples/sample1.jpg')

# Python's dhash process
gray = img.convert('L')
small = gray.resize((9, 8), Image.Resampling.LANCZOS)
pixels = np.array(small)

print("Python 9x8 grayscale pixels for dhash:")
for i in range(8):
    row = pixels[i]
    print(f"Row {i}: {row}")

print(f"\nPython differences (pixel[j+1] > pixel[j]):")
for i in range(8):
    diffs = pixels[i, 1:] > pixels[i, :-1]
    diff_str = ''.join(['1' if x else '0' for x in diffs])
    print(f"Row {i}: {diff_str}")
    
    # Also print actual comparisons
    print(f"  Comparisons:")
    for j in range(8):
        left = pixels[i, j]
        right = pixels[i, j+1]
        result = right > left
        print(f"    pixel[{j}]={left:3d} < pixel[{j+1}]={right:3d} ? {result}")

# Compute what Go should get
print(f"\n\nGo would compute (pixel[j] < pixel[j+1]):")
for i in range(8):
    diffs = pixels[i, :-1] < pixels[i, 1:]
    diff_str = ''.join(['1' if x else '0' for x in diffs])
    print(f"Row {i}: {diff_str}")
    
print("\nNote: Both should give same result since 'right > left' == 'left < right'")