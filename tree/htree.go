package tree

import (
	"errors"
	"fmt"
	"strings"
	"unicode/utf8"
)

func PrintHr(root Node) {
	fmt.Print(SprintHr(root))
}

func SprintHr(root Node) (s string) {
	for _, line := range lines(root) {
		// ignore runes before root node
		line = string([]rune(line)[2:])
		s += strings.TrimRight(line, " ") + "\n"
	}
	return
}

func lines(root Node) (s []string) {
	data := fmt.Sprintf("%s %v ", BoxHor, root.Data())
	l := len(root.Children())
	if l == 0 {
		s = append(s, data)
		return
	}

	w := utf8.RuneCountInString(data)
	for i, c := range root.Children() {
		for j, line := range lines(c) {
			if i == 0 && j == 0 {
				if l == 1 {
					s = append(s, data+BoxHor+line)
				} else {
					s = append(s, data+BoxDownHor+line)
				}
				continue
			}

			var box string
			if i == l-1 && j == 0 {
				// first line of the last child
				box = BoxUpRight
			} else if i == l-1 {
				box = " "
			} else if j == 0 {
				box = BoxVerRight
			} else {
				box = BoxVer
			}
			s = append(s, strings.Repeat(" ", w)+box+line)
		}
	}
	return
}

// PrintHrn prints the horizontal-newline formatted tree to standard output.
func PrintHrn(root Node) {
	fmt.Print(SprintHrn(root))
}

// SprintHrn returns the horizontal-newline formatted tree.
func SprintHrn(root Node) (s string) {
	return strings.Join(lines2(root), "\n") + "\n"
}

func lines2(root Node) (s []string) {
	s = append(s, fmt.Sprintf("%v", root.Data()))
	l := len(root.Children())
	if l == 0 {
		return
	}

	for i, c := range root.Children() {
		for j, line := range lines2(c) {
			// first line of the last child
			if i == l-1 && j == 0 {
				s = append(s, BoxUpRight+BoxHor+" "+line)
			} else if j == 0 {
				s = append(s, BoxVerRight+BoxHor+" "+line)
			} else if i == l-1 {
				s = append(s, "   "+line)
			} else {
				s = append(s, BoxVer+"  "+line)
			}
		}
	}
	return
}

const (
	BoxVer       = "│"
	BoxHor       = "─"
	BoxVerRight  = "├"
	BoxDownLeft  = "┐"
	BoxDownRight = "┌"
	BoxDownHor   = "┬"
	BoxUpRight   = "└"
	// Gutter is number of spaces between two adjacent child nodes.
	Gutter = 2
)

// ErrDuplicateNode indicates that a duplicate Node (node with same hash) was
// encountered while going through the tree. As of now Sprint/Print and
// SprintWithError/PrintWithError cannot operate on such trees.
//
// This error is returned by SprintWithError/PrintWithError. It's also used
// in Sprint/Print as error for panic for the same case.
//
// FIXME: create internal representation of trees that copies data
var ErrDuplicateNode = errors.New("duplicate node")

// Node represents a node in a tree. Type that satisfies Node must be a hashable type.
type Node interface {
	// Data must return a value representing the node. It is stringified using "%v".
	// If empty, a space is used.
	Data() interface{}
	// Children must return a list of all child nodes of the node.
	Children() []Node
}

type queue struct {
	arr []Node
}

func (q queue) empty() bool {
	return len(q.arr) == 0
}

func (q queue) len() int {
	return len(q.arr)
}

func (q *queue) push(n Node) {
	q.arr = append(q.arr, n)
}

func (q *queue) pop() Node {
	if q.empty() {
		return nil
	}
	ele := q.arr[0]
	q.arr = q.arr[1:]
	return ele
}

func (q *queue) peek() Node {
	if q.empty() {
		return nil
	}
	return q.arr[0]
}

// Print prints the formatted tree to standard output. To handle ErrDuplicateNode use PrintWithError.
func Print(root Node) {
	fmt.Print(Sprint(root))
}

// Sprint returns the formatted tree. To handle ErrDuplicateNode use SprintWithError.
func Sprint(root Node) string {
	parents := map[Node]Node{}
	if err := setParents(parents, root); err != nil {
		panic(err)
	}
	return sprint(parents, root)
}

// PrintWithError prints the formatted tree to standard output.
func PrintWithError(root Node) error {
	s, err := SprintWithError(root)
	if err != nil {
		return err
	}
	fmt.Print(s)
	return nil
}

// SprintWithError returns the formatted tree.
func SprintWithError(root Node) (string, error) {
	parents := map[Node]Node{}
	if err := setParents(parents, root); err != nil {
		return "", err
	}
	return sprint(parents, root), nil
}

func sprint(parents map[Node]Node, root Node) string {
	isLeftMostChild := func(n Node) bool {
		p, ok := parents[n]
		if !ok {
			// root
			return true
		}
		return p.Children()[0] == n
	}

	paddings := map[Node]int{}
	setPaddings(paddings, map[Node]int{}, 0, root)

	q := queue{}
	q.push(root)
	lines := []string{}
	for !q.empty() {
		// line storing branches, and line storing nodes
		branches, nodes := "", ""
		// runes covered
		covered := 0
		qLen := q.len()
		for i := 0; i < qLen; i++ {
			n := q.pop()
			for _, c := range n.Children() {
				q.push(c)
			}

			spaces := paddings[n] - covered
			data := safeData(n)
			nodes += strings.Repeat(" ", spaces) + data

			w := utf8.RuneCountInString(data)
			covered += spaces + w
			current, next := isLeftMostChild(n), isLeftMostChild(q.peek())
			if current {
				branches += strings.Repeat(" ", spaces)
			} else {
				branches += strings.Repeat(BoxHor, spaces)
			}

			if current && next {
				branches += BoxVer
			} else if current {
				branches += BoxVerRight
			} else if next {
				branches += BoxDownLeft
			} else {
				branches += BoxDownHor
			}

			if next {
				branches += strings.Repeat(" ", w-1)
			} else {
				branches += strings.Repeat(BoxHor, w-1)
			}
		}
		lines = append(lines, branches, nodes)
	}

	s := ""
	// ignore first line since it's the branch above root
	for _, line := range lines[1:] {
		s += strings.TrimRight(line, " ") + "\n"

	}
	return s
}

// safeData always returns non-empty representation of n's data. Empty data
// messes up tree structure, and ignoring such node will return incomplete
// tree output (tree without an entire subtree). So it returns a space.
func safeData(n Node) string {
	data := fmt.Sprintf("%v", n.Data())
	if data == "" {
		return " "
	}
	return data
}

// setPaddings sets left padding (distance of a node from the root)
// for each node in the tree.
func setPaddings(paddings map[Node]int, widths map[Node]int, pad int, root Node) {
	for _, c := range root.Children() {
		paddings[c] = pad
		setPaddings(paddings, widths, pad, c)
		pad += width(widths, c)
	}
}

// setParents sets child-parent relationships for the tree rooted
// at root.
func setParents(parents map[Node]Node, root Node) error {
	for _, c := range root.Children() {
		if _, ok := parents[c]; ok {
			return ErrDuplicateNode
		}
		parents[c] = root
		if err := setParents(parents, c); err != nil {
			return err
		}
	}
	return nil
}

// width returns either the sum of widths of it's children or its own
// data length depending on which one is bigger. widths is used in
// memoization.
func width(widths map[Node]int, n Node) int {
	if w, ok := widths[n]; ok {
		return w
	}

	w := utf8.RuneCountInString(safeData(n)) + Gutter
	widths[n] = w
	if len(n.Children()) == 0 {
		return w
	}

	sum := 0
	for _, c := range n.Children() {
		sum += width(widths, c)
	}
	if sum > w {
		widths[n] = sum
		return sum
	}
	return w
}
