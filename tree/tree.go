package tree

import (
	"bytes"
	"fmt"
	"io"
	"reflect"

	color "github.com/multiverse-os/cli/text/ansi/color"
)

type EdgeType string

var (
	EdgeTypeLink EdgeType = color.White("│")
	EdgeTypeMid  EdgeType = color.White("├──")
	EdgeTypeEnd  EdgeType = color.White("└──")
)

func New() Tree {
	return &node{Value: color.SkyBlue(".")}
}

type node struct {
	Root  *node
	Value string
	Nodes []*node
}

func (self *node) AddNode(value string) Tree {
	self.Nodes = append(self.Nodes, &node{
		Root:  self,
		Value: value,
	})
	if self.Root != nil {
		return self.Root
	}
	return self
}

func (self *node) AddBranch(value string) Tree {
	branch := &node{
		Value: value,
	}
	self.Nodes = append(self.Nodes, branch)
	return branch
}

func (self *node) Branch() Tree {
	self.Root = nil
	return self
}

func printNodes(writer io.Writer, level int, levelsEnded []int, nodes []*node) {
	for i, node := range nodes {
		edge := EdgeTypeMid
		if i == len(nodes)-1 {
			levelsEnded = append(levelsEnded, level)
			edge = EdgeTypeEnd
		}
		printValues(writer, level, levelsEnded, edge, node.Value)
		if len(node.Nodes) > 0 {
			printNodes(writer, level+1, levelsEnded, node.Nodes)
		}
	}
}

func printValues(writer io.Writer, level int, levelsEnded []int, edge EdgeType, value string) {
	for i := 0; i < level; i++ {
		if isEnded(levelsEnded, i) {
			fmt.Fprint(writer, "    ")
			continue
		}
		fmt.Fprintf(writer, "%s   ", EdgeTypeLink)
	}
	fmt.Fprintf(writer, "%s %v\n", edge, value)
}

func isEnded(levelsEnded []int, level int) bool {
	for _, l := range levelsEnded {
		if l == level {
			return true
		}
	}
	return false
}
