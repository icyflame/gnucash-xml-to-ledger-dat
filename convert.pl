use strict;
binmode(STDOUT, ":utf8");

use XML::Simple;

# dev + debug
use Data::Dumper;


sub help_text {
    print "
perl script.pl input.xml

input.xml is the GnuCash XML file (after gzip decompression)

Commodity names will be enclosed in double-quotes. This allows them to contain any character
(including white space, digits, colons, etc.)

";
}

my $input_file = shift;
if (!$input_file) {
    help_text;
    die "Must provide input file name: perl script.pl input.xml [sub repl]"
}

if ($input_file eq "--help" || $input_file eq "-h") {
    help_text;
    exit 0;
}

my $parsed = XMLin($input_file, ForceArray => [ "trn:split" ]);

my @accounts_list = @{$parsed->{'gnc:book'}->{'gnc:account'}};

my %accounts_map;

while (my $account = shift(@accounts_list)) {
    $accounts_map{$account->{'act:id'}->{'content'}} = $account;
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
    my $transaction_commodity = $transaction->{'trn:currency'}->{'cmdty:id'};

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

        my $split_amount = $quantity;

        # Bug: In a Stock purchase transaction, the "Edit Exchange Rate" dialog is not shown by
        # GnuCash. This seems to be a limitation within the program. If the quantity for a
        # transaction is zero, then it should attempt to use "split:value" instead. If that is zero
        # too, then this split is a no-op and does not need to be included in the resulting Ledger
        # file.  However, for correctness, I *will* include the "0 COMMODITY" split in the output
        # Ledger file.

        my $value = eval($split->{'split:value'});
        if ($quantity == 0 && $value != 0) {
            $split_amount = $value;
            $commodity = $transaction_commodity;
        }

        # Commodity names can have any character (including white space) if they are enclosed in double quotes.
        #
        # See section 4.5.1 Naming Commodities in the Ledger manual: https://ledger-cli.org/doc/ledger3.pdf
        $commodity = '"' . $commodity . '"';

        # TODO: For Stock purchases and sales, this should also have the price at which stocks were
        # purchased, which will make it possible to use Ledger.
        #
        # See 4.5.2 Buying and Selling Stock in the Ledger manual: https://ledger-cli.org/doc/ledger3.pdf
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
#
# Sample Stock purchase transaction:
#
# 2004/05/01 Stock purchase
#   Assets:Broker                50 AAPL @ $30.00
#   Expenses:Broker:Commissions  $19.95
#   Assets:Broker                $-1,519.95
#
# The amount can be left out of the last posting:
#
# 2004/05/01 Stock purchase
#   Assets:Broker                50 AAPL @ $30.00
#   Expenses:Broker:Commissions  $19.95
#   Assets:Broker
