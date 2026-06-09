package repository

import (
	"fmt"
	"os"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"

	"wallet-service/internal/model"
)

type WalletRepository interface {
	GetWallet(id uuid.UUID) (*model.Wallet, error)
	CreateWallet(id uuid.UUID) error
	UpdateBalance(id uuid.UUID, amount float64) error
}

type PostgresRepo struct {
	db *sqlx.DB
}

func NewPostgresRepo(db *sqlx.DB) *PostgresRepo {
	return &PostgresRepo{db: db}
}

func (r *PostgresRepo) GetWallet(id uuid.UUID) (*model.Wallet, error) {
	wallet := &model.Wallet{}
	err := r.db.Get(wallet, "SELECT * FROM wallets WHERE id = $1", id)
	if err != nil {
		return nil, fmt.Errorf("get wallet: %w", err)
	}
	return wallet, nil
}

func (r *PostgresRepo) CreateWallet(id uuid.UUID) error {
	_, err := r.db.Exec("INSERT INTO wallets (id) VALUES ($1) ON CONFLICT DO NOTHING", id)
	return err
}

func (r *PostgresRepo) UpdateBalance(id uuid.UUID, amount float64) error {
	result, err := r.db.Exec(
		"UPDATE wallets SET balance = balance + $1, updated_at = NOW() WHERE id = $2",
		amount, id,
	)
	if err != nil {
		return fmt.Errorf("update balance: %w", err)
	}
	rows, _ := result.RowsAffected()
	if rows == 0 {
		return fmt.Errorf("wallet not found")
	}
	return nil
}

func (r *PostgresRepo) RunMigrations(sqlFilePath string) error {
	content, err := os.ReadFile(sqlFilePath)
	if err != nil {
		return fmt.Errorf("read migration file: %w", err)
	}
	_, err = r.db.Exec(string(content))
	if err != nil {
		return fmt.Errorf("execute migration: %w", err)
	}
	return nil
}
