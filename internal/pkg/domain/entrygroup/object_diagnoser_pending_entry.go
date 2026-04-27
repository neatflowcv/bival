package entrygroup

type pendingEntryDiagnoser struct{}

func (pendingEntryDiagnoser) Diagnose(group *EntryGroup) []*Issue {
	issues := make([]*Issue, 0)

	for _, entry := range group.plainEntries {
		if !entry.HasPendingMap() {
			continue
		}

		issues = append(issues, newIssue(
			issueCodePendingPlainExists,
			map[string]string{
				"instance": entry.Instance(),
				"index":    entry.Index(),
			},
		))
	}

	for _, entry := range group.instanceEntries {
		if !entry.HasPendingMap() {
			continue
		}

		issues = append(issues, newIssue(
			issueCodePendingInstanceExists,
			map[string]string{
				"instance": entry.Instance(),
				"index":    entry.Index(),
			},
		))
	}

	for _, entry := range group.olhEntries {
		if !entry.HasPendingLog() {
			continue
		}

		issues = append(issues, newIssue(
			issueCodePendingOLHExists,
			nil,
		))
	}

	if len(issues) == 0 {
		return nil
	}

	return issues
}
