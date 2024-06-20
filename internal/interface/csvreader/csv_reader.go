package csv

import (
	"encoding/csv"
	"os"
	"strconv"
	"sync"

	"github.com/jordanlanch/stori-test/internal/core/domain"
)

type CSVReaderInterface interface {
	ReadTransactions() ([]domain.Transaction, error)
}

type CSVReader struct {
	FilePath string
}

func NewCSVReader(filePath string) *CSVReader {
	return &CSVReader{FilePath: filePath}
}

func (r *CSVReader) ReadTransactions() ([]domain.Transaction, error) {
	file, err := os.Open(r.FilePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	reader := csv.NewReader(file)
	records, err := reader.ReadAll()
	if err != nil {
		return nil, err
	}

	var wg sync.WaitGroup
	transactions := make([]domain.Transaction, len(records)-1)
	errors := make(chan error, len(records)-1)

	for i, record := range records[1:] { // Skip header
		wg.Add(1)
		go func(i int, record []string) {
			defer wg.Done()
			id, err := strconv.Atoi(record[0])
			if err != nil {
				errors <- err
				return
			}
			amount, err := strconv.ParseFloat(record[2], 64)
			if err != nil {
				errors <- err
				return
			}
			transactions[i] = domain.Transaction{
				ID:     id,
				Date:   record[1],
				Amount: amount,
			}
		}(i, record)
	}

	wg.Wait()
	close(errors)

	for err := range errors {
		if err != nil {
			return nil, err
		}
	}

	return transactions, nil
}
