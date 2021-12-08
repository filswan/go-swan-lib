package test

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"math"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/filswan/go-swan-lib/client"
	"github.com/filswan/go-swan-lib/client/ipfs"
	"github.com/filswan/go-swan-lib/client/lotus"
	"github.com/filswan/go-swan-lib/client/swan"
	"github.com/filswan/go-swan-lib/client/web"
	"github.com/filswan/go-swan-lib/constants"
	"github.com/filswan/go-swan-lib/logs"
	"github.com/filswan/go-swan-lib/model"
	"github.com/filswan/go-swan-lib/utils"
)

func TestOsCmdClient() {
	result, err := client.ExecOsCmd("ls -l", true)
	logs.GetLogger().Info(result, err)

	result, err = client.ExecOsCmd("pwd", true)
	logs.GetLogger().Info(result, err)

	result, err = client.ExecOsCmd("ls -l | grep common", true)
	logs.GetLogger().Info(result, err)

	words := strings.Fields(result)
	for _, word := range words {
		logs.GetLogger().Info(word)
	}
}

func TestCreateTask() {
	carFiles := []*model.FileDesc{}
	sourceId := 2
	startEpoch := 8888
	carFile := model.FileDesc{
		Uuid:           "",
		SourceFileName: "",
		SourceFilePath: "",
		SourceFileMd5:  "",
		SourceFileSize: 28888,
		CarFileName:    "bafybeie6mpsur5hul7ejcetnsfghisb4fnc7hewykilhul5aycyl23mlde.car",
		CarFilePath:    "",
		CarFileMd5:     "",
		CarFileUrl:     "http://192.168.88.41:5050/ipfs/QmS2dz1D6auP6viubUgfA4k7FFykb9ekefurRh1rSVS9hP",
		CarFileSize:    26666,
		DealCid:        "bafyreihdzwblxxuafpzn7olws62lxpqjc2u37zrer22gawxi3z7zrytrma",
		DataCid:        "bafybeie6mpsur5hul7ejcetnsfghisb4fnc7hewykilhul5aycyl23mlde",
		PieceCid:       "baga6ea4seaqorozvmj4b6lvruspz26gwubg37uesjuuyou5ictv4vhu5sdzliiy",
		MinerFid:       "t024557",
		StartEpoch:     &startEpoch,
		SourceId:       &sourceId,
		Cost:           "131781500000000",
	}
	carFiles = append(carFiles, &carFile)

	isPublic := 1
	bidMode := 1
	fastRetrieval := 1
	task := model.Task{
		TaskName:      "test",
		Type:          "regular",
		IsPublic:      &isPublic,
		Uuid:          "hello-world",
		BidMode:       &bidMode,
		FastRetrieval: &fastRetrieval,
		SourceId:      2,
		Duration:      8888,
	}

	params := map[string]interface{}{
		"task":      task,
		"car_files": carFiles,
	}
	response := web.HttpPost("http://127.0.0.1:8886/tasks/create_task", "eyJ0eXAiOiJKV1QiLCJhbGciOiJIUzI1NiJ9.eyJleHAiOjE2NzA1MTQzMjIsImlhdCI6MTYzODk3ODMyMiwic3ViIjoyNjF9.ZdKyTntKRwzYAbdhshJSNFYFpEr1I8HVTyDMtnFN6GM", params)
	logs.GetLogger().Info(response)

	response = web.HttpPut("http://127.0.0.1:8886/tasks/update_task_by_uuid", "eyJ0eXAiOiJKV1QiLCJhbGciOiJIUzI1NiJ9.eyJleHAiOjE2NzA1MTQzMjIsImlhdCI6MTYzODk3ODMyMiwic3ViIjoyNjF9.ZdKyTntKRwzYAbdhshJSNFYFpEr1I8HVTyDMtnFN6GM", params)
	logs.GetLogger().Info(response)
	//params
}

func TestOsCmdClient1() {
	/*result, err := */ client.ExecOsCmd2Screen("ls -l", true)
	//logs.GetLogger().Info(result, err)

	/*result, err = */
	client.ExecOsCmd2Screen("pwd", true)
	//logs.GetLogger().Info(result, err)

	/*result, err = */
	client.ExecOsCmd2Screen("ls -l | grep x", true)
	//logs.GetLogger().Info(result, err)
}

func TestLotusClient() {
	/*
		currentEpoch := client.LotusGetCurrentEpoch()
		logs.GetLogger().Info("currentEpoch: ", currentEpoch)
		status, message := client.LotusGetDealOnChainStatus("bafyreigbcdmozbfyr5sfipu7xm4fj23r3g2idgk7jibaku4y4r2z4x55bq")
		logs.GetLogger().Info("status: ", status)
		logs.GetLogger().Info("message: ", message)
		message = client.LotusImportData("bafyreiaj7av2qgziwfyvo663a2kjg3n35rvfr2i5r2dyrexxukdbybz7ky", "/tmp/swan-downloads/185/202107/go1.15.5.linux-amd64.tar.gz.car")
		logs.GetLogger().Info("message: ", message)
		message = client.LotusImportData("bafyreia5qflut2hqbwfhhhiybes5uhnx6aehgg3ltvam2aqbkekkyuoboy", "/tmp/swan-downloads/185/202107/go1.15.5.linux-amd64.tar.gz.car")
		logs.GetLogger().Info("message: ", message)
	*/
}

func Test() {
	TestLotusClient()
}

type Todo struct {
	UserID    int    `json:"userId"`
	ID        int    `json:"id"`
	Title     string `json:"title"`
	Completed bool   `json:"completed"`
}

func TestRestApiClient() {
	response := web.HttpGet("https://jsonplaceholder.typicode.com/todos/1", "", "")
	logs.GetLogger().Info(response)

	todo := Todo{1, 2, "lorem ipsum dolor sit amet", true}
	response = web.HttpPostNoToken("https://jsonplaceholder.typicode.com/todos", todo)
	logs.GetLogger().Info(response)

	response = web.HttpPut("https://jsonplaceholder.typicode.com/todos/1", "", todo)
	logs.GetLogger().Info(response)

	title := utils.GetFieldFromJson(response, "title")
	logs.GetLogger().Info(title)

	response = web.HttpDelete("https://jsonplaceholder.typicode.com/todos/1", "", todo)
	logs.GetLogger().Info(response)
}

func TestSwanClient() {
	swanClient, err := swan.SwanGetClient("", "", "", "")
	if err != nil {
		logs.GetLogger().Error(err)
	}

	deals := swanClient.SwanGetOfflineDeals("", "Downloading", "10")
	logs.GetLogger().Info(deals)

	response := swanClient.SwanUpdateOfflineDealStatus(2455, "Downloaded", "test note")
	logs.GetLogger().Info(response)

	response = swanClient.SwanUpdateOfflineDealStatus(2455, "Completed", "test note", "/test/test", "0003222")
	logs.GetLogger().Info(response)

	err = swanClient.SendHeartbeatRequest("")
	if err != nil {
		logs.GetLogger().Error(err)
	}
	logs.GetLogger().Info(response)
}

func TestDecodeJwtToken() {
	a, err := utils.DecodeJwtToken("eyJ0eXAiOiJKV1QiLCJhbGciOiJIUzI1NiJ9.eyJleHAiOjE2MzYxMzQyNjAsImlhdCI6MTYzNjA0Nzg2MCwic3ViIjoieXIwd1VXMzdQRW0xWlV0ZXMtME5WZyJ9.cAcLVH2SeLFykKcBACbInsz0BFyh5eHKBvfybMVoD7Y")
	if err != nil {
		logs.GetLogger().Error(err)
		return
	}
	logs.GetLogger().Error(a)
}

func TestGetCarFile() {
	fileHash, err := ipfs.IpfsUploadFileByWebApi("http://192.168.88.41:5001/api/v0/add?stream-channels=true&pin=true", "/home/peware/swan_dora/srcFiles/gnomad.genomes.v3.1.1.sites.chr22.vcf.bgz_02.car")
	if err != nil {
		logs.GetLogger().Error(err)
		return
	}

	logs.GetLogger().Info("source file hash:", *fileHash)

	cids := []string{
		*fileHash,
	}
	dataCid, err := ipfs.MergeFiles2CarFile("http://192.168.88.41:5001", cids)
	if err != nil {
		logs.GetLogger().Error(err)
		return
	}

	logs.GetLogger().Info("data CID:", *dataCid)
	err = ipfs.Export2CarFile("http://192.168.88.41:5001", *dataCid, "/home/peware/swan_dora/srcFiles/"+*dataCid+".car")
	if err != nil {
		logs.GetLogger().Error(err)
		return
	}
}

func TestGetDeal() {
	httpposturl := "https://eaehi1usrl.execute-api.us-east-1.amazonaws.com/test/link"
	fmt.Println("HTTP JSON POST URL:", httpposturl)

	var jsonData = []byte(`{"id":0,"data":{"deal":58172} }`)
	request, error := http.NewRequest("POST", httpposturl, bytes.NewBuffer(jsonData))
	request.Header.Set("Content-Type", "application/json; charset=UTF-8")

	client := &http.Client{}
	response, error := client.Do(request)
	if error != nil {
		panic(error)
	}
	defer response.Body.Close()

	fmt.Println("response Status:", response.Status)
	fmt.Println("response Headers:", response.Header)
	body, _ := ioutil.ReadAll(response.Body)
	fmt.Println("response Body:", string(body))
}

func TestMergeFile() {
	price := utils.ConvertPrice2AttoFil("0.0001 FIL")
	logs.GetLogger().Info(price)
	price = utils.ConvertPrice2AttoFil("1 fil")
	logs.GetLogger().Info(price)

	cids := []string{
		"QmaLTsfGTynrnbFeG5CmPRqdbYa1E4jbD9GjzkXXugTPfx",
		"QmdTf7TiBpYYv6sf7E9nbQdjLRDZLBvhSNwpGFf1B6zFtz",
	}
	dataCid, err := ipfs.MergeFiles2CarFile("http://127.0.0.1:5001", cids)
	if err != nil {
		logs.GetLogger().Error(err)
		return
	}
	logs.GetLogger().Info(*dataCid)
}

func TestIpfs() {
	hash, err := ipfs.IpfsUploadFileByWebApi("http://192.168.88.41:5001/api/v0/add?stream-channels=true&pin=true", "/Users/dorachen/go/src/go-swan-lib_DoraNebula/go.mod")
	if err != nil {
		logs.GetLogger().Error(err)
		return
	}
	logs.GetLogger().Info(*hash)
}

func TestLotusAuthVerify(apiUrl, token string) {
	auths, err := lotus.LotusAuthVerify(apiUrl, token)
	if err != nil {
		logs.GetLogger().Error(err)
		return
	}
	logs.GetLogger().Info(auths)

	isMatch, err := lotus.LotusCheckAuth(apiUrl, token, constants.LOTUS_AUTH_WRITE)
	if err != nil {
		logs.GetLogger().Error(err)
		return
	}
	logs.GetLogger().Info(isMatch)
}

func Test2Title() {
	result := utils.Convert2Title("StorageDealStaged,funds computed:823728125000000,funds reserved:823728125000000,funds released:823728125000000")
	logs.GetLogger().Info(result)
	result = utils.Convert2Title("abc,    def.    txt    ddd ....a. . . ....")
	logs.GetLogger().Info(result)
	result = utils.FirstLetter2Upper("StorageDealStaged,funds computed:823728125000000,funds reserved:823728125000000,funds released:823728125000000")
	logs.GetLogger().Info(result)
	result = utils.FirstLetter2Upper("abc,    def.    txt    ddd ....a. . . ....")
	logs.GetLogger().Info(result)
	result = utils.FirstLetter2Upper("deal already imported,StorageDealPublish")
	logs.GetLogger().Info(result)
	logs.GetLogger().Info("a"[1:])
}

func TestLotusClientDealInfo() {
	logs.GetLogger().Info(1e18 == math.Pow10(18))
	lotusClient, err := lotus.LotusGetClient("http://192.168.88.41:1234/rpc/v0", "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJBbGxvdyI6WyJyZWFkIiwid3JpdGUiLCJzaWduIiwiYWRtaW4iXX0.-Y4pF34RGOten6YXoau-sEMOWOEeiHwGh9u2lsl4cv8")
	if err != nil {
		logs.GetLogger().Error(err)
		return
	}
	costStatus, err := lotusClient.LotusClientGetDealInfo("bafyreihacunihjlocvxxwgumj5ugqlfry6xpwiwmovqqfmsz475rv7icei")
	if err != nil {
		logs.GetLogger().Error(err)
		return
	}
	logs.GetLogger().Info(*costStatus)
}

func TestLotusClientCalcCommP() {
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

func TestDataCap() {
	swanClient, _ := swan.SwanGetClient("http://192.168.88.41:5002", "yr0wUW37PEm1ZUtes-0NVg", "878da9defc1841dd5ab9f4dcef1ec9af", "")
	isV, _ := swanClient.CheckDatacap("t3u7pumush376xbytsgs5wabkhtadjzfydxxda2vzyasg7cimkcphswrq66j4dubbhwpnojqd3jie6ermpwvvq")
	logs.GetLogger().Info(isV)
}

func TestGenerateFile() {
	utils.GenerateFile("./", "test.txt", 2)
}

func TestGenerateUploadFile() {
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

		carFileHash, err := ipfs.IpfsUploadFileByWebApi(apiUrl, filefullpath)
		if err != nil {
			logs.GetLogger().Error(err)
			return
		}

		logs.GetLogger().Info(*carFileHash)
	default:
		logs.GetLogger().Error("only support subcommand:generate|upload")
	}
}

func TestLotusClientQeryAsk() {
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

func TestRandStr() {
	var letterRunes = []rune("abcdefghijklmnopqrstuvwxyz0123456789")
	var optionChars = "abcdefghijklmnopqrstuvwxyz0123456789"
	for i := 0; i < 100; i++ {
		randStr := utils.RandStringRunes(letterRunes, 6)
		randStr1 := utils.RandString(optionChars, 6)
		logs.GetLogger().Info(randStr, "  ", randStr1)
	}
}

func TestTask() {
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

func TestLevelDb() {
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
