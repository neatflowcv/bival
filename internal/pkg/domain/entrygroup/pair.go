package entrygroup

import (
	"errors"

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

func (p *Pair) Version() (string, bool) {
	if p.plain != nil {
		return p.plain.Instance(), true
	}

	if p.instance != nil {
		return p.instance.Instance(), true
	}

	return "", false
}

func (p *Pair) MTime() (string, bool) {
	if p.plain != nil {
		return p.plain.MTime(), true
	}

	if p.instance != nil {
		return p.instance.MTime(), true
	}

	return "", false
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
