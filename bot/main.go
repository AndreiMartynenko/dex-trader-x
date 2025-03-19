package main

import (
	"context"
	"crypto/ecdsa"
	"fmt"
	"log"
	"math/big"
	"os"
	"strings"
	"time"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/joho/godotenv"
)

func main() {
	// Load .env file
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	// Get Infura API Key and Private Key from .env
	infuraURL := "https://sepolia.infura.io/v3/" + os.Getenv("INFURA_API_KEY")
	privateKey := os.Getenv("PRIVATE_KEY")
	contractAddress := "0x7308C5edb894c5027ed0eAA536dC4EeCAab7237C" // Sepolia contract

	client, err := ethclient.Dial(infuraURL)
	if err != nil {
		log.Fatal(err)
	}
	defer client.Close()

	privateKeyECDSA, err := crypto.HexToECDSA(privateKey)
	if err != nil {
		log.Fatal(err)
	}

	auth, err := bind.NewKeyedTransactorWithChainID(privateKeyECDSA, big.NewInt(11155111)) // Sepolia Chain ID
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("✅ Connected to Ethereum Sepolia Testnet")

	// Correct Sepolia Uniswap & SushiSwap WETH/USDC trading pairs
	uniswapPair := common.HexToAddress("0x3A9D48AB9751398BbFa63ad67599Bb04e4BdF98b")
	sushiswapPair := common.HexToAddress("0xYourSepoliaSushiSwapPairAddress")

	for {
		uniPrice, err := getPrice(client, uniswapPair)
		if err != nil {
			log.Println("Uniswap Error:", err)
			continue
		}

		sushiPrice, err := getPrice(client, sushiswapPair)
		if err != nil {
			log.Println("SushiSwap Error:", err)
			continue
		}

		fmt.Printf("\n🚀 Uniswap: %.2f USDC | 🍣 SushiSwap: %.2f USDC\n", uniPrice, sushiPrice)

		diff := uniPrice - sushiPrice
		if diff < 0 {
			diff = -diff
		}

		percentage := (diff / sushiPrice) * 100
		fmt.Printf("🔍 Difference: %.2f%%\n", percentage)

		if percentage > 0.3 {
			fmt.Println("✅ Profitable arbitrage detected!")
			err := triggerArbitrage(auth, client, contractAddress, privateKeyECDSA)
			if err != nil {
				log.Println("Error executing trade:", err)
			} else {
				fmt.Println("🎉 Arbitrage trade executed successfully!")
			}
		} else {
			fmt.Println("❌ No profitable opportunity currently.")
		}

		time.Sleep(30 * time.Second)
	}
}

func getPrice(client *ethclient.Client, pairAddress common.Address) (float64, error) {
	pairABI, _ := abi.JSON(strings.NewReader(`[{"constant":true,"inputs":[],"name":"getReserves","outputs":[{"internalType":"uint112","name":"_reserve0","type":"uint112"},{"internalType":"uint112","name":"_reserve1","type":"uint112"},{"internalType":"uint32","name":"_blockTimestampLast","type":"uint32"}],"payable":false,"stateMutability":"view","type":"function"}]`))
	callData, _ := pairABI.Pack("getReserves")

	res, err := client.CallContract(context.Background(), ethereum.CallMsg{To: &pairAddress, Data: callData}, nil)
	if err != nil || len(res) == 0 {
		return 0, fmt.Errorf("empty response or error")
	}

	var reserves struct {
		Reserve0, Reserve1 *big.Int
		BlockTimestampLast uint32
	}

	pairABI.UnpackIntoInterface(&reserves, "getReserves", res)

	reserveUSDC := new(big.Float).Quo(new(big.Float).SetInt(reserves.Reserve0), big.NewFloat(1e6))
	reserveWETH := new(big.Float).Quo(new(big.Float).SetInt(reserves.Reserve1), big.NewFloat(1e18))

	price, _ := new(big.Float).Quo(reserveUSDC, reserveWETH).Float64()
	return price, nil
}

func triggerArbitrage(auth *bind.TransactOpts, client *ethclient.Client, contractAddress string, privateKeyECDSA *ecdsa.PrivateKey) error {
	contractABI, err := abi.JSON(strings.NewReader(`[{"inputs":[{"internalType":"address","name":"router1","type":"address"},{"internalType":"address","name":"router2","type":"address"},{"internalType":"address","name":"token1","type":"address"},{"internalType":"address","name":"token2","type":"address"},{"internalType":"uint256","name":"amountIn","type":"uint256"}],"name":"executeTrade","outputs":[],"stateMutability":"nonpayable","type":"function"}]`))
	if err != nil {
		return err
	}

	contractAddr := common.HexToAddress(contractAddress)
	router1 := common.HexToAddress("0xYourSepoliaUniswapRouterAddress")
	router2 := common.HexToAddress("0xYourSepoliaSushiSwapRouterAddress")
	token1 := common.HexToAddress("0xYourSepoliaUSDCAddress")
	token2 := common.HexToAddress("0xYourSepoliaWETHAddress")
	amountIn := big.NewInt(1000000) // 1 USDC (6 decimals)

	inputData, err := contractABI.Pack("executeTrade", router1, router2, token1, token2, amountIn)
	if err != nil {
		return err
	}

	nonce, err := client.PendingNonceAt(context.Background(), auth.From)
	if err != nil {
		return err
	}

	gasPrice, err := client.SuggestGasPrice(context.Background())
	if err != nil {
		return err
	}

	// ✅ Dynamically Get the Chain ID
	chainID, err := client.NetworkID(context.Background())
	if err != nil {
		return fmt.Errorf("failed to get network ID: %v", err)
	}

	// ✅ Create and Sign the Transaction with `chainID`
	tx := types.NewTransaction(nonce, contractAddr, big.NewInt(0), uint64(300000), gasPrice, inputData)
	signedTx, err := types.SignTx(tx, types.NewEIP155Signer(chainID), privateKeyECDSA)
	if err != nil {
		return fmt.Errorf("failed to sign transaction: %v", err)
	}

	// ✅ Send the Signed Transaction
	err = client.SendTransaction(context.Background(), signedTx)
	if err != nil {
		return fmt.Errorf("failed to send transaction: %v", err)
	}

	fmt.Printf("✅ Transaction sent successfully: %s\n", signedTx.Hash().Hex())
	return nil
}
