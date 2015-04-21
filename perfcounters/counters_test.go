package perfcounters

import (
	"testing"
)

func TestNumberOfItems32(t *testing.T) {

	counter := NewNumberOfItems32()

	counter.Increment()

	if counter.String() != "1" {
		t.Error("Expected count of 1.")
	}

	counter.Add(9)

	if counter.String() != "10" {
		t.Error("Expected count of 10.")
	}
}

func TestNumberOfItems64(t *testing.T) {

	counter := NewNumberOfItems64()

	counter.Increment()

	if counter.String() != "1" {
		t.Error("Expected count of 1.")
	}

	counter.Add(9)

	if counter.String() != "10" {
		t.Error("Expected count of 10.")
	}
}

func TestAverageCount32(t *testing.T) {
}

func TestAverageTimer32(t *testing.T) {
}

func TestRateOfCountsPerSecond32(t *testing.T) {
}

func TestCountPerItemInterval32(t *testing.T) {
}
