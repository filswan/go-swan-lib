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
}

func GetDealConfig(verifiedDeal, fastRetrieval, skipConfirmation bool, minerPrice decimal.Decimal, startEpoch, duration int, minerFid, senderWallet, transferType string) *DealConfig {
	dealConfig := &DealConfig{
		VerifiedDeal:     verifiedDeal,
		FastRetrieval:    fastRetrieval,
		SkipConfirmation: skipConfirmation,
		MinerPrice:       minerPrice,
		StartEpoch:       startEpoch,
		MinerFid:         minerFid,
		SenderWallet:     senderWallet,
		Duration:         duration,
		TransferType:     transferType,
	}

	return dealConfig
}
