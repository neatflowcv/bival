package domain

import "time"

type DirPayload struct {
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

func NewDirPayload(
	name string,
	instance string,
	pool int,
	epoch int,
	vEpoch int,
	locator string,
	exists bool,
	tag string,
	flags int,
	category int,
	size int64,
	accountedSize int64,
	appendable bool,
	mTime time.Time,
	eTag string,
	storageClass string,
	contentType string,
	ownerUserID string,
	ownerDisplayName string,
	pendingMaps []*PendingMap,
) *DirPayload {
	return &DirPayload{
		name:             name,
		instance:         instance,
		pool:             pool,
		epoch:            epoch,
		vEpoch:           vEpoch,
		locator:          locator,
		exists:           exists,
		tag:              tag,
		flags:            flags,
		category:         category,
		size:             size,
		accountedSize:    accountedSize,
		appendable:       appendable,
		mTime:            mTime,
		eTag:             eTag,
		storageClass:     storageClass,
		contentType:      contentType,
		ownerUserID:      ownerUserID,
		ownerDisplayName: ownerDisplayName,
		pendingMaps:      pendingMaps,
	}
}

func (p *DirPayload) Name() string {
	if p == nil {
		return ""
	}

	return p.name
}

func (p *DirPayload) Instance() string {
	if p == nil {
		return ""
	}

	return p.instance
}

func (p *DirPayload) Pool() int {
	if p == nil {
		return 0
	}

	return p.pool
}

func (p *DirPayload) Epoch() int {
	if p == nil {
		return 0
	}

	return p.epoch
}

func (p *DirPayload) VersionedEpoch() int {
	if p == nil {
		return 0
	}

	return p.vEpoch
}

func (p *DirPayload) IsVersionMissing() bool {
	if p == nil {
		return false
	}

	return p.pool == -1 && p.epoch == 0
}

func (p *DirPayload) Locator() string {
	if p == nil {
		return ""
	}

	return p.locator
}

func (p *DirPayload) Exists() bool {
	if p == nil {
		return false
	}

	return p.exists
}

func (p *DirPayload) Tag() string {
	if p == nil {
		return ""
	}

	return p.tag
}

func (p *DirPayload) Flags() int {
	if p == nil {
		return 0
	}

	return p.flags
}

func (p *DirPayload) HasMetaParts() bool {
	return p != nil
}

func (p *DirPayload) IsMetaDefault() bool {
	if p == nil {
		return false
	}

	return p.isDefaultObjectMeta() &&
		p.isDefaultAuditMeta() &&
		p.isDefaultContentMeta() &&
		p.isDefaultOwnerMeta()
}

func (p *DirPayload) PendingMaps() []*PendingMap {
	if p == nil {
		return nil
	}

	return p.pendingMaps
}

func (p *DirPayload) isDefaultObjectMeta() bool {
	return p.category == 0 &&
		p.size == 0 &&
		p.accountedSize == 0 &&
		!p.appendable
}

func (p *DirPayload) isDefaultAuditMeta() bool {
	return p.mTime.IsZero() &&
		p.eTag == ""
}

func (p *DirPayload) isDefaultContentMeta() bool {
	return p.storageClass == "" &&
		p.contentType == ""
}

func (p *DirPayload) isDefaultOwnerMeta() bool {
	return p.ownerUserID == "" &&
		p.ownerDisplayName == ""
}
