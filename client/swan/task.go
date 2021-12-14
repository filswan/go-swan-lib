package swan

import (
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/filswan/go-swan-lib/client/web"
	"github.com/filswan/go-swan-lib/constants"
	"github.com/filswan/go-swan-lib/logs"
	"github.com/filswan/go-swan-lib/model"
	"github.com/filswan/go-swan-lib/utils"
)

type SwanServerResponse struct {
	Status  string `json:"status"`
	Message string `json:"message"`
}

func (swanClient *SwanClient) SwanCreateTask(task model.Task, carFiles []*model.FileDesc) (*SwanServerResponse, error) {
	apiUrl := utils.UrlJoin(swanClient.ApiUrl, "tasks/create_task")
	params := map[string]interface{}{
		"task":      task,
		"car_files": carFiles,
	}

	response := web.HttpPost(apiUrl, swanClient.SwanToken, params)

	if response == "" {
		err := fmt.Errorf("no response from:%s", apiUrl)
		return nil, err
	}

	swanServerResponse := &SwanServerResponse{}
	err := json.Unmarshal([]byte(response), swanServerResponse)
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

func (swanClient *SwanClient) SwanUpdateTaskByUuid(task model.Task, carFiles []*model.FileDesc) (*SwanServerResponse, error) {
	apiUrl := utils.UrlJoin(swanClient.ApiUrl, "tasks/create_task")
	params := map[string]interface{}{
		"task":      task,
		"car_files": carFiles,
	}

	response := web.HttpPost(apiUrl, swanClient.SwanToken, params)

	if response == "" {
		err := fmt.Errorf("no response from:%s", apiUrl)
		logs.GetLogger().Error(err)
		return nil, err
	}

	swanServerResponse := &SwanServerResponse{}
	err := json.Unmarshal([]byte(response), swanServerResponse)
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

type GetTaskResult struct {
	Data   GetTaskResultData `json:"data"`
	Status string            `json:"status"`
}

type GetTaskResultData struct {
	Task           []model.Task `json:"task"`
	TotalItems     int          `json:"total_items"`
	TotalTaskCount int          `json:"total_task_count"`
}

func (swanClient *SwanClient) SwanGetTasks(limit *int, status *string) (*GetTaskResult, error) {
	apiUrl := utils.UrlJoin(swanClient.ApiUrl, "tasks")
	filters := ""
	if limit != nil {
		filters = filters + "?limit=" + strconv.Itoa(*limit)
	}

	if status != nil {
		if filters == "" {
			filters = filters + "?"
		} else {
			filters = filters + "&"
		}
		filters = filters + "status=" + *status
	}

	apiUrl = apiUrl + filters

	response := web.HttpGet(apiUrl, swanClient.SwanToken, "")

	if response == "" {
		err := errors.New("failed to get tasks from swan")
		logs.GetLogger().Error(err)
		return nil, err
	}

	getTaskResult := &GetTaskResult{}
	err := json.Unmarshal([]byte(response), getTaskResult)
	if err != nil {
		logs.GetLogger().Error(err)
		return nil, err
	}

	if !strings.EqualFold(getTaskResult.Status, constants.SWAN_API_STATUS_SUCCESS) {
		err := fmt.Errorf("error:%s", getTaskResult.Status)
		logs.GetLogger().Error(err)
		return nil, err
	}

	return getTaskResult, nil
}

func (swanClient *SwanClient) SwanGetAllTasks(status string) ([]model.Task, error) {
	limit := -1
	getTaskResult, err := swanClient.SwanGetTasks(&limit, &status)
	if err != nil {
		logs.GetLogger().Error(err)
		return nil, err
	}
	return getTaskResult.Data.Task, err
}

type GetTaskByUuidResult struct {
	Data   GetTaskByUuidResultData `json:"data"`
	Status string                  `json:"status"`
}
type GetTaskByUuidResultData struct {
	//AverageBid       string              `json:"average_bid"`
	Task             model.Task           `json:"task"`
	CarFiles         []model.CarFile      `json:"car_file"`
	Miner            model.Miner          `json:"miner"`
	Deal             []*model.OfflineDeal `json:"deal"`
	TotalItems       int                  `json:"total_items"`
	TotalTaskCount   int                  `json:"total_task_count"`
	BidCount         int                  `json:"bid_count"`
	DealCompleteRate string               `json:"deal_complete_rate"`
}

func (swanClient *SwanClient) SwanGetTaskByUuid(taskUuid string) (*GetTaskByUuidResult, error) {
	if len(taskUuid) == 0 {
		err := fmt.Errorf("please provide task uuid")
		logs.GetLogger().Error(err)
		return nil, err
	}
	apiUrl := utils.UrlJoin(swanClient.ApiUrl, "tasks", taskUuid)

	response := web.HttpGet(apiUrl, swanClient.SwanToken, "")

	if response == "" {
		err := fmt.Errorf("no response from:%s", apiUrl)
		logs.GetLogger().Error(err)
		return nil, err
	}

	getTaskByUuidResult := &GetTaskByUuidResult{}
	err := json.Unmarshal([]byte(response), getTaskByUuidResult)
	if err != nil {
		logs.GetLogger().Error(err)
		return nil, err
	}

	if !strings.EqualFold(getTaskByUuidResult.Status, constants.SWAN_API_STATUS_SUCCESS) {
		err := fmt.Errorf("error:%s", getTaskByUuidResult.Status)
		logs.GetLogger().Error(err)
		return nil, err
	}

	return getTaskByUuidResult, nil
}

type SwanOfflineDeals4CarFileResult struct {
	Data   SwanOfflineDeals4CarFileResultData `json:"data"`
	Status string                             `json:"status"`
}
type SwanOfflineDeals4CarFileResultData struct {
	CarFile          model.CarFile        `json:"car_file"`
	OfflineDeals     []*model.OfflineDeal `json:"deal"`
	TotalItems       int                  `json:"total_items"`
	TotalTaskCount   int                  `json:"total_task_count"`
	BidCount         int                  `json:"bid_count"`
	DealCompleteRate string               `json:"deal_complete_rate"`
}

func (swanClient *SwanClient) SwanOfflineDeals4CarFile(taskUuid, carFileUrl string) (*SwanOfflineDeals4CarFileResultData, error) {
	if len(taskUuid) == 0 {
		err := fmt.Errorf("please provide task uuid")
		logs.GetLogger().Error(err)
		return nil, err
	}
	if len(carFileUrl) == 0 {
		err := fmt.Errorf("please provide car file url")
		logs.GetLogger().Error(err)
		return nil, err
	}
	apiUrl := fmt.Sprintf("%s?task_uuid=%s&&car_file_url=%s", swanClient.ApiUrl, taskUuid, carFileUrl)

	response := web.HttpGet(apiUrl, swanClient.SwanToken, "")

	if response == "" {
		err := fmt.Errorf("no response from:%s", apiUrl)
		logs.GetLogger().Error(err)
		return nil, err
	}

	swanOfflineDeals4CarFileResult := &SwanOfflineDeals4CarFileResult{}
	err := json.Unmarshal([]byte(response), swanOfflineDeals4CarFileResult)
	if err != nil {
		logs.GetLogger().Error(err)
		return nil, err
	}

	if !strings.EqualFold(swanOfflineDeals4CarFileResult.Status, constants.SWAN_API_STATUS_SUCCESS) {
		err := fmt.Errorf("error:%s", swanOfflineDeals4CarFileResult.Status)
		logs.GetLogger().Error(err)
		return nil, err
	}

	return &swanOfflineDeals4CarFileResult.Data, nil
}
