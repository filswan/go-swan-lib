package client

import (
	"fmt"

	"github.com/filswan/go-swan-lib/constants"
	"github.com/filswan/go-swan-lib/logs"
	"github.com/filswan/go-swan-lib/utils"
)

func IpfsUploadCarFileByWebApi(apiUrl, carFilePath string) (*string, error) {
	response, err := HttpUploadFileByStream(apiUrl, carFilePath)
	if err != nil {
		logs.GetLogger().Error(err)
		return nil, err
	}

	if response == constants.EMPTY_STRING {
		err := fmt.Errorf("no response from %s", apiUrl)
		logs.GetLogger().Error(err)
		return nil, err
	}
	//logs.GetLogger().Info(response)
	carFileHash := utils.GetFieldStrFromJson(response, "Hash")
	//logs.GetLogger().Info(carFileHash)

	if carFileHash == constants.EMPTY_STRING {
		err := fmt.Errorf("cannot get file hash from response:%s", response)
		logs.GetLogger().Error(err)
		return nil, err
	}

	return &carFileHash, nil
}
