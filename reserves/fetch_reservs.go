package reserves

import (
	"math/big"
	"strings"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
)

// ABI for Uniswap/SushiSwap Pairs
const pairABI = `[{"constant":true,"inputs":[],"name":"getReserves","outputs":[{"internalType":"uint112","name":"_reserve0","type":"uint112"},{"internalType":"uint112","name":"_reserve1","type":"uint112"},{"internalType":"uint32","name":"_blockTimestampLast","type":"uint32"}],"payable":false,"stateMutability":"view","type":"function"}]`

func GetReserves(pairAddress common.Address, client *ethclient.Client) (*big.Int, *big.Int, error) {
	parsedABI, _ := abi.JSON(strings.NewReader(pairABI))
	contract := bind.NewBoundContract(pairAddress, parsedABI, client, client, client)

	var reserves []*big.Int
	var result []interface{}
	err := contract.Call(nil, &result, "getReserves")
	if err != nil {
		return nil, nil, err
	}

	reserves = []*big.Int{result[0].(*big.Int), result[1].(*big.Int)}
	return reserves[0], reserves[1], nil
}
