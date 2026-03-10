#!/usr/bin/env python3
from PIL import Image
import numpy as np

# PIL's convert('L') uses the following formula:
# L = R * 299/1000 + G * 587/1000 + B * 114/1000
# Or approximately: 0.299*R + 0.587*G + 0.114*B

# Go's Rgb2Gray uses:
# lum := 0.299*float64(r/257) + 0.587*float64(g/257) + 0.114*float64(b/256)

# The division by 257 vs 256 might cause differences!
# Also, Python uses integer arithmetic with the formula above.

# Let's test with a sample pixel
R, G, B = 100, 150, 200

# Python PIL convert('L') 
# Actually, let's check what PIL really does
img = Image.new('RGB', (1, 1), (R, G, B))
gray_img = img.convert('L')
pixel = gray_img.getpixel((0, 0))
print(f"PIL convert('L') for RGB({R},{G},{B}): {pixel}")

# Manual calculation using PIL's formula
# According to PIL documentation: L = R * 299/1000 + G * 587/1000 + B * 114/1000
manual = int(R * 299/1000 + G * 587/1000 + B * 114/1000)
print(f"Manual calculation: {manual}")

# Go's calculation (approximated in Python)
go_style = 0.299*(R/257) + 0.587*(G/257) + 0.114*(B/256)
go_value = int(go_style * 255)  # Scale to 0-255
print(f"Go-style calculation: {go_value}")
