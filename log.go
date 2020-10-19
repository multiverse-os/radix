package radix

import (
	"time"
)

type Activity uint8

const (
	Insert Activity = iota
	Delete
	Update
)

type ActivityLog struct {
	Timestamp    time.Time
	Activity     Activity
	Transactions []*Transaction
}
