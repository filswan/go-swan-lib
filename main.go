package main

import (
	"github.com/filswan/go-swan-lib/client"
	"github.com/filswan/go-swan-lib/client/swan"
	"github.com/filswan/go-swan-lib/logs"
	"github.com/filswan/go-swan-lib/utils"
)

func main() {
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

func testRandStr() {
	client.IpfsUploadCarFileByWebApi("http://192.168.88.41:5001/api/v0/add?stream-channels=true&pin=true", "/Users/dorachen/go-workspace/src/testGo/go.mod")

	var letterRunes = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

	logs.GetLogger().Info(utils.RandStringRunes(letterRunes, 6))
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
