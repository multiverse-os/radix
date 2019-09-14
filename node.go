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
	Index    int // x coord system
	Depth    int // y
	key      string
	parent   *Node
	children []*Node
	Value    interface{}
}

type ByteKeys []string

func (self ByteKeys) Len() int           { return len(self) }
func (self ByteKeys) Swap(i, j int)      { self[i], self[j] = self[j], self[i] }
func (self ByteKeys) Less(i, j int) bool { return len(self[i]) < len(self[j]) }

func NewNode(key string, value interface{}) *Node {
	return &Node{
		key:   key,
		Value: value,
	}
}

func (self *Node) AddChild(key string, value interface{}) *Node {
	child := &Node{
		key:      key,
		Value:    value,
		children: []*Node{},
		parent:   self,
		Depth:    (self.Depth + 1),
	}
	self.children = append(self.children, child)
	if len(self.children) != 0 {
		self.Type = Branch
	}
	return child
}

// NOTE: This breaks up the keys as new children are added
func (self *Node) SplitKeyAtIndex(index int) *Node {
	prefixKey := self.key[:index]
	suffixKey := self.key[index:]
	value := self.Value
	children := self.children

	self.key = prefixKey
	self.Value = nil
	self.children = []*Node{}
	fmt.Println("self.key:", self.key)

	child := self.AddChild(suffixKey, value)
	child.children = children
	for _, c := range children {
		c.Depth += 1
		if len(c.children) == 0 {
			c.Type = Edge
		}
	}
	child.Value = value
	return child
}

func (self *Node) Walk() (nodes []*Node) {
	nodes = []*Node{self}
	for _, child := range self.children {
		nodes = append(nodes, child)
		if child.Type == Branch {
			nodes = append(nodes, child.Walk()...)
		}
	}
	return nodes
}

func (self *Node) String() string {
	return fmt.Sprintf("["+self.Type.String()+"][key='"+string(self.key)+"', value='%v'][depth='%v'][index='%v']", self.Value, self.Depth, self.Index)
}

func (self *Node) JSON() string {
	return fmt.Sprintf(`{
	'type':`+self.Type.String()+`',
	'depth':`+strconv.Itoa(self.Depth)+`',
	'key':`+string(self.key)+`',
	'value': '%v',
	'parent_key':`+string(self.parent.key)+`',
	'children_count':`+strconv.Itoa(len(self.children))+`',
}`, self.Value)
}
