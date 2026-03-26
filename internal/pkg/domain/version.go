package domain

type Version struct {
	pool  int
	epoch int
}

func NewVersion(pool int, epoch int) *Version {
	return &Version{
		pool:  pool,
		epoch: epoch,
	}
}

func (v *Version) Pool() int {
	return v.pool
}

func (v *Version) Epoch() int {
	return v.epoch
}

func (v *Version) IsMissing() bool {
	return v.pool == -1 && v.epoch == 0
}
