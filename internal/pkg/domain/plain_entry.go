package domain

import "time"

type PlainEntry struct {
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

func NewPlainEntry(p DirEntryParams) *PlainEntry {
	return &PlainEntry{
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

func (e *PlainEntry) Index() string {
	return string(e.index)
}

func (e *PlainEntry) Name() string {
	return e.name
}

func (e *PlainEntry) Instance() string {
	return e.instance
}

func (e *PlainEntry) VersionPool() int {
	return e.pool
}

func (e *PlainEntry) VersionEpoch() int {
	return e.epoch
}

func (e *PlainEntry) Exists() bool {
	return e.exists
}

func (e *PlainEntry) MTime() time.Time {
	return e.meta.mTime
}

func (e *PlainEntry) ETag() string {
	return e.meta.eTag
}

func (e *PlainEntry) Tag() string {
	return e.tag
}

func (e *PlainEntry) Flags() int {
	return e.flags
}

func (e *PlainEntry) HasPendingMap() bool {
	return len(e.pendingMaps) > 0
}

func (e *PlainEntry) Payload() *DirPayload {
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
