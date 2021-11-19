package main

import (
	"math"
	"os"
	"strconv"

	"github.com/filswan/go-swan-lib/client"
	"github.com/filswan/go-swan-lib/client/lotus"
	"github.com/filswan/go-swan-lib/client/swan"
	"github.com/filswan/go-swan-lib/logs"
	"github.com/filswan/go-swan-lib/utils"
)

func main() {
	result := utils.Convert2Title("StorageDealStaged,funds computed:823728125000000,funds reserved:823728125000000,funds released:823728125000000")
	logs.GetLogger().Info(result)
	result = utils.Convert2Title("abc,    def.    txt    ddd ....a. . . ....")
	logs.GetLogger().Info(result)
	logs.GetLogger().Info("a"[1:])
	testLotusClientDealInfo()
	//testGenerateUploadFile()
}

func testLotusClientDealInfo() {
	logs.GetLogger().Info(1e18 == math.Pow10(18))
	lotusClient, err := lotus.LotusGetClient("http://192.168.88.41:1234/rpc/v0", "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJBbGxvdyI6WyJyZWFkIiwid3JpdGUiLCJzaWduIiwiYWRtaW4iXX0.-Y4pF34RGOten6YXoau-sEMOWOEeiHwGh9u2lsl4cv8")
	if err != nil {
		logs.GetLogger().Error(err)
		return
	}
	costStatus, err := lotusClient.LotusClientGetDealInfo("bafyreifrrcveyjcc3vnpvahuus2whngqhndt3a6qnqel2zpyykdk6xtspm")
	if err != nil {
		logs.GetLogger().Error(err)
		return
	}
	logs.GetLogger().Info(*costStatus)
}

func testLotusClientCalcCommP() {
	lotusClient, err := lotus.LotusGetClient("http://192.168.88.41:1234/rpc/v0", "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJBbGxvdyI6WyJyZWFkIiwid3JpdGUiLCJzaWduIiwiYWRtaW4iXX0.-Y4pF34RGOten6YXoau-sEMOWOEeiHwGh9u2lsl4cv8")
	if err != nil {
		logs.GetLogger().Error(err)
		return
	}
	pieceCid := lotusClient.LotusClientCalcCommP("~/swan-gh/gotest/test/sc/fs3test/volume_1637010138995215.car")
	if pieceCid != nil {
		logs.GetLogger().Info(*pieceCid)
	} else {
		logs.GetLogger().Error("piece CID is nil")
	}
}

func testDataCap() {
	swanClient, _ := swan.SwanGetClient("http://192.168.88.41:5002", "yr0wUW37PEm1ZUtes-0NVg", "878da9defc1841dd5ab9f4dcef1ec9af", "")
	isV, _ := swanClient.CheckDatacap("t3u7pumush376xbytsgs5wabkhtadjzfydxxda2vzyasg7cimkcphswrq66j4dubbhwpnojqd3jie6ermpwvvq")
	logs.GetLogger().Info(isV)
}

func testGenerateFile() {
	utils.GenerateFile("./", "test.txt", 2)
}

func testGenerateUploadFile() {
	if len(os.Args) <= 1 {
		logs.GetLogger().Error("please provide subcommand:generate|upload")
		return
	}
	switch os.Args[1] {
	case "generate":
		logs.GetLogger().Println("usage:swan-lib generate filepath filename filesizeInGigabyte")
		if len(os.Args) < 5 {
			logs.GetLogger().Error("not enough arguments")
			return
		}
		filepath := os.Args[2]
		filename := os.Args[3]
		filesizeInGigabyte, err := strconv.ParseInt(os.Args[4], 10, 64)
		if err != nil {
			logs.GetLogger().Error(err)
			return
		}

		utils.GenerateFile(filepath, filename, filesizeInGigabyte)
	case "upload":
		logs.GetLogger().Println("usage:swan-lib upload apiUrl filefullpath")
		if len(os.Args) < 4 {
			logs.GetLogger().Error("not enough arguments")
			return
		}

		apiUrl := os.Args[2]
		filefullpath := os.Args[3]

		if !utils.IsFileExistsFullPath(filefullpath) {
			logs.GetLogger().Error(filefullpath, " not exists")
			return
		}

		carFileHash, err := client.IpfsUploadCarFileByWebApi(apiUrl, filefullpath)
		if err != nil {
			logs.GetLogger().Error(err)
			return
		}

		logs.GetLogger().Info(*carFileHash)
	default:
		logs.GetLogger().Error("only support subcommand:generate|upload")
	}
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
