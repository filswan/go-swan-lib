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

type SwanClient struct {
	ApiUrl      string
	SwanToken   string
	ApiKey      string
	AccessToken string
}

func GetClient(apiUrl, apiKey, accessToken, swanToken string) (*SwanClient, error) {
	if len(apiUrl) == 0 {
		err := fmt.Errorf("api url is required")
		logs.GetLogger().Error(err)
		return nil, err
	}

	swanClient := &SwanClient{
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
