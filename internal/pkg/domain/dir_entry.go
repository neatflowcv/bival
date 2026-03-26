package domain

type DirEntry struct {
	kind    string
	index   []byte
	payload *DirPayload
}

func NewDirEntry(kind string, index []byte, payload *DirPayload) *DirEntry {
	return &DirEntry{
		kind:    kind,
		index:   index,
		payload: payload,
	}
}

func (e *DirEntry) hasPendingMap() bool {
	return len(e.payload.pendingMaps) > 0
}
