package domain

type OLHEntryParams struct {
	Kind    string
	Index   []byte
	Payload *OLHPayload
}

type OLHEntry struct {
	kind    string
	index   []byte
	payload *OLHPayload
}

func NewOLHEntry(p OLHEntryParams) *OLHEntry {
	return &OLHEntry{
		kind:    p.Kind,
		index:   p.Index,
		payload: p.Payload,
	}
}

func (e *OLHEntry) Name() string {
	return e.payload.key.name
}

func (e *OLHEntry) HasPendingLog() bool {
	return len(e.payload.pendingLogs) > 0
}

func (e *OLHEntry) Payload() *OLHPayload {
	if e == nil {
		return nil
	}

	return e.payload
}
