package domain

type InstanceEntry struct {
	entry *DirEntry
}

func NewInstanceEntry(entry *DirEntry) *InstanceEntry {
	return &InstanceEntry{
		entry: entry,
	}
}

func (e *InstanceEntry) Name() string {
	return e.entry.key.name
}

func (e *InstanceEntry) HasPendingMap() bool {
	return e.entry.hasPendingMap()
}

func (e *InstanceEntry) Payload() *DirPayload {
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
