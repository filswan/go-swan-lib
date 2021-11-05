package client

import (
	"fmt"
	"strings"

	"github.com/filswan/go-swan-lib/constants"
	"github.com/filswan/go-swan-lib/logs"
	"github.com/filswan/go-swan-lib/utils"
)

func IpfsUploadCarFile(carFilePath string) (*string, error) {
	cmd := "ipfs add " + carFilePath
	logs.GetLogger().Info(cmd)

	result, err := ExecOsCmd(cmd, false)

	if err != nil {
		logs.GetLogger().Error(err)
		return nil, err
	}

	if result == "" {
		err := fmt.Errorf("cmd(%s) result is empty", cmd)
		logs.GetLogger().Error(err)
		return nil, err
	}

	words := strings.Fields(result)
	if len(words) < 2 {
		err := fmt.Errorf("cmd(%s) result(%s) does not have enough fields", cmd, result)
		logs.GetLogger().Error(err)
		return nil, err
	}

	carFileHash := words[1]

	return &carFileHash, nil
}

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
	logs.GetLogger().Info(response)
	carFileHash := utils.GetFieldStrFromJson(response, "Hash")
	logs.GetLogger().Info(carFileHash)

	if carFileHash == constants.EMPTY_STRING {
		err := fmt.Errorf("cannot get file hash from response:%s", response)
		logs.GetLogger().Error(err)
		return nil, err
	}

	return &carFileHash, nil
}
