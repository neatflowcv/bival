package domain

import (
	"errors"
	"fmt"
)

var (
	errEntrySetMissingPlain = errors.New("plain entries require at least one plain entry")
	errEntrySetInstanceOnly = errors.New("instance entries require an olh entry")
	errEntrySetOLHOnly      = errors.New("olh entries require an instance entry")
	errEntrySetTooManyOLH   = errors.New("olh entries must be 0 or 1")
	errEntrySetInvalidTotal = errors.New("total entries must be 1 or 4 plus an even offset")
)

type EntryRegistry struct {
	sets map[string]*entrySet
}

type entrySet struct {
	name string

	plainCount    int
	instanceCount int
	olhCount      int
}

func NewEntryRegistry() *EntryRegistry {
	return &EntryRegistry{
		sets: make(map[string]*entrySet),
	}
}

func (r *EntryRegistry) Add(entry *Entry) {
	if entry == nil {
		return
	}

	set := r.getOrCreateSet(entry.Name())
	set.add(entry)
}

func (r *EntryRegistry) Validate() error {
	var errs []error

	for _, set := range r.sets {
		err := set.validate()
		if err != nil {
			errs = append(errs, err)
		}
	}

	return errors.Join(errs...)
}

func (r *EntryRegistry) getOrCreateSet(name string) *entrySet {
	set := r.sets[name]
	if set != nil {
		return set
	}

	set = &entrySet{
		name:          name,
		plainCount:    0,
		instanceCount: 0,
		olhCount:      0,
	}
	r.sets[name] = set

	return set
}

func (s *entrySet) add(entry *Entry) {
	if entry == nil {
		return
	}

	switch entry.Kind() {
	case KindPlain:
		s.plainCount++
	case KindInstance:
		s.instanceCount++
	case KindOLH:
		s.olhCount++
	}
}

func (s *entrySet) validate() error {
	if s.plainCount == 0 {
		return fmt.Errorf("%w: %q", errEntrySetMissingPlain, s.name)
	}

	if s.instanceCount > 0 && s.olhCount == 0 {
		return fmt.Errorf("%w: %q", errEntrySetInstanceOnly, s.name)
	}

	if s.olhCount > 0 && s.instanceCount == 0 {
		return fmt.Errorf("%w: %q", errEntrySetOLHOnly, s.name)
	}

	if s.olhCount > 1 {
		return fmt.Errorf("%w: %q", errEntrySetTooManyOLH, s.name)
	}

	totalCount := s.plainCount + s.instanceCount + s.olhCount
	if totalCount != 1 && (totalCount < 4 || totalCount%2 != 0) {
		return fmt.Errorf("%w: %q total=%d", errEntrySetInvalidTotal, s.name, totalCount)
	}

	return nil
}
