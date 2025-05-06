package slide_test

import (
	"testing"
	"time"
)

func generateRFC3389FromString(t *testing.T, timestamp string) time.Time {
	t.Helper()

	expectedTime, err := time.Parse(time.RFC3339, timestamp)
	if err != nil {
		t.Fatalf("error during test setup - converting timestamp string to time.Time (%s) failed: %s", timestamp, err.Error())
	}

	return expectedTime
}
