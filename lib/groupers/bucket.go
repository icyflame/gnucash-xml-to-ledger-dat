package groupers

import "github.com/icyflame/gnucash-xml-to-ledger-dat/lib/parsers/ledger"

type Bucket struct {
	Name         string
	Transactions []ledger.RegisterTransaction
}
