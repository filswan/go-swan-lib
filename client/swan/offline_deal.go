package swan

import (
	"encoding/json"
	"fmt"
	"net/url"
	"strconv"
	"strings"

	"github.com/filswan/go-swan-lib/client"
	"github.com/filswan/go-swan-lib/constants"
	"github.com/filswan/go-swan-lib/logs"
	"github.com/filswan/go-swan-lib/model"
)

const GET_OFFLINEDEAL_LIMIT_DEFAULT = 50

type GetOfflineDealResponse struct {
	Data   GetOfflineDealData `json:"data"`
	Status string             `json:"status"`
}

type GetOfflineDealData struct {
	Deal []model.OfflineDeal `json:"deal"`
}

type UpdateOfflineDealResponse struct {
	Data   UpdateOfflineDealData `json:"data"`
	Status string                `json:"status"`
}

type UpdateOfflineDealData struct {
	Deal    model.OfflineDeal `json:"deal"`
	Message string            `json:"message"`
}

func (swanClient *SwanClient) SwanGetOfflineDeals(minerFid, status string, limit ...string) []model.OfflineDeal {
	err := swanClient.SwanGetJwtTokenUp3Times()
	if err != nil {
		logs.GetLogger().Error(err)
		return nil
	}

	rowLimit := strconv.Itoa(GET_OFFLINEDEAL_LIMIT_DEFAULT)
	if len(limit) > 0 {
		rowLimit = limit[0]
	}

	urlStr := swanClient.ApiUrl + "/offline_deals/" + minerFid + "?deal_status=" + status + "&limit=" + rowLimit + "&offset=0"
	response := client.HttpGet(urlStr, swanClient.SwanToken, "")
	getOfflineDealResponse := GetOfflineDealResponse{}
	err = json.Unmarshal([]byte(response), &getOfflineDealResponse)
	if err != nil {
		logs.GetLogger().Error(err)
		return nil
	}

	if !strings.EqualFold(getOfflineDealResponse.Status, constants.SWAN_API_STATUS_SUCCESS) {
		logs.GetLogger().Error("Get offline deal with status ", status, " failed")
		return nil
	}

	return getOfflineDealResponse.Data.Deal
}

func (swanClient *SwanClient) SwanUpdateOfflineDealStatus(dealId int, status string, statusInfo ...string) bool {
	err := swanClient.SwanGetJwtTokenUp3Times()
	if err != nil {
		logs.GetLogger().Error(err)
		return false
	}

	if len(status) == 0 {
		logs.GetLogger().Error("Please provide status")
		return false
	}

	apiUrl := swanClient.ApiUrl + "/my_miner/deals/" + strconv.Itoa(dealId)

	params := url.Values{}
	params.Add("status", status)

	if len(statusInfo) > 0 && len(statusInfo[0]) > 0 {
		params.Add("note", statusInfo[0])
	}

	if len(statusInfo) > 1 && len(statusInfo[1]) > 0 {
		params.Add("file_path", statusInfo[1])
	}

	if len(statusInfo) > 2 && len(statusInfo[2]) > 0 {
		params.Add("file_size", statusInfo[2])
	}

	if len(statusInfo) > 3 && len(statusInfo[3]) > 0 {
		params.Add("cost", statusInfo[3])
	}

	response := client.HttpPut(apiUrl, swanClient.SwanToken, strings.NewReader(params.Encode()))

	updateOfflineDealResponse := &UpdateOfflineDealResponse{}
	err = json.Unmarshal([]byte(response), updateOfflineDealResponse)
	if err != nil {
		logs.GetLogger().Error(err)
		return false
	}

	if !strings.EqualFold(updateOfflineDealResponse.Status, constants.SWAN_API_STATUS_SUCCESS) {
		err := fmt.Errorf("deal(id=%d),failed to update offline deal status to %s,%s", dealId, status, updateOfflineDealResponse.Data.Message)
		logs.GetLogger().Error(err)
		return false
	}

	return true
}
