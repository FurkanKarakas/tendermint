package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
)

// Prints the information about application blockchain interface (ABCI) such as last block height, etc.
func main() {
	portNr := flag.String("portNr", "26657", "Port numbers: 26657, 26660, 26662, and 26664. Default 26657")
	flag.Parse()
	resp, err := http.Get("http://127.0.0.1:" + *portNr + "/abci_info?")
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusOK {
		bodybytes, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			panic(err)
		}
		bodyString := string(bodybytes)
		fmt.Println(bodyString)
	}
}
