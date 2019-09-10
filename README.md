<img src="https://avatars2.githubusercontent.com/u/24763891?s=400&u=c1150e7da5667f47159d433d8e49dad99a364f5f&v=4"  width="256px" height="256px" align="right" alt="Multiverse OS Logo">

## Multiverse OS: `radix` trie library
**URL** [multiverse-os.org](https://multiverse-os.org)

A simple prefix sorting radix tree designed specifically for working with words, providing fuzzy search functionality, prefix filtering for autocomplete.

### Examples
A simple example using the `.String()` function provided which prints out the
current state of the trie. 

``` go
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
	tree.Add("romanus", 11)
	tree.Add("tuber", 44)
	tree.Add("romulus", 4)
	tree.Add("ruber", 8)
	tree.Add("tubular", 99)
	tree.Add("rubens", 9)
	tree.Add("tub", 3)
	tree.Add("rubicon", 44)
	tree.Add("rubicundus", 71)

	fmt.Printf("%s", tree.String())
}
```
