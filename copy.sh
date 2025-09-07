#!/bin/bash

usage() {
    cat <<EOF >&2
$0 path-to-gnucash-file [output-file]

    output-file is an optional argument. If not specified, the output ledger
    file will be written to $(pwd)/output.dat.
EOF
}

if [[ -z "$1" || "$1" == "-h" || "$1" == "--help" ]]; then
    usage
    exit 0
fi

GNUCASH_FILE="$1"
if [[ ! -f "$GNUCASH_FILE" ]]; then
    echo "ERROR: Input GnuCash file must be a file" >&2
    usage
    exit 42
fi

OUTPUT_FILE="$2"
if [[ -z "$OUTPUT_FILE" ]]; then
    OUTPUT_FILE=$(pwd)/output.dat;
fi

if [[ -f "$OUTPUT_FILE" ]]; then
    echo "ERROR: Output file must not already exist." >&2
    usage
    exit 43
fi

CWD=$(pwd)
pushd /tmp

rm -vf input.xml && \
    cp -v "$GNUCASH_FILE" input.xml.gz && \
    gzip -d input.xml.gz && \
    perl ${CWD}/convert.pl input.xml > $OUTPUT_FILE && \
    rm -fv input.xml

if [[ -f "${CWD}/commodities.dat" ]]; then
    cat "${CWD}/commodities.dat" >> $OUTPUT_FILE
fi

popd
