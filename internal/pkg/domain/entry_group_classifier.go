package domain

type EntryGroupSpecification interface {
	IsSatisfiedBy(group *EntryGroup) bool
}

type EntryGroupClassifier struct {
	unversionedSpec EntryGroupSpecification
	versionedSpec   EntryGroupSpecification
}

func NewEntryGroupClassifier() EntryGroupClassifier {
	return EntryGroupClassifier{
		unversionedSpec: unversionedObjectSpecification{},
		versionedSpec:   versionedObjectSpecification{},
	}
}

func (c EntryGroupClassifier) Classify(group *EntryGroup) ObjectKind {
	if c.unversionedSpec.IsSatisfiedBy(group) {
		return UnversionedObject
	}

	if c.versionedSpec.IsSatisfiedBy(group) {
		return VersionedObject
	}

	return UnknownObject
}

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

type versionedObjectSpecification struct{}

func (versionedObjectSpecification) IsSatisfiedBy(group *EntryGroup) bool {
	return group.OLHCount() == 1 &&
		group.InstanceCount()+1 == group.PlainCount()
}
