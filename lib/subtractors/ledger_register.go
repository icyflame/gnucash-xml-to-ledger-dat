package subtractors

import (
	"sort"

	"github.com/icyflame/gnucash-xml-to-ledger-dat/lib/parsers/ledger"
)

// RegisterSubtractor handles set subtraction operations on Ledger register transactions.
type RegisterSubtractor struct{}

// New creates a new RegisterSubtractor instance.
func New() *RegisterSubtractor {
	return &RegisterSubtractor{}
}

// Subtract returns the set difference A-B.
func (rs *RegisterSubtractor) Subtract(setA, setB map[string][]ledger.RegisterTransaction) []ledger.RegisterTransaction {
	var diff []ledger.RegisterTransaction

	for key, txnsA := range setA {
		txnsB, exists := setB[key]

		if !exists {
			diff = append(diff, txnsA...)
			continue
		}

		if len(txnsA) != len(txnsB) {
			diff = append(diff, txnsA...)
			continue
		}

		if !rs.allAmountsMatch(txnsA, txnsB) {
			diff = append(diff, txnsA...)
		}
	}

	return diff
}

func (rs *RegisterSubtractor) allAmountsMatch(txnsA, txnsB []ledger.RegisterTransaction) bool {
	if len(txnsA) != len(txnsB) {
		return false
	}

	sort.Sort(ledger.RegisterTransactionSlice(txnsA))
	sort.Sort(ledger.RegisterTransactionSlice(txnsB))

	for i := range txnsA {
		if txnsA[i].Amount != txnsB[i].Amount {
			return false
		}
	}

	return true
}
