package pairs

import (
	"fmt"
	"strings"

	"github.com/ethereum/go-ethereum/ethclient"
)

func FetchUniswapPairs(client *ethclient.Client) map[string]string {
	// Uniswap Pairs
	uniswapPairs := map[string]string{
		"0xB4e16d0168e52d35CaCD2c6185b44281Ec28C9Dc": "USDC/WETH",
		"0xA478c2975Ab1Ea89e8196811F51A7B7Ade33eB11": "DAI/WETH",
		"0xBb2b8038a1640196FbE3e38816F3e67Cba72D940": "WBTC/WETH",
		"0xd3d2E2692501A5c9Ca623199D38826e513033a17": "UNI/WETH",
		"0xa2107faF4604931f62f3C5ac62c1fE0896f4c4b9": "LINK/WETH",
	}

	fmt.Println("Uniswap Pairs Fetched:", uniswapPairs)
	return uniswapPairs
}

// Fetch SushiSwap pairs
func FetchSushiSwapPairs(client *ethclient.Client) map[string]string {
	// SushiSwap Pairs
	sushiswapPairs := map[string]string{
		"0x397ff1542f962076d0bfe58ea045ffa2d347aca0": "USDC/WETH",
		"0xc3d3a24a77afd62f1f5d2f4bc6ca2c1c3a317790": "DAI/WETH",
		"0xceff51756c56ceffca006cd410b03ffc46dd3a58": "WBTC/WETH",
		"0xd71ecff9342a5ced620049e616c5035f1db98620": "UNI/WETH",
		"0x3954a503bf87f49443af1e37e732ab2a0456a41c": "LINK/WETH",
	}

	fmt.Println("SushiSwap Pairs Fetched:", sushiswapPairs)
	return sushiswapPairs
}

// Find common pairs
// Find common pairs by comparing names instead of addresses
func FindCommonPairs(uniswapPairs, sushiswapPairs map[string]string) map[string]string {
	commonPairs := make(map[string]string)

	fmt.Println("\nFinding Common Pairs...")

	for uniAddr, uniName := range uniswapPairs {
		for sushiAddr, sushiName := range sushiswapPairs {
			if strings.EqualFold(uniName, sushiName) {
				commonPairs[uniAddr+"_"+sushiAddr] = uniName
				fmt.Printf("Common Pair Found: %s (%s)\n", uniName, uniAddr+" / "+sushiAddr)
			}
		}
	}

	fmt.Printf("Found %d common pairs!\n", len(commonPairs))
	return commonPairs
}
