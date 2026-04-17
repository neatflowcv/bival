package domain

type DirEntryParams struct {
	Kind        string
	Index       []byte
	Name        string
	Instance    string
	VersionInfo *DirVersionInfo
	Locator     string
	Exists      bool
	Tag         string
	Flags       int
	Meta        *Meta
	PendingMaps []*PendingMap
}
