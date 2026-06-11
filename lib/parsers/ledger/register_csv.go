package ledger

import (
	"encoding/csv"
	"fmt"
	"os"
	"slices"
)

// RegisterTransaction represents a single transaction from a Ledger register CSV export.
type RegisterTransaction struct {
	TxnIdx      string // Transaction index
	Date        string // YYYY-MM-DD format
	Code        string // Transaction code
	Description string // Transaction description
	Account     string // Account name
	Amount      string // Amount with commodity (e.g., "INR -3500")
	Total       string // Running total with commodity
}

// RegisterCSVParser parses a Ledger register CSV file.
type RegisterCSVParser struct {
	transactions []RegisterTransaction
}

// New creates a new RegisterCSVParser.
func New() *RegisterCSVParser {
	return &RegisterCSVParser{
		transactions: []RegisterTransaction{},
	}
}

// Parse reads and parses a Ledger register CSV file.
// The CSV is expected to have the format:
//
//	"txnidx","date","code","description","account","amount","total"
func (p *RegisterCSVParser) Parse(filePath string) error {
	file, err := os.Open(filePath)
	if err != nil {
		return fmt.Errorf("failed to open CSV file: %w", err)
	}
	defer file.Close()

	reader := csv.NewReader(file)

	// Read and validate header
	header, err := reader.Read()
	if err != nil {
		return fmt.Errorf("failed to read CSV header: %w", err)
	}

	if !isValidHeader(header) {
		return fmt.Errorf("invalid CSV header: expected [txnidx, date, code, description, account, amount, total], got %v", header)
	}

	// Read all records
	records, err := reader.ReadAll()
	if err != nil {
		return fmt.Errorf("failed to read CSV records: %w", err)
	}

	// Parse records into RegisterTransaction objects
	for _, record := range records {
		txn := RegisterTransaction{
			TxnIdx:      record[0],
			Date:        record[1],
			Code:        record[2],
			Description: record[3],
			Account:     record[4],
			Amount:      record[5],
			Total:       record[6],
		}

		p.transactions = append(p.transactions, txn)
	}

	return nil
}

// GetTransactions returns the parsed transactions.
func (p *RegisterCSVParser) GetTransactions() []RegisterTransaction {
	return p.transactions
}

// isValidHeader validates that the CSV header matches the expected format.
func isValidHeader(header []string) bool {
	expectedHeader := []string{"txnidx", "date", "code", "description", "account", "amount", "total"}
	return slices.Equal(header, expectedHeader)
}
