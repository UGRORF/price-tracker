package main

import (
	"log"
	"os"

	"github.com/UGRORF/price-tracker/internal/api"
	"github.com/UGRORF/price-tracker/internal/api/handlers"
	"github.com/UGRORF/price-tracker/internal/config"
	"github.com/UGRORF/price-tracker/internal/repository/postgres"
	"github.com/UGRORF/price-tracker/internal/service/auth"
)

func main() {
	logger := log.New(os.Stdout, "PRICE-TRACKER: ", log.LstdFlags|log.Lshortfile)

	cfg, err := config.Load()
	if err != nil {
		logger.Fatal("Failed to load config:", err)
	}

	db, err := postgres.NewDB(&cfg.Database)
	if err != nil {
		logger.Fatal("Failed to connect to database:", err)
	}
	defer db.Close()

	// Handlers
	userRepo := postgres.NewUserRepo(db)
	handlers.InitUserRepo(userRepo)

	storeRepo := postgres.NewStoreRepo(db)
	handlers.InitStoreRepo(storeRepo)

	productRepo := postgres.NewProductRepo(db)
	handlers.InitProductRepo(productRepo)

	offerRepo := postgres.NewOfferRepo(db)
	handlers.InitOfferRepo(offerRepo)

	bucketRepo := postgres.NewBucketRepo(db)
	handlers.InitBucketRepo(bucketRepo)

	jwtService := auth.NewJWTService()
	authService := auth.NewService(userRepo, jwtService)
	authHandler := handlers.NewAuthHandler(authService)

	logger.Printf("Connected to database %s on %s:%d",
		cfg.Database.Name, cfg.Database.Host, cfg.Database.Port)

	srv := api.NewServer(logger, authHandler, jwtService)

	if err := srv.Start(); err != nil {
		logger.Fatal("Server failed:", err)
	}
}
