package domain

type Plain struct {
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
	mTime            string
	eTag             string
	storageClass     string
	contentType      string
	ownerUserID      string
	ownerDisplayName string
	pendingMaps      []*PendingMap
}

func NewPlain(p DirEntryParams) *Plain {
	return &Plain{
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

func (e *Plain) Index() string {
	return string(e.index)
}

func (e *Plain) Name() string {
	return e.name
}

func (e *Plain) Instance() string {
	return e.instance
}

func (e *Plain) VersionPool() int {
	return e.pool
}

func (e *Plain) VersionEpoch() int {
	return e.epoch
}

func (e *Plain) VersionedEpoch() int {
	return e.vEpoch
}

func (e *Plain) Exists() bool {
	return e.exists
}

func (e *Plain) MTime() string {
	return e.mTime
}

func (e *Plain) ETag() string {
	return e.eTag
}

func (e *Plain) Tag() string {
	return e.tag
}

func (e *Plain) Flags() int {
	return e.flags
}

func (e *Plain) HasPendingMap() bool {
	return len(e.pendingMaps) > 0
}

func (e *Plain) IsPlaceholder() bool {
	return e.hasPlaceholderIdentity() &&
		e.hasPlaceholderVersion() &&
		e.hasPlaceholderState() &&
		e.hasPlaceholderMeta()
}

func (e *Plain) hasPlaceholderIdentity() bool {
	return string(e.index) == e.name &&
		e.instance == ""
}

func (e *Plain) hasPlaceholderVersion() bool {
	return e.pool == -1 &&
		e.epoch == 0 &&
		e.vEpoch == 0
}

func (e *Plain) hasPlaceholderState() bool {
	return !e.exists &&
		e.locator == "" &&
		e.tag == "" &&
		e.flags == 8 &&
		len(e.pendingMaps) == 0
}

func (e *Plain) hasPlaceholderMeta() bool {
	return e.category == 0 &&
		e.size == 0 &&
		e.accountedSize == 0 &&
		!e.appendable &&
		e.mTime == "0.000000" &&
		e.eTag == "" &&
		e.storageClass == "" &&
		e.contentType == "" &&
		e.ownerUserID == "" &&
		e.ownerDisplayName == ""
}
