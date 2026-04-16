package domain

type OLHEntryParams struct {
	Kind           string
	Index          []byte
	Name           string
	Instance       string
	DeleteMarker   bool
	PendingRemoval bool
	Exists         bool
	Epoch          int
	PendingLogs    []*PendingLog
	Tag            string
}

type OLHEntry struct {
	kind           string
	index          []byte
	name           string
	instance       string
	deleteMarker   bool
	pendingRemoval bool
	exists         bool
	epoch          int
	pendingLogs    []*PendingLog
	tag            string
}

func NewOLHEntry(p OLHEntryParams) *OLHEntry {
	return &OLHEntry{
		kind:           p.Kind,
		index:          p.Index,
		name:           p.Name,
		instance:       p.Instance,
		deleteMarker:   p.DeleteMarker,
		pendingRemoval: p.PendingRemoval,
		exists:         p.Exists,
		epoch:          p.Epoch,
		pendingLogs:    p.PendingLogs,
		tag:            p.Tag,
	}
}

func (e *OLHEntry) Name() string {
	return e.name
}

func (e *OLHEntry) Instance() string {
	return e.instance
}

func (e *OLHEntry) HasPendingLog() bool {
	return len(e.pendingLogs) > 0
}
