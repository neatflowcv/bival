package domain

type DirVersionInfo struct {
	version        *Version
	versionedEpoch int
}

func NewDirVersionInfo(version *Version, versionedEpoch int) *DirVersionInfo {
	return &DirVersionInfo{
		version:        version,
		versionedEpoch: versionedEpoch,
	}
}

func (i *DirVersionInfo) Version() *Version {
	return i.version
}

func (i *DirVersionInfo) VersionedEpoch() int {
	return i.versionedEpoch
}
