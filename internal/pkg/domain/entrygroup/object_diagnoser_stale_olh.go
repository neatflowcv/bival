package entrygroup

import (
	"errors"
	"strconv"
	"time"

	"github.com/neatflowcv/bival/internal/pkg/domain"
)

const staleOLHThreshold = 7 * 24 * time.Hour

var (
	errMissingOLH          = errors.New(missingOLHReason)
	errInvalidOLH          = errors.New(invalidOLHReason)
	errInvalidOLHReference = errors.New(invalidOLHReferenceReason)
)

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

func buildVersionedOLH(olhEntries []*domain.OLH, instanceEntries []*domain.Instance) (*domain.OLH, error) {
	olh, reason := singleValidOLHEntry(olhEntries)
	switch reason {
	case "":
	case missingOLHReason:
		return nil, errMissingOLH
	default:
		return nil, errInvalidOLH
	}

	if !hasValidOLHReference(olh, instanceEntries) {
		return nil, errInvalidOLHReference
	}

	return olh, nil
}

func hasValidOLHReference(olhEntry *domain.OLH, instanceEntries []*domain.Instance) bool {
	instanceSet, instanceSetOK := instanceNameSet(instanceEntries)
	if !instanceSetOK {
		return false
	}

	referencedInstance := olhEntry.Instance()
	_, exists := instanceSet[referencedInstance]

	return exists
}

func singleValidOLHEntry(entries []*domain.OLH) (*domain.OLH, string) {
	if len(entries) == 0 {
		return nil, missingOLHReason
	}

	if len(entries) != 1 {
		return nil, invalidOLHReason
	}

	entry := entries[0]
	if entry == nil || entry.Name() == "" {
		return nil, invalidOLHReason
	}

	if entry.HasPendingLog() {
		return nil, invalidOLHReason
	}

	return entry, ""
}

func isStaleInstance(entry *domain.Instance, now time.Time) bool {
	mtime, err := time.Parse(time.RFC3339Nano, entry.MTime())
	if err != nil {
		return false
	}

	return now.Sub(mtime) >= staleOLHThreshold
}
