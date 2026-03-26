package domain

type PlainEntry struct {
	entry *DirEntry
}

func NewPlainEntry(entry *DirEntry) *PlainEntry {
	return &PlainEntry{
		entry: entry,
	}
}

func (e *PlainEntry) Name() string {
	return e.entry.payload.key.name
}

func (e *PlainEntry) HasPendingMap() bool {
	return e.entry.hasPendingMap()
}
