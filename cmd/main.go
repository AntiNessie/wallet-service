package main

import (
	"fmt"
	"net/http"
	"os"
	"time"
	"wallet-service/internal/logger"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"

	"wallet-service/internal/config"
	"wallet-service/internal/handler"
	"wallet-service/internal/repository"
	"wallet-service/internal/service"
)

func main() {

	logger.Init()
	cfg := config.Load()

	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		cfg.DBHost, cfg.DBPort, cfg.DBUser, cfg.DBPassword, cfg.DBName)

	var db *sqlx.DB
	var err error
	for i := 0; i < 10; i++ {
		db, err = sqlx.Connect("postgres", dsn)
		if err == nil {
			break
		}
		logger.Log.Info("Waiting for database...", "attempt", i+1)
		time.Sleep(2 * time.Second)
	}
	if err != nil {
		logger.Log.Error("Failed to connect to database", "error", err)
		os.Exit(1)
	}
	defer db.Close()

	repo := repository.NewPostgresRepo(db)
	svc := service.NewWalletService(repo)
	h := handler.NewWalletHandler(svc)

	if err := repo.RunMigrations("migrations/001_init.sql"); err != nil {
		logger.Log.Error("Failed to run migrations", "error", err)
		os.Exit(1)
	}

	http.HandleFunc("/api/v1/wallet", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}
		h.HandleOperation(w, r)
	})

	http.HandleFunc("/api/v1/wallets/", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}
		h.HandleGetBalance(w, r)
	})

	logger.Log.Info("Server starting", "port", cfg.ServerPort)
	if err := http.ListenAndServe(":"+cfg.ServerPort, nil); err != nil {
		logger.Log.Error("Failed", "error", err)
		os.Exit(1)
	}
}
