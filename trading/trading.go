package trading

import (
	"context"
	"dex-trader-x/config"
	"dex-trader-x/uniswap"
	"fmt"
	"log"
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
)

// resolvePath dynamically determines the token path based on the given pair
func resolvePath(pair string) []common.Address {
	switch pair {
	case "ETH/USDC":
		return []common.Address{
			common.HexToAddress("0xC02aaa39b223FE8D0A0e5C4F27eAD9083C756Cc2"), // WETH
			common.HexToAddress("0xA0b86991c6218b36c1d19d4a2e9eb0ce3606eb48"), // USDC
		}
	case "ETH/DAI":
		return []common.Address{
			common.HexToAddress("0xC02aaa39b223FE8D0A0e5C4F27eAD9083C756Cc2"), // WETH
			common.HexToAddress("0x6B175474E89094C44Da98b954EedeAC495271d0F"), // DAI
		}
	default:
		log.Fatalf("❌ Unknown trading pair: %s", pair)
		return nil
	}
}

// ExecuteTrade handles token swaps on a decentralized exchange
func ExecuteTrade(pair string, buyExchange, sellExchange string, client *ethclient.Client) {
	fmt.Printf("⚡ Executing trade on %s -> %s for pair %s\n", buyExchange, sellExchange, pair)

	// Load private key
	privateKey, err := crypto.HexToECDSA(config.PrivateKey)
	if err != nil {
		log.Fatalf("❌ Invalid private key: %v", err)
	}

	// Create a transactor with the given chain ID
	auth, err := bind.NewKeyedTransactorWithChainID(privateKey, config.ChainID)
	if err != nil {
		log.Fatalf("❌ Error creating transactor: %v", err)
	}

	// Determine the correct router address
	var router common.Address
	if buyExchange == "Uniswap" {
		router = config.UniswapRouter
	} else if buyExchange == "SushiSwap" {
		router = config.SushiSwapRouter
	} else {
		log.Fatalf("❌ Unknown buy exchange: %s", buyExchange)
		return
	}

	// Initialize the router contract
	routerContract, err := uniswap.NewUniswapV2Router02(router, client) // Adjust to match the actual constructor in the uniswap package
	if err != nil {
		log.Fatalf("❌ Error loading router contract: %v", err)
	}

	// Dynamically resolve token path based on the trading pair
	path := resolvePath(pair)
	amountIn := big.NewInt(1000000000000000000)     // Example: 1 WETH (18 decimals)
	amountOutMin := big.NewInt(0)                   // Minimal slippage (adjustable)
	deadline := big.NewInt(time.Now().Unix() + 300) // Transaction deadline: 5 minutes

	// Ensure gas price is set
	gasPrice, err := client.SuggestGasPrice(context.Background())
	if err != nil {
		log.Fatalf("❌ Error getting gas price: %v", err)
	}
	auth.GasPrice = gasPrice

	// Perform the token swap
	tx, err := routerContract.SwapExactTokensForTokens(
		auth, amountIn, amountOutMin, path, config.WalletAddress, deadline,
	)
	if err != nil {
		log.Fatalf("❌ Failed to swap tokens via %s router: %v", buyExchange, err)
	}

	fmt.Printf("✅ Trade executed! Transaction Hash: %s\n", tx.Hash().Hex())
}
