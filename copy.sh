#!/bin/bash

GNUCASH_FILE="$1"

if [[ -z "$GNUCASH_FILE" ]]; then
    echo "ERROR: ./script gnucash-file-path"
    exit 42
fi

rm -f input.xml && \
    cp "$GNUCASH_FILE" input.xml.gz && \
    gzip -d input.xml.gz && \
    rm -f simple.dat && \
    perl simple.pl input.xml 4385 MERC > simple.dat && \
    rm input.xml

if [[ -f "commodities.dat" ]]; then
    cat commodities.dat >> simple.dat
fi
