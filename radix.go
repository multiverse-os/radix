package radix

import (
	"sync"
)

// TODO: All non-edge nodes technically are associated or are the root of a
// subtree.
// TODO: Need a solid way to build a tree from any node
// TODO: Should probably have longest branch length
//LongestBranch int

type Tree struct {
	mutex     *sync.RWMutex
	Root      *Node
	Children  []*Node
	Height    int // Or MaxDepth
	NodeCount int // Total
}

// TODO: Would prefer a proper build tree from node function
func New() *Tree {
	root := &Node{
		Type:     Branch,
		Depth:    -1,
		Children: []*Node{},
	}
	return &Tree{
		Root:      root,
		mutex:     new(sync.RWMutex),
		Children:  root.Children,
		Height:    0,
		NodeCount: 0,
	}
}

func (self *Tree) Keys() (keys []string) {
	for _, child := range self.Children {
		keys = append(keys, child.Key)
	}
	return append(keys, self.Root.Key)
}

func (self *Tree) NodesAtDepth(depth int) []*Node {
	nodes := []*Node{}
	for _, node := range self.Children {
		if node.Depth == depth {
			nodes = append(nodes, node)
		}
	}
	return nodes
}

func (self *Tree) KeyExists(key string) bool {
	for _, child := range self.Children {
		if child.Key == key {
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
		self.Children = append(self.Children, node)
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

func (self *Tree) AddEdge(node *Node) {
	self.Children = append(self.Children, node)
}

// TODO: This is a gross function, needs recursive
func (self *Tree) AddBranch(child *Node) (subtree *Tree) {
	subtree = self.AddBranch(child)
	for _, subchild := range child.Children {
		if subchild.Type == Edge {
			subtree.AddEdge(subchild)
		} else {
			subtree = subtree.AddBranch(subchild)
			for _, subchild := range subchild.Children {
				if subchild.Type == Edge {
					subtree.AddEdge(subchild)
				} else {
					subtree = subtree.AddBranch(subchild)
					for _, subchild := range subchild.Children {
						if subchild.Type == Edge {
							subtree.AddEdge(subchild)
						} else {
							subtree.AddBranch(subchild)
						}
					}
				}
			}
		}
	}
	return subtree
}

func (self *Node) IsBranch() bool { return len(self.Children) != 0 }
func (self *Node) IsEdge() bool   { return len(self.Children) == 0 }

// TODO: Need to actually write a recursive function, this was just used for
// development.
func (self *Tree) String() string {
	tree := &Tree{}
	for _, node := range self.Root.Children {
		if node.IsBranch() {
			tree = tree.AddBranch(node)
			for _, child := range node.Children {
				if child.IsEdge() {
					tree.AddEdge(child)
				}
			}
		} else {
			tree.AddEdge(node)
		}
	}
	return tree.String()
}
