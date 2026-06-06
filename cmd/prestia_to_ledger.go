package cmd

import (
	"fmt"
	"os"

	"github.com/icyflame/gnucash-xml-to-ledger-dat/lib/parsers/prestia"
	"github.com/icyflame/gnucash-xml-to-ledger-dat/internal/transformers"
	"github.com/spf13/cobra"
)

var prestiaToLedgerCmd = &cobra.Command{
	Use:   "prestia-to-ledger <input-prestia-csv> [output-ledger]",
	Short: "Convert Prestia CSV to Ledger DAT format",
	Long: `Convert Prestia credit card transaction CSV file to Ledger's dat file.

All transactions will use the generic commodity "COMM" and counterparty account "Expenses".
Empty records and transactions with invalid amounts are automatically skipped.`,
	Args: cobra.RangeArgs(1, 2),
	RunE: runPrestiaConvert,
}

func runPrestiaConvert(cmd *cobra.Command, args []string) error {
	inputFile := args[0]

	// Open and parse the Prestia CSV file (use stdin if inputFile is "-")
	var file *os.File
	var err error
	if inputFile == "-" {
		file = os.Stdin
	} else {
		file, err = os.Open(inputFile)
		if err != nil {
			return fmt.Errorf("failed to open input file: %w", err)
		}
		defer file.Close()
	}

	records, err := prestia.ParsePrestia(file)
	if err != nil {
		return fmt.Errorf("failed to parse input file: %w", err)
	}

	// Create transformer and convert to Ledger format
	transformer := transformers.NewPrestiaLedgerTransformer(records)

	var output *os.File
	if len(args) == 2 {
		var err error
		output, err = os.Create(args[1])
		if err != nil {
			return fmt.Errorf("failed to create output file: %w", err)
		}
		defer output.Close()
	} else {
		output = os.Stdout
	}

	if err := transformer.Write(output); err != nil {
		return fmt.Errorf("failed to write output: %w", err)
	}

	return nil
}
