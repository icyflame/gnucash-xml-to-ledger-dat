package presenters

import (
	"fmt"
	"io"

	"github.com/icyflame/gnucash-xml-to-ledger-dat/lib/parsers/ledger"
)

func PresentRegisterDiff(w io.Writer, setA, setB []ledger.RegisterTransaction) {
	fmt.Fprintln(w, "\n- === Transactions in File 1 but not in File 2 ===", "Count: ", len(setA))
	for _, txn := range setA {
		fmt.Fprintf(w, "- %s | %s | %s | %s\n", txn.Date, txn.Amount, txn.Account, txn.Description)
	}

	fmt.Fprintln(w, "\n+ === Transactions in File 2 but not in File 1 ===", "Count: ", len(setB))
	for _, txn := range setB {
		fmt.Fprintf(w, "+ %s | %s | %s | %s\n", txn.Date, txn.Amount, txn.Account, txn.Description)
	}
}
