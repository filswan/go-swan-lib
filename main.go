package main

import (
	"math"

	"github.com/filswan/go-swan-lib/client"
	"github.com/filswan/go-swan-lib/client/lotus"
	"github.com/filswan/go-swan-lib/client/swan"
	"github.com/filswan/go-swan-lib/logs"
	"github.com/filswan/go-swan-lib/utils"
)

func main() {
	testLotusClientQeryAsk()
}

func testLotusClientQeryAsk() {
	logs.GetLogger().Info(1e18 == math.Pow10(18))
	minerFid := "t03354"
	lotusClient, err := lotus.LotusGetClient("http://192.168.88.41:1234/rpc/v0", "")
	if err != nil {
		logs.GetLogger().Error(err)
		return
	}

	minerConf, err := lotusClient.LotusClientQueryAsk(minerFid)
	if err != nil {
		logs.GetLogger().Error(err)
		return
	}

	logs.GetLogger().Info(minerConf.Price)
	logs.GetLogger().Info(minerConf.VerifiedPrice)
	logs.GetLogger().Info(minerConf.MinPieceSize)
	logs.GetLogger().Info(minerConf.MaxPieceSize)

	price, verifiedPrice, maxSize, minSize := lotusClient.LotusGetMinerConfig(minerFid)
	logs.GetLogger().Info(*price)
	logs.GetLogger().Info(*verifiedPrice)
	logs.GetLogger().Info(*maxSize)
	logs.GetLogger().Info(*minSize)
}

func testRandStr() {
	var letterRunes = []rune("abcdefghijklmnopqrstuvwxyz0123456789")
	var optionChars = "abcdefghijklmnopqrstuvwxyz0123456789"
	for i := 0; i < 100; i++ {
		randStr := utils.RandStringRunes(letterRunes, 6)
		randStr1 := utils.RandString(optionChars, 6)
		logs.GetLogger().Info(randStr, "  ", randStr1)
	}
}

func testTask() {
	swanClient, err := swan.SwanGetClient("", "", "", "")
	if err != nil {
		logs.GetLogger().Error(err)
		return
	}

	tasks, err := swanClient.SwanGetAssignedTasks()
	if err != nil {
		logs.GetLogger().Error(err)
		return
	}

	for _, task := range tasks {
		logs.GetLogger().Info(task.Uuid, " ", task.TaskFileName)
	}

	utils.DecodeJwtToken("")
}

func testIpfs() {
	client.IpfsUploadCarFileByWebApi("http://192.168.88.41:5001/api/v0/add?stream-channels=true&pin=true", "/Users/dorachen/go-workspace/src/testGo/go.mod")
}

func testLevelDb() {
	leveldbFile := "~/.swan/client/leveldbfile"
	leveldbKey := "test"
	err := client.LevelDbPut(leveldbFile, leveldbKey, "hello")
	if err != nil {
		logs.GetLogger().Error(err)
	}

	data, err := client.LevelDbGet(leveldbFile, leveldbKey)
	if err != nil {
		logs.GetLogger().Error(err)
	}
	datastr := string(data)
	if err != nil {
		logs.GetLogger().Error(err)
	}

	logs.GetLogger().Info(datastr)
}
