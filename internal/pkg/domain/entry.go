package domain

import (
	"errors"
	"fmt"
)

var errInvalidKind = errors.New("invalid entry kind")
var errEmptyName = errors.New("entry name is required")

type Entry struct {
	Kind     Kind
	Name     string
	Instance string
}

func NewEntry(kind Kind, name string, instance string) (*Entry, error) {
	if !kind.IsValid() {
		return nil, fmt.Errorf("%w: %q", errInvalidKind, kind)
	}

	if name == "" {
		return nil, errEmptyName
	}

	return &Entry{
		Kind:     kind,
		Name:     name,
		Instance: instance,
	}, nil
}
