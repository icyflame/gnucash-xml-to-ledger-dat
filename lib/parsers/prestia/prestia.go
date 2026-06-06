package prestia

import (
	"encoding/csv"
	"io"
	"log/slog"
	"strings"
	"time"
)

// PrestiaRecord represents a parsed transaction from a Prestia CSV file
type PrestiaRecord struct {
	Date        string
	Description string
	Amount      string
	Quantity    string
	Unit        string
	Total       string
	Notes       string
}

// ParsePrestia parses a Prestia CSV file and returns all valid transaction records.
// Only rows with exactly 7 fields that start with a valid date [0-9]{4}/[0-9]{2}/[0-9]{2} are parsed.
func ParsePrestia(r io.Reader) ([]PrestiaRecord, error) {
	var records []PrestiaRecord
	reader := csv.NewReader(r)
	reader.FieldsPerRecord = -1 // Allow varying field counts

	for {
		fields, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			// Skip lines with parse errors
			continue
		}

		// Only process lines with exactly 7 fields
		if len(fields) != 7 {
			slog.Warn("invalid line in Prestia statement > number of fields mismatch, want 7", slog.Int("num_fields", len(fields)), slog.Any("fields", fields))
			continue
		}

		// Check if first field is a valid date
		date := strings.TrimSpace(fields[0])
		if _, err := time.Parse("2006/01/02", date); err != nil {
			slog.Warn("invalid line in Prestia statement > first field is not a date in YYYY/MM/DD format", slog.String("first_field", fields[0]), slog.Any("parse_error", err))
			continue
		}

		record := PrestiaRecord{
			Date:        date,
			Description: strings.TrimSpace(fields[1]),
			Amount:      strings.TrimSpace(fields[2]),
			Quantity:    strings.TrimSpace(fields[3]),
			Unit:        strings.TrimSpace(fields[4]),
			Total:       strings.TrimSpace(fields[5]),
			Notes:       strings.TrimSpace(fields[6]),
		}
		records = append(records, record)
	}

	return records, nil
}
