package domain

type DirPayload struct {
	key         *Key
	versionInfo *DirVersionInfo
	state       *DirState
	meta        *Meta
	pendingMaps []*PendingMap
}

func NewDirPayload(
	key *Key,
	versionInfo *DirVersionInfo,
	state *DirState,
	meta *Meta,
	pendingMaps []*PendingMap,
) *DirPayload {
	return &DirPayload{
		key:         key,
		versionInfo: versionInfo,
		state:       state,
		meta:        meta,
		pendingMaps: pendingMaps,
	}
}

func (p *DirPayload) Key() *Key {
	if p == nil {
		return nil
	}

	return p.key
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
