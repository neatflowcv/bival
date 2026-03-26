package domain

type Meta struct {
	objectSpec  *ObjectSpec
	auditInfo   *AuditInfo
	contentInfo *ContentInfo
	owner       *Owner
}

func NewMeta(
	objectSpec *ObjectSpec,
	auditInfo *AuditInfo,
	contentInfo *ContentInfo,
	owner *Owner,
) *Meta {
	return &Meta{
		objectSpec:  objectSpec,
		auditInfo:   auditInfo,
		contentInfo: contentInfo,
		owner:       owner,
	}
}

func (m *Meta) IsDefault() bool {
	return m.objectSpec != nil &&
		m.auditInfo != nil &&
		m.contentInfo != nil &&
		m.owner != nil &&
		m.objectSpec.IsDefault() &&
		m.auditInfo.IsDefault() &&
		m.contentInfo.IsDefault() &&
		m.owner.IsDefault()
}
