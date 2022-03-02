package trueclock

import (
	"github.com/facebook/time/ntp/chrony"
	"time"
)


func dispersionAt(tracking chrony.Tracking, systemTime time.Time, maxClockError float64) float64 {
	elapsed := systemTime.Sub(tracking.RefTime).Seconds()
	errorRate := (maxClockError + tracking.SkewPPM + tracking.ResidFreqPPM) * 1e-6
	return tracking.RootDispersion + elapsed * errorRate
}
