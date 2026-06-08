package service

import (
	"errors"
	"fmt"
	"github.com/google/uuid"
	"wallet-service/internal/model"
	"wallet-service/internal/repository"
)

var (
	ErrInsufficientFunds = errors.New("insufficient funds")
	ErrInvalidOperation  = errors.New("invalid operation type")
)

type WalletService struct {
	repo *repository.PostgresRepo
}

func NewWalletService(repo *repository.PostgresRepo) *WalletService {
	return &WalletService{repo: repo}
}

func (s *WalletService) GetBalance(id uuid.UUID) (*model.WalletResponse, error) {
	wallet, err := s.repo.GetWallet(id)
	if err != nil {
		return nil, fmt.Errorf("get balance: %w", err)
	}
	return &model.WalletResponse{
		WalletID: wallet.ID,
		Balance:  wallet.Balance,
	}, nil
}

func (s *WalletService) ProcessOperation(req *model.WalletRequest) error {

	if req.OperationType != "DEPOSIT" && req.OperationType != "WITHDRAW" {
		return ErrInvalidOperation
	}

	if req.Amount <= 0 {
		return errors.New("amount must be positive")
	}

	if req.OperationType == "WITHDRAW" {
		wallet, err := s.repo.GetWallet(req.WalletID)
		if err != nil {
			if err.Error() == "sql: no rows in result set" {
				s.repo.CreateWallet(req.WalletID)
				return ErrInsufficientFunds
			}
			return fmt.Errorf("get wallet: %w", err)
		}
		if wallet.Balance < req.Amount {
			return ErrInsufficientFunds
		}
	} else {
		s.repo.CreateWallet(req.WalletID)
	}

	amount := req.Amount
	if req.OperationType == "WITHDRAW" {
		amount = -amount
	}

	if err := s.repo.UpdateBalance(req.WalletID, amount); err != nil {
		return fmt.Errorf("update balance: %w", err)
	}

	return nil
}
