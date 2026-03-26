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
	return group.PlainCount() == 1 &&
		group.InstanceCount() == 0 &&
		group.OLHCount() == 0
}

type versionedObjectSpecification struct{}

func (versionedObjectSpecification) IsSatisfiedBy(group *EntryGroup) bool {
	return group.OLHCount() == 1 &&
		group.InstanceCount()+1 == group.PlainCount()
}
