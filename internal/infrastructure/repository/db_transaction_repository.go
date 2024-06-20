package repository

import (
	"context"
	"crypto/sha256"
	"encoding/csv"
	"encoding/hex"
	"os"

	"github.com/jordanlanch/stori-test/internal/core/domain"
	csvreader "github.com/jordanlanch/stori-test/internal/interface/csvreader"
	"gorm.io/gorm"
)

type DBTransactionRepository struct {
	db        *gorm.DB
	csvReader csvreader.CSVReaderInterface
}

func NewDBTransactionRepository(db *gorm.DB, filePath string) *DBTransactionRepository {
	return &DBTransactionRepository{db: db, csvReader: csvreader.NewCSVReader(filePath)}
}

func (r *DBTransactionRepository) GetAllTransactions(ctx context.Context) ([]domain.Transaction, error) {
	return r.csvReader.ReadTransactions()
}

func (r *DBTransactionRepository) SaveTransactions(ctx context.Context, transactions []domain.Transaction) error {
	// Create a slice without IDs for insertion
	transactionsWithoutIDs := make([]domain.Transaction, len(transactions))
	for i, t := range transactions {
		transactionsWithoutIDs[i] = domain.Transaction{
			Date:   t.Date,
			Amount: t.Amount,
		}
	}
	return r.db.WithContext(ctx).Create(&transactionsWithoutIDs).Error
}

func (r *DBTransactionRepository) GetCSVHash() (string, error) {
	file, err := os.Open(r.csvReader.(*csvreader.CSVReader).FilePath)
	if err != nil {
		return "", err
	}
	defer file.Close()

	hash := sha256.New()
	if _, err := hash.Write([]byte(r.csvReader.(*csvreader.CSVReader).FilePath)); err != nil {
		return "", err
	}

	reader := csv.NewReader(file)
	records, err := reader.ReadAll()
	if err != nil {
		return "", err
	}

	for _, record := range records {
		for _, field := range record {
			if _, err := hash.Write([]byte(field)); err != nil {
				return "", err
			}
		}
	}

	return hex.EncodeToString(hash.Sum(nil)), nil
}
