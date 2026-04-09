package entry

type OLHParmas struct {
	Index           []byte
	Name            string
	Version         string
	DeleteMarker    bool
	Epoch           uint64
	PendingVersions []string
	Tag             string
	Exists          bool
	PendingRemoval  bool
}

type OLH struct {
	index           []byte
	name            string
	version         string
	deleteMarker    bool
	epoch           uint64
	pendingVersions []string
	tag             string
	exists          bool
	pendingRemoval  bool
}

func NewOLH(p OLHParmas) *OLH {
	return &OLH{
		index:           p.Index,
		name:            p.Name,
		version:         p.Version,
		deleteMarker:    p.DeleteMarker,
		epoch:           p.Epoch,
		pendingVersions: p.PendingVersions,
		tag:             p.Tag,
		exists:          p.Exists,
		pendingRemoval:  p.PendingRemoval,
	}
}
