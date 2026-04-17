package entrygroup

type versionedEntryKeyRule struct{}

func (versionedEntryKeyRule) Check(group *EntryGroup) error {
	pairedPlainEntries, collectIssue := collectVersionedPlainEntries(group.PlainEntries())
	if collectIssue == nil {
		_, reason := buildPlainEntryMap(pairedPlainEntries)
		if reason != "" {
			return errDuplicateVersionedEntryKey
		}

		_, reason = buildInstanceEntryMap(group.InstanceEntries())
		if reason != "" {
			return errDuplicateVersionedEntryKey
		}
	}

	return nil
}
