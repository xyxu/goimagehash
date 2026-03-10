#!/usr/bin/env python3
import inspect
import imagehash

# Get the source code of average_hash function
print("=== Python imagehash average_hash source ===")
try:
    source = inspect.getsource(imagehash.average_hash)
    print(source)
except Exception as e:
    print(f"Could not get source: {e}")

print("\n=== Python imagehash dhash source ===")
try:
    source = inspect.getsource(imagehash.dhash)
    print(source)
except Exception as e:
    print(f"Could not get source: {e}")

print("\n=== Checking imagehash module location ===")
print(f"imagehash module: {imagehash.__file__}")

# Try to read the source file directly
import os
module_path = imagehash.__file__
if module_path.endswith('.pyc'):
    module_path = module_path[:-1]  # Remove .pyc to get .py

if os.path.exists(module_path):
    print(f"\nReading source file: {module_path}")
    with open(module_path, 'r') as f:
        lines = f.readlines()
        # Look for average_hash function
        in_function = False
        function_lines = []
        for i, line in enumerate(lines):
            if 'def average_hash' in line:
                in_function = True
            if in_function:
                function_lines.append(line)
                if line.strip() == '' and i > 0 and 'def ' in lines[i-1]:
                    # Empty line after function definition
                    break
        print("".join(function_lines[:50]))  # First 50 lines of function