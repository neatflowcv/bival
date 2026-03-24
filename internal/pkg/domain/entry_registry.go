package domain

import (
	"errors"
	"fmt"
	"time"
)

var (
	errEntrySetNonVersionedPlain = errors.New("non-versioned set must contain exactly 1 plain entry")
	errEntrySetVersionedOLH      = errors.New("versioned set must contain exactly 1 olh entry")
	errEntrySetVersionedPlain    = errors.New(
		"versioned set must contain exactly 1 head plain entry plus 1 plain entry per instance entry",
	)
	errEntrySetOLHInstance = errors.New("versioned set olh must reference an existing instance entry")
	errEntrySetOLHLatest   = errors.New("versioned set olh must reference the latest instance entry")
)

type EntryRegistry struct {
	sets map[string]*entrySet
}

type entrySet struct {
	name string

	plainCount    int
	instanceCount int
	olhCount      int
	olhInstance   string
	instanceMTime map[string]time.Time
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
		olhInstance:   "",
		instanceMTime: make(map[string]time.Time),
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
		s.instanceMTime[entry.Instance()] = entry.MTime()
	case KindOLH:
		s.olhCount++
		s.olhInstance = entry.Instance()
	}
}

func (s *entrySet) validate() error {
	if !s.isVersioningObject() {
		if !s.isValidNonVersionedObject() {
			return fmt.Errorf(
				"%w: %q plain=%d instance=%d olh=%d",
				errEntrySetNonVersionedPlain,
				s.name,
				s.plainCount,
				s.instanceCount,
				s.olhCount,
			)
		}

		return nil
	}

	if !s.isValidVersionedObject() {
		if s.olhCount != 1 {
			return fmt.Errorf(
				"%w: %q plain=%d instance=%d olh=%d",
				errEntrySetVersionedOLH,
				s.name,
				s.plainCount,
				s.instanceCount,
				s.olhCount,
			)
		}

		return fmt.Errorf(
			"%w: %q plain=%d instance=%d olh=%d",
			errEntrySetVersionedPlain,
			s.name,
			s.plainCount,
			s.instanceCount,
			s.olhCount,
		)
	}

	if !s.hasValidOLHInstance() {
		return fmt.Errorf(
			"%w: %q olh_instance=%q",
			errEntrySetOLHInstance,
			s.name,
			s.olhInstance,
		)
	}

	if !s.isLatestOLHInstance() {
		return fmt.Errorf(
			"%w: %q olh_instance=%q latest_instance=%q",
			errEntrySetOLHLatest,
			s.name,
			s.olhInstance,
			s.latestInstance(),
		)
	}

	return nil
}

func (s *entrySet) isVersioningObject() bool {
	return s.instanceCount > 0 || s.olhCount > 0
}

func (s *entrySet) isValidNonVersionedObject() bool {
	return s.plainCount == 1 && s.instanceCount == 0 && s.olhCount == 0
}

func (s *entrySet) isValidVersionedObject() bool {
	// Versioned objects keep one extra plain entry for the head object.
	return s.olhCount == 1 && s.plainCount == s.instanceCount+1
}

func (s *entrySet) hasValidOLHInstance() bool {
	_, ok := s.instanceMTime[s.olhInstance]

	return ok
}

func (s *entrySet) isLatestOLHInstance() bool {
	olhMTime, ok := s.instanceMTime[s.olhInstance]
	if !ok {
		return false
	}

	for _, instanceMTime := range s.instanceMTime {
		if instanceMTime.After(olhMTime) {
			return false
		}
	}

	return true
}

func (s *entrySet) latestInstance() string {
	var (
		latestInstance string
		latestMTime    time.Time
	)

	for instance, instanceMTime := range s.instanceMTime {
		if latestInstance == "" || instanceMTime.After(latestMTime) {
			latestInstance = instance
			latestMTime = instanceMTime
		}
	}

	return latestInstance
}
