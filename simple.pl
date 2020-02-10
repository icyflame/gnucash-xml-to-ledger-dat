use strict;
binmode(STDOUT, ":utf8");

use XML::Simple;

# dev + debug
use Data::Dumper;

my $parsed = XMLin("input.gnucash");

my @transactions = @{$parsed->{'gnc:book'}->{'gnc:transaction'}};

while (my $transaction = shift(@transactions)) {
    print join("\t",
        substr($transaction->{'trn:date-posted'}->{'ts:date'}, 0, 10),
        $transaction->{'trn:currency'}->{'cmdty:id'},
        $transaction->{'trn:description'},
    );

    print "\n";
}

# print Dumper $parsed->{'gnc:book'}->{'gnc:transaction'};
# print Dumper @{$parsed->{'gnc:book'}};

# Sample Ledger transaction:
# 2015/10/12 Exxon
#     Expenses:Auto:Gas         $10.00
#     Liabilities:MasterCard   $-10.00


