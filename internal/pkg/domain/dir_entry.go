package domain

type DirEntryParams struct {
	Kind        string
	Index       []byte
	Key         *Key
	VersionInfo *DirVersionInfo
	State       *DirState
	Meta        *Meta
	PendingMaps []*PendingMap
}
