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
