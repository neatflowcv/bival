package domain

import "time"

type Meta struct {
	objectSpec       *ObjectSpec
	mTime            time.Time
	eTag             string
	storageClass     string
	contentType      string
	ownerUserID      string
	ownerDisplayName string
}

func NewMeta(
	objectSpec *ObjectSpec,
	mTime time.Time,
	eTag string,
	storageClass string,
	contentType string,
	ownerUserID string,
	ownerDisplayName string,
) *Meta {
	return &Meta{
		objectSpec:       objectSpec,
		mTime:            mTime,
		eTag:             eTag,
		storageClass:     storageClass,
		contentType:      contentType,
		ownerUserID:      ownerUserID,
		ownerDisplayName: ownerDisplayName,
	}
}

func (m *Meta) IsDefault() bool {
	return m.objectSpec != nil &&
		m.objectSpec.IsDefault() &&
		m.mTime.IsZero() &&
		m.eTag == "" &&
		m.storageClass == "" &&
		m.contentType == "" &&
		m.ownerUserID == "" &&
		m.ownerDisplayName == ""
}

func (m *Meta) HasParts() bool {
	return m != nil &&
		m.objectSpec != nil
}
