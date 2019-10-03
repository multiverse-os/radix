package radix

import (
	"fmt"
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

func (self *Tree) Prefix(key string) (bool, []*Node) {
	nodes := []*Node{}
	fmt.Println("iterating through ", len(self.Nodes), " nodes:", self.Nodes)
	for _, child := range self.Nodes {
		fmt.Println("comparing:")
		fmt.Println(" child.Key            :", child.Key)
		if len(child.Key) < len(key) && child.Key == key[:len(child.Key)] {
			fmt.Println("  key[:len(child.Key)]:", key[:len(child.Key)])
			nodes = append(nodes, child)
			if child.Type == Branch {
				nodes = append(nodes, child.Walk()...)
			}

		}
	}
	return (len(nodes) > 0), nodes
}
