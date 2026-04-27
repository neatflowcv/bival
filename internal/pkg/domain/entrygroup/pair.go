package entrygroup

import (
	"errors"
	"slices"

	"github.com/neatflowcv/bival/internal/pkg/domain"
)

var (
	errMissingMatchingPlain    = errors.New(missingMatchingPlainReason)
	errMissingMatchingInstance = errors.New(missingMatchingInstanceReason)
	errMismatchedVersionPair   = errors.New(mismatchedVersionPairReason)
)

type Pair struct {
	plain    *domain.Plain
	instance *domain.Instance
}

func NewPair(plain *domain.Plain, instance *domain.Instance) *Pair {
	return &Pair{
		plain:    plain,
		instance: instance,
	}
}

func (p *Pair) Plain() *domain.Plain {
	return p.plain
}

func (p *Pair) Instance() *domain.Instance {
	return p.instance
}

func (p *Pair) MTime() string {
	if p.plain != nil {
		return p.plain.MTime()
	}

	if p.instance != nil {
		return p.instance.MTime()
	}

	return ""
}

func NewPairsByGroup(group *EntryGroup) ([]*Pair, error) {
	versionMap := map[string]struct{}{}
	plains := slices.DeleteFunc(group.PlainEntries(), func(entry *domain.Plain) bool {
		return entry.IsPlaceholder()
	})
	plainMap := map[string]*domain.Plain{}

	for _, plain := range plains {
		versionMap[plain.Instance()] = struct{}{}
		plainMap[plain.Instance()] = plain
	}

	instances := group.InstanceEntries()
	instanceMap := map[string]*domain.Instance{}

	for _, instance := range instances {
		versionMap[instance.Instance()] = struct{}{}
		instanceMap[instance.Instance()] = instance
	}

	var pairs []*Pair

	for version := range versionMap {
		pair := NewPair(plainMap[version], instanceMap[version])
		pairs = append(pairs, pair)
	}

	return pairs, nil
}

func composeVersionedPairs(
	plainByKey map[string]*domain.Plain,
	instanceByKey map[string]*domain.Instance,
) (map[string]*Pair, error) {
	pairCount := max(len(instanceByKey), len(plainByKey))

	pairs := make(map[string]*Pair, pairCount)
	errs := make([]error, 0)

	for key, plainEntry := range plainByKey {
		pair := pairByInstance(pairs, key)
		pair.plain = plainEntry
	}

	for key, instanceEntry := range instanceByKey {
		pair := pairByInstance(pairs, key)
		pair.instance = instanceEntry
	}

	errs = append(errs, validateVersionedPairs(pairs)...)

	if len(errs) > 0 {
		return nil, errors.Join(errs...)
	}

	for key, pair := range pairs {
		pairs[key] = NewPair(pair.plain, pair.instance)
	}

	return pairs, nil
}

func validateVersionedPairs(pairs map[string]*Pair) []error {
	errs := make([]error, 0)

	for _, pair := range pairs {
		switch {
		case pair.plain == nil:
			errs = append(errs, errMissingMatchingPlain)
		case pair.instance == nil:
			errs = append(errs, errMissingMatchingInstance)
		case !domain.IsVersionPair(pair.plain, pair.instance):
			errs = append(errs, errMismatchedVersionPair)
		}
	}

	return errs
}

func pairByInstance(pairs map[string]*Pair, instance string) *Pair {
	pair := pairs[instance]
	if pair != nil {
		return pair
	}

	pair = &Pair{
		plain:    nil,
		instance: nil,
	}
	pairs[instance] = pair

	return pair
}
