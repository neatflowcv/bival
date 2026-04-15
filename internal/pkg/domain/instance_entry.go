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
	return e.entry.payload.key.name
}

func (e *InstanceEntry) HasPendingMap() bool {
	return e.entry.hasPendingMap()
}

func (e *InstanceEntry) Payload() *DirPayload {
	if e == nil || e.entry == nil {
		return nil
	}

	return e.entry.payload
}
