package main

import (
	"fmt"

	radix "github.com/multiverse-os/cli/radix"
)

func main() {
	fmt.Println("radix trie example")
	fmt.Println("==================")

	tree := radix.New()
	tree.Add("romane", 112)
	tree.Add("romaoos", 11)
	tree.Add("tuber", 44)
	tree.Add("romulus", 4)
	//tree.Add("ruber", 8)
	//tree.Add("tubular", 99)
	//tree.Add("rubens", 9)
	//tree.Add("tub", 3)
	//tree.Add("rubicon", 44)
	//tree.Add("rubicundus", 71)

	//fmt.Printf("%s", tree.String())

	nodes := tree.Root.Walk()
	fmt.Println("nodes:", nodes)
	for _, node := range nodes {
		fmt.Println("node:", string(node.Key()))
	}

	//keys, values := tree.PrefixSearch("rom")
	//fmt.Println("keys:", keys)
	//fmt.Println("values:", values)

}
