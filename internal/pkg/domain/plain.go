package domain

import "time"

var _ MyEntry = (*Plain)(nil)

type Plain struct {
	Idx              []byte
	Name             string
	Instance         string
	Pool             int
	Epoch            int
	Locator          string
	Exists           bool
	Category         int
	Size             int
	Mtime            time.Time
	Etag             string
	StorageClass     string
	Owner            string
	OwnerDisplayName string
	ContentType      string
	AccountedSize    int
	UserData         string
	Appendable       bool
	Tag              string
	Flags            int
	PendingMap       []any
	VersionedEpoch   int
}

func (p *Plain) Validate() error {
	if p.Name == "" {
		return errEmptyName
	}

	return nil
}
