#!/usr/bin/env python3
from PIL import Image
import imagehash

# Load the sample image
img = Image.open("../_examples/sample1.jpg")

# Compute hashes
dhash_val = imagehash.dhash(img)
print(f"DHash (sample1.jpg): {dhash_val}")

whash_val = imagehash.whash(img)
print(f"WHash (sample1.jpg): {whash_val}")

colorhash_val = imagehash.colorhash(img, binbits=3)
print(f"ColorHash (sample1.jpg, binbits=3): {colorhash_val}")

# For crop-resistant hash, we need to be careful with parameters
crophash_val = imagehash.crop_resistant_hash(img)
print(f"CropResistantHash (sample1.jpg): {crophash_val}")