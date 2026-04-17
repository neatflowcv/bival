package domain

import "time"

type PendingMap struct {
	key       string
	state     int
	timestamp time.Time
	op        int
}

func NewPendingMap(key string, state int, timestamp time.Time, op int) *PendingMap {
	return &PendingMap{
		key:       key,
		state:     state,
		timestamp: timestamp,
		op:        op,
	}
}
