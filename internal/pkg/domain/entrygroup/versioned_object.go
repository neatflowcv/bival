package entrygroup

import (
	"errors"

	"github.com/neatflowcv/bival/internal/pkg/domain"
)

var errPairAlreadyFull = errors.New("pair is already full")

type Pair struct {
	Plain    *domain.Plain
	Instance *domain.Instance
}

func (p *Pair) isFull() bool {
	return p.Plain != nil &&
		p.Instance != nil
}

func collectVersionedPlainEntries(entries []*domain.Plain) ([]*domain.Plain, error) {
	pairedPlainEntries := make([]*domain.Plain, 0, len(entries))

	var (
		placeholder      *domain.Plain
		invalidHeadCount int
		errs             []error
	)

	for _, entry := range entries {
		if isVersionedHeadCandidate(entry) && !entry.IsPlaceholder() {
			invalidHeadCount++

			continue
		}

		if !entry.IsPlaceholder() {
			pairedPlainEntries = append(pairedPlainEntries, entry)

			continue
		}

		if placeholder != nil {
			errs = append(errs, errors.New(duplicateVersionedHeadReason))

			continue
		}

		placeholder = entry
	}

	if placeholder == nil {
		if invalidHeadCount > 0 {
			errs = append(errs, errors.New(invalidVersionedHeadReason))
		} else {
			errs = append(errs, errors.New(missingVersionedHeadReason))
		}
	}

	if len(errs) > 0 {
		return nil, errors.Join(errs...)
	}

	return pairedPlainEntries, nil
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
			errs = append(errs, errors.New(missingMatchingPlainReason))
		case pair.Instance == nil:
			errs = append(errs, errors.New(missingMatchingInstanceReason))
		case !domain.IsVersionPair(pair.Plain, pair.Instance):
			errs = append(errs, errors.New(mismatchedVersionPairReason))
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

func buildVersionedOLH(olhEntries []*domain.OLH, instanceEntries []*domain.Instance) (*domain.OLH, error) {
	olh, reason := singleValidOLHEntry(olhEntries)
	switch reason {
	case "":
	case missingOLHReason:
		return nil, errors.New(missingOLHReason)
	default:
		return nil, errors.New(invalidOLHReason)
	}

	if !hasValidOLHReference(olh, instanceEntries) {
		return nil, errors.New(invalidOLHReferenceReason)
	}

	return olh, nil
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
