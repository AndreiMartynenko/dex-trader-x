package pairs

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

// Uniswap & SushiSwap Factory ABI
const factoryABI = `[{"constant":true,"inputs":[],"name":"allPairsLength","outputs":[{"internalType":"uint256","name":"","type":"uint256"}],"payable":false,"stateMutability":"view","type":"function"},
{"constant":true,"inputs":[{"internalType":"uint256","name":"","type":"uint256"}],"name":"allPairs","outputs":[{"internalType":"address","name":"","type":"address"}],"payable":false,"stateMutability":"view","type":"function"}]`

func FetchPairs(factoryAddress string, client *ethclient.Client, exchange string) ([]common.Address, error) {
	parsedABI, err := abi.JSON(strings.NewReader(factoryABI))
	if err != nil {
		return nil, fmt.Errorf("Error parsing ABI: %v", err)
	}

	contract := bind.NewBoundContract(common.HexToAddress(factoryAddress), parsedABI, client, client, client)

	var pairsCount *big.Int
	var result []interface{}
	err = contract.Call(nil, &result, "allPairsLength")
	if err != nil {
		return nil, fmt.Errorf("Error fetching total pairs: %v", err)
	}
	pairsCount = result[0].(*big.Int)

	fmt.Printf("✅ Found %d pairs at factory %s (%s)\n", pairsCount, factoryAddress, exchange)

	// Inside FetchPairs function, before making each request
	time.Sleep(200 * time.Millisecond) // Adds a 200ms delay per request
	var pairs []common.Address
	for i := int64(0); i < 10 && i < pairsCount.Int64(); i++ { // Limit to 10 pairs
		var pair common.Address
		var result []interface{}
		err = contract.Call(nil, &result, "allPairs", big.NewInt(i))
		if err != nil {
			log.Printf("Error fetching pair at index %d: %v", i, err)
			continue
		}
		pair = result[0].(common.Address)
		pairs = append(pairs, pair)
		fmt.Printf("%s Pair %d: %s\n", exchange, i+1, pair.Hex())
	}
	return pairs, nil
}
