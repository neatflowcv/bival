package entrygroup

import (
	"strconv"
	"time"

	"github.com/neatflowcv/bival/internal/pkg/domain"
)

const staleOLHThreshold = 7 * 24 * time.Hour

type staleOLHDiagnoser struct {
	now time.Time
}

func (d staleOLHDiagnoser) Diagnose(group *EntryGroup) []*Issue {
	olh, pairs, instanceByKey := olhContext(group)
	if olh == nil || pairs == nil || instanceByKey == nil {
		return nil
	}

	referencedInstance, ok := instanceByKey[olh.Instance()]
	if !ok || !isStaleInstance(referencedInstance, d.now) {
		return nil
	}

	referencedVersion := olh.Instance()

	if olh.DeleteMarker() {
		if len(pairs) == 0 {
			return nil
		}

		return []*Issue{
			newIssue(
				issueCodeStaleDeleteMarkerOLH,
				map[string]string{
					"referenced_version": referencedVersion,
					"pair_count":         strconv.Itoa(len(pairs)),
				},
			),
		}
	}

	if len(pairs) <= 1 {
		return nil
	}

	return []*Issue{
		newIssue(
			issueCodeStaleOLHReference,
			map[string]string{
				"referenced_version": referencedVersion,
				"pair_count":         strconv.Itoa(len(pairs)),
			},
		),
	}
}

func olhContext(group *EntryGroup) (*domain.OLH, map[string]*Pair, map[string]*domain.Instance) {
	olh, _ := buildVersionedOLH(group.OLHEntries(), group.InstanceEntries())
	if olh == nil {
		return nil, nil, nil
	}

	pairedPlainEntries, _ := collectVersionedPlainEntries(group.PlainEntries())
	if pairedPlainEntries == nil {
		return nil, nil, nil
	}

	plainByKey, reason := buildPlainEntryMap(pairedPlainEntries)
	if reason != "" {
		return nil, nil, nil
	}

	instanceByKey, reason := buildInstanceEntryMap(group.InstanceEntries())
	if reason != "" {
		return nil, nil, nil
	}

	pairs, _ := composeVersionedPairs(plainByKey, instanceByKey)
	if pairs == nil {
		return nil, nil, nil
	}

	return olh, pairs, instanceByKey
}

func isStaleInstance(entry *domain.Instance, now time.Time) bool {
	mtime, err := time.Parse(time.RFC3339Nano, entry.MTime())
	if err != nil {
		return false
	}

	return now.Sub(mtime) >= staleOLHThreshold
}
