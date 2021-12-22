package swan

import (
	"encoding/json"
	"fmt"
	"net/url"
	"strconv"
	"strings"

	"github.com/filswan/go-swan-lib/client/web"
	"github.com/filswan/go-swan-lib/constants"
	"github.com/filswan/go-swan-lib/logs"
	"github.com/filswan/go-swan-lib/model"
	"github.com/filswan/go-swan-lib/utils"
)

type MinerResponse struct {
	Status  string      `json:"status"`
	Message string      `json:"message"`
	Data    model.Miner `json:"data"`
}

func (swanClient *SwanClient) GetMiner(minerFid string) (*MinerResponse, error) {
	apiUrl := swanClient.ApiUrl + "/miner/info/" + minerFid

	response := web.HttpGetNoToken(apiUrl, "")
	if response == "" {
		err := fmt.Errorf("no response from %s", apiUrl)
		logs.GetLogger().Error(err)
		return nil, err
	}

	minerResponse := &MinerResponse{}
	err := json.Unmarshal([]byte(response), minerResponse)
	if err != nil {
		logs.GetLogger().Error(err)
		return nil, err
	}

	if !strings.EqualFold(minerResponse.Status, constants.SWAN_API_STATUS_SUCCESS) {
		err := fmt.Errorf("status:%s, message:%s", minerResponse.Status, minerResponse.Message)
		logs.GetLogger().Error(err)
		return nil, err

	}

	return minerResponse, nil
}

func (swanClient *SwanClient) UpdateMinerBidConf(minerFid string, confMiner model.Miner) {
	err := swanClient.GetJwtTokenUp3Times()
	if err != nil {
		logs.GetLogger().Error(err)
		return
	}

	minerResponse, err := swanClient.GetMiner(minerFid)
	if err != nil {
		logs.GetLogger().Error(err)
		return
	}

	if minerResponse == nil || minerResponse.Status != constants.SWAN_API_STATUS_SUCCESS {
		logs.GetLogger().Error("Error: Get miner information failed")
		return
	}

	miner := minerResponse.Data

	if miner.BidMode == confMiner.BidMode &&
		miner.ExpectedSealingTime == confMiner.ExpectedSealingTime &&
		miner.StartEpoch == confMiner.StartEpoch &&
		miner.AutoBidTaskPerDay == confMiner.AutoBidTaskPerDay {
		logs.GetLogger().Info("No changes in bid configuration")
		return
	}

	logs.GetLogger().Info("Begin updating bid configuration")
	apiUrl := swanClient.ApiUrl + "/miner/info"

	params := url.Values{}
	params.Add("miner_fid", minerFid)
	params.Add("bid_mode", strconv.Itoa(confMiner.BidMode))
	params.Add("expected_sealing_time", strconv.Itoa(confMiner.ExpectedSealingTime))
	params.Add("start_epoch", strconv.Itoa(confMiner.StartEpoch))
	params.Add("auto_bid_task_per_day", strconv.Itoa(confMiner.AutoBidTaskPerDay))

	response := web.HttpPost(apiUrl, swanClient.SwanToken, strings.NewReader(params.Encode()))

	minerResponse = &MinerResponse{}
	err = json.Unmarshal([]byte(response), minerResponse)
	if err != nil {
		logs.GetLogger().Error(err)
		return
	}

	if !strings.EqualFold(minerResponse.Status, constants.SWAN_API_STATUS_SUCCESS) {
		logs.GetLogger().Error("Error: failed to update bid configuration.", minerResponse.Message)
		return
	}

	logs.GetLogger().Info("Bid configuration updated.")
}

func (swanClient *SwanClient) SendHeartbeatRequest(minerFid string) error {
	err := swanClient.GetJwtTokenUp3Times()
	if err != nil {
		logs.GetLogger().Error(err)
		return err
	}

	apiUrl := swanClient.ApiUrl + "/heartbeat"
	params := url.Values{}
	params.Add("miner_id", minerFid)

	response := web.HttpPost(apiUrl, swanClient.SwanToken, strings.NewReader(params.Encode()))

	if strings.Contains(response, "fail") {
		err := fmt.Errorf("failed to send heartbeat")
		logs.GetLogger().Error(err)
		return err
	}

	status := utils.GetFieldStrFromJson(response, "status")
	message := utils.GetFieldStrFromJson(response, "message")
	msg := fmt.Sprintf("status:%s, message:%s", status, message)
	logs.GetLogger().Info(msg)
	return nil
}
