package market

import (
	"encoding/json"
	"log"
	"time"

	"github.com/gorilla/websocket"
	"github.com/ruturajlaad/TradeExecEngine/internal/trade"
)

type MarketPrice struct {
	Symbol string  `json:"s"`
	Price  float64 `json:"p,string"`
}

func StartMarketFeed(priceChan chan<- float64, tradeService *trade.Service) {
	wsURL := "wss://stream.binance.com:9443/ws/btcusdt@trade"

	c, _, err := websocket.DefaultDialer.Dial(wsURL, nil)
	if err != nil {
		log.Fatal("Websocket dial error:", err)
	}
	defer c.Close()
	log.Println("Connected to Binance BTC/USDT feed!")

	var (
		lastBroadcast time.Time
		latestPrice   float64
	)

	for {
		_, msg, err := c.ReadMessage()
		if err != nil {
			log.Println("⚠️ Read error:", err)
			time.Sleep(2 * time.Second)
			continue
		}

		var tradeTick MarketPrice
		if err := json.Unmarshal(msg, &tradeTick); err != nil {
			continue
		}

		latestPrice = tradeTick.Price

		// ✅ Emit once every 200ms (to reduce spam)
		if time.Since(lastBroadcast) >= 200*time.Millisecond {
			select {
			case priceChan <- latestPrice:
				lastBroadcast = time.Now()

				// ✅ Run order checks in background (non-blocking)
				go func(p float64) {
					tradeService.CheckOrdersForExecution("BTCUSDT", p)
				}(latestPrice)

			default:

			}
		}
	}
}
