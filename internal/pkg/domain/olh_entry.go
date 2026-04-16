package domain

type OLHEntryParams struct {
	Kind        string
	Index       []byte
	Key         *Key
	State       *OLHState
	Epoch       int
	PendingLogs []*PendingLog
	Tag         string
}

type OLHEntry struct {
	kind        string
	index       []byte
	key         *Key
	state       *OLHState
	epoch       int
	pendingLogs []*PendingLog
	tag         string
}

func NewOLHEntry(p OLHEntryParams) *OLHEntry {
	return &OLHEntry{
		kind:        p.Kind,
		index:       p.Index,
		key:         p.Key,
		state:       p.State,
		epoch:       p.Epoch,
		pendingLogs: p.PendingLogs,
		tag:         p.Tag,
	}
}

func (e *OLHEntry) Name() string {
	return e.key.Name()
}

func (e *OLHEntry) Instance() string {
	return e.key.Instance()
}

func (e *OLHEntry) HasPendingLog() bool {
	return len(e.pendingLogs) > 0
}
