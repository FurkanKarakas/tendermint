package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"math/rand"
	"net/http"
	"strings"
	"sync"
	"time"
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
func sendTX(wg *sync.WaitGroup, portNr *string, TXSize *int, myChan chan int) {
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
		if strings.Contains(bodyString, "error") {
			fmt.Println(bodyString)
			if strings.Contains(bodyString, "timed out waiting for tx to be included in a block") {
				myChan <- 1
			}
		} else {
			myChan <- 1
		}
	}
}

// Counter counts the total valid transactions
func Counter(throughput *int, myChan chan int) {
	for {
		*throughput += <-myChan
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
	throughput := 0
	myChan := make(chan int)
	go Counter(&throughput, myChan)
	t1 := time.Now()

	for i := 0; i < *TXNr; i++ {
		wg.Add(1)
		go sendTX(&wg, portNr, TXSize, myChan)
		fmt.Println("Sleeping", *TXTime, "milliseconds...", i+1)
		time.Sleep(time.Duration(*TXTime) * time.Millisecond)
	}

	wg.Wait()
	t2 := time.Now()
	timediff := t2.Sub(t1)

	time.Sleep(1 * time.Millisecond)

	throughputReal := float64(throughput) / float64(*TXNr) * 100
	TXperS := float64(throughput) / timediff.Seconds()

	fmt.Printf("Valid TX ratio: %.2f%%\n", throughputReal)
	fmt.Printf("Total Throughput: %.2f\n", TXperS)
}
