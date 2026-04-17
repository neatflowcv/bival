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
	return e.entry.key.name
}

func (e *PlainEntry) Instance() string {
	return e.entry.key.instance
}

func (e *PlainEntry) VersionPool() int {
	return e.entry.versionInfo.Version().Pool()
}

func (e *PlainEntry) VersionEpoch() int {
	return e.entry.versionInfo.Version().Epoch()
}

func (e *PlainEntry) Exists() bool {
	return e.entry.state.exists
}

func (e *PlainEntry) MTime() time.Time {
	return e.entry.meta.auditInfo.mTime
}

func (e *PlainEntry) ETag() string {
	return e.entry.meta.auditInfo.eTag
}

func (e *PlainEntry) Tag() string {
	return e.entry.state.tag
}

func (e *PlainEntry) Flags() int {
	return e.entry.state.flags
}

func (e *PlainEntry) HasPendingMap() bool {
	return e.entry.hasPendingMap()
}

func (e *PlainEntry) Payload() *DirPayload {
	if e == nil || e.entry == nil {
		return nil
	}

	return NewDirPayload(
		e.entry.key,
		e.entry.versionInfo,
		e.entry.state,
		e.entry.meta,
		e.entry.pendingMaps,
	)
}
