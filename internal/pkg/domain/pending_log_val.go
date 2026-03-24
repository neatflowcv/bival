package domain

type PendingLogVal struct {
	epoch        int
	op           string
	opTag        string
	key          *Key
	deleteMarker bool
}

func NewPendingLogVal(epoch int, op string, opTag string, key *Key, deleteMarker bool) *PendingLogVal {
	return &PendingLogVal{
		epoch:        epoch,
		op:           op,
		opTag:        opTag,
		key:          key,
		deleteMarker: deleteMarker,
	}
}
