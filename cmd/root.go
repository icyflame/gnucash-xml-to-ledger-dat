package cmd

import (
	"fmt"
	"os"

	"github.com/icyflame/gnucash-xml-to-ledger-dat/internal/parser"
	"github.com/icyflame/gnucash-xml-to-ledger-dat/internal/writer"
	"github.com/spf13/cobra"
)

// TODO: Don't use a global variable for this flag. It can be passed down to the parsers, and other
// objects which are created.
var verbose bool

var rootCmd = &cobra.Command{
	Use:   "gnucash-converter <input.xml> [output.dat]",
	Short: "Convert GNUCash XML to Ledger DAT format",
	Long: `Convert GNUCash's uncompressed XML file to Ledger's dat file.

Commodity names will be enclosed in double-quotes. This allows them to contain
any character (including white space, digits, colons, etc.)`,
	Args: cobra.RangeArgs(1, 2),
	RunE: runConvert,
}

func init() {
	rootCmd.Flags().BoolVarP(&verbose, "verbose", "v", false, "Enable verbose debug output to stderr")
}

func Execute() error {
	return rootCmd.Execute()
}

func runConvert(cmd *cobra.Command, args []string) error {
	inputFile := args[0]

	p := parser.New()
	p.Verbose = verbose

	if err := p.Parse(inputFile); err != nil {
		return fmt.Errorf("failed to parse input file: %w", err)
	}

	w := writer.New(p, verbose)

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

	if err := w.Write(output); err != nil {
		return fmt.Errorf("failed to write output: %w", err)
	}

	return nil
}
