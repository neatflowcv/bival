package domain

type InstanceEntry struct {
	kind        string
	index       []byte
	name        string
	instance    string
	pool        int
	epoch       int
	vEpoch      int
	locator     string
	exists      bool
	tag         string
	flags       int
	meta        *Meta
	pendingMaps []*PendingMap
}

func NewInstanceEntry(p DirEntryParams) *InstanceEntry {
	return &InstanceEntry{
		kind:        p.Kind,
		index:       p.Index,
		name:        p.Name,
		instance:    p.Instance,
		pool:        p.Pool,
		epoch:       p.Epoch,
		vEpoch:      p.VEpoch,
		locator:     p.Locator,
		exists:      p.Exists,
		tag:         p.Tag,
		flags:       p.Flags,
		meta:        p.Meta,
		pendingMaps: p.PendingMaps,
	}
}

func (e *InstanceEntry) Name() string {
	return e.name
}

func (e *InstanceEntry) Instance() string {
	return e.instance
}

func (e *InstanceEntry) VersionPool() int {
	return e.pool
}

func (e *InstanceEntry) VersionEpoch() int {
	return e.epoch
}

func (e *InstanceEntry) VersionedEpoch() int {
	return e.vEpoch
}

func (e *InstanceEntry) HasPendingMap() bool {
	return len(e.pendingMaps) > 0
}

func (e *InstanceEntry) Payload() *DirPayload {
	return NewDirPayload(
		e.name,
		e.instance,
		e.pool,
		e.epoch,
		e.vEpoch,
		e.locator,
		e.exists,
		e.tag,
		e.flags,
		e.meta,
		e.pendingMaps,
	)
}
