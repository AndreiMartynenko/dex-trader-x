package config

import (
	"fmt"
	"log"
	"math/big"
	"os"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/joho/godotenv"
)

var (
	AlchemyURL       string
	InfuraURL        string
	UniswapFactory   string
	SushiSwapFactory string

	PrivateKey    string
	WalletAddress common.Address
	ChainID       = big.NewInt(1) // Ethereum Mainnet
)
var (
	UniswapRouter   = common.HexToAddress("0x7a250d5630b4cf539739df2c5dacb4c659f2488d") // Uniswap Router
	SushiSwapRouter = common.HexToAddress("0xd9e1CE17F2641f24aE83637ab66a2cca9C378B9F") // SushiSwap Router
)

// LoadEnv
func LoadEnv() error {
	err := godotenv.Load()
	if err != nil {
		log.Println("Warning: .env file not found, using system environment variables.")
	}

	AlchemyURL = os.Getenv("ALCHEMY_URL")
	InfuraURL = os.Getenv("INFURA_URL")

	UniswapFactory = os.Getenv("UNISWAP_FACTORY")
	SushiSwapFactory = os.Getenv("SUSHISWAP_FACTORY")

	PrivateKey = os.Getenv("WALLET_PRIVATE_KEY")
	if PrivateKey == "" {
		log.Fatal("WALLET_PRIVATE_KEY is missing in .env")
	}

	privateKeyECDSA, err := crypto.HexToECDSA(PrivateKey)
	if err != nil {
		log.Fatal("Invalid Private Key:", err)
	}

	WalletAddress = crypto.PubkeyToAddress(privateKeyECDSA.PublicKey)

	if UniswapFactory == "" || SushiSwapFactory == "" {
		log.Fatal("Missing UNISWAP_FACTORY or SUSHISWAP_FACTORY in environment variables!")
	}

	if AlchemyURL == "" && InfuraURL == "" {
		log.Fatal("Missing ALCHEMY_URL and INFURA_URL in environment variables. At least one is required!")
	}

	return nil
}
func GetEthereumClient() (*ethclient.Client, error) {
	var client *ethclient.Client
	var err error

	if AlchemyURL != "" {
		client, err = ethclient.Dial(AlchemyURL)
		if err == nil {
			fmt.Println("Connected to Ethereum via Alchemy")
			return client, nil
		}
		log.Println("Alchemy connection failed. Trying Infura...")
	}

	if InfuraURL != "" {
		client, err = ethclient.Dial(InfuraURL)
		if err == nil {
			fmt.Println("Connected to Ethereum via Infura")
			return client, nil
		}
		log.Println("Both Alchemy and Infura failed to connect!")
	}

	return nil, fmt.Errorf("unable to connect to Ethereum via Alchemy or Infura")
}
