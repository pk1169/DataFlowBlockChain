package consensus


// State notes the state of consenesus
type State struct {
	ViewID 		int64
	MsgLogs  	*MsgLogs
}

// MsgLogs logs the msg for consensus
types MsgLogs


// Stage marks the present stage of consensus
type Stage int

const (
	Idle  Stage = iota
	PrePrepared
	Prepared
	Committed
)

// the number of nodes can be tolerated to be attacked
const f = 1

// la
func CreateState(viewID int64, lastSequenceID  int64)