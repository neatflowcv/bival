package domain

type OLHPayload struct {
	key         *Key
	state       *OLHState
	epoch       int
	pendingLogs []*PendingLog
	tag         string
}

func NewOLHPayload(
	key *Key,
	state *OLHState,
	epoch int,
	pendingLogs []*PendingLog,
	tag string,
) *OLHPayload {
	return &OLHPayload{
		key:         key,
		state:       state,
		epoch:       epoch,
		pendingLogs: pendingLogs,
		tag:         tag,
	}
}
