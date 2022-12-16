package utils

import "time"

// FormatTimeRFC3339 Format time according to RFC3339Nano
func FormatTimeRFC3339(t *time.Time) (s string) {
	if t == nil {
		return
	}

	if t.Nanosecond() == 0 {
		return t.Format("2022-12-15T15:04:05.000000000Z07:00")
	}

	return t.Format(time.RFC3339Nano)
}
