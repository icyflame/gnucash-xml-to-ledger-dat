package presenters

import (
	"fmt"
	"io"

	"github.com/icyflame/gnucash-xml-to-ledger-dat/lib/groupers"
)

func PresentRegisterDiff(w io.Writer, sideABuckets, sideBBuckets []groupers.Bucket) {
	if len(sideABuckets) != len(sideBBuckets) {
		panic(fmt.Sprintf("logic error: sideA (%d) and sideB (%d) must be of the same size", len(sideABuckets), len(sideBBuckets)))
	}

	for i := range len(sideABuckets) {
		setABucket := sideABuckets[i]
		setBBucket := sideBBuckets[i]
		setA := setABucket.Transactions
		setB := setBBucket.Transactions

		if max(len(setA), len(setB)) == 0 {
			continue
		}

		fmt.Fprintf(w, "\n- === %s === Count: %d\n", setABucket.Name, len(setA))
		for _, txn := range setA {
			fmt.Fprintf(w, "- %s | %s | %s | %s\n", txn.Date, txn.Amount, txn.Account, txn.Description)
		}

		fmt.Fprintf(w, "\n- === %s === Count: %d\n", setBBucket.Name, len(setB))
		for _, txn := range setB {
			fmt.Fprintf(w, "+ %s | %s | %s | %s\n", txn.Date, txn.Amount, txn.Account, txn.Description)
		}
	}
}
