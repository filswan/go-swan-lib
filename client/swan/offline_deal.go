package swan

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/filswan/go-swan-lib/client/web"
	"github.com/filswan/go-swan-lib/constants"
	"github.com/filswan/go-swan-lib/logs"
	"github.com/filswan/go-swan-lib/model"
	"github.com/filswan/go-swan-lib/utils"
)

const GET_OFFLINEDEAL_LIMIT_DEFAULT = 50

type GetOfflineDealsByStatusParams struct {
	DealStatus string  `json:"status"`
	ForMiner   bool    `json:"for_miner"`
	TaskUuid   *string `json:"task_uuid"`
	SourceId   *int    `json:"source_id"`
	MinerFid   *string `json:"miner_fid"`
	PageNum    *int    `json:"page_num"`
	PageSize   *int    `json:"page_size"`
}

type GetOfflineDealsByStatusResponse struct {
	Data struct {
		OfflineDeals []*model.OfflineDeal `json:"offline_deals"`
	} `json:"data"`
	Status string `json:"status"`
}

func (swanClient *SwanClient) GetOfflineDealsByStatus(params GetOfflineDealsByStatusParams) ([]*model.OfflineDeal, error) {
	if utils.IsStrEmpty(&params.DealStatus) {
		err := fmt.Errorf("deal status is required")
		logs.GetLogger().Error(err)
		return nil, err
	}

	err := swanClient.GetJwtTokenUp3Times()
	if err != nil {
		logs.GetLogger().Error(err)
		return nil, err
	}

	apiUrl := utils.UrlJoin(swanClient.ApiUrl, "offline_deals/get_by_status")
	response, err := web.HttpGet(apiUrl, swanClient.SwanToken, params)
	if err != nil {
		logs.GetLogger().Error(err)
		return nil, err
	}

	getOfflineDealsByStatusResponse := GetOfflineDealsByStatusResponse{}
	err = json.Unmarshal([]byte(response), &getOfflineDealsByStatusResponse)
	if err != nil {
		logs.GetLogger().Error(err)
		return nil, err
	}

	if !strings.EqualFold(getOfflineDealsByStatusResponse.Status, constants.SWAN_API_STATUS_SUCCESS) {
		err := fmt.Errorf("get offline deal with status:%s failed", params.DealStatus)
		logs.GetLogger().Error(err)
		return nil, err
	}

	return getOfflineDealsByStatusResponse.Data.OfflineDeals, nil
}

type UpdateOfflineDealParams struct {
	DealId       int     `json:"id"`
	DealCid      *string `json:"deal_cid"`
	FilePath     *string `json:"file_path"`
	Status       string  `json:"status"`
	StartEpoch   *int    `json:"start_epoch"`
	Note         *string `json:"note"`
	Cost         *string `json:"cost"`
	ChainDealId  int64   `json:"chain_deal_id"`
	ClientAddr   string  `json:"client_addr"`
	AllocationID uint64  `json:"allocation_id"`
}

// for public and auto-bid task
func (swanClient *SwanClient) UpdateOfflineDeal(params UpdateOfflineDealParams) error {
	err := swanClient.GetJwtTokenUp3Times()
	if err != nil {
		logs.GetLogger().Error(err)
		return err
	}

	if len(params.Status) == 0 {
		err := fmt.Errorf("status is invalid")
		logs.GetLogger().Error(err)
		return err
	}

	if params.DealId <= 0 {
		err := fmt.Errorf("deal id is invalid")
		logs.GetLogger().Error(err)
		return err
	}

	apiUrl := utils.UrlJoin(swanClient.ApiUrl, "offline_deals/update_offline_deal")

	response, err := web.HttpPut(apiUrl, swanClient.SwanToken, params)
	if err != nil {
		logs.GetLogger().Error(err)
		return err
	}

	swanServerResponse := &SwanServerResponse{}
	err = json.Unmarshal([]byte(response), swanServerResponse)
	if err != nil {
		logs.GetLogger().Error(err)
		return err
	}

	if !strings.EqualFold(swanServerResponse.Status, constants.SWAN_API_STATUS_SUCCESS) {
		err := fmt.Errorf("deal(id=%d),failed to update offline deal status to %s,%s", params.DealId, params.Status, swanServerResponse.Message)
		logs.GetLogger().Error(err)
		return err
	}

	return nil
}

// for public and non auto-bid task
func (swanClient *SwanClient) CreateOfflineDeals(fileDescs []*model.FileDesc) (*SwanServerResponse, error) {
	apiUrl := utils.UrlJoin(swanClient.ApiUrl, "offline_deals/create_offline_deals")
	response, err := web.HttpPost(apiUrl, swanClient.SwanToken, fileDescs)
	if err != nil {
		logs.GetLogger().Error(err)
		return nil, err
	}

	swanServerResponse := &SwanServerResponse{}
	err = json.Unmarshal(response, swanServerResponse)
	if err != nil {
		logs.GetLogger().Error(err)
		return nil, err
	}

	if !strings.EqualFold(swanServerResponse.Status, constants.SWAN_API_STATUS_SUCCESS) {
		err := fmt.Errorf("error:%s,%s", swanServerResponse.Status, swanServerResponse.Message)
		logs.GetLogger().Error(err)
		return nil, err
	}

	return swanServerResponse, nil
}

func (swanClient *SwanClient) GetDealListByTaskUuid(taskUuId string, pageNum int) (*DealListByTaskUuIdResp, error) {
	apiUrl := fmt.Sprintf("%s/tasks/%s?limit=100&offset=%d", swanClient.ApiUrl, taskUuId, 100*pageNum)
	response, err := web.HttpGet(apiUrl, swanClient.SwanToken, "")
	if err != nil {
		logs.GetLogger().Error(err)
		return nil, err
	}

	dealListByTaskUuIdResp := DealListByTaskUuIdResp{}
	err = json.Unmarshal([]byte(response), &dealListByTaskUuIdResp)
	if err != nil {
		logs.GetLogger().Error(err)
		return nil, err
	}

	if !strings.EqualFold(dealListByTaskUuIdResp.Status, constants.SWAN_API_STATUS_SUCCESS) {
		err := fmt.Errorf("get dealList by taskUuid status:%s failed", dealListByTaskUuIdResp.Status)
		logs.GetLogger().Error(err)
		return nil, err
	}
	return &dealListByTaskUuIdResp, nil
}

type DealListByTaskUuIdResp struct {
	Data struct {
		Deal []struct {
			DealCid     string `json:"deal_cid"`
			MinerFid    string `json:"miner_fid"`
			Note        string `json:"note"`
			PayloadCid  string `json:"payload_cid"`
			PieceCid    string `json:"piece_cid"`
			StartEpoch  int    `json:"start_epoch"`
			Status      string `json:"status"`
			ChainDealId int64  `json:"chain_deal_id"`
		} `json:"deal"`
		DealCompleteRate int `json:"deal_complete_rate"`
		Task             struct {
			BidMode        int    `json:"bid_mode"`
			CuratedDataset string `json:"curated_dataset"`
			Description    string `json:"description"`
			Duration       int    `json:"duration"`
			Status         string `json:"status"`
			TaskName       string `json:"task_name"`
			Type           string `json:"type"`
			Uuid           string `json:"uuid"`
		} `json:"task"`
		TotalDealCount int `json:"total_deal_count"`
		TotalItems     int `json:"total_items"`
	} `json:"data"`
	Status string `json:"status"`
}
