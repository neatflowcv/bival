package domain

var _ MyEntry = (*OLH)(nil)

type OLH struct {
	Idx            []byte
	Name           string
	Instance       string
	DeleteMarker   bool
	Epoch          int
	PendingLog     []any
	Tag            string
	Exists         bool
	PendingRemoval bool
}

func (o *OLH) Validate() error {
	if o.Name == "" {
		return errEmptyName
	}

	return nil
}
