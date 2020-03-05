package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"math/rand"
	"net/http"
	"time"

	"github.com/prometheus/common/log"
)

var randChars = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")

// Generates a random string for transactions.
func randSeq(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = randChars[rand.Intn(len(randChars))]
	}
	return string(b)
}

func main() {
	rand.Seed(time.Now().UnixNano())

	portNr := flag.String("portNr", "26657", "Port numbers: 26657, 26660, 26662, and 26664. Default 26657")
	TXNr := flag.Int("TXNr", 100, "Number of transactions to send. Default: 100")
	TXKeySize := flag.Int("kSize", 100, "Specify the size of the transaction (key) in bytes. Default is 100 bytes.")
	TXValueSize := flag.Int("vSize", 100, "Specify the size of the transaction (value) in bytes. Default is 100 bytes.")
	TXTime := flag.Int("TXTime", 1000, "Time spent between transactions. In milliseconds.")
	flag.Parse()

	for i := 0; i < *TXNr; i++ {
		resp, err := http.Get("http://127.0.0.1:" + *portNr + "/broadcast_tx_commit?tx=\"" + randSeq(*TXKeySize) +
			"=" + randSeq(*TXValueSize) + "\"")
		if err != nil {
			panic(err)
		}
		defer resp.Body.Close()
		if resp.StatusCode == http.StatusOK {
			bodyBytes, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				panic(err)
			}
			bodyString := string(bodyBytes)
			log.Info(bodyString)
		}
		fmt.Println("Sleeping", *TXTime, "milliseconds...", i+1)
		time.Sleep(time.Duration(*TXTime) * time.Millisecond)
	}
}
