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
	// Both values being nil indicates a programmer error, so panic intentionally.
	if plain == nil && instance == nil {
		panic("entrygroup.NewPair: plain and instance cannot both be nil")
	}

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

func (p *Pair) Version() string {
	if p.plain != nil {
		return p.plain.Instance()
	}

	return p.instance.Instance()
}

func (p *Pair) MTime() string {
	if p.plain != nil {
		return p.plain.MTime()
	}

	return p.instance.MTime()
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
