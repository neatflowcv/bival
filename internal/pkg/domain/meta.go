package domain

import "time"

type Meta struct {
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
}

func NewMeta(
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
) *Meta {
	return &Meta{
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
	}
}

func (m *Meta) IsDefault() bool {
	return m.category == 0 &&
		m.size == 0 &&
		m.accountedSize == 0 &&
		!m.appendable &&
		m.mTime.IsZero() &&
		m.eTag == "" &&
		m.storageClass == "" &&
		m.contentType == "" &&
		m.ownerUserID == "" &&
		m.ownerDisplayName == ""
}

func (m *Meta) HasParts() bool {
	return m != nil
}
