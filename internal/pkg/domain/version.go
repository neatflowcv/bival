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
