package adapter

import (
	"fmt"
	"time"

	"github.com/rismapa/go-banking-auth/domain"

	"github.com/jmoiron/sqlx"
)

type AuthRepository interface {
	SaveToken(userID, token, expiresAt string) error
	GetAccountByUsername(username string) (*domain.Account, error)
	GetTokenExpiration(token string) (time.Time, error)
}

type AccountRepositoryDB struct {
	DB *sqlx.DB
}

func NewAccountRepositoryDB(db *sqlx.DB) *AccountRepositoryDB {
	return &AccountRepositoryDB{DB: db}
}

func (a *AccountRepositoryDB) SaveToken(userID, token, expiresAt string) error {
	query := "INSERT INTO tokens (user_id, token, expires_at) VALUES (?, ?, ?) ON DUPLICATE KEY UPDATE token = ?, expires_at = ?"
	_, err := a.DB.Exec(query, userID, token, expiresAt, token, expiresAt)
	if err != nil {
		return fmt.Errorf("failed to save token: %v", err)
	}

	return nil
}

func (a *AccountRepositoryDB) GetAccountByUsername(username string) (*domain.Account, error) {
	var account domain.Account
	query := "SELECT id, customer_id, username, password, balance, currency, status FROM accounts WHERE username = ?"
	err := a.DB.Get(&account, query, username)
	if err != nil {
		if account == (domain.Account{}) {
			return nil, fmt.Errorf("no accounts found")
		}
		return nil, fmt.Errorf("database error: %v", err)
	}

	return &account, nil
}

func (r *AccountRepositoryDB) GetTokenExpiration(token string) (time.Time, error) {
	var expiresAtStr string
	query := "SELECT expires_at FROM tokens WHERE token = ?"
	err := r.DB.Get(&expiresAtStr, query, token)
	if err != nil {
		return time.Time{}, fmt.Errorf("failed to fetch token expiration: %v", err)
	}

	// Parse the expiration time string into a time.Time object
	expiresAt, err := time.Parse("2006-01-02 15:04:05", expiresAtStr)
	if err != nil {
		return time.Time{}, fmt.Errorf("failed to parse expires_at: %v", err)
	}

	return expiresAt, nil
}
