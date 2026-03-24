package domain

type ContentInfo struct {
	storageClass string
	contentType  string
}

func NewContentInfo(storageClass string, contentType string) *ContentInfo {
	return &ContentInfo{
		storageClass: storageClass,
		contentType:  contentType,
	}
}
