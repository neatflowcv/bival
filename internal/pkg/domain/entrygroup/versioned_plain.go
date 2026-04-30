package entrygroup

import (
	"errors"

	"github.com/neatflowcv/bival/internal/pkg/domain"
)

const (
	missingVersionedHeadReason   = "missing versioned head"
	duplicateVersionedHeadReason = "duplicate versioned head"
	invalidVersionedHeadReason   = "invalid versioned head"
)

var (
	errDuplicateVersionedHead = errors.New(duplicateVersionedHeadReason)
	errInvalidVersionedHead   = errors.New(invalidVersionedHeadReason)
	errMissingVersionedHead   = errors.New(missingVersionedHeadReason)
)

func isVersionedHeadCandidate(entry *domain.Plain) bool {
	return entry.Index() == entry.Name() &&
		entry.Instance() == ""
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
			errs = append(errs, errDuplicateVersionedHead)

			continue
		}

		placeholder = entry
	}

	if placeholder == nil {
		if invalidHeadCount > 0 {
			errs = append(errs, errInvalidVersionedHead)
		} else {
			errs = append(errs, errMissingVersionedHead)
		}
	}

	if len(errs) > 0 {
		return nil, errors.Join(errs...)
	}

	return pairedPlainEntries, nil
}
