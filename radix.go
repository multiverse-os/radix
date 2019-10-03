package radix

import (
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
			Children: []*Node{},
		},
		mutex:     new(sync.RWMutex),
		Keys:      []string{},
		Edges:     []*Node{},
		Height:    0,
		NodeCount: 0,
	}
}

func (self *Tree) NodesAtDepth(depth int) []*Node {
	nodes := []*Node{}
	for _, node := range self.Nodes {
		if node.Depth == depth {
			nodes = append(nodes, node)
		}
	}
	return nodes
}

func (self *Tree) KeyExists(key string) bool {
	for _, existingKey := range self.Keys {
		if existingKey == key {
			return true
		}
	}
	return false
}

func (self *Tree) Add(key string, value interface{}) *Node {
	if len(key) == 0 {
		return &Node{}
	}
	self.mutex.Lock()
	defer self.mutex.Unlock()
	node := self.Root.add(key, value)

	if !self.KeyExists(key) && node != nil {
		if len(node.Children) == 0 {
			node.Type = Edge
			if self.Height < node.Depth {
				self.Height = node.Depth
			}
		} else {
			for _, child := range node.Children {
				child.Depth = child.Parent.Depth + 1
			}
		}
		// NOTE: Cache
		self.Nodes = append(self.Nodes, node)
		self.Keys = append(self.Keys, key)
		return node
	} else {
		return nil
	}
}

// TODO: Remaining issues:
//         1) Doesn't properly set depth in complex edge cases
func (self *Node) add(key string, value interface{}) *Node {
	if len(self.Children) != 0 {
		for _, child := range self.Children {
			commonPrefixIndex := longestCommonPrefixIndex(child.Key, key)
			if commonPrefixIndex == -1 {
				continue
			} else if commonPrefixIndex > 0 {
				if len(child.Key) > commonPrefixIndex {
					child.SplitKeyAtIndex(commonPrefixIndex)
				} else if len(child.Key) == commonPrefixIndex {
					return child.add(key[commonPrefixIndex:], value)
				}
				if len(key) == commonPrefixIndex {
					child.Value = value
					return child
				} else {
					return child.AddChild(key[commonPrefixIndex:], value)
				}
			}
		}
	}
	if self.Parent != nil {
		self.Depth = self.Parent.Depth + 1
	}
	return self.AddChild(key, value)
}

// TODO: Need to actually write a recursive function, this was just used for
// development.
func (self *Tree) String() string {
	tree := tree.New()
	for _, node := range self.Root.Children {
		if node.Type == Edge {
			tree.AddNode(node)
		} else {
			branch := tree.AddBranch(node)
			for _, child := range node.Children {
				if child.Type == Edge {
					branch.AddNode(child)
				} else {
					subbranch := branch.AddBranch(child)
					for _, subchild := range child.Children {
						if subchild.Type == Edge {
							subbranch.AddNode(subchild)
						} else {
							subsubbranch := subbranch.AddBranch(subchild)
							for _, subsubchild := range subchild.Children {
								if subsubchild.Type == Edge {
									subsubbranch.AddNode(subsubchild)
								} else {
									subsubsubbranch := subsubbranch.AddBranch(subsubchild)
									for _, subsubsubchild := range subsubchild.Children {
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
	return tree.String()
}
