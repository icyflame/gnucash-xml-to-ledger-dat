package cmd

import (
	"github.com/spf13/cobra"
)

var diffLedgerCmd = &cobra.Command{
	Use:   "diff-ledger",
	Short: "Compare 2 Ledger files using various subcommand outputs",
}

func init() {
	diffLedgerCmd.AddCommand(diffLedgerRegisterCmd)
}
