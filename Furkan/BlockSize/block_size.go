package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
)

// DecodedHTTPResponse stores the HTTP response from the server for the ABCI Info URL
type DecodedHTTPResponse struct {
	Jsonrpc string     `json:"jsonrpc"`
	ID      int        `json:"id"`
	Result  NestedJSON `json:"result"`
}

// NestedJSON stores the nested JSON object. Relevant for the DecodedHTTPResponse struct.
type NestedJSON struct {
	Response NestedJSON2 `json:"response"`
}

// NestedJSON2 stores the 3rd layer. Relevant for the NestedJSON struct.
type NestedJSON2 struct {
	Data             string `json:"data"`
	Version          string `json:"version"`
	AppVersion       string `json:"app_version"`
	LastBlockHeight  string `json:"last_block_height"`
	LastBlockAppHash string `json:"last_block_app_hash"`
}

// BlockHeight stores the HTTP response of block height URL in a struct.
type BlockHeight struct {
	Result BlockHeight1 `json:"result"`
}

//BlockHeight1 is relevant, i.e. a nested JSON inside BlockHeight.
type BlockHeight1 struct {
	Block BlockHeight2 `json:"block"`
}

//BlockHeight2 is relevant, i.e. a nested JSON inside BlockHeight1.
type BlockHeight2 struct {
	Data BlockHeight3 `json:"data"`
}

//BlockHeight3 is relevant, i.e. a nested JSON inside BlockHeight2.
type BlockHeight3 struct {
	TXs []string `json:"txs"`
}

// getJSON decodes the HTTP JSON response.
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

// getBlockSizeinTX returns the size of the block in number of Transactions in it
func getBlockSizeinTX(url string) int {
	r, err := http.Get(url)
	if err != nil {
		panic(err)
	}
	defer r.Body.Close()

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		panic(err)
	}
	var target BlockHeight
	jsonErr := json.Unmarshal(body, &target)
	if jsonErr != nil {
		panic(jsonErr)
	}
	return len(target.Result.Block.Data.TXs)

}

func main() {
	portNr := flag.String("portNr", "26657", "Port numbers: 26657, 26660, 26662, and 26664. Default 26657")
	flag.Parse()
	ABCIInfo := new(DecodedHTTPResponse)
	getJSON("http://127.0.0.1:"+*portNr+"/abci_info?", ABCIInfo)
	lastHeight := ABCIInfo.Result.Response.LastBlockHeight
	fmt.Println("Last block height is:", lastHeight)
	lastHeightasInt, err := strconv.Atoi(lastHeight)
	if err != nil {
		panic(err)
	}

	for i := 1; i <= lastHeightasInt; i++ {
		size := getBlockSizeinTX("http://127.0.0.1:" + *portNr + "/block?height=" + strconv.Itoa(i))
		fmt.Println("The number of transactions in the block with height", i, "is:", size)
	}
}
