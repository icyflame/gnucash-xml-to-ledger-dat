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
    echo "ERROR: Output file must be specified" >&2
    usage
    exit 44
fi

cp -v "$GNUCASH_FILE" /tmp/input.xml.gz && \
    gzip -d /tmp/input.xml.gz && \
    perl /src/convert.pl /tmp/input.xml > $OUTPUT_FILE && \
    rm -fv /tmp/input.xml
