package trade

import (
	"errors"
	"time"
)

type Service struct {
	repo *Repository
}

func NewService(r *Repository) *Service {
	return &Service{repo: r}
}

// filling placeorder
func (s *Service) PlaceOrder(o *Order) error {
	o.Status = "pending"
	o.CreatedAt = time.Now()
	return s.repo.CreateOrder(o)
}

func (s *Service) CheckOrdersForExecution(symbol string, currentPrice float64) error {
	orders, err := s.repo.GetPendingOrders()

	if err != nil {
		return err
	}
	for _, order := range orders {
		if order.Symbol != symbol {
			continue
		}

		shouldExec := false
		if order.Side == "buy" && currentPrice <= order.LimitPrice {
			shouldExec = true
		} else if order.Side == "sell" && currentPrice >= order.LimitPrice {
			shouldExec = true
		}

		if shouldExec {
			pos := &Position{
				OrderID:    &order.ID,
				UserID:     order.UserID,
				Symbol:     order.Symbol,
				Side:       order.Side,
				EntryPrice: currentPrice,
				Contracts:  order.Contracts,
				Leverage:   order.Leverage,
				StopLoss:   order.StopLoss,
				TakeProfit: order.TakeProfit,
				Status:     "open",
				OpenedAt:   time.Now(),
			}
			if err := s.repo.CreatePosition(pos); err != nil {
				return err
			}
			s.repo.MarkOrderExecuted(order.ID)
		}
	}
	return nil
}
func (s *Service) ClosePosition(id int, exitPrice float64) error {
	posList, err := s.repo.GetOpenPositions()

	if err != nil {
		return err
	}

	for _, p := range posList {
		if p.ID == id {
			var pnl float64
			if p.Side == "buy" {
				pnl = (exitPrice - p.EntryPrice) * p.Contracts * float64(p.Leverage)
			} else if p.Side == "sell" {
				pnl = (p.EntryPrice - exitPrice) * p.Contracts * float64(p.Leverage)
			} else {
				return errors.New("Invalid Side!")
			}
			if err := s.repo.ClosePosition(p.ID, exitPrice, pnl); err != nil {
				return err
			}

			// 2️⃣ Update user balance with the realized PnL
			user, err := s.repo.GetUserByID(p.UserID)
			if err != nil {
				return err
			}
			newBalance := user.Balance + pnl
			if err := s.repo.UpdateUserBalance(p.UserID, newBalance); err != nil {
				return err
			}

			return nil //
		}
	}

	return errors.New("Position not found!")
}
func (s *Service) GetUserBalance(userID int) (float64, error) {
	return s.repo.GetUserBalance(userID)
}
