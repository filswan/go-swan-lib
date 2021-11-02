package test

import (
	"strings"

	"github.com/filswan/go-swan-lib/logs"
	"github.com/filswan/go-swan-lib/utils"

	"github.com/filswan/go-swan-lib/client"
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
	response := client.HttpGet("https://jsonplaceholder.typicode.com/todos/1", "", "")
	logs.GetLogger().Info(response)

	todo := Todo{1, 2, "lorem ipsum dolor sit amet", true}
	response = client.HttpPostNoToken("https://jsonplaceholder.typicode.com/todos", todo)
	logs.GetLogger().Info(response)

	response = client.HttpPut("https://jsonplaceholder.typicode.com/todos/1", "", todo)
	logs.GetLogger().Info(response)

	title := utils.GetFieldFromJson(response, "title")
	logs.GetLogger().Info(title)

	response = client.HttpDelete("https://jsonplaceholder.typicode.com/todos/1", "", todo)
	logs.GetLogger().Info(response)
}
