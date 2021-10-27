package main

import (
	leveldbapi "github.com/filswan/go-swan-lib/client"
	"github.com/filswan/go-swan-lib/logs"
)

func main() {
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
