package lotus

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/filswan/go-swan-lib/client"
	"github.com/filswan/go-swan-lib/logs"
	"github.com/filswan/go-swan-lib/utils"
)

const (
	LOTUS_MARKET_GET_ASK               = "Filecoin.MarketGetAsk"
	LOTUS_MARKET_IMPORT_DATA           = "Filecoin.MarketImportDealData"
	LOTUS_MARKET_LIST_INCOMPLETE_DEALS = "Filecoin.MarketListIncompleteDeals"
)

type LotusMarket struct {
	ApiUrl       string
	AccessToken  string
	ClientApiUrl string
}

type MarketGetAsk struct {
	LotusJsonRpcResult
	Result *MarketGetAskResult `json:"result"`
}

type MarketGetAskResult struct {
	Ask MarketGetAskResultAsk
}
type MarketGetAskResultAsk struct {
	Price         string
	VerifiedPrice string
	MinPieceSize  int
	MaxPieceSize  int
	Miner         string
	Timestamp     int
	Expiry        int
	SeqNo         int
}

func GetLotusMarket(apiUrl, accessToken, clientApiUrl string) (*LotusMarket, error) {
	if len(apiUrl) == 0 {
		err := fmt.Errorf("lotus market api url is required")
		logs.GetLogger().Error(err)
		return nil, err
	}

	lotusMarket := &LotusMarket{
		ApiUrl:       apiUrl,
		AccessToken:  accessToken,
		ClientApiUrl: clientApiUrl,
	}

	return lotusMarket, nil
}

//"lotus client query-ask " + minerFid
func (lotusMarket *LotusMarket) LotusMarketGetAsk() *MarketGetAskResultAsk {
	var params []interface{}

	jsonRpcParams := LotusJsonRpcParams{
		JsonRpc: LOTUS_JSON_RPC_VERSION,
		Method:  LOTUS_MARKET_GET_ASK,
		Params:  params,
		Id:      LOTUS_JSON_RPC_ID,
	}

	//here the api url should be miner's api url, need to change later on
	response := client.HttpGetNoToken(lotusMarket.ApiUrl, jsonRpcParams)
	if response == "" {
		return nil
	}

	marketGetAsk := &MarketGetAsk{}
	err := json.Unmarshal([]byte(response), marketGetAsk)
	if err != nil {
		logs.GetLogger().Error(err)
		return nil
	}

	if marketGetAsk.Result == nil {
		return nil
	}

	return &marketGetAsk.Result.Ask
}

type DealCid struct {
	DealCid string `json:"/"`
}

type MarketListIncompleteDeals struct {
	Id      int           `json:"id"`
	JsonRpc string        `json:"jsonrpc"`
	Result  []Deal        `json:"result"`
	Error   *JsonRpcError `json:"error"`
}

type Deal struct {
	State       int     `json:"State"`
	Message     string  `json:"Message"`
	ProposalCid DealCid `json:"ProposalCid"`
}

func (lotusMarket *LotusMarket) LotusGetDeals() []Deal {
	var params []interface{}
	jsonRpcParams := LotusJsonRpcParams{
		JsonRpc: LOTUS_JSON_RPC_VERSION,
		Method:  LOTUS_MARKET_LIST_INCOMPLETE_DEALS,
		Params:  params,
		Id:      LOTUS_JSON_RPC_ID,
	}

	logs.GetLogger().Info("Get deal list from ", lotusMarket.ApiUrl)
	response := client.HttpGet(lotusMarket.ApiUrl, lotusMarket.AccessToken, jsonRpcParams)
	logs.GetLogger().Info("Got deal list from ", lotusMarket.ApiUrl)
	deals := &MarketListIncompleteDeals{}
	err := json.Unmarshal([]byte(response), deals)
	if err != nil {
		logs.GetLogger().Error(err)
		return nil
	}

	return deals.Result
}

func (lotusMarket *LotusMarket) LotusGetDealOnChainStatusFromDeals(deals []Deal, dealCid string) (string, string) {
	if len(deals) == 0 {
		logs.GetLogger().Error("Deal list is empty.")
		return "", ""
	}

	lotusClient, err := LotusGetClient(lotusMarket.ClientApiUrl, "")
	if err != nil {
		logs.GetLogger().Error(err)
		return "", ""
	}
	for _, deal := range deals {
		if deal.ProposalCid.DealCid != dealCid {
			continue
		}

		status := lotusClient.LotusGetDealStatus(deal.State)
		msg := fmt.Sprintf("deal:%s,%s", dealCid, status)
		if deal.Message != "" {
			msg = msg + "," + deal.Message
		}
		logs.GetLogger().Info(msg)
		return status, deal.Message
	}

	logs.GetLogger().Error("Did not find your deal:", dealCid, " in the returned list.")

	return "", ""
}

//"lotus-miner storage-deals list -v | grep -a " + dealCid
func (lotusMarket *LotusMarket) LotusGetDealOnChainStatus(dealCid string) (string, string) {
	deals := lotusMarket.LotusGetDeals()
	status, message := lotusMarket.LotusGetDealOnChainStatusFromDeals(deals, dealCid)
	return status, message
}

func (lotusMarket *LotusMarket) LotusImportData(dealCid string, filepath string) error {
	var params []interface{}
	getDealInfoParam := DealCid{DealCid: dealCid}
	params = append(params, getDealInfoParam)
	params = append(params, filepath)

	jsonRpcParams := LotusJsonRpcParams{
		JsonRpc: LOTUS_JSON_RPC_VERSION,
		Method:  LOTUS_MARKET_IMPORT_DATA,
		Params:  params,
		Id:      LOTUS_JSON_RPC_ID,
	}

	response := client.HttpPost(lotusMarket.ApiUrl, lotusMarket.AccessToken, jsonRpcParams)
	if response == "" {
		err := fmt.Errorf("no response, please check your market api url:%s and access token", lotusMarket.ApiUrl)
		logs.GetLogger().Error(err)
		return err
	}
	//logs.GetLogger().Info(response)

	errorInfo := utils.GetFieldMapFromJson(response, "error")

	if errorInfo == nil {
		return nil
	}

	//logs.GetLogger().Error(errorInfo)
	errCode := int(errorInfo["code"].(float64))
	errMsg := errorInfo["message"].(string)
	err := fmt.Errorf("error code:%d message:%s", errCode, errMsg)
	//logs.GetLogger().Error(err)
	if strings.Contains(response, "(need 'write')") {
		logs.GetLogger().Error("please check your access token, it should have write access")
		logs.GetLogger().Error(err)
	}
	return err
}
