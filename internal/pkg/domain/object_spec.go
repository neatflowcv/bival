package domain

type ObjectSpec struct {
	category      int
	size          int64
	accountedSize int64
	appendable    bool
}

func NewObjectSpec(category int, size int64, accountedSize int64, appendable bool) *ObjectSpec {
	return &ObjectSpec{
		category:      category,
		size:          size,
		accountedSize: accountedSize,
		appendable:    appendable,
	}
}
