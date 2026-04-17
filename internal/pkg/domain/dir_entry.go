package domain

type DirEntryParams struct {
	Kind        string
	Index       []byte
	Name        string
	Instance    string
	VersionInfo *DirVersionInfo
	State       *DirState
	Meta        *Meta
	PendingMaps []*PendingMap
}
