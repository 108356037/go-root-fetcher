package queue

import (
	"github.com/108356037/torn-root-fetcher/builder"
	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
)

func TornRootSingleContractQueryInjector(contract *common.Address, rootSize int, injectChan chan *ethereum.CallMsg, triggerChan <-chan string) {
	for t := range triggerChan {
		if t == "" {
			continue
		} else {
			for i := 0; i < rootSize; i++ {
				injectChan <- builder.RootCalldataSingleTx(contract, i)
			}
		}
	}
}

func TornRootMutipleContractsQueryInjector(contracts [](*common.Address), rootSizes []int, injectChan chan *ethereum.CallMsg, triggerChan <-chan string) {
	for t := range triggerChan {
		if t == "" {
			continue
		} else {
			txs := builder.RootCalldataMultipleTx(contracts, rootSizes)
			for _, tx := range txs {
				injectChan <- tx
			}
		}
	}
}
