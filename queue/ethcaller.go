package queue

import (
	"context"
	"sync"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/log"
)

var (
	wg sync.WaitGroup
	// Counter int
)

func EthCaller(ctx context.Context, client ethclient.Client, injectChan <-chan *ethereum.CallMsg, processChan chan []byte) {
	for data := range injectChan {
		_data := data
		// wg.Add(1)

		go func() {
			// defer wg.Done()
			root, err := client.CallContract(ctx, *_data, nil)
			if err != nil {
				log.Warn(err.Error())
			}
			// Counter++
			processChan <- append(root, _data.To.Bytes()...)
		}()

		// wg.Wait()
	}
}
