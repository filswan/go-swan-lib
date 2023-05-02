package boost

import (
	"context"
	"fmt"
	"github.com/docker/go-units"
	boostapi "github.com/filecoin-project/boost/api"
	"github.com/filecoin-project/boost/storagemarket/types"
	"github.com/filecoin-project/boost/storagemarket/types/dealcheckpoints"
	"github.com/filecoin-project/go-jsonrpc"
	"github.com/filecoin-project/go-state-types/abi"
	"github.com/filecoin-project/lotus/build"
	chain_type "github.com/filecoin-project/lotus/chain/types"
	"github.com/filswan/go-swan-lib/model"
	"github.com/google/uuid"
	"github.com/ipfs/go-cid"
	"github.com/pkg/errors"
	"golang.org/x/xerrors"
	"net/http"
	"time"
)

type Client struct {
	stub boostapi.BoostStruct
}

func NewClient(authToken, apiUrl string) (*Client, jsonrpc.ClientCloser, error) {
	var headers http.Header
	if authToken != "" {
		headers = http.Header{"Authorization": []string{"Bearer " + authToken}}
	} else {
		headers = nil
	}

	var apiSub boostapi.BoostStruct
	closer, err := jsonrpc.NewMergeClient(context.Background(), "ws://"+apiUrl+"/rpc/v0", "Filecoin",
		[]interface{}{&apiSub.Internal}, headers)
	if err != nil {
		return nil, nil, errors.Wrap(err, "connecting with boost failed")
	}

	return &Client{
		stub: apiSub,
	}, closer, nil
}

func (pc *Client) GetDealInfoByDealUuid(ctx context.Context, dealUuid string) (*model.ProviderDealState, error) {
	dealUid, err := uuid.Parse(dealUuid)
	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("dealUuid=[%s] parse failed", dealUid))
	}
	boostDeal, err := pc.stub.BoostDeal(ctx, dealUid)
	if err != nil {
		return nil, err
	}
	var pds model.ProviderDealState
	pds.DealUuid = boostDeal.DealUuid.String()
	pds.IsOffline = boostDeal.IsOffline
	pds.DealDataRoot = boostDeal.DealDataRoot.String()
	pds.ChainDealID = uint64(boostDeal.ChainDealID)
	pds.PublishCID = boostDeal.PublishCID.String()
	pds.DealStatus = statusMessage(&types.DealStatusResponse{
		DealUUID:  boostDeal.DealUuid,
		Error:     boostDeal.Err,
		IsOffline: boostDeal.IsOffline,
		DealStatus: &types.DealStatus{
			Error:  boostDeal.Err,
			Status: boostDeal.Checkpoint.String(),
		},
	})
	pds.Err = boostDeal.Err
	return &pds, nil
}

func (pc *Client) OfflineDealWithData(ctx context.Context, dealUuid, filePath string) (*model.ProviderDealRejectionInfo, error) {
	dealUid, err := uuid.Parse(dealUuid)
	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("dealUuid=[%s] parse failed", dealUid))
	}
	offlineDealWithData, err := pc.stub.BoostOfflineDealWithData(ctx, dealUid, filePath, false)
	if err != nil {
		return nil, err
	}
	return &model.ProviderDealRejectionInfo{
		Accepted: offlineDealWithData.Accepted,
		Reason:   offlineDealWithData.Reason,
	}, nil
}

func (pc *Client) OfflineDealWithDataByMarket(ctx context.Context, proposalCid, filePath string) error {
	propCid, err := cid.Decode(proposalCid)
	if err != nil {
		return fmt.Errorf("could not parse '%s' as deal proposal cid", proposalCid)
	}
	return pc.stub.MarketImportDealData(ctx, propCid, filePath)
}

func (pc *Client) MarketSetAsk(ctx context.Context, price, verifiedPrice, minPieceSize, maxPieceSize string) error {
	pri, err := chain_type.ParseFIL(price)
	if err != nil {
		return err
	}

	vpri, err := chain_type.ParseFIL(verifiedPrice)
	if err != nil {
		return err
	}

	min, err := units.RAMInBytes(minPieceSize)
	if err != nil {
		return xerrors.Errorf("cannot parse min-piece-size to quantity of bytes: %w", err)
	}

	if min < 256 {
		return xerrors.New("minimum piece size (w/bit-padding) is 256B")
	}

	max, err := units.RAMInBytes(maxPieceSize)
	if err != nil {
		return xerrors.Errorf("cannot parse max-piece-size to quantity of bytes: %w", err)
	}
	dur, err := time.ParseDuration("720h0m0s")
	if err != nil {
		return xerrors.Errorf("cannot parse duration: %w", err)
	}

	qty := dur.Seconds() / float64(build.BlockDelaySecs)

	return pc.stub.MarketSetAsk(ctx, chain_type.BigInt(pri), chain_type.BigInt(vpri), abi.ChainEpoch(qty), abi.PaddedPieceSize(min), abi.PaddedPieceSize(max))
}

func (pc *Client) GetDealsConsiderOfflineStorageDeals(ctx context.Context) (bool, error) {
	return pc.stub.DealsConsiderOfflineStorageDeals(ctx)
}

func statusMessage(resp *types.DealStatusResponse) string {
	switch resp.DealStatus.Status {
	case dealcheckpoints.Accepted.String():
		if resp.IsOffline {
			return "Awaiting Offline Data Import"
		}
	case dealcheckpoints.Transferred.String():
		return "Ready to Publish"
	case dealcheckpoints.Published.String():
		return "Awaiting Publish Confirmation"
	case dealcheckpoints.PublishConfirmed.String():
		return "Adding to Sector"
	case dealcheckpoints.AddedPiece.String():
		return "Announcing"
	case dealcheckpoints.IndexedAndAnnounced.String():
		return "Sealing"
	case dealcheckpoints.Complete.String():
		if resp.DealStatus.Error != "" {
			return "Error: " + resp.DealStatus.Error
		}
		return "Expired"
	}
	return resp.DealStatus.Status
}
