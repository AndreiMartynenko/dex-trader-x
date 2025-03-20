package arbitrage

import (
	"fmt"
	"log"
	"math/big"

	"dex-trader-x/reserves"

	"github.com/ethereum/go-ethereum/ethclient"
)

// CheckArbitrage scans for profitable opportunities
func CheckArbitrage(pairs map[string]string, client *ethclient.Client) {
	fmt.Println("\n🔍 Checking for Arbitrage Opportunities...")

	for pair, sushiPair := range pairs {
		uniReserve0, uniReserve1, err := reserves.GetReserves(pair, client)
		if err != nil {
			log.Printf("Skipping pair %s due to error fetching Uniswap reserves\n", pair)
			continue
		}

		sushiReserve0, sushiReserve1, err := reserves.GetReserves(sushiPair, client)
		if err != nil {
			log.Printf("Skipping pair %s due to error fetching SushiSwap reserves\n", sushiPair)
			continue
		}

		// Calculate Uniswap and SushiSwap prices
		uniPrice := new(big.Float).Quo(new(big.Float).SetInt(uniReserve1), new(big.Float).SetInt(uniReserve0))
		sushiPrice := new(big.Float).Quo(new(big.Float).SetInt(sushiReserve1), new(big.Float).SetInt(sushiReserve0))

		priceDiff := new(big.Float).Sub(uniPrice, sushiPrice)
		fmt.Printf("💱 Pair: %s | Uniswap Price: %f | SushiSwap Price: %f | Difference: %f\n",
			pair, uniPrice, sushiPrice, priceDiff)

		// Arbitrage condition: If price difference is > 1% and covers gas fees
		gasCost := big.NewFloat(0.001)                                // Example: 0.001 ETH gas cost
		threshold := new(big.Float).Mul(uniPrice, big.NewFloat(0.01)) // 1% threshold
		if priceDiff.Cmp(threshold) > 0 && priceDiff.Cmp(gasCost) > 0 {
			fmt.Printf("🚀 Arbitrage Opportunity Found! Buy from SushiSwap, Sell on Uniswap\n")
		}
	}
}
