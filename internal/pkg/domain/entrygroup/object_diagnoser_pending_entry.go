package entrygroup

type pendingEntryDiagnoser struct{}

func (pendingEntryDiagnoser) Diagnose(group *EntryGroup) []*Issue {
	if !group.HasPendingEntries() {
		return nil
	}

	return []*Issue{
		newIssue(
			issueCodePendingEntryExists,
			nil,
		),
	}
}
