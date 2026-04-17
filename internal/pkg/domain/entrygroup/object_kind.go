package entrygroup

type ObjectKind string

const (
	UnversionedObject ObjectKind = "unversioned_object"
	VersionedObject   ObjectKind = "versioned_object"
)
