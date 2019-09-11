package radix

import "errors"

type terminate bool
type walkerFunc func(*Node, int, bool, bool, int) terminate

// The all-important building block
// The child bytes is a bit map where 0-26 is A-Z and 27 - 32 is squashed into 0-9
// bitmask: Should probably be renamed to reflect that it is bytes version of key
type Node struct {
	key        []byte
	bitMask    uint32
	childbytes int32
	parent     *Node
	children   []*Node
	Value      interface{}
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

func (self *Node) NewChild(key []byte) *Node {
	child := &Node{
		key:        key,
		childbytes: 0,
		parent:     self,
		bitMask:    genBitMask(key),
	}
	self.children = append(self.children, child)
	return child
}

func (self *Node) Break(index int) (*Node, error) {
	if index > len(self.Key()) {
		return nil, errors.New("Index exceeds key length")
	}
	prefixKey := self.Key()[:index]
	suffixKey := self.Key()[index:]
	value := self.Value
	children := self.Children()

	// Set the vars, move children and add the child
	self.key = prefixKey
	self.Value = nil
	self.children = make([]*Node, 0)

	child := self.NewChild(suffixKey)
	child.children = children
	child.Value = value

	// Rebuild the child bit mask (contain itself and it's children)
	child.OrBitMask(genBitMask(child.Key()))
	for _, childsChild := range child.Children() {
		child.OrBitMask(childsChild.BitMask())
	}

	self.OrBitMask(genBitMask(suffixKey))

	return self, nil
}

func (self *Node) WalkDepthFirst(walk walkerFunc, depth int) {
	isFirst := true
	currentChildCount := len(self.Children())
	for i, childNode := range self.Children() {
		isLast := false
		childCount := len(childNode.Children())

		isLast = (i == currentChildCount-1)
		// TODO: Is this isFirst and isLast required because it seems clunky
		stop := walk(childNode, depth, isFirst, isLast, childCount)

		isFirst = false
		if stop == terminate(true) {
			return
		}
		childNode.WalkDepthFirst(walk, depth+1)
	}
}
