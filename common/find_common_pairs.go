package common

import (
	"fmt"

	"github.com/ethereum/go-ethereum/common"
)

func FindCommonPairs(uniswapPairs, sushiPairs []common.Address) []common.Address {
	commonPairs := []common.Address{}
	lookup := make(map[string]bool)

	for _, pair := range uniswapPairs {
		lookup[pair.Hex()] = true
	}

	for _, pair := range sushiPairs {
		if lookup[pair.Hex()] {
			commonPairs = append(commonPairs, pair)
		}
	}

	return commonPairs
}

func PrintCommonPairs(commonPairs []common.Address) {
	fmt.Println("\nCommon Pairs Between Uniswap & SushiSwap:")
	for _, pair := range commonPairs {
		fmt.Println(pair.Hex())
	}
}
