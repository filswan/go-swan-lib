package main

import (
	leveldbapi "github.com/filswan/go-swan-lib/client"
	"github.com/filswan/go-swan-lib/logs"
)

func main() {
	leveldbFile := "~/.swan/client/leveldbfile"
	leveldbKey := "test"
	err := leveldbapi.Put(leveldbFile, leveldbKey, "hello")
	if err != nil {
		logs.GetLogger().Error(err)
	}

	data, err := leveldbapi.Get(leveldbFile, leveldbKey)
	if err != nil {
		logs.GetLogger().Error(err)
	}
	datastr := string(data)
	if err != nil {
		logs.GetLogger().Error(err)
	}

	logs.GetLogger().Info(datastr)
}
