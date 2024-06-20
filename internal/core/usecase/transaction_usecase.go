package usecase

import (
	"context"
	"fmt"
	"sort"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/jordanlanch/stori-test/internal/core/domain"
	"golang.org/x/time/rate"
)

type TransactionUseCase interface {
	ProcessTransactions(ctx context.Context) error
}

type TransactionRepository interface {
	GetAllTransactions(ctx context.Context) ([]domain.Transaction, error)
	SaveTransactions(ctx context.Context, transactions []domain.Transaction) error
	GetCSVHash() (string, error)
}

type CacheRepository interface {
	Get(ctx context.Context, key string) ([]domain.Transaction, error)
	Set(ctx context.Context, key string, value []domain.Transaction) error
}

type EmailService interface {
	SendEmail(ctx context.Context, templatePath string, data interface{}) error
}

type transactionUseCaseImpl struct {
	DBRepo        TransactionRepository
	CacheRepo     CacheRepository
	Email         EmailService
	RedisClient   *redis.Client
	CacheMutex    sync.Mutex
	RateLimiter   *rate.Limiter
	Timeout       time.Duration
	CacheDuration time.Duration
}

func NewTransactionUseCase(dbRepo TransactionRepository, cacheRepo CacheRepository, email EmailService, redisClient *redis.Client, rateLimit int, timeoutSec int, cacheDuration int) TransactionUseCase {
	return &transactionUseCaseImpl{
		DBRepo:        dbRepo,
		CacheRepo:     cacheRepo,
		Email:         email,
		RedisClient:   redisClient,
		RateLimiter:   rate.NewLimiter(rate.Every(time.Second), rateLimit), // rateLimit requests per second
		Timeout:       time.Duration(timeoutSec) * time.Second,
		CacheDuration: time.Duration(cacheDuration) * time.Second,
	}
}

func (uc *transactionUseCaseImpl) ProcessTransactions(ctx context.Context) error {
	if !uc.RateLimiter.Allow() {
		return fmt.Errorf("too many requests")
	}

	ctx, cancel := context.WithTimeout(ctx, uc.Timeout)
	defer cancel()

	hash, err := uc.DBRepo.GetCSVHash()
	if err != nil {
		return err
	}

	transactions, err := uc.CacheRepo.Get(ctx, hash)
	if err == nil && transactions != nil {
		summary := uc.generateHTMLSummary(transactions)
		return uc.Email.SendEmail(ctx, "./internal/infrastructure/email/templates/summary_template.html", summary)
	}

	transactions, err = uc.DBRepo.GetAllTransactions(ctx)
	if err != nil {
		return err
	}

	err = uc.DBRepo.SaveTransactions(ctx, transactions)
	if err != nil {
		return err
	}

	err = uc.CacheRepo.Set(ctx, hash, transactions)
	if err != nil {
		return err
	}

	summary := uc.generateHTMLSummary(transactions)
	return uc.Email.SendEmail(ctx, "./internal/infrastructure/email/templates/summary_template.html", summary)
}

func (uc *transactionUseCaseImpl) generateHTMLSummary(transactions []domain.Transaction) map[string]interface{} {
	var totalBalance float64
	monthlyTransactions := make(map[string][]domain.Transaction)

	for _, t := range transactions {
		totalBalance += t.Amount
		parts := strings.Split(t.Date, "/")
		if len(parts) != 2 {
			continue
		}
		month, err := strconv.Atoi(parts[0])
		if err != nil {
			continue
		}
		monthName := time.Month(month).String()
		monthlyTransactions[monthName] = append(monthlyTransactions[monthName], t)
	}

	months := make([]string, 0, len(monthlyTransactions))
	for month := range monthlyTransactions {
		months = append(months, month)
	}
	sort.Slice(months, func(i, j int) bool {
		mi, _ := time.Parse("January", months[i])
		mj, _ := time.Parse("January", months[j])
		return mi.Month() < mj.Month()
	})

	summary := map[string]interface{}{
		"TotalBalance": totalBalance,
		"MonthlyData":  make([]map[string]interface{}, 0),
	}

	for _, month := range months {
		txns := monthlyTransactions[month]
		var creditSum, debitSum float64
		var creditCount, debitCount int

		var wg sync.WaitGroup
		var mu sync.Mutex

		for _, t := range txns {
			wg.Add(1)
			go func(t domain.Transaction) {
				defer wg.Done()
				mu.Lock()
				if t.Amount > 0 {
					creditSum += t.Amount
					creditCount++
				} else {
					debitSum += t.Amount
					debitCount++
				}
				mu.Unlock()
			}(t)
		}
		wg.Wait()

		monthData := map[string]interface{}{
			"Month":         month,
			"Transactions":  len(txns),
			"DebitSum":      debitSum,
			"DebitCount":    debitCount,
			"CreditSum":     creditSum,
			"CreditCount":   creditCount,
			"AverageDebit":  0.0,
			"AverageCredit": 0.0,
		}

		if debitCount > 0 {
			monthData["AverageDebit"] = debitSum / float64(debitCount)
		}
		if creditCount > 0 {
			monthData["AverageCredit"] = creditSum / float64(creditCount)
		}

		summary["MonthlyData"] = append(summary["MonthlyData"].([]map[string]interface{}), monthData)
	}

	return summary
}
