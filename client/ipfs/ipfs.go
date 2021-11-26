package ipfs

import (
	"context"
	"fmt"
	"os"

	"github.com/filswan/go-swan-lib/client"
	"github.com/filswan/go-swan-lib/constants"
	"github.com/filswan/go-swan-lib/logs"
	"github.com/filswan/go-swan-lib/utils"
	files "github.com/ipfs/go-ipfs-files"
	ipfsClient "github.com/ipfs/go-ipfs-http-client"
)

func IpfsUploadFileByWebApi(apiUrl, filefullpath string) (*string, error) {
	response, err := client.HttpUploadFileByStream(apiUrl, filefullpath)
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

func IpfsCreateCarFile(apiUrl, srcFilesPath string) error {
	/*api, err := ipfsClient.NewPathApi(apiUrl)
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

	logs.GetLogger().Info("path CID:", path.Cid())
	// Output the resulting CID
	logs.GetLogger().Info("result CID:", path.Root().String())

	return nil
}
