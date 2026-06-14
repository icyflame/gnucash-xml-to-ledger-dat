package groupers

import (
	"time"

	"github.com/icyflame/gnucash-xml-to-ledger-dat/lib/parsers/ledger"
)

// Monthly groups transactions by month for both sides A and B.
// It returns two slices of Buckets with the same length, one for side A and one for side B.
// Each Bucket represents transactions from a given month.
// Assumes transactions in both slices are already sorted by date.
func Monthly(a, b []ledger.RegisterTransaction) ([]Bucket, []Bucket) {
	startDate, endDate := dateRange(a, b)
	bucketStart := func(d time.Time) time.Time {
		return time.Date(d.Year(), d.Month(), 1, 0, 0, 0, 0, time.UTC)
	}
	nextBucketStart := func(d time.Time) time.Time {
		return d.AddDate(0, 1, 0)
	}
	bucketName := func(suffix string) func(time.Time) string {
		return func(d time.Time) string {
			return d.Format("2006-01") + " " + suffix
		}
	}
	sideA := groupByBucket(a, startDate, endDate, bucketStart, nextBucketStart, bucketName("A-side"))
	sideB := groupByBucket(b, startDate, endDate, bucketStart, nextBucketStart, bucketName("B-side"))
	return sideA, sideB
}
