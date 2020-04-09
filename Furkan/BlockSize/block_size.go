package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
)

// Foo stores the decoded HTTP response
type Foo struct {
	Jsonrpc string     `json:"jsonrpc"`
	ID      int        `json:"id"`
	Result  NestedJSON `json:"result"`
}

// NestedJSON stores the nested JSON object.
type NestedJSON struct {
	Response NestedJSON2 `json:"response"`
}

// NestedJSON2 stores the 3rd layer.
type NestedJSON2 struct {
	Data             string `json:"data"`
	Version          string `json:"version"`
	AppVersion       string `json:"app_version"`
	LastBlockHeight  string `json:"last_block_height"`
	LastBlockAppHash string `json:"last_block_app_hash"`
}

// getJson decodes the HTTP JSON response.
func getJSON(url string, target interface{}) {
	r, err := http.Get(url)
	if err != nil {
		panic(err)
	}
	defer r.Body.Close()

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		panic(err)
	}
	jsonErr := json.Unmarshal(body, target)
	if jsonErr != nil {
		panic(jsonErr)
	}
}

// getblockSize returns the size of the JSON string from the HTTP response in bytes
func getBlockSize(url string) int {
	var bodyString string
	resp, err := http.Get(url)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusOK {
		bodyBytes, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			panic(err)
		}
		bodyString = string(bodyBytes)
	}
	return len(bodyString)
}

func main() {
	portNr := flag.String("portNr", "26657", "Port numbers: 26657, 26660, 26662, and 26664. Default 26657")
	height := flag.String("height", "", "Height of the block to look at. Default is the last height")
	flag.Parse()
	if *height == "" {
		fmt.Println("No height given. Taking the last block height...")
		foo1 := new(Foo)
		getJSON("http://127.0.0.1:"+*portNr+"/abci_info?", foo1)
		*height = foo1.Result.Response.LastBlockHeight
		fmt.Println("Last block height is:", *height)
	} else {
		fmt.Println("Height argument is given:", *height)
	}

	size := getBlockSize("http://127.0.0.1:" + *portNr + "/block?height=" + *height)
	fmt.Println("Size of the block is:", size, "bytes")

}
