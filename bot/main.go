package main

import (
	"context"
	"fmt"
	"log"
	"math/big"
	"strings"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
)

const infuraURL = "https://mainnet.infura.io/v3/670cff15b19a47d6800fe3d5ef36940a"

func main() {
	client, err := ethclient.Dial(infuraURL)
	if err != nil {
		log.Fatal(err)
	}
	defer client.Close()

	fmt.Println("✅ Connected to Ethereum Mainnet")

	uniswapPair := common.HexToAddress("0xB4e16d0168e52d35CaCD2c6185b44281Ec28C9Dc") // USDC/WETH pair

	price, err := getPrice(client, uniswapPair)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("🚀 Uniswap WETH/USDC Price: 1 WETH = %.2f USDC\n", price)
}

func getPrice(client *ethclient.Client, pairAddress common.Address) (float64, error) {
	pairABI, err := abi.JSON(strings.NewReader(`[{"constant":true,"inputs":[],"name":"getReserves","outputs":[{"internalType":"uint112","name":"_reserve0","type":"uint112"},{"internalType":"uint112","name":"_reserve1","type":"uint112"},{"internalType":"uint32","name":"_blockTimestampLast","type":"uint32"}],"payable":false,"stateMutability":"view","type":"function"}]`))
	if err != nil {
		return 0, err
	}

	callData, err := pairABI.Pack("getReserves")
	if err != nil {
		return 0, err
	}

	res, err := client.CallContract(context.Background(), ethereum.CallMsg{
		To:   &pairAddress,
		Data: callData,
	}, nil)
	if err != nil {
		return 0, err
	}

	if len(res) == 0 {
		return 0, fmt.Errorf("empty response from contract call")
	}

	var reserves struct {
		Reserve0           *big.Int
		Reserve1           *big.Int
		BlockTimestampLast uint32
	}

	err = pairABI.UnpackIntoInterface(&reserves, "getReserves", res)
	if err != nil {
		return 0, err
	}

	// USDC has 6 decimals, WETH has 18 decimals
	reserveUSDC := new(big.Float).Quo(new(big.Float).SetInt(reserves.Reserve0), big.NewFloat(1e6))
	reserveWETH := new(big.Float).Quo(new(big.Float).SetInt(reserves.Reserve1), big.NewFloat(1e18))

	price := new(big.Float).Quo(reserveUSDC, reserveWETH)

	finalPrice, _ := price.Float64()

	return finalPrice, nil
}
