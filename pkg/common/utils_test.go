package common

import (
	"testing"
	"time"
)

func TestTimeCorrectedDeviationReturnsZeroForEmptyValues(t *testing.T) {
	if got := TimeCorrectedDeviation(nil); got != 0 {
		t.Fatalf("TimeCorrectedDeviation(nil) = %v, want 0", got)
	}
}

func TestTimeCorrectedDeviationReturnsZeroForSingleValue(t *testing.T) {
	values := []time.Duration{1500 * time.Millisecond}
	if got := TimeCorrectedDeviation(values); got != 0 {
		t.Fatalf("TimeCorrectedDeviation(single value) = %v, want 0", got)
	}
}
