package radix

import (
	"fmt"
	"sync"

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
// Clone makes a copy of an existing trie.
// Items stored in both tries become shared, obviously.

// VisitSubtree works much like Visit, but it only visits nodes matching prefix.

type Tree struct {
	//Size      int // Byte Size
	mutex     *sync.RWMutex
	Root      *Node
	Keys      []string
	Edges     []*Node
	Nodes     []*Node
	Height    int // Or MaxDepth
	NodeCount int // Total
}

func New() *Tree {
	return &Tree{
		Root: &Node{
			Type:     Branch,
			Depth:    -1,
			Index:    0,
			children: []*Node{},
		},
		mutex:     new(sync.RWMutex),
		Keys:      []string{},
		Edges:     []*Node{},
		Height:    0,
		NodeCount: 0,
	}
}

func (self *Tree) Add(key string, value interface{}) *Node {
	if key == "" {
		return &Node{}
	}
	self.mutex.Lock()
	defer self.mutex.Unlock()

	fmt.Println("tree:", self.String())
	fmt.Println("root children?:", len(self.Root.children))

	node := self.Root.add(key, value)

	if node != nil && len(node.children) == 0 {
		node.Type = Edge
		if self.Height < node.Depth {
			self.Height = node.Depth
		}
	}
	return node
}

// TODO: Remaining issues:
// 1) Doesn't properly set depth in complex edge cases
// 2) Doesn't always set Edge correctly
// 3) If a key has items below it, the key gets stored as "" below itself.
func (self *Node) add(key string, value interface{}) *Node {
	//fmt.Println("Merging in key:", key)
	if len(self.children) != 0 {
		for _, child := range self.children {
			commonPrefixIndex := longestCommonPrefixIndex(child.key, key)
			if commonPrefixIndex == -1 {
				continue
			} else if commonPrefixIndex > 0 {
				if len(child.key) > commonPrefixIndex {
					child.SplitKeyAtIndex(commonPrefixIndex)
				} else if len(child.key) == commonPrefixIndex {
					return child.add(key[commonPrefixIndex:], value)
				}
				return child.AddChild(key[commonPrefixIndex:], value)
			}
		}
	}
	return self.AddChild(key, value)
}

func (self *Tree) String() string {
	tree := tree.New()

	fmt.Println("Height of the tree is:", self.Height)
	fmt.Println("children of root node:", len(self.Root.children))
	for _, node := range self.Root.children {
		if node.Type == Edge {
			tree.AddNode(node)
		} else {
			branch := tree.AddBranch(node)
			for _, child := range node.children {
				if child.Type == Edge {
					branch.AddNode(child)
				} else {
					subbranch := branch.AddBranch(child)
					for _, subchild := range child.children {
						if subchild.Type == Edge {
							subbranch.AddNode(subchild)
						} else {
							subsubbranch := subbranch.AddBranch(subchild)
							for _, subsubchild := range subchild.children {
								if subsubchild.Type == Edge {
									subsubbranch.AddNode(subsubchild)
								} else {
									subsubsubbranch := subsubbranch.AddBranch(subsubchild)
									for _, subsubsubchild := range subsubchild.children {
										if subsubsubchild.Type == Edge {
											subsubsubbranch.AddNode(subsubsubchild)
										} else {
											subsubsubbranch.AddBranch(subsubsubchild)
										}
									}
								}
							}
						}
					}
				}
			}
		}
	}
	//}

	//tree.AddNode(fmt.Sprintf("['LEAF':{'key':'%s', 'value':'%v'}]", string(node.Key()), node.Value))
	//tree.AddNode(fmt.Sprintf("[{'key':'%s'}]", string(node.Key())))

	return tree.String()
}
