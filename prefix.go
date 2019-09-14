package radix

import (
	"strings"
)

func longestCommonPrefixIndex(key1, key2 string) int {
	for i := len(key2); i >= 1; i-- {
		subKey2 := key2[:i]
		splitStrings := strings.SplitAfter(key1, subKey2)
		if len(splitStrings) == 2 {
			return len(splitStrings[0])
		}
	}
	return -1
}
