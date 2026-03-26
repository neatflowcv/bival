package domain

type InstanceEntry struct {
	entry *DirEntry
}

func NewInstanceEntry(entry *DirEntry) *InstanceEntry {
	return &InstanceEntry{
		entry: entry,
	}
}

func (e *InstanceEntry) HasPendingMap() bool {
	return e.entry.hasPendingMap()
}
