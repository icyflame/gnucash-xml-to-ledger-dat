package presenters

import (
	"encoding/csv"

	"github.com/icyflame/gnucash-xml-to-ledger-dat/lib/parsers/ledger"
)

func PresentTransactionsAsDiff(w *csv.Writer, sideA, sideB []ledger.RegisterTransaction) {
	for i := range max(len(sideB), len(sideA)) {
		var row []string

		var left, right *ledger.RegisterTransaction
		if i < len(sideA) {
			left = &sideA[i]
		}

		if i < len(sideB) {
			right = &sideB[i]
		}

		if left == nil {
			row = append(row, "", "", "", "")
		} else {
			row = append(row, left.Date, left.Amount, left.Account, left.Description)
		}

		row = append(row, "")

		if right == nil {
			row = append(row, "", "", "", "")
		} else {
			row = append(row, right.Date, right.Amount, right.Account, right.Description)
		}

		w.Write(row)
	}
}
