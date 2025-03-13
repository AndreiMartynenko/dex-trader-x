package exchange

import (
    "context"
    "fmt"
    "log"
    "math/big"

    "github.com/AndreiMartynenko/dex-trader-x/config"
    "github.com/ethereum/go-ethereum/rpc"
)

const uniswapV3PoolAddress = "0x8ad599c3A0FF1De082011EFDDC58f1908eb6e6D8" // ETH/USDC Pool

// FetchSqrtPriceX96 fetches Uniswap V3 `sqrtPriceX96` from the ETH/USDC pool
func FetchSqrtPriceX96() (*big.Int, error) {
    infuraURL := "https://mainnet.infura.io/v3/" + config.InfuraAPIKey
    fmt.Println("✅ Connecting to Infura at:", infuraURL) // Debugging log

    client, err := rpc.Dial(infuraURL)
    if err != nil {
        log.Fatal("❌ Failed to connect to Infura:", err)
    }

    data := "0x3850c7bd" // Function signature for `slot0()`
    var result string
    err = client.CallContext(context.Background(), &result, "eth_call", map[string]interface{}{
        "to":   uniswapV3PoolAddress,
        "data": data,
    }, "latest")

    if err != nil {
        log.Fatal("❌ Error fetching Uniswap price:", err)
    }

    // Extract sqrtPriceX96
    sqrtPriceX96Hex := result[2:66] // Get first 64 hex characters (32 bytes)
    sqrtPriceX96 := new(big.Int)
    sqrtPriceX96.SetString(sqrtPriceX96Hex, 16)

    fmt.Println("✅ Extracted sqrtPriceX96:", sqrtPriceX96.String())
    return sqrtPriceX96, nil
}

// ConvertSqrtPriceToETHPrice converts sqrtPriceX96 to ETH price in USDC
func ConvertSqrtPriceToETHPrice(sqrtPriceX96 *big.Int) *big.Float {
    // Convert sqrtPriceX96 to big.Float
    sqrtPrice := new(big.Float).SetPrec(256).SetInt(sqrtPriceX96)

    // Step 1: Square the sqrtPriceX96
    squaredPrice := new(big.Float).Mul(sqrtPrice, sqrtPrice)
    fmt.Println("✅ Squared Price:", squaredPrice.Text('f', 6))

    // Step 2: Divide by 2^192
    divisor := new(big.Float).SetPrec(256).SetInt(new(big.Int).Exp(big.NewInt(2), big.NewInt(192), nil))
    priceFloat := new(big.Float).Quo(squaredPrice, divisor)
    fmt.Println("✅ Price After Division by 2^192:", priceFloat.Text('f', 6))

    // Step 3: Adjust for USDC Decimals (USDC has 6 decimal places)
    usdcDecimals := new(big.Float).SetFloat64(1e6) // 1 USDC = 10^6
    priceFloat = new(big.Float).Quo(priceFloat, usdcDecimals)
    fmt.Println("✅ Adjusted ETH Price (USDC Decimals):", priceFloat.Text('f', 6))

    // Step 4: **Fix the Inversion**
    // If ETH price is still too low (<1 USDC per ETH), it’s stored as USDC/ETH → Invert it
    threshold := big.NewFloat(1) // If ETH price < $1, it's likely stored as USDC/ETH
    if priceFloat.Cmp(threshold) < 0 {
        priceFloat = new(big.Float).Quo(big.NewFloat(1), priceFloat)
        fmt.Println("✅ Fixed Inverted Price (ETH/USDC):", priceFloat.Text('f', 6))
    }

    return priceFloat
}