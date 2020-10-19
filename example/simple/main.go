package main

import (
	"fmt"
	"strings"

	radix "../../../radix"
)

func main() {
	fmt.Println("radix trie example")
	fmt.Println("==================")

	tree := radix.New()
	tree.Add("romane", 112)
	fmt.Println("Adding romane")
	tree.Add("romul", 8)
	fmt.Println("Adding romul")
	tree.Add("romanous", 44)
	fmt.Println("Adding romanous")
	tree.Add("tubular", 99)
	fmt.Println("Adding tubular")
	tree.Add("rubens", 9)
	fmt.Println("Adding rubens")
	tree.Add("romaoos", 11)
	fmt.Println("Adding romaoos")
	tree.Add("ruber", 19)
	fmt.Println("Adding ruber")
	tree.Add("tub", 3)
	fmt.Println("Adding tub")
	tree.Add("tuber", 44)
	fmt.Println("Adding tuber")
	tree.Add("rubicon", 44)
	fmt.Println("Adding rubicon")
	tree.Add("rubicundus", 71)
	fmt.Println("Adding rubicundus")

	fmt.Println("=================================================")
	fmt.Println("tree:", tree.String())

	fmt.Println("=================================================")
	fmt.Println("edges:")
	for _, edge := range tree.Nodes {
		if len(edge.Children) == 0 {
			fmt.Println("fullKey:", strings.Join(edge.FullKey(), ""))
		}
	}
	fmt.Println("=================================================")
	ok, nodes := tree.Prefix("rom")
	if ok {
		for _, node := range nodes {
			fmt.Println("node:", node.String())
		}
	}

	fmt.Println("=================================================")
	//fmt.Println("keys:", keys)
	//fmt.Println("values:", values)

}
