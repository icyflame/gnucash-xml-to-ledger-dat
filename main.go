package main

import (
	"os"

	"github.com/icyflame/gnucash-xml-to-ledger-dat/cmd"
)

func main() {
	if err := cmd.Execute(); err != nil {
		os.Exit(1)
	}
}
