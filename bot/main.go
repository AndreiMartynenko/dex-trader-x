package main

// import (
// 	"context"
// 	"fmt"
// 	"log"
// 	"math/big"
// 	"os"
// 	"strings"
// 	"time"

// 	"github.com/ethereum/go-ethereum"
// 	"github.com/ethereum/go-ethereum/accounts/abi"
// 	"github.com/ethereum/go-ethereum/accounts/abi/bind"
// 	"github.com/ethereum/go-ethereum/common"
// 	"github.com/ethereum/go-ethereum/core/types"
// 	"github.com/ethereum/go-ethereum/crypto"
// 	"github.com/ethereum/go-ethereum/ethclient"
// 	"github.com/joho/godotenv"
// )

// func main() {
// 	// Load .env file
// 	err := godotenv.Load()
// 	if err != nil {
// 		log.Fatal("Error loading .env file")
// 	}

// 	// Get Infura API Key and Private Key from .env
// 	infuraURL := "https://sepolia.infura.io/v3/" + os.Getenv("INFURA_API_KEY")
// 	privateKey := os.Getenv("PRIVATE_KEY")

// 	client, err := ethclient.Dial(infuraURL)
// 	if err != nil {
// 		log.Fatal("Failed to connect to Ethereum node:", err)
// 	}
// 	defer client.Close()

// 	privateKeyECDSA, err := crypto.HexToECDSA(privateKey)
// 	if err != nil {
// 		log.Fatal("Invalid private key:", err)
// 	}

// 	// auth, err := bind.NewKeyedTransactorWithChainID(privateKeyECDSA, big.NewInt(11155111)) // Sepolia Chain ID
// 	auth, err := bind.NewKeyedTransactorWithChainID(privateKeyECDSA, big.NewInt(31337)) // Anvil Chain ID

// 	if err != nil {
// 		log.Fatal("Failed to create transactor:", err)
// 	}

// 	fmt.Println("✅ Connected to Sepolia Testnet")

// 	// Updated Pair Addresses for Sepolia
// 	uniswapPair := common.HexToAddress("0x5AcC9C21F85A4c6482554b20F67b6BdB17BfCda3")
// 	sushiswapPair := common.HexToAddress("0x18fE0d3AdF63cDbf7FDE6B69D77D2DFe0CA89eC1")

// 	for {
// 		uniPrice, err := getPrice(client, uniswapPair)
// 		if err != nil {
// 			log.Printf("Uniswap Error: %v (Pair Address: %s)\n", err, uniswapPair.Hex())
// 			continue
// 		}

// 		sushiPrice, err := getPrice(client, sushiswapPair)
// 		if err != nil {
// 			log.Printf("SushiSwap Error: %v (Pair Address: %s)\n", err, sushiswapPair.Hex())
// 			continue
// 		}

// 		fmt.Printf("\n🚀 Uniswap V2 (UNI/WETH): %.6f\n", uniPrice)
// 		fmt.Printf("🍣 SushiSwap V2 (UNI/WETH): %.6f\n", sushiPrice)

// 		diff := uniPrice - sushiPrice
// 		if diff < 0 {
// 			diff = -diff
// 		}

// 		percentage := (diff / sushiPrice) * 100
// 		fmt.Printf("🔍 Difference: %.6f%%\n", percentage)

// 		if percentage > 0.3 {
// 			fmt.Println("✅ Profitable arbitrage detected!")
// 			err := triggerArbitrage(auth, client, uniswapPair, sushiswapPair, uniPrice, sushiPrice)
// 			if err != nil {
// 				log.Println("Error executing trade:", err)
// 			} else {
// 				fmt.Println("🎉 Arbitrage trade executed successfully!")
// 			}
// 		} else {
// 			fmt.Println("❌ No profitable opportunity currently.")
// 		}

// 		time.Sleep(30 * time.Second)
// 	}
// }

// func getPrice(client *ethclient.Client, pairAddress common.Address) (float64, error) {
// 	pairABI, _ := abi.JSON(strings.NewReader(`[{"constant":true,"inputs":[],"name":"getReserves","outputs":[
//         {"internalType":"uint112","name":"_reserve0","type":"uint112"},
//         {"internalType":"uint112","name":"_reserve1","type":"uint112"},
//         {"internalType":"uint32","name":"_blockTimestampLast","type":"uint32"}
//     ],"payable":false,"stateMutability":"view","type":"function"}]`))

// 	callData, _ := pairABI.Pack("getReserves")

// 	res, err := client.CallContract(context.Background(), ethereum.CallMsg{To: &pairAddress, Data: callData}, nil)
// 	if err != nil {
// 		return 0, fmt.Errorf("Error calling getReserves: %v", err)
// 	}

// 	if len(res) == 0 {
// 		return 0, fmt.Errorf("Empty response from getReserves at address: %s", pairAddress.Hex())
// 	}

// 	fmt.Printf("Raw response bytes: %x\n", res)

// 	var reserves struct {
// 		Reserve0, Reserve1 *big.Int
// 		BlockTimestampLast uint32
// 	}

// 	err = pairABI.UnpackIntoInterface(&reserves, "getReserves", res)
// 	if err != nil {
// 		return 0, fmt.Errorf("Error unpacking reserves: %v", err)
// 	}

// 	reserveUNI := new(big.Float).SetInt(reserves.Reserve0)
// 	reserveWETH := new(big.Float).SetInt(reserves.Reserve1)

// 	price := new(big.Float).Quo(reserveWETH, reserveUNI)
// 	priceFloat, _ := price.Float64()
// 	return priceFloat, nil
// }

// func triggerArbitrage(auth *bind.TransactOpts, client *ethclient.Client, uniswapPair, sushiswapPair common.Address, uniPrice, sushiPrice float64) error {
// 	contractABI, err := abi.JSON(strings.NewReader(`[{"inputs":[{"internalType":"address","name":"router1","type":"address"},{"internalType":"address","name":"router2","type":"address"},{"internalType":"address","name":"token1","type":"address"},{"internalType":"address","name":"token2","type":"address"},{"internalType":"uint256","name":"amountIn","type":"uint256"}],"name":"executeTrade","outputs":[],"stateMutability":"nonpayable","type":"function"}]`))
// 	if err != nil {
// 		return err
// 	}

// 	contractAddr := common.HexToAddress("0x5C7f71C15d45Af5d9C76B5a5Ee4Bc1c0973C6A5D") // Updated Arbitrage Contract for Sepolia

// 	router1 := common.HexToAddress("0x7a250d5630B4cf539739df2C5dAcb4c659F2488D") // Uniswap V2 Router
// 	router2 := common.HexToAddress("0xd9e1Ce17F2641f24AE83637ab66a2CCA9C378B9F") // SushiSwap Router
// 	token1 := common.HexToAddress("0x3e4B6140A546db48Ff5bF09b0B022Af1Eb3A5680")  // UNI Token for Sepolia
// 	token2 := common.HexToAddress("0xC02Aa39b223FE8d0a0e5c4F27Ead9083C756Cc2")   // WETH Token
// 	amountIn := big.NewInt(1000000000000000000)                                  // 1 UNI (18 decimals)

// 	inputData, err := contractABI.Pack("executeTrade", router1, router2, token1, token2, amountIn)
// 	if err != nil {
// 		return err
// 	}

// 	nonce, err := client.PendingNonceAt(context.Background(), auth.From)
// 	if err != nil {
// 		return err
// 	}

// 	gasPrice, err := client.SuggestGasPrice(context.Background())
// 	if err != nil {
// 		return err
// 	}

// 	tx := types.NewTransaction(nonce, contractAddr, big.NewInt(0), uint64(300000), gasPrice, inputData)

// 	signedTx, err := auth.Signer(auth.From, tx)
// 	if err != nil {
// 		return err
// 	}

// 	err = client.SendTransaction(context.Background(), signedTx)
// 	if err != nil {
// 		return err
// 	}

// 	fmt.Printf("✅ Transaction sent successfully: %s\n", signedTx.Hash().Hex())
// 	return nil
// }
