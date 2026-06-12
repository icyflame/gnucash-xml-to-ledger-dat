package groupers

import (
	"log/slog"
	"time"

	"github.com/icyflame/gnucash-xml-to-ledger-dat/lib/parsers/ledger"
)

// Monthly groups transactions by month for both sides A and B.
// It returns two slices of Buckets with the same length, one for side A and one for side B.
// Each Bucket represents transactions from a given month.
// Assumes transactions in both slices are already sorted by date.
func Monthly(a, b []ledger.RegisterTransaction) ([]Bucket, []Bucket) {
	// Find the start date (minimum of dates in both sets) - first transaction in sorted order
	var startDate time.Time
	if len(a) > 0 {
		d, _ := time.Parse("2006-01-02", a[0].Date)
		startDate = d
	}
	if len(b) > 0 {
		d, _ := time.Parse("2006-01-02", b[0].Date)
		if startDate.IsZero() || d.Before(startDate) {
			startDate = d
		}
	}

	// Find the end date (maximum of dates in both sets) - last transaction in sorted order
	var endDate time.Time
	if len(a) > 0 {
		d, _ := time.Parse("2006-01-02", a[len(a)-1].Date)
		endDate = d
	}
	if len(b) > 0 {
		d, _ := time.Parse("2006-01-02", b[len(b)-1].Date)
		if d.After(endDate) {
			endDate = d
		}
	}

	// Group A side transactions by month
	sideA := groupByMonth(a, startDate, endDate, "A-side")
	// Group B side transactions by month
	sideB := groupByMonth(b, startDate, endDate, "B-side")

	return sideA, sideB
}

// groupByMonth groups transactions by month, creating a bucket for each month between start and end dates.
// Assumes transactions are already sorted by date.
func groupByMonth(txns []ledger.RegisterTransaction, startDate, endDate time.Time, bucketNameSuffix string) []Bucket {
	var buckets []Bucket
	var txnIndex int

	// Create a bucket for each month between start and end date
	currentDate := startDate
	for !currentDate.After(endDate) {
		monthStart := time.Date(currentDate.Year(), currentDate.Month(), 1, 0, 0, 0, 0, time.UTC)
		nextMonthStart := monthStart.AddDate(0, 1, 0)

		var monthTransactions []ledger.RegisterTransaction
		for txnIndex < len(txns) {
			txnDate, err := time.Parse("2006-01-02", txns[txnIndex].Date)
			if err != nil {
				slog.Error("transaction date could not be pared", slog.String("date", txns[txnIndex].Date))
				txnIndex++
				continue
			}

			if txnDate.Before(monthStart) {
				txnIndex++
				continue
			}

			if txnDate.Before(nextMonthStart) {
				monthTransactions = append(monthTransactions, txns[txnIndex])
				txnIndex++
			} else {
				break
			}
		}

		bucketName := monthStart.Format("2006-01")
		buckets = append(buckets, Bucket{
			Name:         bucketName + " " + bucketNameSuffix,
			Transactions: monthTransactions,
		})

		currentDate = nextMonthStart
	}

	return buckets
}
