# GnuCash to Ledger

> A Golang script to convert a GNUCash file into a Ledger file

You can run the `main.go` file to convert GnuCash's uncompressed XML file to a Ledger journal. This
script was originally written using Perl, but maintaining Perl dependencies is neither fun nor
something that I am interested in anymore. So, the script has been re-written to Golang; [primarily
using AI].[^1]

The file formats for [`ledger`] and [`hledger`] are nearly identical; so, the output from this script
can be used with either program.

## Usage

``` shell
# Decompress the GnuCash file and write to an input XML file
zcat input.gnucash > input.xml

# Use this script to convert the input into a Ledger file
go run main.go gnucash-to-ledger input.xml output.dat

# Confirm that the output is valid
hledger -f output.dat bal
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
