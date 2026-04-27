package entrygroup

import (
	"strconv"
	"time"
)

const (
	issueCodePendingEntryExists      = "entry.pending.exists"
	issueCodeTooManyVersionedEntries = "entry.versioned.count.exceeded"
	issueCodeMissingVersionedHead    = "versioned.head.missing"
	issueCodeDuplicateVersionedHead  = "versioned.head.duplicate"
	issueCodeInvalidVersionedHead    = "versioned.head.invalid"
	issueCodeDuplicateEntryKey       = "versioned.entry_key.duplicate"
	issueCodeInvalidEntryKey         = "versioned.entry_key.invalid"
	issueCodeMissingMatchingPlain    = "pair.plain.missing"
	issueCodeMissingMatchingInstance = "pair.instance.missing"
	issueCodeMismatchedVersionPair   = "pair.version.mismatched"
	issueCodeMissingOLH              = "olh.missing"
	issueCodeInvalidOLH              = "olh.invalid"
	issueCodeInvalidOLHReference     = "olh.reference.invalid"
	issueCodeStaleOLHReference       = "olh.reference.stale"
	issueCodeStaleDeleteMarkerOLH    = "olh.delete_marker.stale"
)

type Diagnoser interface {
	Diagnose(group *EntryGroup) []*Issue
}

func newObjectDiagnosers() []Diagnoser {
	return []Diagnoser{
		headDiagnoser{},
		entryKeyDiagnoser{},
		pairDiagnoser{},
		olhDiagnoser{},
		staleOLHDiagnoser{now: time.Now()},
	}
}

func diagnose(group *EntryGroup, diagnosers []Diagnoser) []*Issue {
	issues := make([]*Issue, 0)

	for _, diagnoser := range diagnosers {
		issues = append(issues, diagnoser.Diagnose(group)...)
	}

	if len(issues) == 0 {
		return nil
	}

	return issues
}

func (g *EntryGroup) Issues() []*Issue {
	issues := make([]*Issue, 0)

	if g.HasPendingEntries() {
		issues = append(issues, newIssue(
			issueCodePendingEntryExists,
			nil,
		))
	}

	if g.isUnversionedObject() {
		if len(issues) == 0 {
			return nil
		}

		return issues
	}

	if count := g.versionedEntryCount(); count > maxVersionedEntryCount {
		issues = append(issues, newIssue(
			issueCodeTooManyVersionedEntries,
			map[string]string{
				"count":   strconv.Itoa(count),
				"maximum": strconv.Itoa(maxVersionedEntryCount),
			},
		))
	}

	issues = append(issues, diagnose(g, newObjectDiagnosers())...)
	if len(issues) == 0 {
		return nil
	}

	return issues
}
