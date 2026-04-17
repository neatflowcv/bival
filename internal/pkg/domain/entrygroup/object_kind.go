package entrygroup

type ObjectKind string

const (
	UnversionedObjectKind ObjectKind = "unversioned_object"
	VersionedObjectKind   ObjectKind = "versioned_object"
)
