package domain

import "time"

type Instance struct {
	kind             string
	index            []byte
	name             string
	instance         string
	pool             int
	epoch            int
	vEpoch           int
	locator          string
	exists           bool
	tag              string
	flags            int
	category         int
	size             int64
	accountedSize    int64
	appendable       bool
	mTime            time.Time
	eTag             string
	storageClass     string
	contentType      string
	ownerUserID      string
	ownerDisplayName string
	pendingMaps      []*PendingMap
}

func NewInstance(p DirEntryParams) *Instance {
	return &Instance{
		kind:             p.Kind,
		index:            p.Index,
		name:             p.Name,
		instance:         p.Instance,
		pool:             p.Pool,
		epoch:            p.Epoch,
		vEpoch:           p.VEpoch,
		locator:          p.Locator,
		exists:           p.Exists,
		tag:              p.Tag,
		flags:            p.Flags,
		category:         p.Category,
		size:             p.Size,
		accountedSize:    p.AccountedSize,
		appendable:       p.Appendable,
		mTime:            p.MTime,
		eTag:             p.ETag,
		storageClass:     p.StorageClass,
		contentType:      p.ContentType,
		ownerUserID:      p.OwnerUserID,
		ownerDisplayName: p.OwnerDisplayName,
		pendingMaps:      p.PendingMaps,
	}
}

func (e *Instance) Name() string {
	return e.name
}

func (e *Instance) Instance() string {
	return e.instance
}

func (e *Instance) VersionPool() int {
	return e.pool
}

func (e *Instance) VersionEpoch() int {
	return e.epoch
}

func (e *Instance) VersionedEpoch() int {
	return e.vEpoch
}

func (e *Instance) HasPendingMap() bool {
	return len(e.pendingMaps) > 0
}
