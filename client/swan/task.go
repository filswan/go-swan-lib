package swan

import (
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/filswan/go-swan-lib/client/web"
	"github.com/filswan/go-swan-lib/constants"
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
		return nil, err
	}

	if !strings.EqualFold(swanServerResponse.Status, constants.SWAN_API_STATUS_SUCCESS) {
		err := fmt.Errorf("error:%s,%s", swanServerResponse.Status, swanServerResponse.Message)
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
		return nil, err
	}

	swanServerResponse := &SwanServerResponse{}
	err := json.Unmarshal([]byte(response), swanServerResponse)
	if err != nil {
		return nil, err
	}

	if !strings.EqualFold(swanServerResponse.Status, constants.SWAN_API_STATUS_SUCCESS) {
		err := fmt.Errorf("error:%s,%s", swanServerResponse.Status, swanServerResponse.Message)
		return nil, err
	}

	return swanServerResponse, nil
}

type GetTaskByUuidResult struct {
	Data   GetTaskByUuidResultData `json:"data"`
	Status string                  `json:"status"`
}
type GetTaskByUuidResultData struct {
	//AverageBid       string              `json:"average_bid"`
	Task             model.Task          `json:"task"`
	Miner            model.Miner         `json:"miner"`
	Deal             []model.OfflineDeal `json:"deal"`
	TotalItems       int                 `json:"total_items"`
	TotalTaskCount   int                 `json:"total_task_count"`
	BidCount         int                 `json:"bid_count"`
	DealCompleteRate string              `json:"deal_complete_rate"`
}

func (swanClient *SwanClient) SwanGetTaskByUuid(taskUuid string) (*GetTaskByUuidResult, error) {
	if len(taskUuid) == 0 {
		err := fmt.Errorf("please provide task uuid")
		return nil, err
	}
	apiUrl := utils.UrlJoin(swanClient.ApiUrl, "tasks", taskUuid)

	response := web.HttpGet(apiUrl, swanClient.SwanToken, "")

	if response == "" {
		err := fmt.Errorf("no response from:%s", apiUrl)
		return nil, err
	}

	getTaskByUuidResult := &GetTaskByUuidResult{}
	err := json.Unmarshal([]byte(response), getTaskByUuidResult)
	if err != nil {
		return nil, err
	}

	if !strings.EqualFold(getTaskByUuidResult.Status, constants.SWAN_API_STATUS_SUCCESS) {
		err := fmt.Errorf("error:%s", getTaskByUuidResult.Status)
		return nil, err
	}

	return getTaskByUuidResult, nil
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
		return nil, err
	}

	getTaskResult := &GetTaskResult{}
	err := json.Unmarshal([]byte(response), getTaskResult)
	if err != nil {
		return nil, err
	}

	if !strings.EqualFold(getTaskResult.Status, constants.SWAN_API_STATUS_SUCCESS) {
		err := fmt.Errorf("error:%s", getTaskResult.Status)
		return nil, err
	}

	return getTaskResult, nil
}

func (swanClient *SwanClient) SwanGetAllTasks(status string) ([]model.Task, error) {
	getTaskResult, err := swanClient.SwanGetTasks(nil, &status)
	if err != nil {
		return nil, err
	}

	if len(getTaskResult.Data.Task) == 0 {
		return nil, nil
	}

	getTaskResult, err = swanClient.SwanGetTasks(&getTaskResult.Data.TotalTaskCount, &status)
	if err != nil {
		return nil, err
	}

	if len(getTaskResult.Data.Task) == 0 {
		return nil, nil
	}

	result := []model.Task{}

	for _, task := range getTaskResult.Data.Task {
		if task.Status == constants.TASK_STATUS_ASSIGNED && task.MinerFid != "" {
			result = append(result, task)
		}
	}

	return result, nil
}
