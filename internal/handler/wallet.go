package handler

import (
	"encoding/json"
	"errors"
	"github.com/google/uuid"
	"log"
	"net/http"
	"strings"
	"wallet-service/internal/model"
	"wallet-service/internal/service"
)

type WalletHandler struct {
	service *service.WalletService
}

func NewWalletHandler(service *service.WalletService) *WalletHandler {
	return &WalletHandler{service: service}

}

func (h *WalletHandler) HandleOperation(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	var req model.WalletRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "invalid request body"})
		return
	}
	err := h.service.ProcessOperation(&req)
	if err == nil {
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
		return
	}

	switch {
	case errors.Is(err, service.ErrInvalidOperation):
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "invalid operation type"})
		return
	case errors.Is(err, service.ErrInsufficientFunds):
		w.WriteHeader((http.StatusPaymentRequired))
		json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
	default:
		log.Printf("Internal error: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "internal service error"})
	}
}

func (h *WalletHandler) HandleGetBalance(w http.ResponseWriter, r *http.Request) {

	idStr := strings.TrimPrefix(r.URL.Path, "/api/v1/wallets/")
	id, err := uuid.Parse(idStr)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "invalid wallet id"})
		return
	}

	wallet, err := h.service.GetBalance(id)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "wallet not found"})
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(wallet)

}
