package domain

type DirVersionInfo struct {
	pool           int
	epoch          int
	versionedEpoch int
}

func NewDirVersionInfo(pool int, epoch int, versionedEpoch int) *DirVersionInfo {
	return &DirVersionInfo{
		pool:           pool,
		epoch:          epoch,
		versionedEpoch: versionedEpoch,
	}
}

func (i *DirVersionInfo) Pool() int {
	return i.pool
}

func (i *DirVersionInfo) Epoch() int {
	return i.epoch
}

func (i *DirVersionInfo) VersionedEpoch() int {
	return i.versionedEpoch
}

func (i *DirVersionInfo) IsMissing() bool {
	return i.pool == -1 && i.epoch == 0
}
