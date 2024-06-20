package repository

import (
	"context"
	"crypto/sha256"
	"encoding/csv"
	"encoding/hex"
	"os"
	"testing"

	"github.com/jordanlanch/stori-test/internal/core/domain"
	csvreader "github.com/jordanlanch/stori-test/internal/interface/csvreader"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// MockCSVReader is a mock of the CSVReaderInterface
type MockCSVReader struct {
	mock.Mock
}

func (m *MockCSVReader) ReadTransactions() ([]domain.Transaction, error) {
	args := m.Called()
	return args.Get(0).([]domain.Transaction), args.Error(1)
}

func createTestDB() (*gorm.DB, error) {
	return gorm.Open(sqlite.Open("file::memory:?cache=shared"), &gorm.Config{})
}

func TestGetAllTransactions(t *testing.T) {
	mockReader := new(MockCSVReader)
	repo := &DBTransactionRepository{
		db:        nil,
		csvReader: mockReader,
	}

	expectedTransactions := []domain.Transaction{
		{Date: "1/1", Amount: 100},
		{Date: "1/2", Amount: -50},
	}
	mockReader.On("ReadTransactions").Return(expectedTransactions, nil)

	transactions, err := repo.GetAllTransactions(context.Background())
	assert.NoError(t, err)
	assert.Equal(t, expectedTransactions, transactions)

	mockReader.AssertExpectations(t)
}

func TestSaveTransactions(t *testing.T) {
	db, err := createTestDB()
	assert.NoError(t, err)

	repo := &DBTransactionRepository{
		db:        db,
		csvReader: nil,
	}

	err = db.AutoMigrate(&domain.Transaction{})
	assert.NoError(t, err)

	transactions := []domain.Transaction{
		{Date: "1/1", Amount: 100},
		{Date: "1/2", Amount: -50},
	}

	err = repo.SaveTransactions(context.Background(), transactions)
	assert.NoError(t, err)

	var savedTransactions []domain.Transaction
	err = db.Find(&savedTransactions).Error
	assert.NoError(t, err)
	assert.Len(t, savedTransactions, 2)
	assert.Equal(t, transactions[0].Date, savedTransactions[0].Date)
	assert.Equal(t, transactions[0].Amount, savedTransactions[0].Amount)
	assert.Equal(t, transactions[1].Date, savedTransactions[1].Date)
	assert.Equal(t, transactions[1].Amount, savedTransactions[1].Amount)
}

func createTempCSVFile(content string) (string, error) {
	file, err := os.CreateTemp("", "testcsv")
	if err != nil {
		return "", err
	}
	defer file.Close()

	_, err = file.WriteString(content)
	if err != nil {
		return "", err
	}

	return file.Name(), nil
}

func TestGetCSVHash(t *testing.T) {
	content := "Date,Amount\n1/1,100\n1/2,-50\n"
	filePath, err := createTempCSVFile(content)
	assert.NoError(t, err)
	defer os.Remove(filePath)

	repo := &DBTransactionRepository{
		db:        nil,
		csvReader: csvreader.NewCSVReader(filePath),
	}

	expectedHash := sha256.New()
	expectedHash.Write([]byte(filePath))

	file, err := os.Open(filePath)
	assert.NoError(t, err)
	defer file.Close()

	reader := csv.NewReader(file)
	records, err := reader.ReadAll()
	assert.NoError(t, err)
	for _, record := range records {
		for _, field := range record {
			expectedHash.Write([]byte(field))
		}
	}
	expectedHashStr := hex.EncodeToString(expectedHash.Sum(nil))

	hash, err := repo.GetCSVHash()
	assert.NoError(t, err)
	assert.Equal(t, expectedHashStr, hash)
}
