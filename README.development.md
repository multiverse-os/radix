
## Concepts


#### URI Concepts

If file name is of the form `tree://<ID>`, the entire tree, indicating it is the root. 

If file name is of the form `subtree://<ID>`, is the subtree, and can enable downloading speicifcally the subtree. (*Is subtree required if we are using tree?*)

If file name is of the form `treenode://<ID>`, is the tree node

## Organization Methods

#### Hashmaps
There should be a variety of ways of organizing the data, and it should be kept
in sync. For example `node[id]` should be a map that returns the nodes, and
`tree[id]` should be a map that returns the trees or subtrees. 

#### Edge-tracking
At all times the edges should be tracked and kept together, in a list or some
sort. 

*Edges should store the entire branch length.*

