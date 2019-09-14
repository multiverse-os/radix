package radix

type edges []edge

func (self edges) Len() int           { return len(self) }
func (self edges) Less(i, j int) bool { return self[i].key < self[j].key }
func (self edges) Swap(i, j int)      { self[i], self[j] = self[j], self[i] }
func (self edges) Sort()              { sort.Sort(self) }
