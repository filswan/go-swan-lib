package swan

import (
	"fmt"
	"net/url"
	"strings"

	"github.com/filswan/go-swan-lib/client"
	"github.com/filswan/go-swan-lib/constants"
	"github.com/filswan/go-swan-lib/logs"
	"github.com/filswan/go-swan-lib/utils"
)

func (swanClient *SwanClient) CheckDatacap(wallet string) (bool, error) {
	apiUrl := swanClient.ApiUrl + "/tools/check_datacap?address=" + wallet
	params := url.Values{}

	response := client.HttpGetNoToken(apiUrl, strings.NewReader(params.Encode()))

	if response == "" {
		err := fmt.Errorf("no response from:%s", apiUrl)
		logs.GetLogger().Error(err)
		return false, err
	}

	status := utils.GetFieldStrFromJson(response, "status")
	message := utils.GetFieldStrFromJson(response, "message")

	if strings.EqualFold(status, constants.SWAN_API_STATUS_SUCCESS) {
		return true, nil
	}

	if strings.EqualFold(message, constants.WALLET_NON_VERIFIED_MESSAGE) {
		return false, nil
	}

	err := fmt.Errorf("status:%s,message:%s", status, message)
	logs.GetLogger().Error(err)
	return false, err
}
