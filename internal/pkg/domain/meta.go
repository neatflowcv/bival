package domain

type Meta struct {
	objectSpec       *ObjectSpec
	auditInfo        *AuditInfo
	contentInfo      *ContentInfo
	ownerUserID      string
	ownerDisplayName string
}

func NewMeta(
	objectSpec *ObjectSpec,
	auditInfo *AuditInfo,
	contentInfo *ContentInfo,
	ownerUserID string,
	ownerDisplayName string,
) *Meta {
	return &Meta{
		objectSpec:       objectSpec,
		auditInfo:        auditInfo,
		contentInfo:      contentInfo,
		ownerUserID:      ownerUserID,
		ownerDisplayName: ownerDisplayName,
	}
}

func (m *Meta) IsDefault() bool {
	return m.objectSpec != nil &&
		m.auditInfo != nil &&
		m.contentInfo != nil &&
		m.objectSpec.IsDefault() &&
		m.auditInfo.IsDefault() &&
		m.contentInfo.IsDefault() &&
		m.ownerUserID == "" &&
		m.ownerDisplayName == ""
}

func (m *Meta) HasParts() bool {
	return m != nil &&
		m.objectSpec != nil &&
		m.auditInfo != nil &&
		m.contentInfo != nil
}
