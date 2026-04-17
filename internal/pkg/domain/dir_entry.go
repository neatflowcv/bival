package domain

type DirEntryParams struct {
	Kind             string
	Index            []byte
	Name             string
	Instance         string
	Pool             int
	Epoch            int
	VEpoch           int
	Locator          string
	Exists           bool
	Tag              string
	Flags            int
	Category         int
	Size             int64
	AccountedSize    int64
	Appendable       bool
	MTime            string
	ETag             string
	StorageClass     string
	ContentType      string
	OwnerUserID      string
	OwnerDisplayName string
	PendingMaps      []*PendingMap
}
