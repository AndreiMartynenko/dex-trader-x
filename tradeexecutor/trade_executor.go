package main

import (
	"context"
	"crypto/ecdsa"
	"fmt"
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"

	"dex-trader-x/uniswap"
)

const (
	UniswapRouterAddress   = "0xE592427A0AEce92De3Edee1F18E0157C05861564"
	SushiswapRouterAddress = "0xd9e1CE17F2641F24aE83637ab66a2CCA9C378B9F"
)

func ExecuteSwap(client *ethclient.Client, routerAddress string, privateKeyHex string, tokenIn, tokenOut common.Address, amountIn *big.Int) error {
	privateKey, err := crypto.HexToECDSA(privateKeyHex)
	if err != nil {
		return fmt.Errorf("failed to load private key: %v", err)
	}

	publicKey := privateKey.Public()
	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		return fmt.Errorf("error casting public key to ECDSA")
	}
	fromAddress := crypto.PubkeyToAddress(*publicKeyECDSA)

	nonce, err := client.PendingNonceAt(context.Background(), fromAddress)
	if err != nil {
		return fmt.Errorf("failed to get nonce: %v", err)
	}

	gasPrice, err := client.SuggestGasPrice(context.Background())
	if err != nil {
		return fmt.Errorf("failed to get gas price: %v", err)
	}

	amountOutMin := big.NewInt(0)
	path := []common.Address{tokenIn, tokenOut}
	deadline := big.NewInt(time.Now().Unix() + 300)

	router, err := uniswap.NewUniswap(common.HexToAddress(routerAddress), client)
	if err != nil {
		return fmt.Errorf("failed to initialize Uniswap router: %v", err)
	}

	auth, err := bind.NewKeyedTransactorWithChainID(privateKey, big.NewInt(1)) // 1 = Ethereum Mainnet
	if err != nil {
		return fmt.Errorf("failed to create transactor: %v", err)
	}

	auth.Nonce = big.NewInt(int64(nonce))
	auth.GasPrice = gasPrice
	auth.GasLimit = uint64(250000) // Set gas limit

	tx, err := router.SwapExactTokensForTokens(auth, amountIn, amountOutMin, path, fromAddress, deadline)
	if err != nil {
		return fmt.Errorf("swap execution failed: %v", err)
	}

	fmt.Printf("Trade executed! TX Hash: %s\n", tx.Hash().Hex())
	return nil
}
