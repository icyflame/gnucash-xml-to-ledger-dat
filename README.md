# Personal Finance Transformer

> A Cobra CLI to transform various personal finance file types (GnuCash, Ledger format, credit card
> statement CSV files)

You can run subcommands within the `main.go` file to convert GnuCash's uncompressed XML file to a
Ledger journal. This script was originally written using Perl, but maintaining Perl dependencies is
neither fun nor something that I am interested in anymore. So, the script has been re-written to
Golang; [primarily using AI].[^1]

The file formats for [`ledger`] and [`hledger`] are nearly identical. So, the output from this
script can be used with either program.

## Usage

### GnuCash to Ledger

``` shell
# Decompress the GnuCash file and write to an input XML file
zcat input.gnucash > input.xml

# Use this script to convert the input into a Ledger file
go run main.go gnucash-to-ledger input.xml output.dat

# Confirm that the output is valid
hledger -f output.dat bal
```

### Prestia to Ledger

``` sh
$ cat ./TestData/TestBook-Prestia/202509.csv ./TestData/TestBook-Prestia/202510.csv | go run main.go prestia-to-ledger -  | head -10
2026/06/12 12:40:25 WARN invalid line in Prestia statement > number of fields mismatch, want 7 num_fields=3 fields="[ジョナサン\u3000太郎\u3000様 1234-00**-****-**** プレスティアＶｉｓａゴールド]"
2026/06/12 12:40:25 WARN invalid line in Prestia statement > number of fields mismatch, want 7 num_fields=3 fields="[エリザベス\u3000太郎\u3000様 1235-11**-****-**** ＧｏｏｇｌｅＰａｙ／ｉＤ]"
2026/06/12 12:40:25 WARN invalid line in Prestia statement > first field is not a date in YYYY/MM/DD format first_field="" parse_error="parsing time \"\" as \"2006/01/02\": cannot parse \"\" as \"2006\""
2026/06/12 12:40:25 WARN invalid line in Prestia statement > number of fields mismatch, want 7 num_fields=3 fields="[ジョナサン\u3000太郎\u3000様 1234-00**-****-**** プレスティアＶｉｓａゴールド]"
2026/06/12 12:40:25 WARN invalid line in Prestia statement > number of fields mismatch, want 7 num_fields=3 fields="[エリザベス\u3000太郎\u3000様 1235-11**-****-**** ＧｏｏｇｌｅＰａｙ／ｉＤ]"
2026/06/12 12:40:25 WARN invalid line in Prestia statement > first field is not a date in YYYY/MM/DD format first_field="" parse_error="parsing time \"\" as \"2006/01/02\": cannot parse \"\" as \"2006\""
2025-08-20  飲食店
  Liabilities:Prestia Card  "JPY"  -600
  Expenses

2025-08-25  飲食店
  Liabilities:Prestia Card  "JPY"  -750
  Expenses

2025-09-04  飲食店
  Liabilities:Prestia Card  "JPY"  -900

```

### Diff Ledger Registers

``` sh
$ go run main.go diff-ledger register \
    <(dock_hledger -f ./TestData/TestBook-Prestia/prestia-statement-jpy.dat register --begin 2025-08 --end 2025-08-31 'Liabilities' -O csv | head -n10) \
    <(dock_hledger -f ./TestData/TestBook-Ledger/TestBook.ledger.dat register --begin 2025-08-01 --end 2025-08-31 'Liabilities' -O csv | head -n10)
Error: currency mismatch between file 1 and file 2: file 1 currency: JPY, file 2 currency: INR
Usage:
  personal-finance-transformer diff-ledger register <csv-file-1> <csv-file-2> [flags]

Flags:
  -h, --help   help for register

exit status 1

$ go run main.go diff-ledger register \
    <(dock_hledger -f ./TestData/TestBook-Prestia/prestia-statement-inr.dat register --begin 2025-08 --end 2025-08-31 'Liabilities' -O csv | head -n10) \
    <(dock_hledger -f ./TestData/TestBook-Ledger/TestBook.ledger.dat register --begin 2025-08-01 --end 2025-08-31 'Liabilities' -O csv | head -n10)
Parsed 3 transactions from /proc/self/fd/11
Parsed 6 transactions from /proc/self/fd/12

=== Transactions in File 1 but not in File 2 === Count:  0

=== Transactions in File 2 but not in File 1 === Count:  3
2025-08-02 | INR -3500 | Liabilities:Credit Card | Groceries
2025-08-20 | INR -750 | Liabilities:Credit Card | Electricity
2025-08-25 | INR -800 | Liabilities:Credit Card | T-shirts

```


## Testing

When making changes to the script, confirm that the changes that are made in the TestData are
logical.

``` shell
# Convert "./TestData/TestBook-GnuCash/TestBook.gnucash" into a Ledger file
# using the converter script
$ zcat TestData/TestBook-GnuCash/TestBook.gnucash | go run main.go gnucash-to-ledger - > /tmp/converted.dat

# Confirm that actual and expected are identical
$ diff ./TestData/TestBook-Ledger/TestBook.ledger.dat /tmp/converted.dat

# Oneliner:
$ diff TestData/TestBook-Ledger/TestBook.ledger.dat <(zcat TestData/TestBook-GnuCash/TestBook.gnucash | go run main.go gnucash-to-ledger -)

```

## Recipes

### Don't install anything

These oneliners will run Ledger with a read-only file system connected to the current working
directory. The read-only file system ensures that the container can not change any of the files that
are in the current working directory. This might be useful to you if you store the GnuCash and
Ledger files in the same directory.

```sh
hledger () {
        docker run --rm --volume $(pwd):/data:ro -w /data dastapov/hledger hledger $@
}

ledger () {
        docker run --rm --volume $(pwd):/data:ro -w /data dcycle/ledger:1 $@
}
```

### Display numbers in the correct format

Commodity directives can be used to ensure that numbers are always displayed in
the appropriate format.

```ledger
$ cat commodities.dat
commodity JPY 1,000,000.
commodity EUR 1,000,000.00
commodity USD 1,000,000.00
commodity INR 10,00,000.

; Enter the appropriate exchange rate for the various commodities here
P 2020-02-08 EUR 0.001 USD
```

### Generate common accounting reports

```sh
hledger -f simple.dat incomestatement
hledger -f simple.dat balancesheet
hledger -f simple.dat cashflow
```

### Display all balances in a single currency at present exchange rate

```sh
# You _must_ have the exchange rate in USD for all the commodities in the
# commodities.dat file
hledger -f simple.dat bal Assets -Y --value=now,USD
```

## License

Code in this repository is licensed under MIT.

Copyright (C) 2020  Siddharth Kannan <mail@siddharthkannan.in>

[1]: https://gist.github.com/nonducor/ddc97e787810d52d067206a592a35ea7
[primarily using AI]: ./ai-assisted/sessions/rewrite-in-go
[`ledger`]: https://ledger-cli.org/
[`hledger`]: https://hledger.org/index.html

[^1]: OpenCode with Claude Sonnet 4.5
