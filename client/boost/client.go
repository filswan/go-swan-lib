package boost

import (
	"bufio"
	"context"
	"encoding/hex"
	"encoding/json"
	"fmt"
	clinode "github.com/filecoin-project/boost/cli/node"
	cliutil "github.com/filecoin-project/boost/cli/util"
	"github.com/filecoin-project/boost/cmd"
	"github.com/filecoin-project/boost/storagemarket/types"
	"github.com/filecoin-project/go-address"
	cborutil "github.com/filecoin-project/go-cbor-util"
	"github.com/filecoin-project/go-state-types/abi"
	"github.com/filecoin-project/go-state-types/big"
	"github.com/filecoin-project/go-state-types/builtin/v9/market"
	"github.com/filecoin-project/lotus/api"
	apiclient "github.com/filecoin-project/lotus/api/client"
	chaintypes "github.com/filecoin-project/lotus/chain/types"
	"github.com/filswan/go-swan-lib/client/lotus"
	"github.com/filswan/go-swan-lib/constants"
	"github.com/filswan/go-swan-lib/logs"
	"github.com/filswan/go-swan-lib/model"
	"github.com/filswan/go-swan-lib/utils"
	"github.com/google/uuid"
	"github.com/ipfs/go-cid"
	inet "github.com/libp2p/go-libp2p/core/network"
	"github.com/mitchellh/go-homedir"
	"github.com/shopspring/decimal"
	"net/url"
	"os"
	"strings"
)

const DealProtocolv120 = "/fil/storage/mk/1.2.0"

type Client struct {
	lotus       *lotus.LotusClient
	FullNodeApi string
	ClientRepo  string
}

func (client *Client) WithUrl(fullNodeApi string) (*Client, error) {
	client.FullNodeApi = fullNodeApi
	apiInfo := cliutil.ParseApiInfo(client.FullNodeApi)
	addr, err := apiInfo.DialArgs("v1")
	if err != nil {
		logs.GetLogger().Error("parse fullNodeApi failed: %w", err)
		return nil, err
	}
	client.lotus = &lotus.LotusClient{
		ApiUrl:      addr,
		AccessToken: string(apiInfo.Token),
	}
	return client, nil
}

func (client *Client) WithClient(lotus *lotus.LotusClient) *Client {
	client.lotus = lotus
	u, err := url.Parse(lotus.ApiUrl)
	if err != nil {
		logs.GetLogger().Error("parse lotus ApiUrl failed: %w", err)
		return nil
	}
	client.FullNodeApi = fmt.Sprintf("%s:/ip4/%s/tcp/%s/http", lotus.AccessToken, u.Hostname(), u.Port())
	return client
}

func GetClient(clientRepo string) *Client {
	if len(clientRepo) == 0 {
		panic("boost repo is required")
	}
	return &Client{
		ClientRepo: clientRepo,
	}
}

func (client *Client) InitRepo(repoPath string) error {
	sdir, err := homedir.Expand(repoPath)
	if err != nil {
		return err
	}
	os.Mkdir(sdir, 0755) //nolint:errcheck

	n, err := clinode.Setup(repoPath)
	if err != nil {
		logs.GetLogger().Error("setup node failed: %w", err)
		return err
	}

	walletAddr, err := n.Wallet.GetDefault()
	if err != nil {
		return err
	}
	logs.GetLogger().Infof("default wallet set: %s", walletAddr)
	return nil
}

func (client *Client) WalletImport(inputData []byte) error {
	ctx := context.Background()
	n, err := clinode.Setup(client.ClientRepo)
	if err != nil {
		logs.GetLogger().Error("setup node failed: %w", err)
		return err
	}

	var ki chaintypes.KeyInfo
	data, err := hex.DecodeString(strings.TrimSpace(string(inputData)))
	if err != nil {
		return err
	}

	if err := json.Unmarshal(data, &ki); err != nil {
		return err
	}

	_, err = n.Wallet.WalletImport(ctx, &ki)
	if err != nil {
		logs.GetLogger().Error("wallet import failed: %w", err)
		return err
	}
	logs.GetLogger().Infof("wallet import successfully")
	return nil
}

func (client *Client) StartDeal(dealConfig *model.DealConfig) (string, error) {
	minerPrice, err := client.lotus.CheckDealConfig(dealConfig)
	if err != nil {
		logs.GetLogger().Error(err)
		return "", err
	}

	pieceSize, sectorSize := utils.CalculatePieceSize(dealConfig.FileSize, true)
	cost := utils.CalculateRealCost(sectorSize, *minerPrice)

	epochPrice := cost.Mul(decimal.NewFromFloat(constants.LOTUS_PRICE_MULTIPLE_1E18))

	if !dealConfig.SkipConfirmation {
		logs.GetLogger().Info("Do you confirm to submit the deal?")
		logs.GetLogger().Info("Press Y/y to continue, other key to quit")
		reader := bufio.NewReader(os.Stdin)
		response, err := reader.ReadString('\n')
		if err != nil {
			logs.GetLogger().Error(err)
			return "", err
		}

		response = strings.TrimRight(response, "\n")

		if !strings.EqualFold(response, "Y") {
			logs.GetLogger().Info("Your input is ", response, ". Now give up submit the deal.")
			return "", err
		}
	}

	dealConfig.PieceCid = strings.Trim(dealConfig.PieceCid, " ")

	dealParam := DealParam{
		Provider:      dealConfig.MinerFid,
		Commp:         dealConfig.PieceCid,
		PieceSize:     uint64(pieceSize),
		CarSize:       uint64(dealConfig.FileSize),
		PayloadCid:    dealConfig.PayloadCid,
		StartEpoch:    int(dealConfig.StartEpoch),
		Duration:      dealConfig.Duration,
		StoragePrice:  int(epochPrice.BigInt().Int64()),
		Verified:      dealConfig.VerifiedDeal,
		FastRetrieval: dealConfig.FastRetrieval,
		Wallet:        dealConfig.SenderWallet,
	}

	dealUuid, err := client.sendDealToMiner(dealParam)
	if err != nil {
		logs.GetLogger().Error(err)
		return "", err
	}
	return dealUuid, nil
}

func (client *Client) sendDealToMiner(dealP DealParam) (string, error) {
	ctx := context.Background()
	n, err := clinode.Setup(client.ClientRepo)
	if err != nil {
		return "", err
	}

	ainfo := cliutil.ParseApiInfo(client.FullNodeApi)
	addr, err := ainfo.DialArgs("v1")
	if err != nil {
		logs.GetLogger().Error("parse fullNodeApi failed: %w", err)
		return "", err
	}

	fullNode, closer, err := apiclient.NewFullNodeRPCV1(context.Background(), addr, ainfo.AuthHeader())
	if err != nil {
		return "", fmt.Errorf("cant setup fullnode connection: %w", err)
	}
	defer closer()

	walletAddr, err := n.GetProvidedOrDefaultWallet(ctx, dealP.Wallet)
	if err != nil {
		return "", err
	}

	logs.GetLogger().Warn("selected wallet: ", walletAddr)

	maddr, err := address.NewFromString(dealP.Provider)
	if err != nil {
		return "", err
	}

	addrInfo, err := cmd.GetAddrInfo(ctx, fullNode, maddr)
	if err != nil {
		return "", err
	}

	logs.GetLogger().Warn("found storage provider ", "id:", addrInfo.ID, ", multiaddrs: ", addrInfo.Addrs, ", minerID:", maddr)

	if err := n.Host.Connect(ctx, *addrInfo); err != nil {
		return "", fmt.Errorf("failed to connect to peer %s: %w", addrInfo.ID, err)
	}

	x, err := n.Host.Peerstore().FirstSupportedProtocol(addrInfo.ID, DealProtocolv120)
	if err != nil {
		return "", fmt.Errorf("getting protocols for peer %s: %w", addrInfo.ID, err)
	}

	if len(x) == 0 {
		return "", fmt.Errorf("boost client cannot make a deal with storage provider %s because it does not support protocol version 1.2.0", maddr)
	}

	dealUuid := uuid.New()

	commp := dealP.Commp
	pieceCid, err := cid.Parse(commp)
	if err != nil {
		return "", fmt.Errorf("parsing commp '%s': %w", commp, err)
	}

	pieceSize := dealP.PieceSize
	if pieceSize == 0 {
		return "", fmt.Errorf("must provide piece-size parameter for CAR url")
	}

	payloadCidStr := dealP.PayloadCid
	rootCid, err := cid.Parse(payloadCidStr)
	if err != nil {
		return "", fmt.Errorf("dealUuid: %s, parsing payload cid %s: %w", dealUuid.String(), payloadCidStr, err)
	}

	carFileSize := dealP.CarSize
	if dealP.CarSize == 0 {
		return "", fmt.Errorf("size of car file cannot be 0")
	}

	transfer := types.Transfer{
		Size: carFileSize,
	}

	var providerCollateral abi.TokenAmount
	if dealP.ProviderCollateral != 0 {
		providerCollateral = abi.NewTokenAmount(int64(dealP.ProviderCollateral))
	} else {
		bounds, err := fullNode.StateDealProviderCollateralBounds(ctx, abi.PaddedPieceSize(pieceSize), dealP.Verified, chaintypes.EmptyTSK)
		if err != nil {
			return "", fmt.Errorf("dealUuid: %s, node error getting collateral bounds: %w", dealUuid.String(), err)
		}

		providerCollateral = big.Div(big.Mul(bounds.Min, big.NewInt(6)), big.NewInt(5)) // add 20%
	}

	var startEpoch abi.ChainEpoch
	if dealP.StartEpoch != 0 {
		startEpoch = abi.ChainEpoch(dealP.StartEpoch)
	} else {
		tipset, err := fullNode.ChainHead(ctx)
		if err != nil {
			return "", fmt.Errorf("getting chain head: %w", err)
		}

		head := tipset.Height()
		logs.GetLogger().Debug("current block height", "number", head)
		startEpoch = head + abi.ChainEpoch(5760) // head + 2 days
	}

	// Create a deal proposal to storage provider using deal protocol v1.2.0 format
	dealProposal, err := dealProposal(ctx, n, walletAddr, rootCid, abi.PaddedPieceSize(pieceSize), pieceCid, maddr, startEpoch, dealP.Duration, dealP.Verified, providerCollateral, abi.NewTokenAmount(int64(dealP.StoragePrice)))
	if err != nil {
		return "", fmt.Errorf("dealUuid: %s, failed to create a deal proposal: %w", dealUuid.String(), err)
	}

	dealParams := types.DealParams{
		DealUUID:           dealUuid,
		ClientDealProposal: *dealProposal,
		DealDataRoot:       rootCid,
		IsOffline:          true,
		Transfer:           transfer,
	}

	logs.GetLogger().Debug("about to submit deal proposal", "uuid", dealUuid.String())

	s, err := n.Host.NewStream(ctx, addrInfo.ID, DealProtocolv120)
	if err != nil {
		return "", fmt.Errorf("failed to open stream to peer %s: %w", addrInfo.ID, err)
	}
	defer s.Close()

	var resp types.DealResponse
	if err := doRpc(ctx, s, &dealParams, &resp); err != nil {
		return "", fmt.Errorf("send proposal rpc: %w", err)
	}

	if !resp.Accepted {
		return "", fmt.Errorf("deal proposal rejected: %s", resp.Message)
	}

	logs.GetLogger().Infof("dealUuid: %s, The deal proposal has been sent to the storage provider, the deal info is as follows: ", dealUuid.String())
	out := map[string]interface{}{
		"dealUuid":           dealUuid.String(),
		"provider":           maddr.String(),
		"clientWallet":       walletAddr.String(),
		"payloadCid":         rootCid.String(),
		"commp":              dealProposal.Proposal.PieceCID.String(),
		"startEpoch":         dealProposal.Proposal.StartEpoch.String(),
		"endEpoch":           dealProposal.Proposal.EndEpoch.String(),
		"providerCollateral": dealProposal.Proposal.ProviderCollateral.String(),
	}
	return dealUuid.String(), cmd.PrintJson(out)
}

func dealProposal(ctx context.Context, n *clinode.Node, clientAddr address.Address, rootCid cid.Cid, pieceSize abi.PaddedPieceSize, pieceCid cid.Cid, minerAddr address.Address, startEpoch abi.ChainEpoch, duration int, verified bool, providerCollateral abi.TokenAmount, storagePrice abi.TokenAmount) (*market.ClientDealProposal, error) {
	endEpoch := startEpoch + abi.ChainEpoch(duration)
	// deal proposal expects total storage price for deal per epoch, therefore we
	// multiply pieceSize * storagePrice (which is set per epoch per GiB) and divide by 2^30
	storagePricePerEpochForDeal := big.Div(big.Mul(big.NewInt(int64(pieceSize)), storagePrice), big.NewInt(int64(1<<30)))
	l, err := market.NewLabelFromString(rootCid.String())
	if err != nil {
		return nil, err
	}
	proposal := market.DealProposal{
		PieceCID:             pieceCid,
		PieceSize:            pieceSize,
		VerifiedDeal:         verified,
		Client:               clientAddr,
		Provider:             minerAddr,
		Label:                l,
		StartEpoch:           startEpoch,
		EndEpoch:             endEpoch,
		StoragePricePerEpoch: storagePricePerEpochForDeal,
		ProviderCollateral:   providerCollateral,
	}

	buf, err := cborutil.Dump(&proposal)
	if err != nil {
		return nil, err
	}

	sig, err := n.Wallet.WalletSign(ctx, clientAddr, buf, api.MsgMeta{Type: api.MTDealProposal})
	if err != nil {
		return nil, fmt.Errorf("wallet sign failed: %w", err)
	}

	return &market.ClientDealProposal{
		Proposal:        proposal,
		ClientSignature: *sig,
	}, nil
}

func doRpc(ctx context.Context, s inet.Stream, req interface{}, resp interface{}) error {
	errc := make(chan error)
	go func() {
		if err := cborutil.WriteCborRPC(s, req); err != nil {
			errc <- fmt.Errorf("failed to send request: %w", err)
			return
		}

		if err := cborutil.ReadCborRPC(s, resp); err != nil {
			errc <- fmt.Errorf("failed to read response: %w", err)
			return
		}

		errc <- nil
	}()

	select {
	case err := <-errc:
		return err
	case <-ctx.Done():
		return ctx.Err()
	}
}

type DealParam struct {
	Provider           string `json:"provider"`            // storage provider on-chain address. Required
	Commp              string `json:"commp"`               // commp of the CAR file. Required
	PieceSize          uint64 `json:"piece_size"`          // size of the CAR file as a padded piece. Required
	CarSize            uint64 `json:"car_size"`            // size of the CAR file. Required
	PayloadCid         string `json:"payload_cid"`         // root CID of the CAR file. Required
	StartEpoch         int    `json:"start_epoch"`         // start epoch by when the deal should be proved by provider on-chain. default: current chain head + 2 days
	Duration           int    `json:"duration"`            // duration of the deal in epochs. default is 2880 * 180 == 180 days  518400
	ProviderCollateral int    `json:"provider_collateral"` // deal collateral that storage miner must put in escrow; if empty, the min collateral for the given piece size will be used
	StoragePrice       int    `json:"storage_price"`       // storage price in attoFIL per epoch per GiB. default 1
	Verified           bool   `json:"verified"`            // whether the deal funds should come from verified client data-cap. default true
	FastRetrieval      bool   `json:"fast_retrieval"`      // indicates that data should be available for fast retrieval. default true
	Wallet             string `json:"wallet"`              // wallet address to be used to initiate the deal
}
