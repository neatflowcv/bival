package domain

type OLHState struct {
	deleteMarker   bool
	pendingRemoval bool
	exists         bool
}

func NewOLHState(deleteMarker bool, pendingRemoval bool, exists bool) *OLHState {
	return &OLHState{
		deleteMarker:   deleteMarker,
		pendingRemoval: pendingRemoval,
		exists:         exists,
	}
}
