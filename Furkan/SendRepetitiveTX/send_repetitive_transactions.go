package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"math/rand"
	"net/http"
	"sync"
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

// Sends transactions
func sendTX(wg *sync.WaitGroup, portNr *string, TXSize *int) {
	defer wg.Done() // Decrement by 1 after function returns.
	resp, err := http.Get("http://127.0.0.1:" + *portNr + "/broadcast_tx_commit?tx=\"" + randSeq(*TXSize) + "\"")
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
}

func main() {
	rand.Seed(time.Now().UnixNano())

	portNr := flag.String("portNr", "26657", "Port numbers: 26657, 26660, 26662, and 26664. Default 26657")
	TXNr := flag.Int("TXNr", 100, "Number of transactions to send. Default: 100")
	TXSize := flag.Int("TXSize", 100, "Specify the size of the transaction in bytes. Default is 100 bytes.")
	TXTime := flag.Int("TXTime", 1000, "Time spent between transactions. In milliseconds.")
	flag.Parse()

	var wg sync.WaitGroup

	for i := 0; i < *TXNr; i++ {
		wg.Add(1)
		go sendTX(&wg, portNr, TXSize)
		fmt.Println("Sleeping", *TXTime, "milliseconds...", i+1)
		time.Sleep(time.Duration(*TXTime) * time.Millisecond)
	}

	wg.Wait()
}
