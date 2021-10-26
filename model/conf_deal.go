package model

import (
	"github.com/shopspring/decimal"
)

type ConfDeal struct {
	SwanApiUrl              string
	SwanApiKey              string
	SwanAccessToken         string
	SenderWallet            string
	MaxPrice                decimal.Decimal
	VerifiedDeal            bool
	FastRetrieval           bool
	SkipConfirmation        bool
	MinerPrice              decimal.Decimal
	StartEpoch              int
	StartEpochIntervalHours int
	OutputDir               string
	MinerFid                *string
	MetadataJsonPath        *string
}
