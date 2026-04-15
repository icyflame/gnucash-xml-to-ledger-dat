# GNUCash XML to Ledger Converter (Go)

A Go implementation of the GNUCash XML to Ledger DAT converter. This is a pure rewrite of the original Perl script with identical behavior.

## Overview

This tool converts GNUCash's uncompressed XML files to Ledger's dat format. The output format is compatible with both Ledger and hledger.

## Project Structure

```
go/
├── cmd/
│   └── root.go                 # Cobra CLI command
├── internal/
│   ├── parser/
│   │   └── gnucash.go         # XML parsing & account hierarchy
│   └── writer/
│       └── ledger.go          # Ledger format output
├── main.go                    # Entry point
├── go.mod                     # Go module
├── go.sum                     # Dependencies
└── README.md                  # This file
```

## Building

```bash
go build -o gnucash-converter
```

## Usage

```bash
# Convert to stdout (useful for piping)
./gnucash-converter input.xml

# Convert to file
./gnucash-converter input.xml output.dat

# With verbose debug output to stderr
./gnucash-converter -v input.xml output.dat
./gnucash-converter --verbose input.xml

# Show help
./gnucash-converter --help
```

## Testing

The input file must be an uncompressed GNUCash XML file. If you have a `.gnucash` file (which is gzip-compressed), decompress it first:

```bash
# Decompress the test file
gunzip -c ../TestData/TestBook-GnuCash/TestBook.gnucash > /tmp/TestBook.xml

# Run the converter
go build -o gnucash-converter
./gnucash-converter /tmp/TestBook.xml /tmp/actual.dat

# Compare with expected output (should show no differences)
diff /tmp/actual.dat ../TestData/TestBook-Ledger/TestBook.ledger.dat
```

## Key Features

- ✅ **XML Parsing** - Handles GNUCash XML with proper namespace support
- ✅ **Account Hierarchy** - Recursive building of full account paths (e.g., "Assets:Bank:Checking")
- ✅ **Fraction Parsing** - Converts "400000/100" format to decimal values
- ✅ **Ledger Output** - Exact format match with the original Perl version
- ✅ **Cobra CLI** - Clean command-line interface with flags
- ✅ **Verbose Mode** - Debug output to stderr showing account names, quantities, and values
- ✅ **Edge Cases** - Handles empty descriptions and zero-quantity splits
- ✅ **UTF-8 Support** - Proper character encoding

## Implementation Details

### Output Format

- Date format: `YYYY-MM-DD` (e.g., `2025-09-01`)
- Transaction format: `DATE  DESCRIPTION`
- Split format: `  ACCOUNT_NAME  "COMMODITY"  AMOUNT`
- Two spaces separate each field
- Blank line between transactions
- All commodity names are enclosed in double-quotes

### Special Handling

- **Empty descriptions**: Replaced with "Empty description (TRANSACTION_ID)"
- **Zero-quantity splits**: Uses the transaction value field instead
- **Account hierarchy**: Built recursively from parent references, skipping ROOT type

### Code Statistics

- **~260 lines of Go code** (original Perl: ~154 lines)
- **7 files total**
- **Clean, readable, maintainable**
- **Type-safe with explicit error handling**

## Deferred Features

The following features are marked as TODO in the original Perl script and are not implemented:

- Stock purchase price notation (`@ $30.00`)
- Transaction second line as comments

## Dependencies

- `github.com/spf13/cobra` - CLI framework
- Go standard library for XML parsing and file I/O

## License

Code in this directory follows the same license as the parent repository (MIT).

Copyright (C) 2020  Siddharth Kannan <mail@siddharthkannan.in>
