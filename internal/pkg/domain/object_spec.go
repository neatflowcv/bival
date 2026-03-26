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

func (s *ObjectSpec) IsDefault() bool {
	return s.category == 0 &&
		s.size == 0 &&
		s.accountedSize == 0 &&
		!s.appendable
}
