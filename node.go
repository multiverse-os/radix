package radix

import (
	"errors"
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
	Root NodeType = iota
	Edge
	Connector
)

func (self NodeType) String() string {
	switch self {
	case Root:
		return "root"
	case Edge:
		return "edge"
	case Connector:
		return "connector"
	default:
		return ""
	}
}

// The all-important building block
// The child bytes is a bit map where 0-26 is A-Z and 27 - 32 is squashed into 0-9
// bitmask: Should probably be renamed to reflect that it is bytes version of key
type Node struct {
	Type     NodeType
	Index    int // x coord system
	Depth    int // y
	key      []byte
	bitMask  uint32
	parent   *Node
	children []*Node
	Value    interface{}
}

// NOTE: The benefit of doing this instead of just making the variables public
// is that this allows them to essentially be read-only, which becomes more
// important when dealing with pointers.
func (self *Node) Key() []byte                      { return self.key }
func (self *Node) Parent() *Node                    { return self.parent }
func (self *Node) Children() []*Node                { return self.children }
func (self *Node) OrBitMask(bitMask uint32)         { self.bitMask |= bitMask }
func (self *Node) IsBitMaskSet(bitMask uint32) bool { return bitMaskContains(self.bitMask, bitMask) }
func (self *Node) BitMask() uint32                  { return self.bitMask }

func NewNode(key []byte, value interface{}) *Node {
	return &Node{
		bitMask: genBitMask(key),
		key:     key,
		Value:   value,
	}
}

func (self *Node) NewChild(key []byte) *Node {
	child := &Node{
		key:     key,
		parent:  self,
		Depth:   (self.Depth + 1),
		bitMask: genBitMask(key),
	}
	self.children = append(self.children, child)
	return child
}

// NOTE: This breaks up the keys as new children are added
func (self *Node) SplitNode(index int) (*Node, error) {
	if index > len(self.Key()) {
		return nil, errors.New("Index exceeds key length")
	}

	prefixKey := self.Key()[:index]
	suffixKey := self.Key()[index:]
	fmt.Println("prefixKey:", string(prefixKey))
	fmt.Println("suffixKey:", string(suffixKey))
	value := self.Value
	children := self.Children()

	// Set the vars, move children and add the child
	// TODO: This should be mostly put inside of a newChild like function
	self.key = prefixKey
	self.Value = nil
	self.Type = Connector
	self.children = make([]*Node, 0)

	child := self.NewChild(suffixKey)
	child.parent = self
	child.children = children
	child.Value = value
	child.Type = Edge

	// Rebuild the child bit mask (contain itself and it's children)
	child.OrBitMask(genBitMask(child.Key()))
	for _, childsChild := range child.Children() {
		child.OrBitMask(childsChild.BitMask())
	}

	self.OrBitMask(genBitMask(suffixKey))

	return self, nil
}

func (self *Node) Walk() (nodes []*Node) {
	nodes = []*Node{self}
	for _, child := range self.Children() {
		nodes = append(nodes, child)
		fmt.Println("added node:", string(child.Key()))
		if child.Type == Connector {
			nodes = append(nodes, child.Walk()...)
		}
	}
	fmt.Println("[radix] found [", len(nodes), "] nodes in the tree")
	return nodes
}

func (self *Node) String() string {
	return fmt.Sprintf("["+self.Type.String()+"][key='"+string(self.Key())+"', value='%v']", self.Value)
}

func (self *Node) JSON() string {
	return fmt.Sprintf(`{
	'type':`+self.Type.String()+`',
	'depth':`+strconv.Itoa(self.Depth)+`',
	'key':`+string(self.Key())+`',
	'value': '%v',
	'parent_key':`+string(self.parent.Key())+`',
	'children_count':`+strconv.Itoa(len(self.children))+`',
}`, self.Value)
}
