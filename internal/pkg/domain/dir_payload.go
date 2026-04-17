package domain

type DirPayload struct {
	name        string
	instance    string
	versionInfo *DirVersionInfo
	state       *DirState
	meta        *Meta
	pendingMaps []*PendingMap
}

func NewDirPayload(
	name string,
	instance string,
	versionInfo *DirVersionInfo,
	state *DirState,
	meta *Meta,
	pendingMaps []*PendingMap,
) *DirPayload {
	return &DirPayload{
		name:        name,
		instance:    instance,
		versionInfo: versionInfo,
		state:       state,
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

func (p *DirPayload) State() *DirState {
	if p == nil {
		return nil
	}

	return p.state
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
