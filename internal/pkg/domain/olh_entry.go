package domain

type OLHEntry struct {
	kind    string
	index   []byte
	payload *OLHPayload
}

func NewOLHEntry(kind string, index []byte, payload *OLHPayload) *OLHEntry {
	return &OLHEntry{
		kind:    kind,
		index:   index,
		payload: payload,
	}
}

func (e *OLHEntry) Name() string {
	return e.payload.key.name
}

func (e *OLHEntry) HasPendingLog() bool {
	return len(e.payload.pendingLogs) > 0
}
