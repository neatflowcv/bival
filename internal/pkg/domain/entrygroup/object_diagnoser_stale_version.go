package entrygroup

import (
	"errors"
	"slices"
	"time"

	"github.com/neatflowcv/bival/internal/pkg/domain"
)

const (
	invalidOLHReason  = "invalid olh"
	staleOLHThreshold = 7 * 24 * time.Hour
)

var (
	errInvalidOLH = errors.New(invalidOLHReason)
)

type staleVersionDiagnoser struct {
	now time.Time
}

func (d staleVersionDiagnoser) Diagnose(group *EntryGroup) []*Issue {
	olh, err := extractOLH(group)
	if err != nil {
		return nil
	}

	pairs := NewPairsByGroup(group)

	olhPair, ok := pairs.PairByVersion(olh.Instance())
	if !ok {
		return nil
	}

	if !isStale(olhPair, d.now) {
		return nil
	}

	items := pairs.Items()
	if !olhPair.IsSoftDeleted() {
		items = slices.DeleteFunc(pairs.Items(), func(pair *Pair) bool {
			return pair.Version() == olh.Instance()
		})
	}

	var issues []*Issue

	for _, pair := range items {
		issues = append(issues, newIssue(issueCodeStaleVersion, map[string]string{
			"version": pair.Version(),
		}))
	}

	return issues
}

func extractOLH(group *EntryGroup) (*domain.OLH, error) {
	olhEntries := group.OLHEntries()
	if len(olhEntries) != 1 {
		return nil, errInvalidOLH
	}

	return olhEntries[0], nil
}

func isStale(pair *Pair, now time.Time) bool {
	mtime, err := time.Parse(time.RFC3339Nano, pair.MTime())
	if err != nil {
		return false
	}

	return now.Sub(mtime) >= staleOLHThreshold
}
