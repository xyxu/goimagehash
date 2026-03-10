#!/usr/bin/env python3
from PIL import Image
import numpy as np

# Test with a specific RGB color
R, G, B = 100, 150, 200

# Create 1x1 image
img = Image.new('RGB', (1, 1), (R, G, B))

# PIL convert('L')
gray_img = img.convert('L')
pil_value = gray_img.getpixel((0, 0))
print(f"PIL convert('L') for RGB({R},{G},{B}): {pil_value}")

# PIL's formula (from documentation)
# L = R * 299/1000 + G * 587/1000 + B * 114/1000
pil_formula = int(R * 299/1000 + G * 587/1000 + B * 114/1000)
print(f"PIL formula calculation: {pil_formula}")

# Go's formula (approximated)
# Go gets RGBA values in range 0-65535, divides by 257 to get 0-255
# Then: 0.299*r + 0.587*g + 0.114*b
r_go = R * 257  # Scale 0-255 to 0-65535
g_go = G * 257
b_go = B * 257

# Go divides by 257 first
r_norm = r_go / 257
g_norm = g_go / 257  
b_norm = b_go / 257

go_value = 0.299 * r_norm + 0.587 * g_norm + 0.114 * b_norm
print(f"Go-style calculation (exact): {go_value}")
print(f"Go-style rounded: {int(go_value)}")

# Actually, Go uses integer division: r/257
r_div = r_go // 257
g_div = g_go // 257
b_div = b_go // 257
go_int_value = 0.299 * r_div + 0.587 * g_div + 0.114 * b_div
print(f"Go with integer division: {go_int_value}")
print(f"Go with integer division rounded: {int(go_int_value)}")

# Test with actual image pixels from sample1.jpg
print("\n=== Testing with actual image ===")
img = Image.open('../_examples/sample1.jpg')

# Get a few sample pixels
pixels = []
for x in [0, 100, 200]:
    for y in [0, 100, 200]:
        if x < img.width and y < img.height:
            r, g, b = img.getpixel((x, y))
            pixels.append((r, g, b))

print(f"Sample pixels (R,G,B):")
for i, (r, g, b) in enumerate(pixels[:5]):
    print(f"  Pixel {i}: ({r},{g},{b})")
    
    # PIL value
    test_img = Image.new('RGB', (1, 1), (r, g, b))
    pil_val = test_img.convert('L').getpixel((0, 0))
    
    # Go calculation
    r_go = r * 257
    g_go = g * 257
    b_go = b * 257
    r_div = r_go // 257
    g_div = g_go // 257
    b_div = b_go // 257
    go_val = 0.299 * r_div + 0.587 * g_div + 0.114 * b_div
    
    print(f"    PIL: {pil_val}, Go: {go_val:.1f}, diff: {abs(pil_val - go_val):.1f}")