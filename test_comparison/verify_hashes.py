#!/usr/bin/env python3
"""
Generate reference hash values from Python imagehash for comparison with Go implementation.
"""
from PIL import Image
import imagehash
import json
import os

def generate_reference_hashes(image_path):
    """Generate reference hash values for all algorithms."""
    img = Image.open(image_path)
    
    results = {
        "image": os.path.basename(image_path),
        "hashes": {}
    }
    
    # Average hash
    ahash = imagehash.average_hash(img)
    results["hashes"]["average"] = str(ahash)
    
    # Difference hash
    dhash = imagehash.dhash(img)
    results["hashes"]["difference"] = str(dhash)
    
    # Perception hash
    phash = imagehash.phash(img)
    results["hashes"]["perception"] = str(phash)
    
    # Wavelet hash
    whash = imagehash.whash(img)
    results["hashes"]["wavelet"] = str(whash)
    
    # Color hash (binbits=3)
    colorhash_val = imagehash.colorhash(img, binbits=3)
    results["hashes"]["color"] = str(colorhash_val)
    
    # Color hash (binbits=1) for debugging
    colorhash1 = imagehash.colorhash(img, binbits=1)
    results["hashes"]["color_binbits1"] = str(colorhash1)
    
    # Crop-resistant hash (with default parameters)
    crophash = imagehash.crop_resistant_hash(img)
    results["hashes"]["crop_resistant"] = str(crophash)
    
    return results

def main():
    # Test with sample images
    sample_dir = "../_examples"
    samples = ["sample1.jpg", "sample2.jpg", "sample3.jpg", "sample4.jpg"]
    
    all_results = {}
    for sample in samples:
        path = os.path.join(sample_dir, sample)
        if os.path.exists(path):
            print(f"Processing {sample}...")
            all_results[sample] = generate_reference_hashes(path)
    
    # Save to JSON file
    with open("python_reference_hashes.json", "w") as f:
        json.dump(all_results, f, indent=2)
    
    print(f"\nSaved reference hashes to python_reference_hashes.json")
    
    # Also print summary
    print("\nReference hash values (sample1.jpg):")
    sample1 = all_results.get("sample1.jpg", {})
    for algo, hash_val in sample1.get("hashes", {}).items():
        print(f"  {algo:20} {hash_val}")

if __name__ == "__main__":
    main()