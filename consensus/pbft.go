package consensus

type PBFT interface {
	StartConsensus()
	PrePrePare()
	Prepare()
	Commit()
}

type Consensus struct{

}