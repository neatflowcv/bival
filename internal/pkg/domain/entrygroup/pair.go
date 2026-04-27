package entrygroup

import (
	"errors"

	"github.com/neatflowcv/bival/internal/pkg/domain"
)

var errPairAlreadyFull = errors.New("pair is already full")

var (
	errMissingMatchingPlain    = errors.New(missingMatchingPlainReason)
	errMissingMatchingInstance = errors.New(missingMatchingInstanceReason)
	errMismatchedVersionPair   = errors.New(mismatchedVersionPairReason)
)

type Pair struct {
	Plain    *domain.Plain
	Instance *domain.Instance
}

func (p *Pair) isFull() bool {
	return p.Plain != nil &&
		p.Instance != nil
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

		err := pair.setPlain(plainEntry)
		if err != nil {
			errs = append(errs, err)
		}
	}

	for key, instanceEntry := range instanceByKey {
		pair := pairByInstance(pairs, key)

		err := pair.setInstance(instanceEntry)
		if err != nil {
			errs = append(errs, err)
		}
	}

	errs = append(errs, validateVersionedPairs(pairs)...)

	if len(errs) > 0 {
		return nil, errors.Join(errs...)
	}

	return pairs, nil
}

func validateVersionedPairs(pairs map[string]*Pair) []error {
	errs := make([]error, 0)

	for _, pair := range pairs {
		switch {
		case pair.Plain == nil:
			errs = append(errs, errMissingMatchingPlain)
		case pair.Instance == nil:
			errs = append(errs, errMissingMatchingInstance)
		case !domain.IsVersionPair(pair.Plain, pair.Instance):
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
		Plain:    nil,
		Instance: nil,
	}
	pairs[instance] = pair

	return pair
}

func (p *Pair) setPlain(entry *domain.Plain) error {
	if p.isFull() {
		return errPairAlreadyFull
	}

	p.Plain = entry

	return nil
}

func (p *Pair) setInstance(entry *domain.Instance) error {
	if p.isFull() {
		return errPairAlreadyFull
	}

	p.Instance = entry

	return nil
}
