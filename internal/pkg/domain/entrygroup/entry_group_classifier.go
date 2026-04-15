package entrygroup

type Specification interface {
	IsSatisfiedBy(group *EntryGroup) bool
}

type Classifier struct {
	unversionedSpec Specification
	versionedSpec   Specification
}

func NewClassifier() Classifier {
	return Classifier{
		unversionedSpec: unversionedObjectSpecification{},
		versionedSpec:   versionedObjectSpecification{},
	}
}

func (c Classifier) Classify(group *EntryGroup) ObjectKind {
	if c.unversionedSpec.IsSatisfiedBy(group) {
		return UnversionedObject
	}

	if c.versionedSpec.IsSatisfiedBy(group) {
		return VersionedObject
	}

	return UnknownObject
}
