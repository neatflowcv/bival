package domain

type PendingLogVal struct {
	epoch        int
	op           string
	opTag        string
	name         string
	instance     string
	deleteMarker bool
}

func NewPendingLogVal(
	epoch int,
	op string,
	opTag string,
	name string,
	instance string,
	deleteMarker bool,
) *PendingLogVal {
	return &PendingLogVal{
		epoch:        epoch,
		op:           op,
		opTag:        opTag,
		name:         name,
		instance:     instance,
		deleteMarker: deleteMarker,
	}
}
