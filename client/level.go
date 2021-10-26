package client

import (
	"github.com/filswan/go-swan-lib/logs"

	"github.com/syndtr/goleveldb/leveldb"
)

func Put(dbFilepath, key, value string) error {
	db, err := leveldb.OpenFile(dbFilepath, nil)
	if err != nil {
		logs.GetLogger().Error(err)
		return err
	}
	defer db.Close()

	err = db.Put([]byte(key), []byte(value), nil)
	if err != nil {
		logs.GetLogger().Error(err)
		return err
	}

	return nil
}

func Get(dbFilepath, key string) ([]byte, error) {
	db, err := leveldb.OpenFile(dbFilepath, nil)
	if err != nil {
		logs.GetLogger().Error(err)
		return nil, err
	}
	defer db.Close()

	data, err := db.Get([]byte(key), nil)
	if err != nil {
		logs.GetLogger().Error(err)
		return nil, err
	}

	return data, nil
}

func Delete(dbFilepath, key, value string) error {
	db, err := leveldb.OpenFile(dbFilepath, nil)
	if err != nil {
		logs.GetLogger().Error(err)
		return err
	}
	defer db.Close()

	err = db.Delete([]byte("key"), nil)
	if err != nil {
		logs.GetLogger().Error(err)
		return err
	}

	return nil
}
