package domain

type PendingLog struct {
	key int
	val []*PendingLogVal
}

func NewPendingLog(key int, val []*PendingLogVal) *PendingLog {
	return &PendingLog{
		key: key,
		val: val,
	}
}
