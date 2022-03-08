package trueclock

import (
	"github.com/facebook/time/ntp/chrony"
	"time"
	"math"
)

type Bounds struct {
	Earliest time.Time
	Latest time.Time
}

func boundsFromTracking(tracking chrony.Tracking) Bounds {
	errorBound := roundFloat64ToNanos(math.Abs(tracking.CurrentCorrection) + tracking.RootDispersion + (tracking.RootDelay / 2.0))

	errorBoundNanos := int64(errorBound * 1000000000.0)

	now := time.Now()
	return Bounds {
		Earliest: now.Add(time.Duration(-errorBoundNanos)),
		Latest: now.Add(time.Duration(errorBoundNanos)),
	}
}

func roundFloat64ToNanos(value float64) float64 {
	return math.Round(value * 1000000000.0) / 1000000000.0
}
