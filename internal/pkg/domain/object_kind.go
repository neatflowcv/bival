package domain

type ObjectKind string

const (
	UnversionedObject ObjectKind = "unversioned_object"
	VersionedObject   ObjectKind = "versioned_object"
	UnknownObject     ObjectKind = "unknown_object"
)
