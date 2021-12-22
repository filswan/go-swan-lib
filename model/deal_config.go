package model

import (
	"github.com/shopspring/decimal"
)

type DealConfig struct {
	VerifiedDeal     bool
	FastRetrieval    bool
	SkipConfirmation bool
	MinerPrice       decimal.Decimal
	StartEpoch       int
	MinerFid         string
	SenderWallet     string
	Duration         int
	TransferType     string
	PayloadCid       string
	PieceCid         string
	FileSize         int64
}
