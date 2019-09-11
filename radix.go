package radix

import (
	tree "github.com/multiverse-os/cli/text/tree"
)

// Terminology
///////////////////////////////////////////////////////////////////////////////
// Root: Base of the Tree data object with a depth of 0
// Leaf: Terminal nodes at the edge of the tree object containing values
// Branch: Inner nodes not containing values

// TODO:TASKS:
///////////////////////////////////////////////////////////////////////////////
// * Add subtree isolation and ability to run commands, so you can search,
// iterate, delete, a subtree
// * Ability to delete a key
// * Track all edges
// * Use levestian to improve functionality
//
//

type Tree struct {
	//Size      int // Byte Size
	Root      *Node
	Keys      [][]byte
	Edges     []*Node
	Height    int // Or MaxDepth
	NodeCount int // Total
}

func New() *Tree {
	return &Tree{
		Root: &Node{
			Type:  Root,
			key:   []byte{0},
			Depth: -1,
			Index: 0,
		},
		Keys:      [][]byte{},
		Edges:     []*Node{},
		Height:    0,
		NodeCount: 0,
	}
}

func (self *Tree) FuzzySearch(query string) ([]string, []interface{}) {
	if len(self.Root.Children()) == 0 || query == "" {
		return []string{}, []interface{}{}
	}
	return self.fuzzySearch(self.byteKey(query),
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
	node, prefix, ok := self.prefixSearch(self.byteKey(query), self.Root, 0, []byte{})
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

	for _, child := range node.Children() {
		keys = append(keys, string(child.Key()))
		values = append(values, child.Value)
	}
	return keys, values
}

func (self *Tree) Remove(key string) (bool, error) {
	// 1) Look up node
	// 2) Remove node
	// 3) Cleanup edges if it was edge
	// 4) Fix counts
	return false, nil
}

func (self *Tree) CacheKey(key []byte) {
	if len(key) > 0 {
		self.Keys = append(self.Keys, key)
		self.NodeCount = len(self.Keys)
	}
}

func (self *Tree) CacheEdge(node *Node) {
	if node.Value != nil {
		self.Edges = append(self.Edges, node)
		node.Type = Edge
	}
}

func (self *Tree) Add(key string, value interface{}) *Node {
	if key == "" {
		return &Node{}
	}
	input := self.byteKey(key)
	bitMask := genBitMask(input)

	leaf := self.add(self.Root, input, bitMask, 0)

	leaf.Value = value
	self.CacheEdge(leaf)
	self.CacheKey(leaf.Key())

	if leaf.Parent() == nil {
		leaf.Depth = 0
		leaf.parent = self.Root
	}
	return leaf
}

func (self *Tree) add(node *Node, input []byte, bitMask uint32, depth int) *Node {
	if len(input) == 0 {
		return node
	} else if len(node.Children()) == 0 {
		newChild := node.NewChild(input)
		self.CacheKey(newChild.Key())
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
					// NOTE: This fixed the issue with the first node breaking
					child.Value = node.Value
					child.SplitNode(i)

					if len(input[i:]) > 0 {
						newNode := child.NewChild(input[i:])
						self.CacheKey(newNode.Key())

						return newNode
					} else {
						return child
					}
				} else {
					if childIndex+1 < len(node.Children()) {
						break
					}

					newNode := node.NewChild(input[i:])
					self.CacheKey(newNode.Key())
					return newNode
				}
			}
		}
	}
	//fmt.Println("Node: bottom ", string(node.Key()), "]")
	return node
}

func (self *Tree) String() string {
	tree := tree.New()

	//tree.AddNode(fmt.Sprintf("['LEAF':{'key':'%s', 'value':'%v'}]", string(node.Key()), node.Value))
	//tree.AddNode(fmt.Sprintf("[{'key':'%s'}]", string(node.Key())))

	return tree.String()
}

// Because we're converting from utf8 down, we'll max out at 255 on the
// letter's value as to not overflow a byte
func (self *Tree) byteKey(key string) []byte {
	bytes := make([]byte, len(key))
	for i, letter := range key {
		if letter < rune(255) {
			bytes[i] = uint8(letter)
		}
	}
	return bytes
}
