package trade

import "time"

type Order struct {
	ID         int       `db:"id" json:"id"`
	UserID     int       `db:"user_id" json:"user_id"`
	Symbol     string    `db:"symbol" json:"symbol"`
	Side       string    `db:"side" json:"side"` // buy or sell
	LimitPrice float64   `db:"limit_price" json:"limit_price"`
	Contracts  float64   `db:"contracts" json:"contracts"`
	Leverage   int       `db:"leverage" json:"leverage"`
	StopLoss   *float64  `db:"stop_loss" json:"stop_loss,omitempty"`
	TakeProfit *float64  `db:"take_profit" json:"take_profit,omitempty"`
	Status     string    `db:"status" json:"status"`
	CreatedAt  time.Time `db:"created_at" json:"created_at"`
}

type Position struct {
	ID          int        `db:"id" json:"id"`
	OrderID     *int       `db:"order_id" json:"order_id,omitempty"`
	UserID      int        `db:"user_id" json:"user_id"`
	Symbol      string     `db:"symbol" json:"symbol"`
	Side        string     `db:"side" json:"side"`
	EntryPrice  float64    `db:"entry_price" json:"entry_price"`
	Contracts   float64    `db:"contracts" json:"contracts"`
	Leverage    int        `db:"leverage" json:"leverage"`
	StopLoss    *float64   `db:"stop_loss" json:"stop_loss,omitempty"`
	TakeProfit  *float64   `db:"take_profit" json:"take_profit,omitempty"`
	Status      string     `db:"status" json:"status"`
	OpenedAt    time.Time  `db:"opened_at" json:"opened_at"`
	ClosedAt    *time.Time `db:"closed_at" json:"closed_at,omitempty"`
	ExitPrice   *float64   `db:"exit_price" json:"exit_price,omitempty"`
	RealizedPnL float64    `db:"realized_pnl" json:"realized_pnl"`
}

type Tradee struct {
	ID        int       `db:"id" json:"id"`
	UserID    int       `db:"user_id" json:"user_id"`
	Symbol    string    `db:"symbol" json:"symbol"`
	Side      string    `db:"side" json:"side"`
	Price     float64   `db:"price" json:"price"`
	Qty       float64   `db:"qty" json:"qty"`
	PNL       float64   `db:"pnl" json:"pnl"`
	CreatedAt time.Time `db:"created_at" json:"created_at"`
}
