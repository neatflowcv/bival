package domain

type DirEntryParams struct {
	Kind    string
	Index   []byte
	Payload *DirPayload
}

type DirEntry struct {
	kind    string
	index   []byte
	payload *DirPayload
}

func NewDirEntry(p DirEntryParams) *DirEntry {
	return &DirEntry{
		kind:    p.Kind,
		index:   p.Index,
		payload: p.Payload,
	}
}

func (e *DirEntry) hasPendingMap() bool {
	return len(e.payload.pendingMaps) > 0
}

func (e *DirEntry) indexString() string {
	return string(e.index)
}
