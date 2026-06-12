package subtractors

import (
	"slices"

	"github.com/icyflame/gnucash-xml-to-ledger-dat/lib/parsers/ledger"
)

// RegisterSubtractor handles set subtraction operations on Ledger register transactions.
type RegisterSubtractor struct{}

// New creates a new RegisterSubtractor instance.
func New() *RegisterSubtractor {
	return &RegisterSubtractor{}
}

// Subtract returns the set difference A-B.
// It returns all transactions from setA that either:
// 1. Don't exist in setB (key not found)
// 2. Exist in setB but with different counts or amounts
func (rs *RegisterSubtractor) Subtract(setA, setB map[string][]ledger.RegisterTransaction) map[string][]ledger.RegisterTransaction {
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
		if !rs.allAmountsMatch(txnsA, txnsB) {
			diff[key] = txnsA
		}
	}

	return diff
}

func transactionSorter(a, b ledger.RegisterTransaction) int {
	if a.Amount < b.Amount {
		return -1
	} else if a.Amount > b.Amount {
		return 1
	}
	return 0
}

// allAmountsMatch verifies that all transactions in both slices have matching amounts.
// Transactions are sorted by Amount before comparison to ensure order-independent matching.
func (rs *RegisterSubtractor) allAmountsMatch(txnsA, txnsB []ledger.RegisterTransaction) bool {
	if len(txnsA) != len(txnsB) {
		return false
	}

	sortedA := slices.Clone(txnsA)
	sortedB := slices.Clone(txnsB)

	slices.SortFunc(sortedA, transactionSorter)
	slices.SortFunc(sortedB, transactionSorter)

	for i := range sortedA {
		if sortedA[i].Amount != sortedB[i].Amount {
			return false
		}
	}

	return true
}
