package entrygroup

import "strconv"

const maxVersionedEntryCount = 8

type versionedEntryCountDiagnoser struct{}

func (versionedEntryCountDiagnoser) Diagnose(group *EntryGroup) []*Issue {
	count := group.versionedEntryCount()
	if count <= maxVersionedEntryCount {
		return nil
	}

	return []*Issue{
		newIssue(
			issueCodeTooManyVersionedEntries,
			map[string]string{
				"count":   strconv.Itoa(count),
				"maximum": strconv.Itoa(maxVersionedEntryCount),
			},
		),
	}
}
