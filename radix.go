package radix

import (
	"fmt"
	"strings"
)

// Terminology
///////////////////////////////////////////////////////////////////////////////
// Root: Base of the Tree data object with a depth of 0
// Leaf: Terminal nodes at the edge of the tree object containing values
// Branch: Inner nodes not containing values

const fuzzyIterationLimit = 2

type Tree struct {
	Root        *Node
	Keys        []string
	Edges       []*Node
	Size        int
	Depth       int
	KeyCount    int
	NodeCount   int // Total
	LeafCount   int
	BranchCount int
}

func New() *Tree {
	return &Tree{
		Root:        &Node{},
		Keys:        []string{},
		Edges:       []*Node{},
		Size:        0,
		Depth:       0,
		LeafCount:   0,
		BranchCount: 0,
	}
}

func (self *Tree) FuzzySearch(query string) ([]string, []interface{}) {
	if len(self.Root.Children()) == 0 || query == "" {
		return []string{}, []interface{}{}
	}
	return self.fuzzySearch(self.keyToBytes(query),
		self.Root,
		0,
		0,
		0,
		[]byte{})
}

func (self *Tree) fuzzySearch(query []byte, node *Node, index int, iteration int, lastIncrement int, found []byte) ([]string, []interface{}) {
	searchBitMask := genBitMask(query[index:])

	if len(node.Children()) == 0 {
		return []string{}, []interface{}{}
	}

	startIndex := index
	for _, child := range node.Children() {
		// Reset index for each iteration of the child
		index = startIndex
		// If this is the case, then somewhere inside the depth of this
		// node there MIGHT exist what we're looking for, or it could
		// be shallow
		if child.IsBitMaskSet(searchBitMask) {
			// Iterate letters
			for _, letter := range child.Key() {
				if index < len(query) {
					if letter == query[index] {
						lastIncrement = iteration
						index++
					}
					// Small optimization, break early
					if index >= len(query) {
						break
					}
				}
				iteration++
			}
		} else {
			// Not set, can't do anything here really
			continue
		}
	}
	return nil, nil
}

func (self *Tree) PrefixSearch(query string) ([]string, []interface{}) {
	if len(self.Root.Children()) == 0 {
		return []string{}, []interface{}{}
	}
	// TODO: Node and prefix decladed and not used...
	node, prefix, ok := self.prefixSearch(self.keyToBytes(query), self.Root, 0, []byte{})
	if !ok {
		return []string{}, []interface{}{}
	}
	return []string{string(prefix)}, []interface{}{node.Value}
}

func (self *Tree) LongestPrefix(query string) (string, bool) {
	if len(self.Root.Children()) == 0 {
		return "", false
	}

	ok := false

	_, prefix, _ := self.prefixSearch(
		[]byte(query),
		self.Root,
		0,
		[]byte{})

	ok = (len(prefix) > 0)

	return string(prefix), ok
}

// Recursively prefix-searches to find the longest prefix that exists
func (self *Tree) prefixSearch(query []byte, node *Node, index int, found []byte) (*Node, []byte, bool) {
	if index+1 > len(query) {
		return node, found, true
	} else if len(node.Children()) == 0 {
		return node, found, false
	}

	for _, child := range node.Children() {
		lettersFound := 0
		searchLetter := query[index]

		for _, letter := range child.Key() {
			// A matching letter has been found
			if searchLetter == letter {
				lettersFound++
				// Otherwise recurse
				if (index+lettersFound) >= len(query) || len(child.Key()) == lettersFound {
					newIndex := index + len(child.Key())
					toAppend := child.Key()
					return self.prefixSearch(query, child, newIndex, append(found, toAppend...))
				}
				if index+lettersFound < len(query) {
					searchLetter = query[index+lettersFound]
					return child, child.Key(), false
				}
			} else {
				break
			}
		}
	}
	return nil, []byte{}, false
}

// The collection will, starting from a given node, recurse and generate
// strings from every leaf
func (self *Tree) collect(node *Node, prefix []byte) ([]string, []interface{}) {
	keys := []string{string(prefix)}
	values := []interface{}{node.Value}

	if len(node.Children()) == 0 {
		return keys, values
	}

	// Recursively append
	for _, child := range node.Children() {
		keys = append(keys, string(child.Key()))
		values = append(values, child.Value)
	}

	return keys, values
}

func (self *Tree) Add(key string, value interface{}) *Node {
	if key == "" {
		return &Node{}
	}
	input := self.keyToBytes(key)
	bitMask := genBitMask(input)
	leaf := self.add(self.Root, input, bitMask, 0)
	self.KeyCount++
	leaf.Value = value
	return leaf
}

func (self *Tree) add(node *Node, input []byte, bitMask uint32, depth int) *Node {
	if len(input) == 0 {
		return node
	} else if len(node.Children()) == 0 {
		self.NodeCount++
		return node.NewChild(input)
	}

	for childIndex, child := range node.Children() {
		for i := 0; i < len(child.Key()); i++ {
			if i > len(input) {
				break
			}

			var inputbyte byte
			if i+1 <= len(input) {
				inputbyte = input[i : i+1][0]
			}

			childbyte := child.Key()[i : i+1][0]
			if childbyte == inputbyte {
				child.OrBitMask(genBitMask(input[i:]))
				if i+1 == len(child.Key()) {
					if len(child.Children()) > 0 {
						return self.add(child, input[i+1:], bitMask, depth+1)
					}
				}
			} else {
				if i > 0 {
					self.NodeCount++
					child.Break(i)
					if len(input[i:]) > 0 {
						self.NodeCount++
						newNode := child.NewChild(input[i:])
						return newNode
					} else {
						// If the break is less than the input then
						// return the child (which is the parent of any
						// new child)
						return child
					}
				} else {
					// If there are more nodes to be seen, continue
					if childIndex+1 < len(node.Children()) {
						break
					}

					// If it's the first letter, just insert to node
					// (not child)
					self.NodeCount++
					newNode := node.NewChild(input[i:])
					return newNode
				}
			}
		}
	}
	return node
}

func (self *Tree) String() string {
	output := "\n"
	first := true

	self.Root.WalkDepthFirst(
		func(node *Node, depth int, firstAtDepth bool, lastAtDepth bool, numChildren int) terminate {

			if !first && firstAtDepth {
				output += strings.Repeat(" ", (depth*3)-3) + "|\n"
			}
			if depth > 0 {
				output += strings.Repeat(" ", (depth*3)-3) + "+- "
			}
			output += fmt.Sprintf("[%s]\n", string(node.Key()))

			first = false
			return terminate(false)
		}, 0)
	return output
}

// Because we're converting from utf8 down, we'll max out at 255 on the
// letter's value as to not overflow a byte
func (self *Tree) keyToBytes(key string) []byte {
	bytes := make([]byte, len(key))
	for i, letter := range key {
		if letter < rune(255) {
			bytes[i] = uint8(letter)
		}
	}
	return bytes
}
