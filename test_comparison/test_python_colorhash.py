#!/usr/bin/env python3
from PIL import Image
import imagehash

img = Image.open("../_examples/sample1.jpg")

colorhash1 = imagehash.colorhash(img, binbits=1)
print(f"ColorHash binbits=1: {colorhash1}")

colorhash3 = imagehash.colorhash(img, binbits=3)
print(f"ColorHash binbits=3: {colorhash3}")

# Also get the binary array
print(f"\nbinbits=1 hash array shape: {colorhash1.hash.shape}")
print(f"binbits=1 hash array: {colorhash1.hash.flatten()}")
