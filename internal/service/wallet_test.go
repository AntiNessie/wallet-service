package service

import (
	"errors"
	"github.com/google/uuid"
	"testing"
	"wallet-service/internal/model"
)

// mockRepo — фейк репа для тестов
type mockRepo struct {
	wallet    *model.Wallet
	getErr    error
	createErr error
	updateErr error
}

func (m *mockRepo) GetWallet(id uuid.UUID) (*model.Wallet, error) {
	return m.wallet, m.getErr
}

func (m *mockRepo) CreateWallet(id uuid.UUID) error {
	return m.createErr
}

func (m *mockRepo) UpdateBalance(id uuid.UUID, amount float64) error {
	return m.updateErr
}

func TestProcessOperation_InvalidType(t *testing.T) {
	mock := &mockRepo{}
	svc := NewWalletService(mock)

	req := &model.WalletRequest{
		WalletID:      uuid.New(),
		OperationType: "TRANSFER", // невалидный тип
		Amount:        100,
	}

	err := svc.ProcessOperation(req)
	if !errors.Is(err, ErrInvalidOperation) {
		t.Errorf("expected ErrInvalidOperation, got %v", err)
	}
}

func TestProcessOperation_InvalidAmount(t *testing.T) {
	mock := &mockRepo{}
	svc := NewWalletService(mock)

	req := &model.WalletRequest{
		WalletID:      uuid.New(),
		OperationType: "DEPOSIT",
		Amount:        -50,
	}

	err := svc.ProcessOperation(req)
	if err == nil {
		t.Error("expected error for negative amount, got nil")
	}
}

func TestProcessOperation_Deposit_Success(t *testing.T) {
	mock := &mockRepo{}
	svc := NewWalletService(mock)

	id := uuid.New()
	req := &model.WalletRequest{
		WalletID:      id,
		OperationType: "DEPOSIT",
		Amount:        500,
	}

	err := svc.ProcessOperation(req)
	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}
}

func TestProcessOperation_Withdraw_InsufficientFunds(t *testing.T) {
	mock := &mockRepo{
		wallet: &model.Wallet{
			ID:      uuid.New(),
			Balance: 50,
		},
	}
	svc := NewWalletService(mock)

	req := &model.WalletRequest{
		WalletID:      uuid.New(),
		OperationType: "WITHDRAW",
		Amount:        100, // больше чем баланс
	}

	err := svc.ProcessOperation(req)
	if !errors.Is(err, ErrInsufficientFunds) {
		t.Errorf("expected ErrInsufficientFunds, got %v", err)
	}
}
