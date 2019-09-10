package radix

// This is specific to 32 bits (8 pairs of hex goodness)
func numBitsSet32(n uint32) int {
	var a uint32 = 0x55555555
	var b uint32 = 0x33333333
	var c uint32 = 0x0F0F0F0F
	var d uint32 = 0x01010101

	n = n - ((n >> 1) & a)
	n = (n & b) + ((n >> 2) & b)
	n = (((n + (n >> 4)) & c) * d) >> 24

	return int(n)
}

// Similar to the 32, except with variables doubled, in case we want to
// increase precision a little later on.
func numBitsSet64(n uint64) int {
	var a uint64 = 0x5555555555555555
	var b uint64 = 0x3333333333333333
	var c uint64 = 0x0F0F0F0F0F0F0F0F
	var d uint64 = 0x0101010101010101

	n = n - ((n >> 1) & a)
	n = (n & b) + ((n >> 2) & b)
	n = (((n + (n >> 4)) & c) * d) >> 56

	return int(n)
}

// Generates a bit mask based on the byte slice, compressing it to 32 bits.
// Priority is given to letters, numbers are compress by half, special
// characters are the final bit.
func genBitMask(str []byte) uint32 {
	mask := uint32(0)
	for _, r := range str {
		setBit := uint32(0)
		if r >= 'a' && r <= 'z' {
			// A-Z and a-z should map the same bits 1 - 26
			setBit = uint32(r - 97) // 'a' is 97
		} else if r >= 'A' && r <= 'Z' {
			// Fit into a-z range by negating to 0 index
			setBit = uint32(r - 65) // 'A' is 65
		} else if r >= '0' && r <= '9' {
			// Half the number and add to start value
			number := uint32(r - 48)
			setBit = (number / 2) + 26
		} else {
			// All other characters (special characters etc.) will appear as character 32
			setBit = 31
		}
		mask |= (1 << setBit)
	}
	return mask
}

func bitMaskContains(haystack, needle uint32) bool {
	return (haystack & needle) == needle
}
