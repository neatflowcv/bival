package domain

type DirState struct {
	locator string
	exists  bool
	tag     string
	flags   int
}

func NewDirState(locator string, exists bool, tag string, flags int) *DirState {
	return &DirState{
		locator: locator,
		exists:  exists,
		tag:     tag,
		flags:   flags,
	}
}

func (s *DirState) Locator() string {
	if s == nil {
		return ""
	}

	return s.locator
}

func (s *DirState) Exists() bool {
	if s == nil {
		return false
	}

	return s.exists
}

func (s *DirState) Tag() string {
	if s == nil {
		return ""
	}

	return s.tag
}

func (s *DirState) Flags() int {
	if s == nil {
		return 0
	}

	return s.flags
}
