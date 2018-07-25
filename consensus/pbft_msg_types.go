package consensus

type RequestMsg struct {
	Timestamp  int64  `json:"timestamp"` 	// 时间戳
	ClientID   string `json:"clientID"`  	// 客户端ID
	Operation  string `json:"operation"`	// 需要达成共识的操作
	SequenceID int64  `json:"sequenceID"`	//
}

type ReplyMsg struct {
	ViewID    int64  `json:"viewID"`
	Timestamp int64  `json:"timestamp"`
	ClientID  string `json:"clientID"`
	NodeID    string `json:"nodeID"`
	Result    string `json:"result"`
}

type PrePrepareMsg struct {
	ViewID     int64       `json:"viewID"`
	SequenceID int64       `json:"sequenceID"`
	Digest     string      `json:"digest"`
	RequestMsg *RequestMsg `json:"requestMsg"`
}

type VoteMsg struct {
	ViewID     int64  `json:"viewID"`
	SequenceID int64  `json:"sequenceID"`
	Digest     string `json:"digest"`
	NodeID     string `json:"nodeID"`
	MsgType           `json:"msgType"`
}

type MsgType int
const (
	PrepareMsg MsgType = iota
	CommitMsg
)
