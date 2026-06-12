package presenters

import (
	"encoding/csv"
	"fmt"

	"github.com/icyflame/gnucash-xml-to-ledger-dat/lib/groupers"
	"github.com/icyflame/gnucash-xml-to-ledger-dat/lib/parsers/ledger"
)

func PresentTransactionsDiffAsTSV(w *csv.Writer, sideABuckets, sideBBuckets []groupers.Bucket) {
	if len(sideABuckets) != len(sideBBuckets) {
		panic(fmt.Sprintf("logic error: sideA (%d) and sideB (%d) must be of the same size", len(sideABuckets), len(sideBBuckets)))
	}

	for bucketCount := range len(sideABuckets) {
		sideA := sideABuckets[bucketCount]
		sideB := sideBBuckets[bucketCount]

		w.Write([]string{
			sideA.Name, "", "", "",
			"",
			sideB.Name, "", "", "",
		})

		w.Write([]string{
			"", "", "", "",
			"",
			"", "", "", "",
		})

		w.Write([]string{
			"Date", "Amount", "Account", "Description",
			"",
			"Date", "Amount", "Account", "Description",
		})

		for i := range max(len(sideB.Transactions), len(sideA.Transactions)) {
			var row []string

			var left, right *ledger.RegisterTransaction
			if i < len(sideA.Transactions) {
				left = &sideA.Transactions[i]
			}

			if i < len(sideB.Transactions) {
				right = &sideB.Transactions[i]
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

		w.Write([]string{
			"", "", "", "",
			"",
			"", "", "", "",
		})
	}
}
