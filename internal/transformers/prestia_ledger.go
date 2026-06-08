package transformers

import (
	"fmt"
	"io"
	"strconv"
	"strings"

	"github.com/icyflame/gnucash-xml-to-ledger-dat/lib/parsers/prestia"
)

// Constants for Prestia to Ledger transformation
const (
	PrestiaCounterpartyAccount       = "Expenses"
	PrestiaDefaultCommodity          = "JPY"
	PrestiaDefaultLiabilitiesAccount = "Liabilities:Prestia Card"
)

type PrestiaLedgerTransformer struct {
	records []prestia.PrestiaRecord
}

func NewPrestiaLedgerTransformer(records []prestia.PrestiaRecord) *PrestiaLedgerTransformer {
	return &PrestiaLedgerTransformer{
		records: records,
	}
}

func (t *PrestiaLedgerTransformer) Write(writer io.Writer) error {
	for _, record := range t.records {
		// Skip empty records (e.g., total rows)
		if strings.TrimSpace(record.Description) == "" {
			continue
		}

		date := t.formatDate(record.Date)
		description := t.formatDescription(record.Description)
		amount, err := t.parseAmount(record.Total)
		if err != nil {
			// Skip records with invalid amounts
			continue
		}

		// Write the transaction header
		fmt.Fprintf(writer, "%s  %s\n", date, description)

		// Write the posting line with the amount
		commodity := fmt.Sprintf("\"%s\"", PrestiaDefaultCommodity)
		fmt.Fprintf(writer, "  %s  %s  %g\n", PrestiaDefaultLiabilitiesAccount, commodity, -1*amount)

		// Write the counterparty posting
		fmt.Fprintf(writer, "  %s\n", PrestiaCounterpartyAccount)

		fmt.Fprintln(writer)
	}

	return nil
}

func (t *PrestiaLedgerTransformer) formatDate(date string) string {
	// Convert YYYY/MM/DD to YYYY-MM-DD
	return strings.ReplaceAll(date, "/", "-")
}

func (t *PrestiaLedgerTransformer) formatDescription(description string) string {
	description = strings.TrimSpace(description)
	if description == "" {
		return "Unknown transaction"
	}
	return description
}

func (t *PrestiaLedgerTransformer) parseAmount(total string) (float64, error) {
	total = strings.TrimSpace(total)
	if total == "" {
		return 0, fmt.Errorf("empty amount")
	}

	amount, err := strconv.ParseFloat(total, 64)
	if err != nil {
		return 0, fmt.Errorf("invalid amount: %w", err)
	}

	return amount, nil
}
