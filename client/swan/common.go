package swan

import (
	"fmt"
	"net/url"
	"strings"

	"github.com/filswan/go-swan-lib/client/web"
	"github.com/filswan/go-swan-lib/constants"
	"github.com/filswan/go-swan-lib/logs"
	"github.com/filswan/go-swan-lib/utils"
)

type TokenAccessInfo struct {
	ApiKey      string `json:"apikey"`
	AccessToken string `json:"access_token"`
}

type SwanClient struct {
	ApiUrlToken string
	ApiUrl      string
	SwanToken   string
	ApiKey      string
	AccessToken string
}

func (swanClient *SwanClient) GetJwtToken() error {
	data := TokenAccessInfo{
		ApiKey:      swanClient.ApiKey,
		AccessToken: swanClient.AccessToken,
	}

	if len(swanClient.ApiUrlToken) == 0 {
		err := fmt.Errorf("api url is required")
		logs.GetLogger().Error(err)
		return err
	}

	if len(data.ApiKey) == 0 {
		err := fmt.Errorf("api key is required")
		logs.GetLogger().Error(err)
		return err
	}

	if len(data.AccessToken) == 0 {
		err := fmt.Errorf("acess token is required")
		logs.GetLogger().Error(err)
		return err
	}

	apiUrl := swanClient.ApiUrlToken + "/user/api_keys/jwt"

	response := web.HttpPostNoToken(apiUrl, data)

	if len(response) == 0 {
		err := fmt.Errorf("no response from swan platform:%s", apiUrl)
		logs.GetLogger().Error(err)
		return err
	}

	if strings.Contains(response, "fail") {
		message := utils.GetFieldStrFromJson(response, "message")
		status := utils.GetFieldStrFromJson(response, "status")
		err := fmt.Errorf("status:%s, message:%s", status, message)
		logs.GetLogger().Error(err)

		if message == "api_key Not found" {
			logs.GetLogger().Error("please check your api key")
		}

		if message == "please provide a valid api token" {
			logs.GetLogger().Error("Please check your access token")
		}

		logs.GetLogger().Info("for more information about how to config, please check https://docs.filswan.com/run-swan-provider/config-swan-provider")

		return err
	}

	jwtToken := utils.GetFieldMapFromJson(response, "data")
	if jwtToken == nil {
		err := fmt.Errorf("error: fail to connect to swan api")
		logs.GetLogger().Error(err)
		return err
	}

	swanClient.SwanToken = jwtToken["jwt"].(string)

	return nil
}

func GetClient(apiUrlToken, apiUrl, apiKey, accessToken, swanToken string) (*SwanClient, error) {
	if len(apiUrl) == 0 {
		err := fmt.Errorf("api url is required")
		logs.GetLogger().Error(err)
		return nil, err
	}

	swanClient := &SwanClient{
		ApiUrlToken: apiUrlToken,
		ApiUrl:      apiUrl,
		ApiKey:      apiKey,
		AccessToken: accessToken,
		SwanToken:   swanToken,
	}

	if swanToken == constants.EMPTY_STRING {
		err := swanClient.GetJwtTokenUp3Times()
		return swanClient, err
	}

	return swanClient, nil
}

func (swanClient *SwanClient) GetJwtTokenUp3Times() error {
	if len(swanClient.ApiUrlToken) == 0 {
		err := fmt.Errorf("api url for token is required")
		logs.GetLogger().Error(err)
		return err
	}

	if len(swanClient.ApiUrl) == 0 {
		err := fmt.Errorf("api url is required")
		logs.GetLogger().Error(err)
		return err
	}

	if len(swanClient.ApiKey) == 0 {
		err := fmt.Errorf("api key is required")
		logs.GetLogger().Error(err)
		return err
	}

	if len(swanClient.AccessToken) == 0 {
		err := fmt.Errorf("access token is required")
		logs.GetLogger().Error(err)
		return err
	}

	var err error
	for i := 0; i < 3; i++ {
		err = swanClient.GetJwtToken()
		if err == nil {
			break
		}
		logs.GetLogger().Error(err)
	}

	if err != nil {
		err = fmt.Errorf("failed to connect to swan platform after trying 3 times")
		logs.GetLogger().Error(err)
		return err
	}

	return nil
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
