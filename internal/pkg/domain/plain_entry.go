package domain

import "time"

type PlainEntry struct {
	kind        string
	index       []byte
	name        string
	instance    string
	versionInfo *DirVersionInfo
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
		versionInfo: p.VersionInfo,
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
	return e.versionInfo.Pool()
}

func (e *PlainEntry) VersionEpoch() int {
	return e.versionInfo.Epoch()
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
		e.versionInfo,
		e.locator,
		e.exists,
		e.tag,
		e.flags,
		e.meta,
		e.pendingMaps,
	)
}
