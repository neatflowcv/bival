package domain

type InstanceEntry struct {
	entry *DirEntry
}

func NewInstanceEntry(entry *DirEntry) *InstanceEntry {
	return &InstanceEntry{
		entry: entry,
	}
}
