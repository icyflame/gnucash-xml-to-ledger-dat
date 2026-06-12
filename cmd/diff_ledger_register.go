package cmd

import (
	"fmt"

	"github.com/icyflame/gnucash-xml-to-ledger-dat/lib/parsers/ledger"
	"github.com/icyflame/gnucash-xml-to-ledger-dat/lib/subtractors"
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

	// Build maps with key "YYYY-MM-DD-AMOUNT" for set operations
	parser1Map := parser1.BuildTransactionMap()
	parser2Map := parser2.BuildTransactionMap()

	// Create a subtractor instance for set differences
	subtractor := subtractors.New()

	// Compute set differences: A-B and B-A
	diffAminusB := subtractor.Subtract(parser1Map, parser2Map)
	diffBminusA := subtractor.Subtract(parser2Map, parser1Map)

	// Output results
	fmt.Println("\n=== Transactions in File 1 but not in File 2 ===")
	for key, txns := range diffAminusB {
		for _, txn := range txns {
			fmt.Printf("%s | %s | %s\n", key, txn.Account, txn.Description)
		}
	}

	fmt.Println("\n=== Transactions in File 2 but not in File 1 ===")
	for key, txns := range diffBminusA {
		for _, txn := range txns {
			fmt.Printf("%s | %s | %s\n", key, txn.Account, txn.Description)
		}
	}

	return nil
}
