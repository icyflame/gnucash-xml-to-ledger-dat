# Implementation Plan: GNUCash XML to Ledger Converter in Go

## Overview

Rewrite the Perl-based GNUCash XML to Ledger converter in Golang using Cobra CLI framework. All new code will be placed in a `go/` directory.

## CLI Interface

```bash
# Basic usage (output to stdout)
gnucash-converter <input.xml>

# Output to file
gnucash-converter <input.xml> <output.dat>

# With verbose debug output
gnucash-converter --verbose <input.xml>
gnucash-converter -v <input.xml> <output.dat>

# Help
gnucash-converter --help
```

## Project Structure

```
go/
├── go.mod
├── go.sum
├── main.go                 # Entry point, Cobra setup
├── cmd/
│   └── root.go            # Root command with convert logic
├── internal/
│   ├── parser/
│   │   └── gnucash.go     # XML parsing & account hierarchy
│   └── writer/
│       └── ledger.go      # Ledger format output
└── README.md
```

## Implementation Tasks

### 1. Setup Go module and dependencies
- Initialize `go.mod` with module name
- Add Cobra dependency

### 2. Create data models (inline in parser package)
- GNUCash XML structure matching the actual XML elements
- Account, Transaction, Split structs with proper XML tags

### 3. Implement XML parser (`internal/parser/gnucash.go`)
- `Parse(filename string)` - read and unmarshal XML
- `BuildAccountMap()` - create ID -> Account lookup
- `GetAccountName(accountID string)` - recursively build full path (e.g., "Assets:Bank:Checking")

### 4. Implement Ledger writer (`internal/writer/ledger.go`)
- `WriteLedger(writer io.Writer, transactions, accounts)` - main output function
- Format: `YYYY/MM/DD  Description`
- Format splits: `  AccountName  "COMMODITY"  amount`
- Handle edge cases:
  - Empty descriptions → use transaction ID
  - Zero quantity splits → use transaction value
  - Quote all commodity names

### 5. Create Cobra CLI (`main.go` + `cmd/root.go`)
- Root command accepts 1-2 positional args
- `--verbose` / `-v` flag for debug output
- `--help` / `-h` built-in from Cobra
- Output to stdout if no second arg, else to file
- Proper error handling and exit codes

### 6. Add README.md for Go version
- Building: `go build -o gnucash-converter`
- Usage examples
- Testing instructions

## Key Features

✅ UTF-8 output support  
✅ Recursive account hierarchy  
✅ Commodity name quoting  
✅ Zero-quantity split handling (use value field)  
✅ Empty description handling  
✅ Debug output to stderr (verbose mode)  
✅ Error handling (file not found, invalid XML, parse errors)  

## Deferred Features (TODO in original)

⏸️ Stock purchase price notation (`@ $30.00`)  
⏸️ Transaction second line as comments  

## Simplifications from Perl version

- **Cleaner error handling**: Go's explicit error returns vs Perl's die
- **Type safety**: Structs instead of hashmaps
- **No regex for HASH detection**: Direct type checking
- **Standard library XML**: No external XML dependency needed

## Testing Approach

Once TestData submodule is initialized:

```bash
go build -o gnucash-converter
./gnucash-converter TestData/TestBook-GnuCash/TestBook.gnucash /tmp/actual.dat
diff /tmp/actual.dat TestData/TestBook-Ledger/TestBook.ledger.dat
```

## Design Decisions

### Positional Arguments
Using positional arguments for input and output files (not flag-based) to keep the interface simple and aligned with the original Perl script.

### Output Destination
- **Default**: Output to stdout (allows piping and redirection)
- **Optional**: Second positional argument specifies output file path

### Debug Output
- **Default**: Clean output only
- **Optional**: `--verbose` / `-v` flag enables debug information to stderr (similar to line 108 in the Perl script)

## Implementation Order

1. Create directory structure and go.mod
2. Define XML parsing structs matching GNUCash format
3. Implement account hierarchy building logic
4. Implement transaction processing logic
5. Implement Ledger format output writer
6. Create main.go with Cobra CLI setup
7. Add help text and usage documentation
8. Test with TestData (if available)
