# dex-trade-x

/tendo-arbitrage-bot  # Project Root
│── /config            # Configuration files (API keys, settings)
│   ├── config.go      # Load Binance, Uniswap, Infura API keys
│── /exchange          # Handles exchange API interactions
│   ├── binance.go     # Fetch prices & execute trades on Binance
│   ├── uniswap.go     # Fetch prices & swap tokens on Uniswap
│── /arbitrage         # Core arbitrage detection logic
│   ├── detect.go      # Identifies arbitrage opportunities
│   ├── execute.go     # Handles trade execution
│── /utils             # Utility functions (logging, error handling, etc.)
│   ├── logger.go      # Logs transactions & errors
│── main.go            # Entry point (runs the bot)
│── go.mod             # Golang dependencies
│── .env               # Stores private API keys (DO NOT COMMIT)
