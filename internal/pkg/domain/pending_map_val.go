package domain

import "time"

type PendingMapVal struct {
	state     int
	timestamp time.Time
	op        int
}

func NewPendingMapVal(state int, timestamp time.Time, op int) *PendingMapVal {
	return &PendingMapVal{
		state:     state,
		timestamp: timestamp,
		op:        op,
	}
}
