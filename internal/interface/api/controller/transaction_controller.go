// internal/interface/api/controller/transaction_controller.go
package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jordanlanch/stori-test/internal/core/usecase"
)

type TransactionController struct {
	UseCase usecase.TransactionUseCase
}

func (ctrl *TransactionController) ProcessTransactions(c *gin.Context) {
	ctx := c.Request.Context()
	err := ctrl.UseCase.ProcessTransactions(ctx)
	if err != nil {
		if err.Error() == "too many requests" {
			c.JSON(http.StatusTooManyRequests, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Transactions processed and email sent"})
}
