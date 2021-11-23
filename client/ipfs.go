package client

import (
	"context"
	"fmt"
	"os"

	"github.com/filswan/go-swan-lib/constants"
	"github.com/filswan/go-swan-lib/logs"
	"github.com/filswan/go-swan-lib/utils"
	files "github.com/ipfs/go-ipfs-files"
	ipfsClient "github.com/ipfs/go-ipfs-http-client"
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

func IpfsCreateCarFile(apiUrl, srcFilesPath string) error {
	/*
		api, err := ipfsClient.NewPathApi(apiUrl)
		if err != nil {
			return err
		}*/
	api, err := ipfsClient.NewLocalApi()
	if err != nil {
		return err
	}

	stat, err := os.Stat(srcFilesPath)
	if err != nil {
		return err
	}
	// This walks the filesystem at /tmp/example/ and create a list of the files / directories we have.
	node, err := files.NewSerialFile(srcFilesPath, true, stat)
	if err != nil {
		return err
	}

	// Add the files / directory to IPFS
	path, err := api.Unixfs().Add(context.Background(), node)
	if err != nil {
		return err
	}

	// Output the resulting CID
	fmt.Println(path.Root().String())

	return nil
}
