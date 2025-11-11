package user

import (
	"time"

	"github.com/ruturajlaad/TradeExecEngine/internal/db"

	"golang.org/x/crypto/bcrypt"
)

type User struct {
	ID           int       `db:"id"	  json:"id"`
	Username     string    `db:"username"	  json:"username"`
	PasswordHash string    `db:"password_hash"`
	Balance      float64   `db:"balance"	  json:"balance"`
	CreatedAt    time.Time `db:"created_at"`
}

func CreateUser(username, password string) error {
	hash, _ := bcrypt.GenerateFromPassword([]byte(password), 10)
	_, err := db.DB.Exec(`INSERT INTO users(username,password_hash,balance,created_at) VALUES($1,$2,$3,$4)`, username, string(hash), 10000.0, time.Now())
	return err
}

func Authenticate(username, password string) (*User, error) {
	var u User

	err := db.DB.Get(&u, "SELECT * FROM users WHERE  username=$1", username)
	if err != nil {
		return nil, err
	}
	if bcrypt.CompareHashAndPassword([]byte(u.PasswordHash), []byte(password)) != nil {
		return nil, err
	}
	return &u, err
}
