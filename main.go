package main

import (
	"github.com/filswan/go-swan-lib/client"
	leveldbapi "github.com/filswan/go-swan-lib/client"
	"github.com/filswan/go-swan-lib/logs"
	"github.com/filswan/go-swan-lib/utils"
)

func main() {
	swanClient, err := client.SwanGetClient("", "", "", "")
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

func testRandStr() {
	client.IpfsUploadCarFileByWebApi("http://192.168.88.41:5001/api/v0/add?stream-channels=true&pin=true", "/Users/dorachen/go-workspace/src/testGo/go.mod")

	var letterRunes = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

	logs.GetLogger().Info(utils.RandStringRunes(letterRunes, 6))
}

func testLevelDb() {
	leveldbFile := "~/.swan/client/leveldbfile"
	leveldbKey := "test"
	err := leveldbapi.LevelDbPut(leveldbFile, leveldbKey, "hello")
	if err != nil {
		logs.GetLogger().Error(err)
	}

	data, err := leveldbapi.LevelDbGet(leveldbFile, leveldbKey)
	if err != nil {
		logs.GetLogger().Error(err)
	}
	datastr := string(data)
	if err != nil {
		logs.GetLogger().Error(err)
	}

	logs.GetLogger().Info(datastr)
}
