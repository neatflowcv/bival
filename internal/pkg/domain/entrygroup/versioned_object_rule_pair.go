package entrygroup

type versionedPairRule struct{}

func (versionedPairRule) Check(group *EntryGroup) error {
	pairedPlainEntries, collectIssue := collectVersionedPlainEntries(group.PlainEntries())
	if collectIssue == nil {
		plainByKey, reason := buildPlainEntryMap(pairedPlainEntries)
		if reason != "" {
			return nil
		}

		instanceByKey, reason := buildInstanceEntryMap(group.InstanceEntries())
		if reason != "" {
			return nil
		}

		_, pairIssue := composeVersionedPairs(plainByKey, instanceByKey)
		if pairIssue != nil {
			return pairIssue
		}
	}

	return nil
}
