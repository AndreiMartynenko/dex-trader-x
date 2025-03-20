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
const pairABI = `[{"constant":true,"inputs":[],"name":"getReserves","outputs":[{"internalType":"uint112","name":"reserve0","type":"uint112"},{"internalType":"uint112","name":"reserve1","type":"uint112"},{"internalType":"uint32","name":"blockTimestampLast","type":"uint32"}],"payable":false,"stateMutability":"view","type":"function"},
{"constant":true,"inputs":[],"name":"token0","outputs":[{"internalType":"address","name":"","type":"address"}],"payable":false,"stateMutability":"view","type":"function"},
{"constant":true,"inputs":[],"name":"token1","outputs":[{"internalType":"address","name":"","type":"address"}],"payable":false,"stateMutability":"view","type":"function"}]`

// GetReserves fetches the reserves of a pair contract
func GetReserves(pairAddress string, client *ethclient.Client) (reserve0, reserve1 *big.Int, err error) {
	pairAddr := common.HexToAddress(pairAddress)

	parsedABI, err := abi.JSON(strings.NewReader(pairABI))
	if err != nil {
		log.Fatalf("Error parsing pair ABI: %v", err)
	}

	contract := bind.NewBoundContract(pairAddr, parsedABI, client, client, client)

	// Inside GetReserves function, before making each request
	time.Sleep(200 * time.Millisecond) // Adds a 200ms delay per request
	var result []interface{}
	err = contract.Call(nil, &result, "getReserves")
	if err != nil {
		log.Printf("Error fetching reserves for pair %s: %v", pairAddress, err)
		return nil, nil, err
	}

	reserve0 = result[0].(*big.Int)
	reserve1 = result[1].(*big.Int)
	fmt.Printf("🔹 Reserves for pair %s: Reserve0 = %s, Reserve1 = %s\n", pairAddress, reserve0.String(), reserve1.String())

	return reserve0, reserve1, nil
}
