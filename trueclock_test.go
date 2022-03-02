package trueclock

import (
	"testing"
	"time"
	"fmt"
	"math"
)

func TestClock(t *testing.T) {
	clock, err := New()
	if err != nil {
		t.Error(err)
	}

	firstBounds := clock.Now()
	firstDiff := firstBounds.Latest.Sub(firstBounds.Earliest)
	fmt.Printf("First measurement uncertainty: %d us\n", firstDiff.Microseconds())

	if !firstBounds.Latest.After(firstBounds.Earliest) {
		t.Errorf("expected latest timestamp to be after earliest")
	}

	fmt.Println("Sleeping for 5 seconds...")
	time.Sleep(time.Second * 5)

	secondBounds := clock.Now()
	secondDiff := secondBounds.Latest.Sub(secondBounds.Earliest)
	fmt.Printf("Second measurement uncertainty: %d us\n", secondDiff.Microseconds())

	if !secondBounds.Latest.After(secondBounds.Earliest) {
		t.Errorf("expected latest timestamp to be after earliest")
	}

	diffEarliest := secondBounds.Earliest.Sub(firstBounds.Earliest)
	diffLatest := secondBounds.Latest.Sub(firstBounds.Latest)

	if math.Abs(float64(diffEarliest.Milliseconds() - 5000)) > 10 {
		t.Errorf("expected earliest timestamps to not have diverged by more than 10 ms")
	}

	if math.Abs(float64(diffLatest.Milliseconds() - 5000)) > 10 {
		t.Errorf("expected latest timestamps to not have diverged by more than 10 ms")
	}
}
