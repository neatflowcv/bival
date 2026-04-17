package domain

type Meta struct {
	objectSpec       *ObjectSpec
	auditInfo        *AuditInfo
	storageClass     string
	contentType      string
	ownerUserID      string
	ownerDisplayName string
}

func NewMeta(
	objectSpec *ObjectSpec,
	auditInfo *AuditInfo,
	storageClass string,
	contentType string,
	ownerUserID string,
	ownerDisplayName string,
) *Meta {
	return &Meta{
		objectSpec:       objectSpec,
		auditInfo:        auditInfo,
		storageClass:     storageClass,
		contentType:      contentType,
		ownerUserID:      ownerUserID,
		ownerDisplayName: ownerDisplayName,
	}
}

func (m *Meta) IsDefault() bool {
	return m.objectSpec != nil &&
		m.auditInfo != nil &&
		m.objectSpec.IsDefault() &&
		m.auditInfo.IsDefault() &&
		m.storageClass == "" &&
		m.contentType == "" &&
		m.ownerUserID == "" &&
		m.ownerDisplayName == ""
}

func (m *Meta) HasParts() bool {
	return m != nil &&
		m.objectSpec != nil &&
		m.auditInfo != nil
}
