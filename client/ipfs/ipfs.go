package ipfs

import (
	"fmt"

	"github.com/filswan/go-swan-lib/client/web"
	"github.com/filswan/go-swan-lib/constants"
	"github.com/filswan/go-swan-lib/utils"
)

func IpfsUploadFileByWebApi(apiUrl, filefullpath string) (*string, error) {
	response, err := web.HttpUploadFileByStream(apiUrl, filefullpath)
	if err != nil {
		//logs.GetLogger().Error(err)
		return nil, err
	}

	if response == constants.EMPTY_STRING {
		err := fmt.Errorf("no response from %s", apiUrl)
		//logs.GetLogger().Error(err)
		return nil, err
	}
	//logs.GetLogger().Info(response)
	fileHash := utils.GetFieldStrFromJson(response, "Hash")
	//logs.GetLogger().Info(carFileHash)

	if fileHash == constants.EMPTY_STRING {
		err := fmt.Errorf("cannot get file hash from response:%s", response)
		//logs.GetLogger().Error(err)
		return nil, err
	}

	return &fileHash, nil
}

func Export2CarFile(apiUrl, fileHash string, carFileFullPath string) {
	apiUrlFull := utils.UrlJoin(apiUrl, "api/v0/dag/export")
	apiUrlFull = apiUrlFull + "?arg=" + fileHash

}
