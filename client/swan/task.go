package swan

import (
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/filswan/go-swan-lib/client"
	"github.com/filswan/go-swan-lib/constants"
	"github.com/filswan/go-swan-lib/logs"
	"github.com/filswan/go-swan-lib/model"
	"github.com/filswan/go-swan-lib/utils"
)

type SwanCreateTaskResponse struct {
	Data    SwanCreateTaskResponseData `json:"data"`
	Status  string                     `json:"status"`
	Message string                     `json:"message"`
}

type SwanCreateTaskResponseData struct {
	Filename string `json:"filename"`
	Uuid     string `json:"uuid"`
}

func (swanClient *SwanClient) SwanCreateTask(task model.Task, csvFilePath string) (*SwanCreateTaskResponse, error) {
	apiUrl := swanClient.ApiUrl + "/tasks"

	params := map[string]string{}
	params["task_name"] = task.TaskName
	params["curated_dataset"] = task.CuratedDataset
	params["description"] = task.Description
	params["is_public"] = strconv.Itoa(*task.IsPublic)

	params["type"] = task.Type

	if task.MinerFid != "" {
		params["miner_id"] = task.MinerFid
	}
	params["fast_retrieval"] = strconv.FormatBool(task.FastRetrievalBool)
	params["bid_mode"] = strconv.Itoa(*task.BidMode)
	params["max_price"] = (*task.MaxPrice).String()
	params["expire_days"] = strconv.Itoa(*task.ExpireDays)
	params["source_id"] = strconv.Itoa(task.SourceId)
	params["duration"] = strconv.Itoa(task.Duration)

	response, err := client.HttpPostFile(apiUrl, swanClient.JwtToken, params, "file", csvFilePath)
	if err != nil {
		logs.GetLogger().Error(err)
		return nil, err
	}

	swanCreateTaskResponse := &SwanCreateTaskResponse{}
	err = json.Unmarshal([]byte(response), swanCreateTaskResponse)
	if err != nil {
		logs.GetLogger().Error(err)
		return nil, err
	}

	if !strings.EqualFold(swanCreateTaskResponse.Status, constants.SWAN_API_STATUS_SUCCESS) {
		err := fmt.Errorf("error:%s,%s", swanCreateTaskResponse.Status, swanCreateTaskResponse.Message)
		logs.GetLogger().Error(err)
		return nil, err
	}

	return swanCreateTaskResponse, nil
}

type GetTaskByUuidResult struct {
	Data   GetTaskResultData `json:"data"`
	Status string            `json:"status"`
}

type GetTaskByUuidResultData struct {
	Task           model.Task `json:"task"`
	TotalItems     int        `json:"total_items"`
	TotalTaskCount int        `json:"total_task_count"`
}

func (swanClient *SwanClient) SwanGetTaskByUuid(uuid string) (*GetTaskByUuidResult, error) {
	apiUrl := swanClient.ApiUrl + "/tasks/" + uuid
	//logs.GetLogger().Info("Getting My swan tasks info")
	response := client.HttpGet(apiUrl, swanClient.JwtToken, "")

	if response == "" {
		err := errors.New("failed to get tasks from swan")
		logs.GetLogger().Error(err)
		return nil, err
	}

	//logs.GetLogger().Info(response)

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

type GetTaskResult struct {
	Data   GetTaskResultData `json:"data"`
	Status string            `json:"status"`
}

type GetTaskResultData struct {
	Task           []model.Task `json:"task"`
	TotalItems     int          `json:"total_items"`
	TotalTaskCount int          `json:"total_task_count"`
}

func (swanClient *SwanClient) SwanGetTasks(limit *int) (*GetTaskResult, error) {
	apiUrl := swanClient.ApiUrl + "/tasks"
	if limit != nil {
		apiUrl = apiUrl + "?limit=" + strconv.Itoa(*limit)
	}
	//logs.GetLogger().Info("Getting My swan tasks info")
	response := client.HttpGet(apiUrl, swanClient.JwtToken, "")

	if response == "" {
		err := errors.New("failed to get tasks from swan")
		logs.GetLogger().Error(err)
		return nil, err
	}

	//logs.GetLogger().Info(response)

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

func (swanClient *SwanClient) SwanGetAssignedTasksByLimit(limit *int) (*GetTaskResult, error) {
	apiUrl := swanClient.ApiUrl + "/tasks?status=Assigned"
	if limit != nil {
		apiUrl = apiUrl + "&limit=" + strconv.Itoa(*limit)
	}
	//logs.GetLogger().Info("Getting My swan tasks info")
	response := client.HttpGet(apiUrl, swanClient.JwtToken, "")

	if response == "" {
		err := errors.New("failed to get tasks from swan")
		logs.GetLogger().Error(err)
		return nil, err
	}

	//logs.GetLogger().Info(response)

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

func (swanClient *SwanClient) SwanGetAssignedTasks() ([]model.Task, error) {
	getTaskResult, err := swanClient.SwanGetAssignedTasksByLimit(nil)
	if err != nil {
		logs.GetLogger().Error(err)
		return nil, err
	}

	if len(getTaskResult.Data.Task) == 0 {
		return nil, nil
	}
	//logs.GetLogger().Info(len(getTaskResult.Data.Task), " ", getTaskResult.Data.TotalTaskCount)

	getTaskResult, err = swanClient.SwanGetAssignedTasksByLimit(&getTaskResult.Data.TotalTaskCount)
	if err != nil {
		logs.GetLogger().Error(err)
		return nil, err
	}

	if len(getTaskResult.Data.Task) == 0 {
		return nil, nil
	}

	//logs.GetLogger().Info(len(getTaskResult.Data.Task), " ", getTaskResult.Data.TotalTaskCount)

	result := []model.Task{}

	for _, task := range getTaskResult.Data.Task {
		if task.Status == constants.TASK_STATUS_ASSIGNED && task.MinerFid != "" {
			//logs.GetLogger().Info("id: ", task.Id, " task:", task.Status, " miner:", *task.MinerFid)
			result = append(result, task)
		}
	}

	return result, nil
}

type GetOfflineDealsByTaskUuidResult struct {
	Data   GetOfflineDealsByTaskUuidResultData `json:"data"`
	Status string                              `json:"status"`
}
type GetOfflineDealsByTaskUuidResultData struct {
	AverageBid       string              `json:"average_bid"`
	BidCount         int                 `json:"bid_count"`
	DealCompleteRate string              `json:"deal_complete_rate"`
	Deal             []model.OfflineDeal `json:"deal"`
	Miner            model.Miner         `json:"miner"`
	Task             model.Task          `json:"task"`
}

func (swanClient *SwanClient) SwanGetOfflineDealsByTaskUuid(taskUuid string) (*GetOfflineDealsByTaskUuidResult, error) {
	if len(taskUuid) == 0 {
		err := fmt.Errorf("please provide task uuid")
		logs.GetLogger().Error(err)
		return nil, err
	}
	apiUrl := swanClient.ApiUrl + "/tasks/" + taskUuid
	logs.GetLogger().Info("Getting My swan tasks info")
	response := client.HttpGet(apiUrl, swanClient.JwtToken, "")

	if response == "" {
		err := errors.New("failed to get tasks from swan")
		logs.GetLogger().Error(err)
		return nil, err
	}
	//logs.GetLogger().Info(response)

	getOfflineDealsByTaskUuidResult := &GetOfflineDealsByTaskUuidResult{}
	err := json.Unmarshal([]byte(response), getOfflineDealsByTaskUuidResult)
	if err != nil {
		logs.GetLogger().Error(err)
		return nil, err
	}

	if !strings.EqualFold(getOfflineDealsByTaskUuidResult.Status, constants.SWAN_API_STATUS_SUCCESS) {
		err := fmt.Errorf("error:%s", getOfflineDealsByTaskUuidResult.Status)
		logs.GetLogger().Error(err)
		return nil, err
	}

	return getOfflineDealsByTaskUuidResult, nil
}

func (swanClient *SwanClient) SwanUpdateTaskByUuid(taskUuid string, minerFid string, csvFilePath string) error {
	apiUrl := swanClient.ApiUrl + "/uuid_tasks/" + taskUuid
	params := map[string]string{}
	params["miner_fid"] = minerFid

	response, err := client.HttpPutFile(apiUrl, swanClient.JwtToken, params, "file", csvFilePath)
	if err != nil {
		logs.GetLogger().Error(err)
		return err
	}

	if response == "" {
		err := fmt.Errorf("no response from :%s", apiUrl)
		logs.GetLogger().Error(err)
		return err
	}

	status := utils.GetFieldStrFromJson(response, "status")
	if !strings.EqualFold(status, constants.SWAN_API_STATUS_SUCCESS) {
		message := utils.GetFieldStrFromJson(response, "message")
		err := fmt.Errorf("access:%s failed, status:%s, message:%s", apiUrl, status, message)
		logs.GetLogger().Error(err)
		return err
	}
	data := utils.GetFieldMapFromJson(response, "data")
	filename := data["filename"]

	msg := fmt.Sprintf("access:%s succeeded, file:%s", apiUrl, filename)
	logs.GetLogger().Info(msg)

	return nil
}

func (swanClient *SwanClient) SwanUpdateAssignedTask(taskUuid, status, csvFilePath string) (*SwanCreateTaskResponse, error) {
	apiUrl := swanClient.ApiUrl + "/tasks/" + taskUuid
	logs.GetLogger().Info("Updating Swan task")
	params := map[string]string{}
	params["status"] = status

	response, err := client.HttpPutFile(apiUrl, swanClient.JwtToken, params, "file", csvFilePath)
	if err != nil {
		logs.GetLogger().Error(err)
		return nil, err
	}

	swanCreateTaskResponse := &SwanCreateTaskResponse{}
	err = json.Unmarshal([]byte(response), swanCreateTaskResponse)
	if err != nil {
		logs.GetLogger().Error(err)
		return nil, err
	}

	if !strings.EqualFold(swanCreateTaskResponse.Status, constants.SWAN_API_STATUS_SUCCESS) {
		err := fmt.Errorf("error:%s,%s", swanCreateTaskResponse.Status, swanCreateTaskResponse.Message)
		logs.GetLogger().Error(err)
		return nil, err
	}

	return swanCreateTaskResponse, nil
}
