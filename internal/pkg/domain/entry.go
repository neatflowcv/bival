package domain

import (
	"errors"
	"fmt"
	"time"
)

var errInvalidKind = errors.New("invalid entry kind")
var errEmptyName = errors.New("entry name is required")

type Entry struct {
	kind     Kind
	name     string
	instance string
	mtime    time.Time
}

func NewEntry(kind Kind, name string, instance string, mtime time.Time) (*Entry, error) {
	if !kind.IsValid() {
		return nil, fmt.Errorf("%w: %q", errInvalidKind, kind)
	}

	if name == "" {
		return nil, errEmptyName
	}

	return &Entry{
		kind:     kind,
		name:     name,
		instance: instance,
		mtime:    mtime,
	}, nil
}

func (e *Entry) Kind() Kind {
	return e.kind
}

func (e *Entry) Name() string {
	return e.name
}

func (e *Entry) Instance() string {
	return e.instance
}

func (e *Entry) MTime() time.Time {
	return e.mtime
}
