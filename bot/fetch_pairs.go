package main

import (
	"fmt"
	"log"
	"math/big"
	"os"
	"strings"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/joho/godotenv"
)

func main() {
	// Load .env file
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file")
	}

	infuraURL := os.Getenv("INFURA_URL")
	if infuraURL == "" {
		log.Fatalf("Error: INFURA_URL not set in .env")
	}

	fmt.Println("🚀 Connecting to Ethereum Mainnet...")
	client, err := ethclient.Dial(infuraURL)
	if err != nil {
		log.Fatalf("Error connecting to Ethereum: %v", err)
	}
	defer client.Close()

	// Fetch Uniswap V2 pairs
	fmt.Println("Fetching Uniswap V2 Pairs...")
	fetchPairs("0x5C69bEe701ef814a2B6a3EDD4B1652CB9cc5aA6f", client, "Uniswap")

	// Fetch SushiSwap V2 pairs
	fmt.Println("\nFetching SushiSwap V2 Pairs...")
	fetchPairs("0xc0aee478e3658e2610c5f7a4a2e1777ce9e4f2ac", client, "SushiSwap")
}

func fetchPairs(factoryAddress string, client *ethclient.Client, name string) {
	parsedABI, err := abi.JSON(strings.NewReader(`[{
      "constant": true,
      "inputs": [],
      "name": "allPairsLength",
      "outputs": [{"internalType": "uint256", "name": "", "type": "uint256"}],
      "payable": false,
      "stateMutability": "view",
      "type": "function"
    },
    {
      "constant": true,
      "inputs": [{"internalType": "uint256", "name": "", "type": "uint256"}],
      "name": "allPairs",
      "outputs": [{"internalType": "address", "name": "", "type": "address"}],
      "payable": false,
      "stateMutability": "view",
      "type": "function"
    }]`))
	if err != nil {
		log.Fatalf("Error parsing ABI: %v", err)
	}

	contract := bind.NewBoundContract(common.HexToAddress(factoryAddress), parsedABI, client, client, client)

	// Get total pairs count
	var pairsCount *big.Int
	var result []interface{}
	err = contract.Call(nil, &result, "allPairsLength")
	if err != nil {
		log.Fatalf("Error fetching total pairs for %s: %v", name, err)
	}
	pairsCount = result[0].(*big.Int)

	fmt.Printf("✅ Found %d pairs at factory %s (%s)\n", pairsCount, factoryAddress, name)

	// Fetch first 50 pairs with delay
	for i := int64(0); i < 10 && i < pairsCount.Int64(); i++ {
		time.Sleep(500 * time.Millisecond) // Delay to prevent rate limiting

		var pairAddress common.Address
		var result []interface{}
		err = contract.Call(nil, &result, "allPairs", big.NewInt(i))
		if err != nil {
			log.Printf("Error fetching pair at index %d from %s: %v", i, name, err)
			continue
		}
		pairAddress = result[0].(common.Address)

		fmt.Printf("%s Pair %d: %s\n", name, i+1, pairAddress.Hex())
	}
}
