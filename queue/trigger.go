package queue

import (
	"context"
	"time"

	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"

	log "github.com/sirupsen/logrus"
)

// NOTE: Needs a way for other goroutine to close ticker chan
func IntervalTrigger(intervalTime time.Duration, triggerChan chan string /*, quitChan chan struct{}, done context.CancelFunc*/) {
	tickerChan := make(chan bool)
	ticker := time.NewTicker(intervalTime)

	for {
		select {
		case <-tickerChan:
			return
		case <-ticker.C:
			log.Info("Fetch data triggered, time: ", time.Now())
			triggerChan <- "triggered"
		}
	}
}

// NOTE: Needs a way for other goroutine to close block header chan & close wss client
func OnNewBlockTrigger(ctx context.Context, ethWssUrl string, triggerChan chan string) {
	wssClient, err := ethclient.Dial(ethWssUrl)
	if err != nil {
		log.Fatal(err.Error())
	}

	headers := make(chan *types.Header)
	sub, err := wssClient.SubscribeNewHead(ctx, headers)
	if err != nil {
		log.Fatalf("Error on subscribing new header chain %s", err.Error())
	}

	for {
		select {
		case err := <-sub.Err():
			log.Fatal(err)

		case header := <-headers:
			log.Infof("Recieved new header: %#v", header.Number)
			triggerChan <- "trigger"
		}
	}
}
