package ledger

import (
	"encoding/csv"
	"fmt"
	"os"
	"regexp"
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

// BuildTransactionMap creates a map with key "YYYY-MM-DD-AMOUNT" for each transaction.
// The amount is extracted without the commodity (numeric value and sign only).
// Multiple transactions with the same key are stored as a slice.
func (p *RegisterCSVParser) BuildTransactionMap() map[string][]RegisterTransaction {
	return buildTransactionMapFromSlice(p.transactions)
}

// buildTransactionMapFromSlice is a helper function that builds a transaction map from a slice.
// This is used by both the parser and can be used independently for creating maps from transaction slices.
func buildTransactionMapFromSlice(transactions []RegisterTransaction) map[string][]RegisterTransaction {
	txnMap := make(map[string][]RegisterTransaction)

	for _, txn := range transactions {
		key := extractTransactionKey(txn)
		txnMap[key] = append(txnMap[key], txn)
	}

	return txnMap
}

// extractTransactionKey extracts the numeric amount from the Amount field and creates a key
// in the format "YYYY-MM-DD-AMOUNT". The amount includes the sign (+ or -) but no commodity.
func extractTransactionKey(txn RegisterTransaction) string {
	// Amount format: "INR -3500" or "INR 50000"
	// Extract the numeric value with sign
	amount := extractAmountValue(txn.Amount)
	return fmt.Sprintf("%s-%s", txn.Date, amount)
}

// extractAmountValue extracts the numeric amount with sign from a string like "INR -3500".
// Returns just the numeric part with sign, e.g., "-3500" or "50000".
func extractAmountValue(amountStr string) string {
	// Use regex to extract the numeric value (with optional sign)
	re := regexp.MustCompile(`(-?\d+(?:\.\d+)?)`)
	matches := re.FindStringSubmatch(amountStr)

	if len(matches) > 0 {
		return matches[1]
	}

	return amountStr // Fallback to original if no match
}
