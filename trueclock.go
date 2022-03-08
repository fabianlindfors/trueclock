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
	tracking := *c.tracking
	defer c.mu.Unlock()

	tracking.RootDispersion = calculateDispersion(tracking)
	return boundsFromTracking(tracking)
}

func calculateDispersion(tracking chrony.Tracking) float64 {
	systemTime := time.Now()
	elapsed := systemTime.Sub(tracking.RefTime).Seconds()
	errorRate := (MaxClockError + tracking.SkewPPM + tracking.ResidFreqPPM) * 1e-6

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

