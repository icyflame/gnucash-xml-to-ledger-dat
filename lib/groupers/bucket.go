package groupers

import (
	"log/slog"
	"time"

	"github.com/icyflame/gnucash-xml-to-ledger-dat/lib/parsers/ledger"
)

type Bucket struct {
	Name         string
	Transactions []ledger.RegisterTransaction
}

// groupByBucket groups transactions into buckets defined by bucketStart, nextBucketStart, and
// bucketName functions. Assumes transactions are already sorted by date.
func groupByBucket(
	txns []ledger.RegisterTransaction,
	startDate, endDate time.Time,
	bucketStart func(time.Time) time.Time,
	nextBucketStart func(time.Time) time.Time,
	bucketName func(time.Time) string,
) []Bucket {
	var buckets []Bucket
	var txnIndex int

	currentDate := startDate
	for !currentDate.After(endDate) {
		bStart := bucketStart(currentDate)
		bEnd := nextBucketStart(bStart)

		var bucketTransactions []ledger.RegisterTransaction
		for txnIndex < len(txns) {
			txnDate, err := time.Parse("2006-01-02", txns[txnIndex].Date)
			if err != nil {
				slog.Error("transaction date could not be parsed", slog.String("date", txns[txnIndex].Date))
				txnIndex++
				continue
			}

			if txnDate.Before(bStart) {
				txnIndex++
				continue
			}

			if txnDate.Before(bEnd) {
				bucketTransactions = append(bucketTransactions, txns[txnIndex])
				txnIndex++
			} else {
				break
			}
		}

		buckets = append(buckets, Bucket{
			Name:         bucketName(bStart),
			Transactions: bucketTransactions,
		})

		currentDate = bEnd
	}

	return buckets
}

// dateRange returns the start and end dates across both sorted transaction slices.
func dateRange(a, b []ledger.RegisterTransaction) (startDate, endDate time.Time) {
	if len(a) > 0 {
		startDate, _ = time.Parse("2006-01-02", a[0].Date)
		endDate, _ = time.Parse("2006-01-02", a[len(a)-1].Date)
	}
	if len(b) > 0 {
		bStart, _ := time.Parse("2006-01-02", b[0].Date)
		bEnd, _ := time.Parse("2006-01-02", b[len(b)-1].Date)
		if startDate.IsZero() || bStart.Before(startDate) {
			startDate = bStart
		}
		if bEnd.After(endDate) {
			endDate = bEnd
		}
	}
	return startDate, endDate
}
