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
