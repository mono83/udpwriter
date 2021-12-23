package influxdb

import "strings"

var sanitizeReplacement = byte('_')

// Sanitize prepares value to be written to InfluxDB
func Sanitize(s string) []byte {
	if len(s) == 0 {
		return []byte{}
	}

	bts := []byte(strings.TrimSpace(s))
	for i, v := range bts {
		if !(v == 46 || (v >= 48 && v <= 57) || (v >= 65 && v <= 90) || (v >= 97 && v <= 122)) {
			bts[i] = sanitizeReplacement
		}
	}

	return bts
}

// SanitizeS prepares value to be written to InfluxDB
func SanitizeS(s string) string {
	return string(Sanitize(s))
}
