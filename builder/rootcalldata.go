package builder

import (
	"fmt"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
)

func RootCalldataSingleTx(contract *common.Address, rootPos int) *ethereum.CallMsg {
	cleanData := "c2b40ae40000000000000000000000000000000000000000000000000000000000000000"

	suffix := fmt.Sprintf("%x", rootPos)
	return &ethereum.CallMsg{
		From: common.HexToAddress("0000000000000000000000000000000000000000"),
		To:   contract,
		Data: common.Hex2Bytes(cleanData[:72-len(suffix)] + suffix),
	}
}

func RootCalldataMultipleTx(contracts [](*common.Address), rootSize []int) [](*ethereum.CallMsg) {
	result := make([]*ethereum.CallMsg, 0)
	cleanData := "c2b40ae40000000000000000000000000000000000000000000000000000000000000000"

	for i := 0; i < len(contracts); i++ {
		for j := 0; j < rootSize[i]; j++ {
			suffix := fmt.Sprintf("%x", j)
			calldata := ethereum.CallMsg{
				From: common.HexToAddress("0000000000000000000000000000000000000000"),
				To:   contracts[i],
				Data: common.Hex2Bytes(cleanData[:72-len(suffix)] + suffix),
			}
			result = append(result, &calldata)
		}
	}
	return result
}
