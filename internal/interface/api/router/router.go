package router

import (
	"github.com/gin-gonic/gin"
	"github.com/jordanlanch/stori-test/internal/interface/api/controller"
)

func SetupRouter(transactionController *controller.TransactionController) *gin.Engine {
	r := gin.Default()
	r.POST("/process-transactions", transactionController.ProcessTransactions)
	return r
}
