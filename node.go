package radix

import (
	"fmt"
	"strconv"
)

// TODO: Consider implementing the root node of the tree as a Node with a tree
// interface. And example can be found in the text/tree library for printing
// trees.
// Path is used to recurse over the tree only visiting nodes
// which are above this node in the tree.
// Subs is used to recurse over the tree only visiting nodes
// which are directly under this node in the tree.

type terminate bool

type NodeType int

const (
	Branch NodeType = iota
	Edge
)

func (self NodeType) String() string {
	switch self {
	case Edge:
		return "edge"
	default: // Branch
		return "branch"
	}
}

type Node struct {
	Type     NodeType
	Depth    int
	Key      string
	Parent   *Node
	Children []*Node
	Value    interface{}
}

type ByteKeys []string

func (self ByteKeys) Len() int           { return len(self) }
func (self ByteKeys) Swap(i, j int)      { self[i], self[j] = self[j], self[i] }
func (self ByteKeys) Less(i, j int) bool { return len(self[i]) < len(self[j]) }

func NewNode(key string, value interface{}) *Node {
	return &Node{
		Key:   key,
		Value: value,
	}
}

func (self *Node) Ancestor(depth int) *Node {
	node := self
	for i := 0; i < depth; i++ {
		if node.Parent != nil {
			node = node.Parent
		} else {
			break
		}
	}
	return node
}

// TODO: This revealed that some of the more complex keys are not getting
// working. Its showing up in the tree output, but the parents are not correct
// and so rebuilding the full key leaves out the middle node of examples like
// "romul" and just ends up being "rul" because its missing "om"
func (self *Node) FullKey() (fullKey []string) {
	for i := (self.Depth + 1); i >= 0; i-- {
		fullKey = append(fullKey, self.Ancestor(i).Key)
	}
	return fullKey
}

func (self *Node) AddChild(key string, value interface{}) *Node {
	child := &Node{
		Key:      key,
		Value:    value,
		Children: []*Node{},
		Parent:   self,
		Depth:    (self.Depth + 1),
	}
	self.Children = append(self.Children, child)
	if len(self.Children) != 0 {
		self.Type = Branch
	}
	return child
}

// NOTE: This breaks up the keys as new Children are added
func (self *Node) SplitKeyAtIndex(index int) *Node {
	prefixKey := self.Key[:index]
	suffixKey := self.Key[index:]
	value := self.Value
	children := self.Children

	self.Key = prefixKey
	self.Value = nil
	self.Children = []*Node{}
	self.Depth = self.Parent.Depth + 1
	child := self.AddChild(suffixKey, value)

	child.Children = children
	if len(child.Children) == 0 {
		child.Type = Edge
	} else {
		for _, c := range children {
			c.Depth = c.Parent.Depth + 1
		}
	}
	child.Value = value
	return child
}

func (self *Node) Walk() (nodes []*Node) {
	nodes = []*Node{self}
	for _, child := range self.Children {
		nodes = append(nodes, child)
		if child.Type == Branch {
			nodes = append(nodes, child.Walk()...)
		}
	}
	return nodes
}

func (self *Node) String() string {
	return fmt.Sprintf("["+self.Type.String()+"][FullKey='"+string(self.Key)+"'][Key='"+string(self.Key)+"', value='%v'][depth='%v']", self.Value, self.Depth)
}

func (self *Node) JSON() string {
	return fmt.Sprintf(`{
	'Type':` + self.Type.String() + `',
	'Depth':` + strconv.Itoa(self.Depth) + `',
	'Key':` + string(self.Key) + `',
	'value': '%v',
}`)
}
