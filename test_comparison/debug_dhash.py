#!/usr/bin/env python3
import numpy as np

# Simulate a simple 3x3 image resized to 4x3 for dhash (hash_size=3)
# Python dhash with hash_size=3 would resize to 4x3
pixels = np.array([
    [10, 20, 30, 40],   # row 0
    [50, 60, 70, 80],   # row 1  
    [90, 100, 110, 120] # row 2
])

print("Pixels array (4x3):")
print(pixels)
print("\nPython dhash comparison: pixels[:, 1:] > pixels[:, :-1]")
diff = pixels[:, 1:] > pixels[:, :-1]
print("Result (3x3 boolean array):")
print(diff)
print("\nFlattened (row-major):")
flat = diff.flatten()
print(flat)
print("\nAs bits (1=True, 0=False):")
bits = [1 if x else 0 for x in flat]
print(bits)
print("\nAs hex (assuming 9 bits):")
# Convert to binary string
binary_str = ''.join(str(b) for b in bits)
hex_val = hex(int(binary_str, 2))[2:]
print(f"0x{hex_val}")