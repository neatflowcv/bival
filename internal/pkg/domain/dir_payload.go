package domain

type DirPayload struct {
	name        string
	instance    string
	versionInfo *DirVersionInfo
	locator     string
	exists      bool
	tag         string
	flags       int
	meta        *Meta
	pendingMaps []*PendingMap
}

func NewDirPayload(
	name string,
	instance string,
	versionInfo *DirVersionInfo,
	locator string,
	exists bool,
	tag string,
	flags int,
	meta *Meta,
	pendingMaps []*PendingMap,
) *DirPayload {
	return &DirPayload{
		name:        name,
		instance:    instance,
		versionInfo: versionInfo,
		locator:     locator,
		exists:      exists,
		tag:         tag,
		flags:       flags,
		meta:        meta,
		pendingMaps: pendingMaps,
	}
}

func (p *DirPayload) Name() string {
	if p == nil {
		return ""
	}

	return p.name
}

func (p *DirPayload) Instance() string {
	if p == nil {
		return ""
	}

	return p.instance
}

func (p *DirPayload) VersionInfo() *DirVersionInfo {
	if p == nil {
		return nil
	}

	return p.versionInfo
}

func (p *DirPayload) Locator() string {
	if p == nil {
		return ""
	}

	return p.locator
}

func (p *DirPayload) Exists() bool {
	if p == nil {
		return false
	}

	return p.exists
}

func (p *DirPayload) Tag() string {
	if p == nil {
		return ""
	}

	return p.tag
}

func (p *DirPayload) Flags() int {
	if p == nil {
		return 0
	}

	return p.flags
}

func (p *DirPayload) Meta() *Meta {
	if p == nil {
		return nil
	}

	return p.meta
}

func (p *DirPayload) PendingMaps() []*PendingMap {
	if p == nil {
		return nil
	}

	return p.pendingMaps
}
