package domain

import "time"

type PlainEntry struct {
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

func NewPlainEntry(p DirEntryParams) *PlainEntry {
	return &PlainEntry{
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

func (e *PlainEntry) Index() string {
	return string(e.index)
}

func (e *PlainEntry) Name() string {
	return e.name
}

func (e *PlainEntry) Instance() string {
	return e.instance
}

func (e *PlainEntry) VersionPool() int {
	return e.pool
}

func (e *PlainEntry) VersionEpoch() int {
	return e.epoch
}

func (e *PlainEntry) VersionedEpoch() int {
	return e.vEpoch
}

func (e *PlainEntry) Exists() bool {
	return e.exists
}

func (e *PlainEntry) MTime() time.Time {
	return e.mTime
}

func (e *PlainEntry) ETag() string {
	return e.eTag
}

func (e *PlainEntry) Tag() string {
	return e.tag
}

func (e *PlainEntry) Flags() int {
	return e.flags
}

func (e *PlainEntry) HasPendingMap() bool {
	return len(e.pendingMaps) > 0
}

func (e *PlainEntry) IsPlaceholder() bool {
	return e.hasPlaceholderIdentity() &&
		e.hasPlaceholderVersion() &&
		e.hasPlaceholderState() &&
		e.hasPlaceholderMeta()
}

func (e *PlainEntry) hasPlaceholderIdentity() bool {
	return string(e.index) == e.name &&
		e.instance == ""
}

func (e *PlainEntry) hasPlaceholderVersion() bool {
	return e.pool == -1 &&
		e.epoch == 0 &&
		e.vEpoch == 0
}

func (e *PlainEntry) hasPlaceholderState() bool {
	return !e.exists &&
		e.locator == "" &&
		e.tag == "" &&
		e.flags == 8 &&
		len(e.pendingMaps) == 0
}

func (e *PlainEntry) hasPlaceholderMeta() bool {
	return e.category == 0 &&
		e.size == 0 &&
		e.accountedSize == 0 &&
		!e.appendable &&
		e.mTime.IsZero() &&
		e.eTag == "" &&
		e.storageClass == "" &&
		e.contentType == "" &&
		e.ownerUserID == "" &&
		e.ownerDisplayName == ""
}
