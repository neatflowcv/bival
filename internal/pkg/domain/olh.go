package domain

type OLHParams struct {
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

type OLH struct {
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

func NewOLH(p OLHParams) *OLH {
	return &OLH{
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

func (e *OLH) Name() string {
	return e.name
}

func (e *OLH) Instance() string {
	return e.instance
}

func (e *OLH) DeleteMarker() bool {
	return e.deleteMarker
}

func (e *OLH) HasPendingLog() bool {
	return len(e.pendingLogs) > 0
}
