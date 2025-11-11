package main

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/ruturajlaad/TradeExecEngine/internal/api"
	"github.com/ruturajlaad/TradeExecEngine/internal/db"
	"github.com/ruturajlaad/TradeExecEngine/internal/market"
	"github.com/ruturajlaad/TradeExecEngine/internal/trade"
)

func main() {

	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, relying on system environment variables")
	}

	db.ConnectDB()

	//live prices channel
	priceChan := make(chan float64)

	tradeRepo := trade.NewRepository(db.DB)
	tradeService := trade.NewService(tradeRepo)
	tradeHandler := trade.NewHandler(tradeService)
	//binance price stream
	go market.StartMarketFeed(priceChan, tradeService)

	r := gin.Default()
	api.SetupRoutes(r, tradeHandler, tradeService)
	api.SetupMarketWebS(r, priceChan)
	r.Run(":8080")
}
