package domain

type DirEntryParams struct {
	Kind        string
	Index       []byte
	Name        string
	Instance    string
	Pool        int
	Epoch       int
	VEpoch      int
	Locator     string
	Exists      bool
	Tag         string
	Flags       int
	Meta        *Meta
	PendingMaps []*PendingMap
}
