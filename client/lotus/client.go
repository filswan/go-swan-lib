package lotus

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/filswan/go-swan-lib/client/web"
	"github.com/filswan/go-swan-lib/constants"
	"github.com/filswan/go-swan-lib/logs"
	"github.com/filswan/go-swan-lib/model"
	"github.com/filswan/go-swan-lib/utils"

	"github.com/shopspring/decimal"
)

const (
	LOTUS_CLIENT_MINER_QUERY     = "Filecoin.ClientMinerQueryOffer"
	LOTUS_CLIENT_QUERY_ASK       = "Filecoin.ClientQueryAsk"
	LOTUS_CLIENT_GET_DEAL_INFO   = "Filecoin.ClientGetDealInfo"
	LOTUS_CLIENT_GET_DEAL_STATUS = "Filecoin.ClientGetDealStatus"
	LOTUS_CHAIN_HEAD             = "Filecoin.ChainHead"
	LOTUS_CLIENT_CALC_COMM_P     = "Filecoin.ClientCalcCommP"
	LOTUS_CLIENT_IMPORT          = "Filecoin.ClientImport"
	LOTUS_CLIENT_GEN_CAR         = "Filecoin.ClientGenCar"
	LOTUS_CLIENT_START_DEAL      = "Filecoin.ClientStartDeal"

	STAGE_RESERVE_FUNDS     = "StorageDealReserveClientFunds"
	STAGE_PROPOSAL_ACCEPTED = "StorageDealProposalAccepted"

	FUNDS_RESERVED = "funds reserved"
	FUNDS_RELEASED = "funds released"
)

type LotusClient struct {
	ApiUrl      string
	AccessToken string
}

type ClientCalcCommP struct {
	LotusJsonRpcResult
	Result *ClientCalcCommPResult `json:"result"`
}

type ClientCalcCommPResult struct {
	Root Cid
	Size int
}
type ClientImport struct {
	LotusJsonRpcResult
	Result *ClientImportResult `json:"result"`
}
type ClientImportResult struct {
	Root     Cid
	ImportID int64
}

func LotusGetClient(apiUrl, accessToken string) (*LotusClient, error) {
	if len(apiUrl) == 0 {
		err := fmt.Errorf("config lotus api_url is required")
		logs.GetLogger().Error(err)
		return nil, err
	}

	lotusClient := &LotusClient{
		ApiUrl:      apiUrl,
		AccessToken: accessToken,
	}

	return lotusClient, nil
}

type ClientMinerQuery struct {
	LotusJsonRpcResult
	Result ClientMinerQueryResult `json:"result"`
}

type ClientMinerQueryResult struct {
	MinerPeer ClientMinerQueryResultPeer
}

type ClientMinerQueryResultPeer struct {
	Address string
	ID      string
}

type ClientDealInfo struct {
	LotusJsonRpcResult
	Result ClientDealResult `json:"result"`
}

type ClientDealResult struct {
	State         int
	Message       string
	DealStages    ClientDealStages
	PricePerEpoch string
	Duration      int
	DealID        int64
	Verified      bool
}

type ClientDealStages struct {
	Stages []ClientDealStage
}
type ClientDealStage struct {
	Name             string
	Description      string
	ExpectedDuration string
	CreatedTime      string
	UpdatedTime      string
	Logs             []ClientDealStageLog
}

type ClientDealStageLog struct {
	Log         string
	UpdatedTime string
}

type ClientDealCostStatus struct {
	CostComputed         string
	ReserveClientFunds   string
	DealProposalAccepted string
	Status               string
	DealId               int64
	Verified             bool
}

func (lotusClient *LotusClient) LotusClientGetDealInfo(dealCid string) (*ClientDealCostStatus, error) {
	var params []interface{}
	cid := Cid{Cid: dealCid}
	params = append(params, cid)

	jsonRpcParams := LotusJsonRpcParams{
		JsonRpc: LOTUS_JSON_RPC_VERSION,
		Method:  LOTUS_CLIENT_GET_DEAL_INFO,
		Params:  params,
		Id:      LOTUS_JSON_RPC_ID,
	}

	response := web.HttpGetNoToken(lotusClient.ApiUrl, jsonRpcParams)

	clientDealInfo := &ClientDealInfo{}
	err := json.Unmarshal([]byte(response), clientDealInfo)
	if err != nil {
		err := fmt.Errorf("deal:%s,%s", dealCid, err.Error())
		logs.GetLogger().Error(err)
		return nil, err
	}

	if clientDealInfo.Error != nil {
		err := fmt.Errorf("deal:%s,code:%d,message:%s", dealCid, clientDealInfo.Error.Code, clientDealInfo.Error.Message)
		logs.GetLogger().Error(err)
		return nil, err
	}

	pricePerEpoch, err := decimal.NewFromString(clientDealInfo.Result.PricePerEpoch)
	if err != nil {
		err := fmt.Errorf("deal:%s,%s", dealCid, err.Error())
		logs.GetLogger().Error(err)
		return nil, err
	}
	duration := decimal.NewFromInt(int64(clientDealInfo.Result.Duration))

	clientDealCostStatus := ClientDealCostStatus{}
	clientDealCostStatus.CostComputed = pricePerEpoch.Mul(duration).String()

	dealStages := clientDealInfo.Result.DealStages.Stages
	for _, stage := range dealStages {
		if strings.EqualFold(stage.Name, STAGE_RESERVE_FUNDS) {
			for _, log := range stage.Logs {
				if strings.Contains(log.Log, FUNDS_RESERVED) {
					clientDealCostStatus.ReserveClientFunds = utils.GetNumStrFromStr(log.Log)
					clientDealCostStatus.ReserveClientFunds = strings.TrimSuffix(clientDealCostStatus.ReserveClientFunds, ">")
				}
			}
		}
		if strings.EqualFold(stage.Name, STAGE_PROPOSAL_ACCEPTED) {
			for _, log := range stage.Logs {
				if strings.Contains(log.Log, FUNDS_RELEASED) {
					clientDealCostStatus.DealProposalAccepted = utils.GetNumStrFromStr(log.Log)
					clientDealCostStatus.DealProposalAccepted = strings.TrimSuffix(clientDealCostStatus.DealProposalAccepted, ">")
				}
			}
		}
	}

	clientDealCostStatus.Status = lotusClient.LotusGetDealStatus(clientDealInfo.Result.State)
	clientDealCostStatus.DealId = clientDealInfo.Result.DealID
	clientDealCostStatus.Verified = clientDealInfo.Result.Verified

	//logs.GetLogger().Info(clientDealCost)
	return &clientDealCostStatus, nil
}

func GetDealCost(dealCost ClientDealCostStatus) string {
	if dealCost.DealProposalAccepted != "" {
		return dealCost.DealProposalAccepted
	}

	if dealCost.ReserveClientFunds != "" {
		return dealCost.ReserveClientFunds
	}

	return dealCost.CostComputed
}

func (lotusClient *LotusClient) LotusClientMinerQuery(minerFid string) (string, error) {
	var params []interface{}
	params = append(params, minerFid)
	params = append(params, nil)
	params = append(params, nil)

	jsonRpcParams := LotusJsonRpcParams{
		JsonRpc: LOTUS_JSON_RPC_VERSION,
		Method:  LOTUS_CLIENT_MINER_QUERY,
		Params:  params,
		Id:      LOTUS_JSON_RPC_ID,
	}

	response := web.HttpGetNoToken(lotusClient.ApiUrl, jsonRpcParams)

	clientMinerQuery := &ClientMinerQuery{}
	err := json.Unmarshal([]byte(response), clientMinerQuery)
	if err != nil {
		err := fmt.Errorf("miner:%s,%s", minerFid, err.Error())
		logs.GetLogger().Error(err)
		return "", err
	}

	if clientMinerQuery.Error != nil {
		err := fmt.Errorf("miner:%s,code:%d,message:%s", minerFid, clientMinerQuery.Error.Code, clientMinerQuery.Error.Message)
		logs.GetLogger().Error(err)
		return "", err
	}

	minerPeerId := clientMinerQuery.Result.MinerPeer.ID
	return minerPeerId, nil
}

type ClientQueryAsk struct {
	LotusJsonRpcResult
	Result ClientQueryAskResult `json:"result"`
}

type ClientQueryAskResult struct {
	Price         string
	VerifiedPrice string
	MinPieceSize  int64
	MaxPieceSize  int64
}
type MinerConfig struct {
	Price         decimal.Decimal
	VerifiedPrice decimal.Decimal
	MinPieceSize  int64
	MaxPieceSize  int64
}

func (lotusClient *LotusClient) LotusClientQueryAsk(minerFid string) (*MinerConfig, error) {
	minerPeerId, err := lotusClient.LotusClientMinerQuery(minerFid)
	if err != nil {
		logs.GetLogger().Error(err)
		return nil, err
	}

	var params []interface{}
	params = append(params, minerPeerId)
	params = append(params, minerFid)

	jsonRpcParams := LotusJsonRpcParams{
		JsonRpc: LOTUS_JSON_RPC_VERSION,
		Method:  LOTUS_CLIENT_QUERY_ASK,
		Params:  params,
		Id:      LOTUS_JSON_RPC_ID,
	}

	response := web.HttpGetNoToken(lotusClient.ApiUrl, jsonRpcParams)

	clientQueryAsk := &ClientQueryAsk{}
	err = json.Unmarshal([]byte(response), clientQueryAsk)
	if err != nil {
		logs.GetLogger().Error(err)
		return nil, err
	}

	if clientQueryAsk.Error != nil {
		err := fmt.Errorf("miner:%s,code:%d,message:%s", minerFid, clientQueryAsk.Error.Code, clientQueryAsk.Error.Message)
		logs.GetLogger().Error(err)
		return nil, err
	}

	price, err := decimal.NewFromString(clientQueryAsk.Result.Price)
	if err != nil {
		logs.GetLogger().Error(err)
		return nil, err
	}

	verifiedPrice, err := decimal.NewFromString(clientQueryAsk.Result.VerifiedPrice)
	if err != nil {
		logs.GetLogger().Error(err)
		return nil, err
	}

	minerConfig := &MinerConfig{
		Price:         price,
		VerifiedPrice: verifiedPrice,
		MinPieceSize:  clientQueryAsk.Result.MinPieceSize,
		MaxPieceSize:  clientQueryAsk.Result.MaxPieceSize,
	}

	return minerConfig, nil
}

func (lotusClient *LotusClient) LotusGetCurrentEpoch() int {
	var params []interface{}

	jsonRpcParams := LotusJsonRpcParams{
		JsonRpc: LOTUS_JSON_RPC_VERSION,
		Method:  LOTUS_CHAIN_HEAD,
		Params:  params,
		Id:      LOTUS_JSON_RPC_ID,
	}

	response := web.HttpPostNoToken(lotusClient.ApiUrl, jsonRpcParams)

	result := utils.GetFieldMapFromJson(response, "result")
	if result == nil {
		logs.GetLogger().Error("Failed to get result from:", lotusClient.ApiUrl)
		return -1
	}

	height := result["Height"]
	if height == nil {
		logs.GetLogger().Error("Failed to get height from:", lotusClient.ApiUrl)
		return -1
	}

	heightFloat := height.(float64)
	return int(heightFloat)
}

//"lotus-miner storage-deals list -v | grep -a " + dealCid
func (lotusClient *LotusClient) LotusGetDealStatus(state int) string {
	var params []interface{}
	params = append(params, state)

	jsonRpcParams := LotusJsonRpcParams{
		JsonRpc: LOTUS_JSON_RPC_VERSION,
		Method:  LOTUS_CLIENT_GET_DEAL_STATUS,
		Params:  params,
		Id:      LOTUS_JSON_RPC_ID,
	}

	response := web.HttpPostNoToken(lotusClient.ApiUrl, jsonRpcParams)

	result := utils.GetFieldStrFromJson(response, "result")
	if result == "" {
		logs.GetLogger().Error("no response from:", lotusClient.ApiUrl)
		return ""
	}

	return result
}

//"lotus client commP " + carFilePath
func (lotusClient *LotusClient) LotusClientCalcCommP(filepath string) *string {
	var params []interface{}
	params = append(params, filepath)

	jsonRpcParams := LotusJsonRpcParams{
		JsonRpc: LOTUS_JSON_RPC_VERSION,
		Method:  LOTUS_CLIENT_CALC_COMM_P,
		Params:  params,
		Id:      LOTUS_JSON_RPC_ID,
	}

	response := web.HttpPost(lotusClient.ApiUrl, lotusClient.AccessToken, jsonRpcParams)
	if response == "" {
		logs.GetLogger().Error("no response from:", lotusClient.ApiUrl)
		return nil
	}

	//logs.GetLogger().Info(response)

	clientCalcCommP := &ClientCalcCommP{}
	err := json.Unmarshal([]byte(response), clientCalcCommP)
	if err != nil {
		logs.GetLogger().Error(err)
		return nil
	}

	if clientCalcCommP.Error != nil {
		err := fmt.Errorf("get piece CID failed for:%s, error code:%d, message:%s", filepath, clientCalcCommP.Error.Code, clientCalcCommP.Error.Message)
		logs.GetLogger().Error(err)
		return nil
	}

	if clientCalcCommP.Result == nil {
		logs.GetLogger().Error("no result from:", lotusClient.ApiUrl)
		return nil
	}

	pieceCid := clientCalcCommP.Result.Root.Cid
	return &pieceCid
}

type ClientFileParam struct {
	Path  string
	IsCAR bool
}

//"lotus client import --car " + carFilePath
func (lotusClient *LotusClient) LotusClientImport(filepath string, isCar bool) (*string, error) {
	var params []interface{}
	clientFileParam := ClientFileParam{
		Path:  filepath,
		IsCAR: isCar,
	}
	params = append(params, clientFileParam)

	jsonRpcParams := LotusJsonRpcParams{
		JsonRpc: LOTUS_JSON_RPC_VERSION,
		Method:  LOTUS_CLIENT_IMPORT,
		Params:  params,
		Id:      LOTUS_JSON_RPC_ID,
	}

	response := web.HttpGet(lotusClient.ApiUrl, lotusClient.AccessToken, jsonRpcParams)
	if response == "" {
		err := fmt.Errorf("lotus import file %s failed, no response from %s", filepath, lotusClient.ApiUrl)
		logs.GetLogger().Error(err)
		return nil, err
	}

	clientImport := &ClientImport{}
	err := json.Unmarshal([]byte(response), clientImport)
	if err != nil {
		logs.GetLogger().Error(err)
		return nil, err
	}

	if clientImport.Error != nil {
		err := fmt.Errorf("lotus import file %s failed, error code:%d, message:%s", filepath, clientImport.Error.Code, clientImport.Error.Message)
		logs.GetLogger().Error(err)
		return nil, err
	}

	if clientImport.Result == nil {
		err := fmt.Errorf("lotus import file %s failed, result is null from %s", filepath, lotusClient.ApiUrl)
		logs.GetLogger().Error(err)
		return nil, err
	}

	dataCid := clientImport.Result.Root.Cid

	return &dataCid, nil
}

//"lotus client generate-car " + srcFilePath + " " + destCarFilePath
func (lotusClient *LotusClient) LotusClientGenCar(srcFilePath, destCarFilePath string, srcFilePathIsCar bool) error {
	var params []interface{}
	clientFileParam := ClientFileParam{
		Path:  srcFilePath,
		IsCAR: srcFilePathIsCar,
	}
	params = append(params, clientFileParam)
	params = append(params, destCarFilePath)

	jsonRpcParams := LotusJsonRpcParams{
		JsonRpc: LOTUS_JSON_RPC_VERSION,
		Method:  LOTUS_CLIENT_GEN_CAR,
		Params:  params,
		Id:      LOTUS_JSON_RPC_ID,
	}

	response := web.HttpGet(lotusClient.ApiUrl, lotusClient.AccessToken, jsonRpcParams)
	if response == "" {
		err := fmt.Errorf("failed to generate car, no response")
		logs.GetLogger().Error(err)
		return err
	}

	//logs.GetLogger().Info(response)
	lotusJsonRpcResult := &LotusJsonRpcResult{}
	err := json.Unmarshal([]byte(response), lotusJsonRpcResult)
	if err != nil {
		logs.GetLogger().Error(err)
		return err
	}

	if lotusJsonRpcResult.Error != nil {
		err := fmt.Errorf("error, code:%d, message:%s", lotusJsonRpcResult.Error.Code, lotusJsonRpcResult.Error.Message)
		logs.GetLogger().Error(err)
		return err
	}

	return nil
}

type ClientStartDealParamData struct {
	TransferType string
	Root         Cid
	PieceCid     Cid
	PieceSize    int
}

type ClientStartDealParam struct {
	Data              ClientStartDealParamData
	Wallet            string
	Miner             string
	EpochPrice        string
	MinBlocksDuration int
	DealStartEpoch    int
	FastRetrieval     bool
	VerifiedDeal      bool
}

type ClientStartDeal struct {
	LotusJsonRpcResult
	Result Cid `json:"result"`
}

func (lotusClient *LotusClient) LotusClientStartDeal(dealConfig model.DealConfig, relativeEpoch int) (*string, *int, error) {
	pieceSize, sectorSize := utils.CalculatePieceSize(dealConfig.FileSize)
	cost := utils.CalculateRealCost(sectorSize, dealConfig.MinerPrice)

	epochPrice := cost.Mul(decimal.NewFromFloat(constants.LOTUS_PRICE_MULTIPLE_1E18))
	startEpoch := dealConfig.StartEpoch - relativeEpoch

	if !dealConfig.SkipConfirmation {
		logs.GetLogger().Info("Do you confirm to submit the deal?")
		logs.GetLogger().Info("Press Y/y to continue, other key to quit")
		reader := bufio.NewReader(os.Stdin)
		response, err := reader.ReadString('\n')
		if err != nil {
			logs.GetLogger().Error(err)
			return nil, nil, err
		}

		response = strings.TrimRight(response, "\n")

		if !strings.EqualFold(response, "Y") {
			logs.GetLogger().Info("Your input is ", response, ". Now give up submit the deal.")
			return nil, nil, nil
		}
	}

	clientStartDealParamData := ClientStartDealParamData{
		TransferType: dealConfig.TransferType, //constants.LOTUS_TRANSFER_TYPE_MANUAL,
		Root: Cid{
			Cid: dealConfig.PayloadCid,
		},
		PieceCid: Cid{
			Cid: dealConfig.PieceCid,
		},
		PieceSize: int(pieceSize),
	}

	clientStartDealParam := ClientStartDealParam{
		Data:              clientStartDealParamData,
		Wallet:            dealConfig.SenderWallet,
		Miner:             dealConfig.MinerFid,
		EpochPrice:        epochPrice.BigInt().String(),
		MinBlocksDuration: dealConfig.Duration,
		DealStartEpoch:    startEpoch,
		FastRetrieval:     dealConfig.FastRetrieval,
		VerifiedDeal:      dealConfig.VerifiedDeal,
	}

	var params []interface{}
	params = append(params, clientStartDealParam)
	logs.GetLogger().Info(utils.ToJson(params))

	jsonRpcParams := LotusJsonRpcParams{
		JsonRpc: LOTUS_JSON_RPC_VERSION,
		Method:  LOTUS_CLIENT_START_DEAL,
		Params:  params,
		Id:      LOTUS_JSON_RPC_ID,
	}

	response := web.HttpGet(lotusClient.ApiUrl, lotusClient.AccessToken, jsonRpcParams)
	if response == "" {
		err := fmt.Errorf("failed to send deal for %s, no response", dealConfig.PayloadCid)
		logs.GetLogger().Error(err)
		return nil, nil, err
	}

	clientStartDeal := &ClientStartDeal{}
	err := json.Unmarshal([]byte(response), clientStartDeal)
	if err != nil {
		logs.GetLogger().Error(err)
		return nil, nil, err
	}

	if clientStartDeal.Error != nil {
		err := fmt.Errorf("error, code:%d, message:%s", clientStartDeal.Error.Code, clientStartDeal.Error.Message)
		logs.GetLogger().Error(err)
		return nil, nil, err
	}

	logs.GetLogger().Info("deal CID:", clientStartDeal.Result.Cid)
	return &clientStartDeal.Result.Cid, &startEpoch, nil
}
