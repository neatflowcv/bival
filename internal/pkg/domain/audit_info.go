package domain

type AuditInfo struct {
	mTime string
	eTag  string
}

func NewAuditInfo(mTime string, eTag string) *AuditInfo {
	return &AuditInfo{
		mTime: mTime,
		eTag:  eTag,
	}
}
