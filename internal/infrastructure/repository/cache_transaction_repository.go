package repository

import (
	"context"
	"encoding/json"
	"sync"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/jordanlanch/stori-test/internal/core/domain"
)

type CacheTransactionRepository struct {
	client        *redis.Client
	cacheDuration time.Duration
	mu            sync.Mutex
}

func NewCacheTransactionRepository(client *redis.Client, cacheDurationSec int) *CacheTransactionRepository {
	return &CacheTransactionRepository{
		client:        client,
		cacheDuration: time.Duration(cacheDurationSec) * time.Second,
	}
}

func (r *CacheTransactionRepository) Get(ctx context.Context, key string) ([]domain.Transaction, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	result, err := r.client.Get(ctx, key).Result()
	if err == redis.Nil {
		return nil, nil
	} else if err != nil {
		return nil, err
	}

	var transactions []domain.Transaction
	err = json.Unmarshal([]byte(result), &transactions)
	if err != nil {
		return nil, err
	}

	return transactions, nil
}

func (r *CacheTransactionRepository) Set(ctx context.Context, key string, transactions []domain.Transaction) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	data, err := json.Marshal(transactions)
	if err != nil {
		return err
	}

	return r.client.Set(ctx, key, data, r.cacheDuration).Err()
}
