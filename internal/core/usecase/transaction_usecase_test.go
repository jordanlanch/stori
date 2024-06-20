package usecase

import (
	"context"
	"errors"
	"testing"

	"github.com/go-redis/redis/v8"
	"github.com/jordanlanch/stori-test/internal/core/domain"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockTransactionRepository struct {
	mock.Mock
}

func (m *MockTransactionRepository) GetAllTransactions(ctx context.Context) ([]domain.Transaction, error) {
	args := m.Called(ctx)
	return args.Get(0).([]domain.Transaction), args.Error(1)
}

func (m *MockTransactionRepository) SaveTransactions(ctx context.Context, transactions []domain.Transaction) error {
	args := m.Called(ctx, transactions)
	return args.Error(0)
}

func (m *MockTransactionRepository) GetCSVHash() (string, error) {
	args := m.Called()
	return args.String(0), args.Error(1)
}

type MockCacheRepository struct {
	mock.Mock
}

func (m *MockCacheRepository) Get(ctx context.Context, key string) ([]domain.Transaction, error) {
	args := m.Called(ctx, key)
	if args.Get(0) != nil {
		return args.Get(0).([]domain.Transaction), args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *MockCacheRepository) Set(ctx context.Context, key string, value []domain.Transaction) error {
	args := m.Called(ctx, key, value)
	return args.Error(0)
}

type MockEmailService struct {
	mock.Mock
}

func (m *MockEmailService) SendEmail(ctx context.Context, templatePath string, data interface{}) error {
	args := m.Called(ctx, templatePath, data)
	return args.Error(0)
}

func TestProcessTransactions(t *testing.T) {
	mockDBRepo := new(MockTransactionRepository)
	mockCacheRepo := new(MockCacheRepository)
	mockEmail := new(MockEmailService)
	redisClient := &redis.Client{}
	rateLimit := 5
	timeoutSec := 5
	cacheDuration := 600

	useCase := NewTransactionUseCase(mockDBRepo, mockCacheRepo, mockEmail, redisClient, rateLimit, timeoutSec, cacheDuration)

	ctx := context.Background()

	transactions := []domain.Transaction{
		{ID: 1, Date: "1/1", Amount: 100},
		{ID: 2, Date: "1/2", Amount: -50},
	}

	mockDBRepo.On("GetCSVHash").Return("hash123", nil)
	mockCacheRepo.On("Get", mock.Anything, "hash123").Return(nil, errors.New("cache miss"))
	mockDBRepo.On("GetAllTransactions", mock.Anything).Return(transactions, nil)
	mockDBRepo.On("SaveTransactions", mock.Anything, transactions).Return(nil)
	mockCacheRepo.On("Set", mock.Anything, "hash123", transactions).Return(nil)
	mockEmail.On("SendEmail", mock.Anything, "./internal/infrastructure/email/templates/summary_template.html", mock.AnythingOfType("map[string]interface {}")).Return(nil)

	err := useCase.ProcessTransactions(ctx)
	assert.NoError(t, err)

	mockDBRepo.AssertExpectations(t)
	mockCacheRepo.AssertExpectations(t)
	mockEmail.AssertExpectations(t)
}

func TestProcessTransactions_RateLimitExceeded(t *testing.T) {
	mockDBRepo := new(MockTransactionRepository)
	mockCacheRepo := new(MockCacheRepository)
	mockEmail := new(MockEmailService)
	redisClient := &redis.Client{}
	rateLimit := 1
	timeoutSec := 5
	cacheDuration := 600

	useCase := NewTransactionUseCase(mockDBRepo, mockCacheRepo, mockEmail, redisClient, rateLimit, timeoutSec, cacheDuration)

	ctx := context.Background()

	transactions := []domain.Transaction{
		{ID: 1, Date: "1/1", Amount: 100},
		{ID: 2, Date: "1/2", Amount: -50},
	}

	mockDBRepo.On("GetCSVHash").Return("hash123", nil)
	mockCacheRepo.On("Get", mock.Anything, "hash123").Return(nil, errors.New("cache miss"))
	mockDBRepo.On("GetAllTransactions", mock.Anything).Return(transactions, nil)
	mockDBRepo.On("SaveTransactions", mock.Anything, transactions).Return(nil)
	mockCacheRepo.On("Set", mock.Anything, "hash123", transactions).Return(nil)
	mockEmail.On("SendEmail", mock.Anything, "./internal/infrastructure/email/templates/summary_template.html", mock.AnythingOfType("map[string]interface {}")).Return(nil)

	// First request should succeed
	err := useCase.ProcessTransactions(ctx)
	assert.NoError(t, err)

	// Second request should exceed rate limit
	err = useCase.ProcessTransactions(ctx)
	assert.Error(t, err)
	assert.Equal(t, "too many requests", err.Error())

	mockDBRepo.AssertExpectations(t)
	mockCacheRepo.AssertExpectations(t)
	mockEmail.AssertExpectations(t)
}

func TestProcessTransactions_CacheHit(t *testing.T) {
	mockDBRepo := new(MockTransactionRepository)
	mockCacheRepo := new(MockCacheRepository)
	mockEmail := new(MockEmailService)
	redisClient := &redis.Client{}
	rateLimit := 5
	timeoutSec := 5
	cacheDuration := 600

	useCase := NewTransactionUseCase(mockDBRepo, mockCacheRepo, mockEmail, redisClient, rateLimit, timeoutSec, cacheDuration)

	ctx := context.Background()

	transactions := []domain.Transaction{
		{ID: 1, Date: "1/1", Amount: 100},
		{ID: 2, Date: "1/2", Amount: -50},
	}

	mockDBRepo.On("GetCSVHash").Return("hash123", nil)
	mockCacheRepo.On("Get", mock.Anything, "hash123").Return(transactions, nil)
	mockEmail.On("SendEmail", mock.Anything, "./internal/infrastructure/email/templates/summary_template.html", mock.AnythingOfType("map[string]interface {}")).Return(nil)

	err := useCase.ProcessTransactions(ctx)
	assert.NoError(t, err)

	mockDBRepo.AssertExpectations(t)
	mockCacheRepo.AssertExpectations(t)
	mockEmail.AssertExpectations(t)
}

func TestProcessTransactions_DBError(t *testing.T) {
	mockDBRepo := new(MockTransactionRepository)
	mockCacheRepo := new(MockCacheRepository)
	mockEmail := new(MockEmailService)
	redisClient := &redis.Client{}
	rateLimit := 5
	timeoutSec := 5
	cacheDuration := 600

	useCase := NewTransactionUseCase(mockDBRepo, mockCacheRepo, mockEmail, redisClient, rateLimit, timeoutSec, cacheDuration)

	ctx := context.Background()

	mockDBRepo.On("GetCSVHash").Return("hash123", nil)
	mockCacheRepo.On("Get", mock.Anything, "hash123").Return(nil, errors.New("cache miss"))
	mockDBRepo.On("GetAllTransactions", mock.Anything).Return(nil, errors.New("db error"))

	err := useCase.ProcessTransactions(ctx)
	assert.Error(t, err)
	assert.Equal(t, "db error", err.Error())

	mockDBRepo.AssertExpectations(t)
	mockCacheRepo.AssertExpectations(t)
	mockEmail.AssertExpectations(t)
}
