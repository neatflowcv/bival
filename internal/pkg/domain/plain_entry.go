package domain

import "time"

type PlainEntry struct {
	kind        string
	index       []byte
	key         *Key
	versionInfo *DirVersionInfo
	state       *DirState
	meta        *Meta
	pendingMaps []*PendingMap
}

func NewPlainEntry(p DirEntryParams) *PlainEntry {
	return &PlainEntry{
		kind:        p.Kind,
		index:       p.Index,
		key:         p.Key,
		versionInfo: p.VersionInfo,
		state:       p.State,
		meta:        p.Meta,
		pendingMaps: p.PendingMaps,
	}
}

func (e *PlainEntry) Index() string {
	return string(e.index)
}

func (e *PlainEntry) Name() string {
	return e.key.name
}

func (e *PlainEntry) Instance() string {
	return e.key.instance
}

func (e *PlainEntry) VersionPool() int {
	return e.versionInfo.Version().Pool()
}

func (e *PlainEntry) VersionEpoch() int {
	return e.versionInfo.Version().Epoch()
}

func (e *PlainEntry) Exists() bool {
	return e.state.exists
}

func (e *PlainEntry) MTime() time.Time {
	return e.meta.auditInfo.mTime
}

func (e *PlainEntry) ETag() string {
	return e.meta.auditInfo.eTag
}

func (e *PlainEntry) Tag() string {
	return e.state.tag
}

func (e *PlainEntry) Flags() int {
	return e.state.flags
}

func (e *PlainEntry) HasPendingMap() bool {
	return len(e.pendingMaps) > 0
}

func (e *PlainEntry) Payload() *DirPayload {
	return NewDirPayload(
		e.key,
		e.versionInfo,
		e.state,
		e.meta,
		e.pendingMaps,
	)
}
