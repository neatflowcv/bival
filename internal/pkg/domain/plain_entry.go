package domain

import "time"

type PlainEntry struct {
	entry *DirEntry
}

func NewPlainEntry(entry *DirEntry) *PlainEntry {
	return &PlainEntry{
		entry: entry,
	}
}

func (e *PlainEntry) Index() string {
	return e.entry.indexString()
}

func (e *PlainEntry) Name() string {
	return e.entry.payload.key.name
}

func (e *PlainEntry) Instance() string {
	return e.entry.payload.key.instance
}

func (e *PlainEntry) VersionPool() int {
	return e.entry.payload.versionInfo.version.pool
}

func (e *PlainEntry) VersionEpoch() int {
	return e.entry.payload.versionInfo.version.epoch
}

func (e *PlainEntry) Exists() bool {
	return e.entry.payload.state.exists
}

func (e *PlainEntry) MTime() time.Time {
	return e.entry.payload.meta.auditInfo.mTime
}

func (e *PlainEntry) ETag() string {
	return e.entry.payload.meta.auditInfo.eTag
}

func (e *PlainEntry) Tag() string {
	return e.entry.payload.state.tag
}

func (e *PlainEntry) Flags() int {
	return e.entry.payload.state.flags
}

func (e *PlainEntry) HasPendingMap() bool {
	return e.entry.hasPendingMap()
}
