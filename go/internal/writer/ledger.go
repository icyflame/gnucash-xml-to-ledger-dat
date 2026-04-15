package writer

import (
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/icyflame/gnucash-xml-to-ledger-dat/internal/parser"
)

type LedgerWriter struct {
	parser  *parser.Parser
	verbose bool
}

func New(p *parser.Parser, verbose bool) *LedgerWriter {
	return &LedgerWriter{
		parser:  p,
		verbose: verbose,
	}
}

func (w *LedgerWriter) Write(writer io.Writer) error {
	transactions := w.parser.GetTransactions()

	for _, txn := range transactions {
		if len(txn.Splits.Split) == 0 {
			continue
		}

		date := w.formatDate(txn.DatePosted)
		description := w.formatDescription(txn.Description, txn.ID)

		fmt.Fprintf(writer, "%s  %s\n", date, description)

		for _, split := range txn.Splits.Split {
			account := w.parser.GetAccount(split.Account)
			if account == nil {
				continue
			}

			accountName := w.parser.GetAccountName(split.Account)
			commodity := account.Commodity.ID
			quantity := parser.ParseFraction(split.Quantity)
			value := parser.ParseFraction(split.Value)

			amount := quantity
			if quantity == 0 && value != 0 {
				amount = value
				commodity = txn.Currency.ID
			}

			if w.verbose {
				fmt.Fprintf(os.Stderr, "%s, %g, %g\n", accountName, quantity, value)
			}

			fmt.Fprintf(writer, "  %s  \"%s\"  %g\n", accountName, commodity, amount)
		}

		fmt.Fprintln(writer)
	}

	return nil
}

func (w *LedgerWriter) formatDate(datePosted string) string {
	// Extract YYYY-MM-DD from "2025-09-01 10:59:00 +0000"
	if len(datePosted) >= 10 {
		return datePosted[:10]
	}
	return datePosted
}

func (w *LedgerWriter) formatDescription(description, txnID string) string {
	description = strings.TrimSpace(description)
	if description == "" {
		return fmt.Sprintf("Empty description (%s)", txnID)
	}
	return description
}
