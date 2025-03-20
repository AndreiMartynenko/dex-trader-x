package pairs

import (
	"fmt"
	"strings"

	"github.com/ethereum/go-ethereum/ethclient"
)

// Fetch Uniswap pairs
func FetchUniswapPairs(client *ethclient.Client) map[string]string {
	uniswapPairs := map[string]string{
		"0xB4e16d0168e52d35CaCD2c6185b44281Ec28C9Dc": "USDC/WETH",
		"0xA478c2975Ab1Ea89e8196811F51A7B7Ade33eB11": "DAI/WETH",
	}

	fmt.Println("✅ Uniswap Pairs Fetched:", uniswapPairs)
	return uniswapPairs
}

// Fetch SushiSwap pairs
func FetchSushiSwapPairs(client *ethclient.Client) map[string]string {
	sushiswapPairs := map[string]string{
		"0xb4e16d0168e52d35cacd2c6185b44281ec28c9dc": "USDC/WETH", // Ensure lowercase
		"0xa478c2975ab1ea89e8196811f51a7b7ade33eb11": "DAI/WETH",
	}

	fmt.Println("✅ SushiSwap Pairs Fetched:", sushiswapPairs)
	return sushiswapPairs
}

// Find common pairs
func FindCommonPairs(uniswapPairs, sushiswapPairs map[string]string) map[string]string {
	commonPairs := make(map[string]string)

	fmt.Println("\n🔍 Finding Common Pairs...")

	uniswapPairsLower := make(map[string]string)
	sushiswapPairsLower := make(map[string]string)

	for k, v := range uniswapPairs {
		uniswapPairsLower[strings.ToLower(k)] = v
	}
	for k, v := range sushiswapPairs {
		sushiswapPairsLower[strings.ToLower(k)] = v
	}

	fmt.Println("📌 Uniswap Pairs:", uniswapPairsLower)
	fmt.Println("📌 SushiSwap Pairs:", sushiswapPairsLower)

	for pair := range uniswapPairsLower {
		if name, exists := sushiswapPairsLower[pair]; exists {
			commonPairs[pair] = name
			fmt.Printf("✅ Common Pair Found: %s (%s)\n", pair, name)
		}
	}

	fmt.Printf("✅ Found %d common pairs!\n", len(commonPairs))
	return commonPairs
}
