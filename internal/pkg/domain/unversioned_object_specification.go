package domain

type unversionedObjectSpecification struct{}

func (unversionedObjectSpecification) IsSatisfiedBy(group *EntryGroup) bool {
	if !hasUnversionedEntryCounts(group) {
		return false
	}

	plainEntries := group.PlainEntries()
	if len(plainEntries) != 1 {
		return false
	}

	return isValidUnversionedPlainEntry(plainEntries[0])
}

func hasUnversionedEntryCounts(group *EntryGroup) bool {
	return group.PlainCount() == 1 &&
		group.InstanceCount() == 0 &&
		group.OLHCount() == 0
}

func isValidUnversionedPlainEntry(entry *PlainEntry) bool {
	return hasValidUnversionedIdentity(entry) &&
		hasValidUnversionedState(entry)
}

func hasValidUnversionedIdentity(entry *PlainEntry) bool {
	return entry.Index() == entry.Name() &&
		entry.Instance() == "" &&
		entry.VersionPool() >= 1 &&
		entry.VersionEpoch() >= 1
}

func hasValidUnversionedState(entry *PlainEntry) bool {
	return entry.Exists() &&
		!entry.MTime().IsZero() &&
		entry.ETag() != "" &&
		entry.Tag() != "" &&
		entry.Flags() == 0
}
