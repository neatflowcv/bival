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
