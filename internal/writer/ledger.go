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

			// Bug: In a Stock purchase transaction, the "Edit Exchange Rate" dialog is not shown by
			// GnuCash. This seems to be a limitation within the program. If the quantity for a
			// transaction is zero, then it should attempt to use "split:value" instead. If that is zero
			// too, then this split is a no-op and does not need to be included in the resulting Ledger
			// file.  However, for correctness, I *will* include the "0 COMMODITY" split in the output
			// Ledger file.
			quantity := parser.ParseFraction(split.Quantity)
			amount := quantity

			value := parser.ParseFraction(split.Value)
			if w.verbose {
				fmt.Fprintf(os.Stderr, "%s, %g, %g\n", accountName, quantity, value)
			}
			if quantity == 0 && value != 0 {
				amount = value
				commodity = txn.Currency.ID
			}

			// Commodity names can have any character (including white space) if they are enclosed in double quotes.
			//
			// See section 4.5.1 Naming Commodities in the Ledger manual: https://ledger-cli.org/doc/ledger3.pdf
			commodity = fmt.Sprintf("\"%s\"", commodity)

			// TODO: For Stock purchases and sales, this should also have the price at which stocks were
			// purchased, which will make it possible to use Ledger to calculate capital gains (maybe)
			//
			// See 4.5.2 Buying and Selling Stock in the Ledger manual: https://ledger-cli.org/doc/ledger3.pdf
			fmt.Fprintf(writer, "  %s  %s  %g\n", accountName, commodity, amount)

			// TODO: Include the second line from the transaction as a comment in the output Ledger file
			// (The second line is visible when the Double line view is enabled in GnuCash)
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

// Sample Ledger transaction:
// 2015/10/12 Exxon
//     Expenses:Auto:Gas         $10.00
//     Liabilities:MasterCard   $-10.00
//
// Sample Stock purchase transaction:
//
// 2004/05/01 Stock purchase
//   Assets:Broker                50 AAPL @ $30.00
//   Expenses:Broker:Commissions  $19.95
//   Assets:Broker                $-1,519.95
//
// The amount can be left out of the last posting:
//
// 2004/05/01 Stock purchase
//   Assets:Broker                50 AAPL @ $30.00
//   Expenses:Broker:Commissions  $19.95
//   Assets:Broker
