package consensus

// PBFT算法的接口，包含了四个步骤：
// 1.开始共识
// 2.预准备
// 3.准备
// 4.提交
type PBFT interface {
	StartConsensus(request *RequestMsg) (*PrePrepareMsg, error)
	PrePrepare(prePrepareMsg *PrePrepareMsg) (*VoteMsg, error)
	Prepare(prepareMsg *VoteMsg) (*VoteMsg, error)
	Commit(commitMsg *VoteMsg) (*ReplyMsg, *RequestMsg, error)
}