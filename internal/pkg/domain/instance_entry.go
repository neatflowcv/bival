package domain

type InstanceEntry struct {
	kind        string
	index       []byte
	key         *Key
	versionInfo *DirVersionInfo
	state       *DirState
	meta        *Meta
	pendingMaps []*PendingMap
}

func NewInstanceEntry(p DirEntryParams) *InstanceEntry {
	return &InstanceEntry{
		kind:        p.Kind,
		index:       p.Index,
		key:         p.Key,
		versionInfo: p.VersionInfo,
		state:       p.State,
		meta:        p.Meta,
		pendingMaps: p.PendingMaps,
	}
}

func (e *InstanceEntry) Name() string {
	return e.key.name
}

func (e *InstanceEntry) EntryKey() *Key {
	return e.key
}

func (e *InstanceEntry) VersionPool() int {
	return e.versionInfo.Version().Pool()
}

func (e *InstanceEntry) VersionEpoch() int {
	return e.versionInfo.Version().Epoch()
}

func (e *InstanceEntry) VersionedEpoch() int {
	return e.versionInfo.VersionedEpoch()
}

func (e *InstanceEntry) HasPendingMap() bool {
	return len(e.pendingMaps) > 0
}

func (e *InstanceEntry) Payload() *DirPayload {
	return NewDirPayload(
		e.key,
		e.versionInfo,
		e.state,
		e.meta,
		e.pendingMaps,
	)
}
