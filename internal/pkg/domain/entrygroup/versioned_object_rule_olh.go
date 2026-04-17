package entrygroup

type versionedOLHRule struct{}

func (versionedOLHRule) Check(group *EntryGroup) error {
	_, err := buildVersionedOLH(group.OLHEntries(), group.InstanceEntries())

	return err
}
