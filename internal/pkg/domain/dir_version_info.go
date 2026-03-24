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
