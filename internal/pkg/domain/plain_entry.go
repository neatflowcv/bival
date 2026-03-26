package domain

type PlainEntry struct {
	entry *DirEntry
}

func NewPlainEntry(entry *DirEntry) *PlainEntry {
	return &PlainEntry{
		entry: entry,
	}
}

func (e *PlainEntry) HasPendingMap() bool {
	return e.entry.hasPendingMap()
}
