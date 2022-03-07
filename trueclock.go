package trueclock

import (
	"github.com/facebook/time/ntp/chrony"
	"net"
	"time"
	"sync"
)

// 1 ppm, same as chrony
const MaxClockError = 1.0

type TrueClock struct {
	chronyConn net.Conn
	chronyClient *chrony.Client
	tracking *chrony.Tracking
	mu sync.RWMutex
}

func New() (*TrueClock, error) {
	chronyConn, err := net.Dial("udp", "[::1]:323")
	if err != nil {
		return nil, err
	}


	chronyClient := chrony.Client {
		Connection: chronyConn,
		Sequence: 1,
	}

	tracking, err := pollTracking(&chronyClient)
	if err != nil {
		return nil, err
	}

	clock := TrueClock {
		chronyConn: chronyConn,
		chronyClient: &chronyClient,
		tracking: tracking,
	}

	clock.startChronyPoller()

	return &clock, nil
}

func (c *TrueClock) Now() Bounds {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.updateRootDispersion()
	return boundsFromTracking(c.tracking)
}

func (c *TrueClock) updateRootDispersion() {
	now := time.Now()
	dispersion := dispersionAt(*c.tracking, now, MaxClockError)
	c.tracking.RootDispersion = dispersion
}

func dispersionAt(tracking chrony.Tracking, systemTime time.Time, maxClockError float64) float64 {
	elapsed := systemTime.Sub(tracking.RefTime).Seconds()
	errorRate := (maxClockError + tracking.SkewPPM + tracking.ResidFreqPPM) * 1e-6
	return tracking.RootDispersion + elapsed * errorRate
}

func (c *TrueClock) startChronyPoller() {
	ticker := time.NewTicker(time.Second)

	go func() {
		for {
			select {
			case <- ticker.C:
				tracking, _ := pollTracking(c.chronyClient)
				c.mu.Lock()
				c.tracking = tracking
				c.mu.Unlock()
			}
		}
	}()
}

func pollTracking(chronyClient *chrony.Client) (*chrony.Tracking, error) {
	req := chrony.NewTrackingPacket()
	resp, err := chronyClient.Communicate(req)
	if err != nil {
		return nil, err
	}

	replyTracking := resp.(*chrony.ReplyTracking)
	return &replyTracking.Tracking, nil
}

