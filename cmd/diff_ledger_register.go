package cmd

import (
	"encoding/csv"
	"fmt"
	"log/slog"
	"os"
	"sort"

	"github.com/icyflame/gnucash-xml-to-ledger-dat/lib/groupers"
	"github.com/icyflame/gnucash-xml-to-ledger-dat/lib/parsers/ledger"
	"github.com/icyflame/gnucash-xml-to-ledger-dat/lib/presenters"
	"github.com/icyflame/gnucash-xml-to-ledger-dat/lib/subtractors"
	"github.com/spf13/cobra"
)

var (
	outputFormat string
	bucket       string
)

var diffLedgerRegisterCmd = &cobra.Command{
	Use:   "register <csv-file-1: A side> <csv-file-2: B side>",
	Short: "Find differences between two Ledger files using the register command output",
	Long: `Find differences between the outputs of 'hledger register -O csv'
run on two Ledger files.

This tool accepts CSV files as input and does not care about their provenance. The input CSV files
MUST match the Ledger register command output.

There are two available output formats can be either plain text or TSV. When the output format is
set to TSV, the output will contain transactions that occur only on the A side on the left, and
those that occur only on the B side on the right.`,
	Args: cobra.ExactArgs(2),
	RunE: runDiffLedgerRegister,
}

func init() {
	diffLedgerRegisterCmd.Flags().StringVarP(&outputFormat, "output-format", "O", "diff", "Output format: diff or tsv")
	diffLedgerRegisterCmd.Flags().StringVarP(&bucket, "bucket", "", "none", "Bucket: none, weekly, monthly")
}

func runDiffLedgerRegister(cmd *cobra.Command, args []string) error {
	csvFile1 := args[0]
	csvFile2 := args[1]

	// Stage 1: Implement parsing logic inside lib/parsers/ledger/register_csv.go which takes a
	// Ledger Register as a CSV file and returns a slice of Ledger transactions
	//
	// File format:
	//
	//    "txnidx","date","code","description","account","amount","total"
	//    "20","2025-08-02","","Groceries","Liabilities:Credit Card","INR -3500","INR -3500"

	parser1 := ledger.New()
	if err := parser1.Parse(csvFile1); err != nil {
		return fmt.Errorf("failed to parse first CSV file: %w", err)
	}

	parser2 := ledger.New()
	if err := parser2.Parse(csvFile2); err != nil {
		return fmt.Errorf("failed to parse second CSV file: %w", err)
	}

	if parser1.GetCurrency() != parser2.GetCurrency() {
		return fmt.Errorf("currency mismatch between file 1 and file 2: file 1 currency: %s, file 2 currency: %s", parser1.GetCurrency(), parser2.GetCurrency())
	}

	transactions1 := parser1.GetTransactions()
	transactions2 := parser2.GetTransactions()

	slog.Info("parsed transactions", slog.Int("count", len(transactions1)), slog.String("file", csvFile1))
	slog.Info("parsed transactions", slog.Int("count", len(transactions2)), slog.String("file", csvFile2))

	// Stage 2: Implement diff logic: Diff should consider the two sides of ledger transactions as Set A and Set B.
	//
	// Diff_A_Side should be the Set subtraction (A-B) and Diff_B_Side should be the Set subtraction
	// (B-A). Both these set differences should be slices that contain a bunch of Ledger transactions.
	//
	// The key for each element in the Set should be "YYYY-MM-DD-AMOUNT" => Note that commodity is
	// not considered here. Retain the sign for + and - amount.

	parser1Map := parser1.BuildTransactionMap()
	parser2Map := parser2.BuildTransactionMap()

	subtractor := subtractors.New()

	diffAminusB := subtractor.Subtract(parser1Map, parser2Map)
	diffBminusA := subtractor.Subtract(parser2Map, parser1Map)

	sort.Sort(ledger.RegisterTransactionSlice(diffAminusB))
	sort.Sort(ledger.RegisterTransactionSlice(diffBminusA))

	// Note: sideA and sideB must contain the same number of elements
	var sideA, sideB []groupers.Bucket
	switch bucket {
	case "weekly":
		// TODO: Implement logic to break up the transactions in diffAminusB and diffBminusA into a
		// slice of slices, where each slice contains transactions from a given week
	case "monthly":
		// TODO: Implement logic to break up the transactions in diffAminusB and diffBminusA into a
		// slice of slices, where each slice contains transactions from a given month
	default:
		sideA = []groupers.Bucket{
			{
				Name:         "A-side Only",
				Transactions: diffAminusB,
			},
		}
		sideB = []groupers.Bucket{
			{
				Name:         "B-side Only",
				Transactions: diffBminusA,
			},
		}
	}

	if len(sideA) != len(sideB) {
		panic(fmt.Sprintf("logic error: sideA (%d) and sideB (%d) must be of the same size", len(sideA), len(sideB)))
	}

	// TODO: Update the presentation functions to accept a slice of slice of transactions, instead
	// of a slice of transactions. The presentation logic should not change much, except for things
	// like putting a new line in between various buckets and such.
	switch outputFormat {
	case "tsv":
		w := csv.NewWriter(os.Stdout)
		// Use tabs as separator to avoid having to deal with double quotes wrapping and such
		w.Comma = '\t'
		defer w.Flush()

		presenters.PresentTransactionsDiffAsTSV(w, sideA, sideB)
	default:
		presenters.PresentRegisterDiff(os.Stdout, sideA, sideB)
	}

	return nil
}
