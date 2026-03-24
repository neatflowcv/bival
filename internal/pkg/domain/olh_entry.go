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
