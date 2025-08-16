# GNUCash's XML to Ledger's Dat File

> A script to convert GNUCash's XML file to Ledger's dat file

[convert.pl](./convert.pl) is the script to convert GnuCash's uncompressed XML
file to a Ledger journal. I run it using Perl 5.26 on a computer running Ubuntu
18.04.

The file formats for ledger and hledger are very similar, I have been able to
generate repors using hledger for the dat file created by the conversion script.

**Note:** There are [other scripts][1] which do the same thing as the Perl
script in this repository.

## Usage

`convert.pl` requires the Perl module `XML::Simple` as a dependency. This script can be run without
installing anything using Docker:

``` sh
# Build the Docker image using the included Dockerfile
$ docker build . -t $(basename $(pwd))

# Use the built Docker image to run the conversion command
$ docker run --name 'gnucash-xml-to-ledger-dat' \
  -it --rm \
  --volume "$GNUCASH_FILE:/input.gnucash:ro" \
  --volume /tmp:/output:rw \
  -w /src $(basename $(pwd)) /input.gnucash /output/output.dat
```

## Recipes

### Don't install anything

These oneliners will run Ledger with a read-only file system connected to the
current working directory. The read-only file system ensures that the container
can not change any of the files that are in the current working directory.

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
