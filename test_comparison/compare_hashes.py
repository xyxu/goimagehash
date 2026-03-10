#!/usr/bin/env python3

go_hash = 0xe14f6f67af0f0b89
py_hash = 0xe89f6f63af4f8b69

print(f"Go hash:  0x{go_hash:016x}")
print(f"Py hash:  0x{py_hash:016x}")
print(f"Go binary: {go_hash:064b}")
print(f"Py binary: {py_hash:064b}")

# XOR to see differences
xor = go_hash ^ py_hash
print(f"\nXOR:       0x{xor:016x}")
print(f"XOR binary: {xor:064b}")
print(f"Different bits: {bin(xor).count('1')}")

# Check if it's a bit reversal
# Let's see if reversing bits makes them match
def reverse_bits_64(n):
    result = 0
    for i in range(64):
        if n & (1 << i):
            result |= 1 << (63 - i)
    return result

go_rev = reverse_bits_64(go_hash)
print(f"\nGo reversed: 0x{go_rev:016x}")
print(f"Py original:  0x{py_hash:016x}")
print(f"Match? {go_rev == py_hash}")

# Check if it's byte reversal
import struct
go_bytes = struct.pack('>Q', go_hash)
py_bytes = struct.pack('>Q', py_hash)
go_rev_bytes = go_bytes[::-1]
go_rev_from_bytes = struct.unpack('>Q', go_rev_bytes)[0]
print(f"\nGo bytes reversed: 0x{go_rev_from_bytes:016x}")
print(f"Match? {go_rev_from_bytes == py_hash}")
