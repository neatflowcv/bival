package entrygroup

type versionedHeadRule struct{}

func (versionedHeadRule) Check(group *EntryGroup) error {
	headCount := 0
	invalidHeadCount := 0

	for _, entry := range group.PlainEntries() {
		if entry.IsPlaceholder() {
			headCount++

			continue
		}

		if isVersionedHeadCandidate(entry) {
			invalidHeadCount++
		}
	}

	switch {
	case headCount == 0 && invalidHeadCount > 0:
		return errInvalidVersionedHead
	case headCount == 0:
		return errMissingVersionedHead
	case headCount > 1:
		return errDuplicateVersionedHead
	default:
		return nil
	}
}
