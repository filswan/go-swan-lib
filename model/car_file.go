package model

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
	DealCids       []string
	DataCid        string
	PieceCid       string
	MinerFids      []string
	StartEpoch     *int
	SourceId       *int
}
