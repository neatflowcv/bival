package entrygroup

import "time"

const (
	issueCodePendingPlainExists      = "plain.pending.exists"
	issueCodePendingInstanceExists   = "instance.pending.exists"
	issueCodePendingOLHExists        = "olh.pending.exists"
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
	issueCodeOutdatedOLHReference    = "olh.reference.outdated"
	issueCodeStaleVersion            = "version.stale"
)

type Diagnoser interface {
	Diagnose(group *EntryGroup) []*Issue
}

func newUnversionedObjectDiagnosers() []Diagnoser {
	return []Diagnoser{
		pendingEntryDiagnoser{},
	}
}

func newVersionedObjectDiagnosers() []Diagnoser {
	return []Diagnoser{
		pendingEntryDiagnoser{},
		headDiagnoser{},
		entryKeyDiagnoser{},
		pairDiagnoser{},
		olhDiagnoser{},
		olhLatestMTimeDiagnoser{},
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
