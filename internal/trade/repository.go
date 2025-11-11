package trade

import "github.com/jmoiron/sqlx"

type Repository struct {
	DB *sqlx.DB
}

type User struct {
	ID       int     `db:"id"`
	Username string  `db:"username"`
	Balance  float64 `db:"balance"`
}

func NewRepository(db *sqlx.DB) *Repository {
	return &Repository{DB: db}
}

//orders
func (r *Repository) CreateOrder(o *Order) error {
	query := `
	INSERT INTO orders (user_id, symbol, side, limit_price, contracts, leverage, stop_loss, take_profit, status)
	VALUES (:user_id, :symbol, :side, :limit_price, :contracts, :leverage, :stop_loss, :take_profit, :status)
	RETURNING id`

	rows, err := r.DB.NamedQuery(query, o)

	if err != nil {
		return err
	}
	defer rows.Close()
	if rows.Next() {
		rows.Scan(&o.ID)
	}
	return nil

}

func (r *Repository) GetPendingOrders() ([]Order, error) {
	var orders []Order
	err := r.DB.Select(&orders, "Select * from orders where status ='pending'")
	return orders, err
}

func (r *Repository) MarkOrderExecuted(id int) error {
	_, err := r.DB.Exec("UPDATE users SET status='executed' WHERE id = $1", id)
	return err
}

func (r *Repository) CreatePosition(p *Position) error {
	query := `INSERT INTO positions (order_id, user_id, symbol, side, entry_price, contracts, leverage, stop_loss, take_profit, status)
	VALUES (:order_id, :user_id, :symbol, :side, :entry_price, :contracts, :leverage, :stop_loss, :take_profit, :status)
	RETURNING id`
	rows, err := r.DB.NamedQuery(query, p)
	if err != nil {
		return err
	}
	defer rows.Close()
	if rows.Next() {
		rows.Scan(&p.ID)
	}
	return nil
}

func (r *Repository) GetOpenPositions() ([]Position, error) {
	var positions []Position
	err := r.DB.Select(&positions, "SELECT * FROM positions where status = 'open'")
	return positions, err
}

func (r *Repository) ClosePosition(id int, exitPrice, realizedPnL float64) error {
	query := `UPDATE positions
	SET status='closed', exit_price=$1, realized_pnl=$2, closed_at=NOW()
	WHERE id=$3`

	_, err := r.DB.Exec(query, exitPrice, realizedPnL, id)

	return err
}

func (r *Repository) GetUserByID(userID int) (*User, error) {
	var u User
	err := r.DB.Get(&u, "SELECT id, username, balance FROM users WHERE id = $1", userID)
	if err != nil {
		return nil, err
	}
	return &u, nil
}

func (r *Repository) UpdateUserBalance(userID int, newBalance float64) error {
	_, err := r.DB.Exec("UPDATE users SET balance = $1 WHERE id = $2", newBalance, userID)
	return err
}

func (r *Repository) GetUserBalance(userID int) (float64, error) {
	var balance float64
	err := r.DB.Get(&balance, "SELECT balance FROM users WHERE id = $1", userID)
	return balance, err
}
