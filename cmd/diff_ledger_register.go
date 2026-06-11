package cmd

import (
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

	// TODO: Implement diff logic for CSV files
	_ = csvFile1
	_ = csvFile2

	return nil
}
