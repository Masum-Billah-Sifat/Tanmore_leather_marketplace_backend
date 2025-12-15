package timeutil

import "time"

// Current UTC time
func NowUTC() time.Time {
	return time.Now().UTC()
}

// Format in ISO-8601 for events, logs, audit trails
func FormatISO8601(t time.Time) string {
	return t.UTC().Format(time.RFC3339)
}
