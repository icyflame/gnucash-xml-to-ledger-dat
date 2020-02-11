# GNUCash's XML to Ledger's Dat File

> A script to convert GNUCash's XML file to Ledger's dat file

**Note:** There are [other scripts][1] which do the same thing as the Perl
script in this repository.

## Recipes

### Don't install anything

```sh
hledger () {
        docker run --rm -v `pwd`:/data -w /data dastapov/hledger hledger $@
}

ledger () {
        docker run --rm -v `pwd`:/data -w /data dcycle/ledger:1 $@
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
