// package config

// import (
// 	"log"
// 	"os"

// 	"github.com/joho/godotenv"
// )

// // Global Config Variables
// var (
// 	InfuraURL        string
// 	PrivateKey       string
// 	UniswapFactory   = "0x5C69bEe701ef814a2B6a3EDD4B1652CB9cc5A"
// 	SushiSwapFactory = "0xc0aee478e3658e2610c5f7a4a2e1777ce9e4f2ac"
// )

// // Load Environment Variables
// func LoadEnv() {
// 	err := godotenv.Load(".env")
// 	if err != nil {
// 		log.Fatalf("Error loading .env file: %v", err)
// 	}

// 	InfuraURL = os.Getenv("INFURA_URL")
// 	PrivateKey = os.Getenv("PRIVATE_KEY")

// 	if InfuraURL == "" || PrivateKey == "" {
// 		log.Fatal("Missing essential environment variables. Check .env file")
// 	}
// }

package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

var (
	InfuraURL        string
	UniswapFactory   string
	SushiSwapFactory string
)

// LoadEnv loads environment variables from the `.env` file
func LoadEnv() error {
	err := godotenv.Load()
	if err != nil {
		log.Printf("Warning: No .env file found or could not be loaded: %v", err)
		// Not a fatal error, return nil
		return nil
	}

	InfuraURL = os.Getenv("INFURA_URL")
	UniswapFactory = os.Getenv("UNISWAP_FACTORY")
	SushiSwapFactory = os.Getenv("SUSHISWAP_FACTORY")

	if InfuraURL == "" || UniswapFactory == "" || SushiSwapFactory == "" {
		log.Println("⚠️ Warning: One or more required environment variables are missing!")
	}

	return nil
}
