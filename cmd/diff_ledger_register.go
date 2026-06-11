package cmd

import (
	"fmt"

	"github.com/icyflame/gnucash-xml-to-ledger-dat/lib/parsers/ledger"
	"github.com/spf13/cobra"
)

var diffLedgerRegisterCmd = &cobra.Command{
	Use:   "register <csv-file-1> <csv-file-2>",
	Short: "Find differences between two Ledger files using the register command output",
	Long: `Find differences between the outputs of 'hledger register -O csv'
run on two Ledger files.

This tool accepts CSV files as input and does not care about their provenance. The input CSV files
MUST match the Ledger register command output.`,
	Args: cobra.ExactArgs(2),
	RunE: runDiffLedgerRegister,
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

	transactions1 := parser1.GetTransactions()
	transactions2 := parser2.GetTransactions()

	fmt.Printf("Parsed %d transactions from %s\n", len(transactions1), csvFile1)
	fmt.Printf("Parsed %d transactions from %s\n", len(transactions2), csvFile2)

	// Stage 2: Implement diff logic: Diff should consider the two sides of ledger transactions as Set A and Set B.
	//
	// Diff_A_Side should be the Set subtraction (A-B) and Diff_B_Side should be the Set subtraction
	// (B-A). Both these set differences should be slices that contain a bunch of Ledger transactions.
	//
	// The key for each element in the Set should be "YYYY-MM-DD-AMOUNT" => Note that commodity is
	// not considered here. Retain the sign for + and - amount.

	// TODO: Stage 2

	_ = transactions1
	_ = transactions2

	return nil
}
