package groupers

import (
	"time"

	"github.com/icyflame/gnucash-xml-to-ledger-dat/lib/parsers/ledger"
)

// Weekly groups transactions by week for both sides A and B.
// It returns two slices of Buckets with the same length, one for side A and one for side B.
// Each Bucket represents transactions from a given week (Monday to Sunday).
// Assumes transactions in both slices are already sorted by date.
func Weekly(a, b []ledger.RegisterTransaction) ([]Bucket, []Bucket) {
	startDate, endDate := dateRange(a, b)
	bucketStart := func(d time.Time) time.Time {
		// Align to the Monday of the week containing d
		weekday := int(d.Weekday())
		if weekday == 0 {
			weekday = 7 // treat Sunday as day 7 so Monday is day 1
		}
		monday := d.AddDate(0, 0, -(weekday - 1))
		return time.Date(monday.Year(), monday.Month(), monday.Day(), 0, 0, 0, 0, time.UTC)
	}
	nextBucketStart := func(d time.Time) time.Time {
		return d.AddDate(0, 0, 7)
	}
	bucketName := func(suffix string) func(time.Time) string {
		return func(d time.Time) string {
			return d.Format("2006-01-02") + " " + suffix
		}
	}
	sideA := groupByBucket(a, startDate, endDate, bucketStart, nextBucketStart, bucketName("A-side"))
	sideB := groupByBucket(b, startDate, endDate, bucketStart, nextBucketStart, bucketName("B-side"))
	return sideA, sideB
}
