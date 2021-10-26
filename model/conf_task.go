package model

type ConfTask struct {
	SwanApiUrl                 string
	SwanApiKey                 string
	SwanAccessToken            string
	PublicDeal                 bool
	BidMode                    int
	VerifiedDeal               bool
	OfflineMode                bool
	FastRetrieval              bool
	MaxPrice                   string
	StorageServerType          string
	WebServerDownloadUrlPrefix string
	ExpireDays                 int
	OutputDir                  string
	InputDir                   string
	TaskName                   *string
	MinerFid                   *string
	Dataset                    *string
	Description                *string
	StartEpoch                 int
	StartEpochIntervalHours    int
}
