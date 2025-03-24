package arbitrage

import (
	"dex-trader-x/reserves"
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/ethclient"
)

func CheckArbitrageOpportunities(commonPairs map[string]string, client *ethclient.Client, uniswapPairs, sushiswapPairs map[string]string) {
	for pairAddress, pairName := range commonPairs {
		fmt.Printf("\nChecking %s (%s) for arbitrage...\n", pairAddress, pairName)

		uniReserve0, uniReserve1, err := reserves.GetReservesFromExchange(pairAddress, client, "uniswap")
		if err != nil {
			fmt.Printf("Error fetching Uniswap reserves for %s: %v\n", pairAddress, err)
			continue
		}
		fmt.Printf("Uniswap %s Reserves: Reserve0 = %s, Reserve1 = %s\n", pairName, uniReserve0.String(), uniReserve1.String())

		sushiReserve0, sushiReserve1, err := reserves.GetReservesFromExchange(pairAddress, client, "sushiswap")
		if err != nil {
			fmt.Printf("Error fetching SushiSwap reserves for %s: %v\n", pairAddress, err)
			continue
		}
		fmt.Printf("SushiSwap %s Reserves: Reserve0 = %s, Reserve1 = %s\n", pairName, sushiReserve0.String(), sushiReserve1.String())

		uniPrice := calculatePrice(uniReserve0, uniReserve1)
		sushiPrice := calculatePrice(sushiReserve0, sushiReserve1)

		fmt.Printf("Uniswap %s Price = %.6f\n", pairName, uniPrice)
		fmt.Printf("SushiSwap %s Price = %.6f\n", pairName, sushiPrice)

		if checkArbitrage(uniPrice, sushiPrice) {
			fmt.Printf("Arbitrage Opportunity Found for %s! \n", pairName)
		} else {
			fmt.Printf("No arbitrage opportunity for %s.\n", pairName)
		}
	}
}

// calculatePrice computes the price as Reserve1 / Reserve0
func calculatePrice(reserve0, reserve1 *big.Int) float64 {
	dec0 := new(big.Float).SetInt(reserve0)
	dec1 := new(big.Float).SetInt(reserve1)

	result := new(big.Float).Quo(dec1, dec0) // reserve1 / reserve0
	price, _ := result.Float64()
	return price
}

func checkArbitrage(price1, price2 float64) bool {
	threshold := 0.005 // 0.5% difference
	diff := (price1 - price2) / price2

	return diff > threshold || diff < -threshold
}
