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

type DirEntry struct {
	kind        string
	index       []byte
	key         *Key
	versionInfo *DirVersionInfo
	state       *DirState
	meta        *Meta
	pendingMaps []*PendingMap
}

func NewDirEntry(p DirEntryParams) *DirEntry {
	return &DirEntry{
		kind:        p.Kind,
		index:       p.Index,
		key:         p.Key,
		versionInfo: p.VersionInfo,
		state:       p.State,
		meta:        p.Meta,
		pendingMaps: p.PendingMaps,
	}
}

func (e *DirEntry) hasPendingMap() bool {
	return len(e.pendingMaps) > 0
}

func (e *DirEntry) indexString() string {
	return string(e.index)
}
