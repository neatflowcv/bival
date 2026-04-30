package entrygroup

func diagnoseObject(group *EntryGroup) []*Issue {
	if group.isUnversionedObject() {
		return diagnoseUnversionedObject(group)
	}

	return diagnoseVersionedObject(group)
}

func diagnoseUnversionedObject(group *EntryGroup) []*Issue {
	return diagnose(group, newUnversionedObjectDiagnosers())
}

func diagnoseVersionedObject(group *EntryGroup) []*Issue {
	return diagnose(group, newVersionedObjectDiagnosers())
}
