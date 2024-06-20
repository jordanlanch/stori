package main

import (
	"fmt"
	"log"

	"github.com/go-redis/redis/v8"
	"github.com/jordanlanch/stori-test/internal/config"
	"github.com/jordanlanch/stori-test/internal/core/usecase"
	"github.com/jordanlanch/stori-test/internal/infrastructure/email"
	"github.com/jordanlanch/stori-test/internal/infrastructure/repository"
	"github.com/jordanlanch/stori-test/internal/interface/api/controller"
	"github.com/jordanlanch/stori-test/internal/interface/api/router"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	env := config.NewEnv(".env")
	if err := env.Validate(); err != nil {
		log.Fatalf("Environment validation failed: %v", err)
	}

	// Setup Redis
	redisOptions := &redis.Options{
		Addr: fmt.Sprintf("%s:%d", env.RedisHost, env.RedisPort),
		DB:   0,
	}

	if env.RedisPassword != "" {
		redisOptions.Password = env.RedisPassword
	}

	redisClient := redis.NewClient(redisOptions)

	// Setup Database
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%d sslmode=disable",
		env.DBHost, env.DBUser, env.DBPassword, env.DBName, env.DBPort)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	// Setup Repository, Services, and UseCase
	dbRepo := repository.NewDBTransactionRepository(db, env.CSVFilePath)
	cacheRepo := repository.NewCacheTransactionRepository(redisClient, env.CacheDurationSec)
	emailService := &email.EmailService{}
	transactionUseCase := usecase.NewTransactionUseCase(dbRepo, cacheRepo, emailService, redisClient, env.RateLimit, env.RedisTimeoutSec, env.CacheDurationSec)
	transactionController := &controller.TransactionController{
		UseCase: transactionUseCase,
	}

	// Setup Router
	r := router.SetupRouter(transactionController)
	r.Run(env.ServerAddress)
}
