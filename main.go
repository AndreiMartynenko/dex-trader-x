package main

import (
	"context"
	"dex-trader-x/arbitrage"
	"dex-trader-x/pairs"
	"dex-trader-x/reserves"
	"dex-trader-x/uniswap"
	"fmt"
	"log"
	"math/big"
	"os"
	"strings"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/joho/godotenv"
)

func main() {
	// Load environment variables
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	ethRPC := os.Getenv("ALCHEMY_URL")
	if ethRPC == "" {
		log.Fatal("Missing ALCHEMY_URL in .env file")
	}

	fmt.Println("Connecting to Ethereum...")
	client, err := ethclient.Dial(ethRPC)
	if err != nil {
		log.Fatalf("Failed to connect to Ethereum: %v", err)
	}
	fmt.Println("Connected to Ethereum!")

	//routerAddress := common.HexToAddress("0x7a250d5630B4cF539739dF2C5dAcb4c659F2488D") // Uniswap V2 Router
	routerAddress := common.HexToAddress("0x7a250d5630B4cF539739dF2C5dAcb4c659F2488D")

	router, err := uniswap.NewUniswap(routerAddress, client)
	if err != nil {
		log.Fatalf("Failed to initialize Uniswap V2 Router: %v", err)
	}

	fmt.Println("\nFetching Uniswap & SushiSwap Pairs...")
	uniswapPairs := pairs.FetchUniswapPairs(client)
	sushiswapPairs := pairs.FetchSushiSwapPairs(client)

	commonPairs := pairs.FindCommonPairs(uniswapPairs, sushiswapPairs)
	if len(commonPairs) == 0 {
		fmt.Println("No common pairs found. Exiting.")
		return
	}

	fmt.Println("\nFetching Liquidity Reserves...")
	reserves.FetchReserves(commonPairs, client)

	fmt.Println("\n🔍 Checking Arbitrage Opportunities...")
	arbitrage.CheckArbitrageOpportunities(commonPairs, client, uniswapPairs, sushiswapPairs)

	executeTrade(client, router)
}

func executeTrade(client *ethclient.Client, router *uniswap.Uniswap) {
	privateKeyHex := os.Getenv("WALLET_PRIVATE_KEY")
	if privateKeyHex == "" {
		log.Fatal("Missing WALLET_PRIVATE_KEY in .env file")
	}

	privateKeyHex = strings.TrimPrefix(privateKeyHex, "0x")
	privateKey, err := crypto.HexToECDSA(privateKeyHex)
	if err != nil {
		log.Fatalf("Failed to load private key: %v", err)
	}

	auth, err := bind.NewKeyedTransactorWithChainID(privateKey, big.NewInt(1)) // Mainnet
	if err != nil {
		log.Fatalf("Failed to create transaction auth: %v", err)
	}

	auth.GasPrice, _ = client.SuggestGasPrice(context.Background())
	auth.GasLimit = uint64(300000)

	tokenIn := common.HexToAddress("0xC02aaA39b223FE8D0A0e5C4F27eAD9083C756Cc2")  // WETH
	tokenOut := common.HexToAddress("0xdAC17F958D2ee523a2206206994597C13D831ec7") // USDT
	amountIn := big.NewInt(1e18)
	amountOutMin := big.NewInt(0)
	path := []common.Address{tokenIn, tokenOut}
	deadline := big.NewInt(time.Now().Unix() + 300)

	tx, err := router.SwapExactTokensForTokens(
		auth,
		amountIn,
		amountOutMin,
		path,
		auth.From,
		deadline,
	)
	if err != nil {
		log.Fatalf("Swap failed: %v", err)
	}

	fmt.Printf("Swap executed! TX: %s\n", tx.Hash().Hex())
}
