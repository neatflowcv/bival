package domain

type versionedObjectSpecification struct{}

func (versionedObjectSpecification) IsSatisfiedBy(group *EntryGroup) bool {
	return group.OLHCount() == 1 &&
		group.InstanceCount()+1 == group.PlainCount()
}
