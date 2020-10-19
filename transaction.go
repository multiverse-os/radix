package radix

import (
	"time"
)

// TODO: By storing the transaction optionally with the activity log, we can
// easily piece together a versioned datastore.
type Transaction struct {
	Parent    *Node
	CreatedAt time.Time
	//Signatures []key.Signature
	Key   []byte
	Value []byte
}
