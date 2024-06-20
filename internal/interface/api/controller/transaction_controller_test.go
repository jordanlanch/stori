package controller

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockTransactionUseCase mocks the TransactionUseCase interface
type MockTransactionUseCase struct {
	mock.Mock
}

func (m *MockTransactionUseCase) ProcessTransactions(ctx context.Context) error {
	args := m.Called(ctx)
	return args.Error(0)
}

func TestTransactionController_ProcessTransactions(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("success", func(t *testing.T) {
		mockUseCase := new(MockTransactionUseCase)
		mockUseCase.On("ProcessTransactions", mock.Anything).Return(nil)

		controller := &TransactionController{UseCase: mockUseCase}
		router := gin.Default()
		router.POST("/process-transactions", controller.ProcessTransactions)

		req, _ := http.NewRequest(http.MethodPost, "/process-transactions", nil)
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		assert.JSONEq(t, `{"message":"Transactions processed and email sent"}`, w.Body.String())

		mockUseCase.AssertExpectations(t)
	})

	t.Run("too many requests", func(t *testing.T) {
		mockUseCase := new(MockTransactionUseCase)
		mockUseCase.On("ProcessTransactions", mock.Anything).Return(errors.New("too many requests"))

		controller := &TransactionController{UseCase: mockUseCase}
		router := gin.Default()
		router.POST("/process-transactions", controller.ProcessTransactions)

		req, _ := http.NewRequest(http.MethodPost, "/process-transactions", nil)
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusTooManyRequests, w.Code)
		assert.JSONEq(t, `{"error":"too many requests"}`, w.Body.String())

		mockUseCase.AssertExpectations(t)
	})

	t.Run("internal server error", func(t *testing.T) {
		mockUseCase := new(MockTransactionUseCase)
		mockUseCase.On("ProcessTransactions", mock.Anything).Return(errors.New("internal error"))

		controller := &TransactionController{UseCase: mockUseCase}
		router := gin.Default()
		router.POST("/process-transactions", controller.ProcessTransactions)

		req, _ := http.NewRequest(http.MethodPost, "/process-transactions", nil)
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
		assert.JSONEq(t, `{"error":"internal error"}`, w.Body.String())

		mockUseCase.AssertExpectations(t)
	})
}
