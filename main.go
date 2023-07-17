package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"
	"time"

	log "github.com/sirupsen/logrus"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"

	"github.com/108356037/torn-root-fetcher/config"
	"github.com/108356037/torn-root-fetcher/queue"
	"github.com/108356037/torn-root-fetcher/redis"
)

func main() {
	log.SetLevel(log.Level(config.LOG_LEVEL))

	redis.Init()

	client, err := ethclient.Dial(config.ETH_URL)
	if err != nil {
		log.Fatal(err.Error())
	}

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	bufferSize := 128
	triggerChan := make(chan string, 1)
	calldataChan := make(chan *ethereum.CallMsg, bufferSize)
	callResultChan := make(chan []byte, bufferSize)

	go queue.OnNewBlockTrigger(ctx, config.ETH_WSS, triggerChan)

	// go queue.TornRootQueryInjector(&config.TORN_01_ETH, 100, calldataChan, triggerChan)
	// go queue.TornRootQueryInjector(&config.TORN_1_ETH, 100, calldataChan, triggerChan)
	// go queue.TornRootQueryInjector(&config.TORN_10_ETH, 100, calldataChan, triggerChan)
	// go queue.TornRootQueryInjector(&config.TORN_100_ETH, 100, calldataChan, triggerChan)

	contracts, rootSizes :=
		[](*common.Address){
			&config.TORN_100_ETH, &config.TORN_10_ETH, &config.TORN_1_ETH, &config.TORN_01_ETH,
			/* &config.TORN_100_100_DAI, &config.TORN_10_000_DAI2 */
		},
		[](int){
			100, 100, 100, 100 /*30, 30*/}

	go queue.TornRootBatchContractQueryInjector(contracts, rootSizes, calldataChan, triggerChan)

	go queue.EthCaller(ctx, *client, calldataChan, callResultChan)

	go queue.TornRootRedisProcessor(callResultChan)

	triggerChan <- "triggered"

	// block until receives os signal
	<-ctx.Done()
	// stop acts like signal.Reset
	stop()

	_, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	log.Info("Cleaning trigger, inject, process channels")
	close(triggerChan)
	close(calldataChan)
	close(callResultChan)

	log.Info("Closing redis")
	redis.Close()
}
