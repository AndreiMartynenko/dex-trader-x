package reserves

import (
	"fmt"
	"log"
	"math/big"
	"strings"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
)

// Uniswap/SushiSwap Pair ABI
const pairABI = `[{"constant":true,"inputs":[],"name":"getReserves","outputs":[{"internalType":"uint112","name":"reserve0","type":"uint112"},{"internalType":"uint112","name":"reserve1","type":"uint112"},{"internalType":"uint32","name":"blockTimestampLast","type":"uint32"}],"payable":false,"stateMutability":"view","type":"function"}]`

// Get reserves for each pair
func FetchReserves(pairs map[string]string, client *ethclient.Client) {
	for pairAddress, pairName := range pairs {
		fmt.Printf("\n🔹 Fetching reserves for %s (%s)...\n", pairAddress, pairName)

		reserve0, reserve1, err := GetReserves(pairAddress, client)
		if err != nil {
			fmt.Printf("❌ Error fetching reserves for %s: %v\n", pairAddress, err)
			continue
		}

		fmt.Printf("✅ %s - Reserve0: %s, Reserve1: %s\n", pairName, reserve0.String(), reserve1.String())
	}
}

// GetReserves fetches reserves from Uniswap/SushiSwap pair contract
func GetReserves(pairAddress string, client *ethclient.Client) (reserve0, reserve1 *big.Int, err error) {
	pairAddr := common.HexToAddress(pairAddress)

	parsedABI, err := abi.JSON(strings.NewReader(pairABI))
	if err != nil {
		log.Fatalf("❌ Error parsing ABI: %v", err)
	}

	contract := bind.NewBoundContract(pairAddr, parsedABI, client, client, client)

	time.Sleep(200 * time.Millisecond) // Avoid rate limiting

	var result []interface{}
	err = contract.Call(nil, &result, "getReserves")
	if err != nil {
		log.Printf("❌ Error fetching reserves for %s: %v", pairAddress, err)
		return nil, nil, err
	}

	reserve0 = result[0].(*big.Int)
	reserve1 = result[1].(*big.Int)
	return reserve0, reserve1, nil

}

// GetReservesFromExchange fetches liquidity reserves for a given pair from a specific exchange
func GetReservesFromExchange(pairAddress string, client *ethclient.Client, exchange string) (*big.Int, *big.Int, error) {
	fmt.Printf("🔹 Fetching reserves for %s (%s)...\n", pairAddress, exchange)

	// Simulating different reserves for Uniswap & SushiSwap
	var reserve0, reserve1 *big.Int
	if exchange == "uniswap" {
		reserve0 = big.NewInt(9900000000000) // Simulated Uniswap Reserve0
		reserve1 = big.NewInt(5000000000000) // Simulated Uniswap Reserve1
	} else if exchange == "sushiswap" {
		reserve0 = big.NewInt(9910000000000) // Simulated SushiSwap Reserve0
		reserve1 = big.NewInt(5050000000000) // Simulated SushiSwap Reserve1
	} else {
		return nil, nil, fmt.Errorf("invalid exchange name")
	}

	return reserve0, reserve1, nil
}
