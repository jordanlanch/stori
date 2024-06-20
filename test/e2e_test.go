package test

import (
	"fmt"
	"math/rand"
	"net/http"
	"os"
	"testing"
	"time"

	"github.com/jordanlanch/stori-test/test/e2e"
)

const (
	statusOK              = http.StatusOK
	statusTooManyRequests = http.StatusTooManyRequests
)

func generateRandomCSV(filePath string, r *rand.Rand) error {
	file, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	file.WriteString("ID,Date,Transaction\n")
	for i := 0; i < 100; i++ {
		id := i
		month := r.Intn(12) + 1 // generate months between 1 and 12
		day := r.Intn(28) + 1   // generate days between 1 and 28 to avoid invalid dates
		date := fmt.Sprintf("%d/%d", month, day)
		amount := r.Float64()*200 - 100
		transaction := fmt.Sprintf("%.2f", amount)
		if amount >= 0 {
			transaction = "+" + transaction
		}
		file.WriteString(fmt.Sprintf("%d,%s,%s\n", id, date, transaction))
	}
	return nil
}

func TestProcessTransactionsOK(t *testing.T) {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	csvFilePath := "../test/transactions.csv"
	err := generateRandomCSV(csvFilePath, r)
	if err != nil {
		t.Fatalf("failed to generate random CSV: %v", err)
	}
	os.Setenv("CSV_FILE_PATH", csvFilePath)

	expect, teardown := e2e.Setup(t, "fixtures/transaction_test")
	defer teardown()

	t.Run("Process Transactions OK", func(t *testing.T) {
		response := expect.POST("/process-transactions").
			Expect()
		response.Status(statusOK)
	})
}

func TestRateLimitExceeded(t *testing.T) {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	csvFilePath := "../test/transactions.csv"
	err := generateRandomCSV(csvFilePath, r)
	if err != nil {
		t.Fatalf("failed to generate random CSV: %v", err)
	}
	os.Setenv("CSV_FILE_PATH", csvFilePath)

	expect, teardown := e2e.Setup(t, "fixtures/transaction_ratelimit_test")
	defer teardown()

	t.Run("Rate Limit Exceeded", func(t *testing.T) {
		for i := 0; i < 11; i++ { // Send 11 requests to exceed rate limit
			response := expect.POST("/process-transactions").
				Expect()
			if i < 10 {
				response.Status(statusOK)
			} else {
				response.Status(statusTooManyRequests)
			}
		}
	})
}
