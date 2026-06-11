package cmd

import (
	"fmt"
	"regexp"
	"slices"

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

	// Build maps with key "YYYY-MM-DD-AMOUNT" for set operations
	set1 := buildTransactionMap(transactions1)
	set2 := buildTransactionMap(transactions2)

	// Compute set differences: A-B and B-A
	diffAminusB := setSubtraction(set1, set2)
	diffBminusA := setSubtraction(set2, set1)

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

// buildTransactionMap creates a map with key "YYYY-MM-DD-AMOUNT" for each transaction.
// The amount is extracted without the commodity (numeric value and sign only).
// Multiple transactions with the same key are stored as a slice.
func buildTransactionMap(transactions []ledger.RegisterTransaction) map[string][]ledger.RegisterTransaction {
	txnMap := make(map[string][]ledger.RegisterTransaction)

	for _, txn := range transactions {
		key := extractTransactionKey(txn)
		txnMap[key] = append(txnMap[key], txn)
	}

	return txnMap
}

// extractTransactionKey extracts the numeric amount from the Amount field and creates a key
// in the format "YYYY-MM-DD-AMOUNT". The amount includes the sign (+ or -) but no commodity.
func extractTransactionKey(txn ledger.RegisterTransaction) string {
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

// setSubtraction returns the set difference A-B.
// It returns all transactions from setA that either:
// 1. Don't exist in setB (key not found)
// 2. Exist in setB but with different counts or amounts
func setSubtraction(setA, setB map[string][]ledger.RegisterTransaction) map[string][]ledger.RegisterTransaction {
	diff := make(map[string][]ledger.RegisterTransaction)

	for key, txnsA := range setA {
		txnsB, exists := setB[key]

		// If key doesn't exist in B, add all transactions from A
		if !exists {
			diff[key] = txnsA
			continue
		}

		// If key exists, check if counts match
		if len(txnsA) != len(txnsB) {
			diff[key] = txnsA
			continue
		}

		// Counts match, verify all amounts are the same
		if !allAmountsMatch(txnsA, txnsB) {
			diff[key] = txnsA
		}
	}

	return diff
}

// allAmountsMatch verifies that all transactions in both slices have matching amounts.
// Transactions are sorted by Amount before comparison to ensure order-independent matching.
func allAmountsMatch(txnsA, txnsB []ledger.RegisterTransaction) bool {
	if len(txnsA) != len(txnsB) {
		return false
	}

	// Create copies to avoid modifying the original slices
	sortedA := slices.Clone(txnsA)
	sortedB := slices.Clone(txnsB)

	// Sort both slices by Amount
	slices.SortFunc(sortedA, func(a, b ledger.RegisterTransaction) int {
		if a.Amount < b.Amount {
			return -1
		} else if a.Amount > b.Amount {
			return 1
		}
		return 0
	})

	slices.SortFunc(sortedB, func(a, b ledger.RegisterTransaction) int {
		if a.Amount < b.Amount {
			return -1
		} else if a.Amount > b.Amount {
			return 1
		}
		return 0
	})

	// Compare sorted slices
	for i := range sortedA {
		if sortedA[i].Amount != sortedB[i].Amount {
			return false
		}
	}

	return true
}
