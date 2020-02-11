use strict;
binmode(STDOUT, ":utf8");

use XML::Simple;

# dev + debug
use Data::Dumper;

my $input_file = shift;
if (!$input_file) {
    die "Must provide input file name: perl script.pl input.xml [sub repl]"
}

if ($input_file eq "--help" || $input_file eq "-h") {
    print "
perl script.pl input.xml [sub repl] [sub repl]

input.xml is the GnuCash XML file (after gzip decompression)

[sub repl] pairs are pairs of regular expressions which will be used to
change the name of the commodity. The transformation s/sub/repl/g will
be run on every single commodity name. You can use this to remove digits
from the commodity name, or change values to readable values.

Eg: To convert AAPL => Apple, you can use this command:

    perl script.pl input.xml AAPL Apple
";

    exit 0;
}

my @replacements;
while (my $sub = shift) {
    my $repl = shift;
    push(@replacements, ($sub, $repl))
}

my $parsed = XMLin($input_file, ForceArray => [ "trn:split" ]);

my @accounts_list = @{$parsed->{'gnc:book'}->{'gnc:account'}};

my %accounts_map;

while (my $account = shift(@accounts_list)) {
    $accounts_map{$account->{'act:id'}->{'content'}} = $account;
}

sub fix_commodity_name {
    my $base = shift;
    for (my $i = 0; $i < @replacements; $i++) {
        my $replacements = @replacements[$i];
        $base =~ s/$replacements[0]/$replacements[1]/g;
    }
    return $base;
}

sub get_account_name {
    my $act = shift;
    my $base = shift;

    my $type = $act->{'act:type'};
    my $name = $act->{'act:name'};

    if ($type eq "ROOT") {
        # Base case
        return $base;
    }

    my $this = sprintf ("%s", $name);
    if (defined $base && $base ne "") {
        # Recurse
        $this = sprintf ("%s:%s", $this, $base);
    }

    if (exists $act->{'act:parent'}) {
        my $account = $accounts_map{$act->{'act:parent'}->{'content'}};
        return get_account_name ($account, $this);
    }

    # Fallback
    return $base
}

my @transactions = @{$parsed->{'gnc:book'}->{'gnc:transaction'}};

while (my $transaction = shift(@transactions)) {
    my $date = substr($transaction->{'trn:date-posted'}->{'ts:date'}, 0, 10);
    my $description = sprintf("%s", $transaction->{'trn:description'});
    my $id = $transaction->{'trn:id'}->{'content'};

    # When description is empty, it looks like this after conversion to a string
    if ($description =~ /HASH\(0x[0-9a-f]+\)/) {
        $description = "Empty description ($id)";
    }

    my @splits = @{$transaction->{'trn:splits'}->{'trn:split'}};

    if (@splits == 0) {
        continue
    }

    printf ("%s  %s\n", $date, $description eq "" ? "Empty description ($id)" : $description);

    while (my $split = shift(@splits)) {
        my $quantity = eval($split->{'split:quantity'});

        my $account = $accounts_map{$split->{'split:account'}->{'content'}};
        my $account_name = get_account_name ($account);
        my $commodity = $account->{'act:commodity'}->{'cmdty:id'};
        $commodity =~ s/(^\d)/S_\1/g;
        $commodity =~ s/:/_/g;

        $commodity = fix_commodity_name $commodity;

        print "  " . join("  ",
            $account_name,
            $commodity,
            $quantity,
        );
        print "\n";
    }

    print "\n";
}

# Sample Ledger transaction:
# 2015/10/12 Exxon
#     Expenses:Auto:Gas         $10.00
#     Liabilities:MasterCard   $-10.00
