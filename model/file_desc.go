package model

type DealInfo struct {
	DealId        int
	DealCid       string
	MinerFid      string
	StartEpoch    int
	Cost          string
	StorageStatus string `json:"StorageStatus,omitempty"`
}
type FileDesc struct {
	Uuid           string
	SourceFileName string
	SourceFilePath string
	SourceFileMd5  string
	SourceFileSize int64
	CarFileName    string
	CarFilePath    string
	CarFileMd5     string
	CarFileUrl     string
	CarFileSize    int64
	PayloadCid     string
	PieceCid       string
	StartEpoch     *int64
	SourceId       *int
	Deals          []*DealInfo
}
