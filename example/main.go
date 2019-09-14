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
	tree.Add("romul", 8)
	tree.Add("romanous", 44)
	tree.Add("tubular", 99)
	tree.Add("rubens", 9)
	tree.Add("romaoos", 11)
	tree.Add("ruber", 19)
	tree.Add("tub", 3)
	tree.Add("tuber", 44)
	tree.Add("rubicon", 44)
	tree.Add("rubicundus", 71)

	fmt.Printf("%s", tree.String())

	//keys, values := tree.PrefixSearch("rom")
	//fmt.Println("keys:", keys)
	//fmt.Println("values:", values)

}
