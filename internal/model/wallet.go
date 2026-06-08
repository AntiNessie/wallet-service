package model

import (
	"github.com/google/uuid"
	"time"
)

type Wallet struct {
	ID        uuid.UUID `json:"id" db:"id"`
	Balance   float64   `json:"balance" db:"balance"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}

type WalletRequest struct {
	WalletID      uuid.UUID `json:"valletId"` // в задании опечатка "valletId", так и мапим
	OperationType string    `json:"operationType"`
	Amount        float64   `json:"amount"`
}

type WalletResponse struct {
	WalletID uuid.UUID `json:"walletId"`
	Balance  float64   `json:"balance"`
}
