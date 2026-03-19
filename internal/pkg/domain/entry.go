package domain

import (
	"errors"
	"fmt"
)

var errInvalidKind = errors.New("invalid entry kind")

type Entry struct {
	Kind Kind
}

func NewEntry(kind Kind) (*Entry, error) {
	if !kind.IsValid() {
		return nil, fmt.Errorf("%w: %q", errInvalidKind, kind)
	}

	return &Entry{
		Kind: kind,
	}, nil
}
