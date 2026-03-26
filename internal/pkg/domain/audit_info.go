package domain

import "time"

type AuditInfo struct {
	mTime time.Time
	eTag  string
}

func NewAuditInfo(mTime time.Time, eTag string) *AuditInfo {
	return &AuditInfo{
		mTime: mTime,
		eTag:  eTag,
	}
}
