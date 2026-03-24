package domain

type PendingMap struct {
	key string
	val *PendingMapVal
}

func NewPendingMap(key string, val *PendingMapVal) *PendingMap {
	return &PendingMap{
		key: key,
		val: val,
	}
}
