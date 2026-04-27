package entrygroup

import "strconv"

func diagnoseObject(group *EntryGroup) []*Issue {
	if group.isUnversionedObject() {
		return diagnoseUnversionedObject(group)
	}

	return diagnoseVersionedObject(group)
}

func diagnoseUnversionedObject(group *EntryGroup) []*Issue {
	return diagnose(group, newUnversionedObjectDiagnosers())
}

func diagnoseVersionedObject(group *EntryGroup) []*Issue {
	issues := make([]*Issue, 0)

	issues = append(issues, diagnoseUnversionedObject(group)...)

	if count := group.versionedEntryCount(); count > maxVersionedEntryCount {
		issues = append(issues, newIssue(
			issueCodeTooManyVersionedEntries,
			map[string]string{
				"count":   strconv.Itoa(count),
				"maximum": strconv.Itoa(maxVersionedEntryCount),
			},
		))
	}

	diagnosers := newVersionedObjectDiagnosers()
	issues = append(issues, diagnose(group, diagnosers[1:])...)

	if len(issues) == 0 {
		return nil
	}

	return issues
}
