package utils

import (
	"time"

	"github.com/sirupsen/logrus"
)

type Timer struct {
	stop chan bool
}

// Start starts work. If the first attempt fails, then returns the error.
// It will trigger the work by the interval until recieving a stop signal.
func (t Timer) Start(f func() error, interval time.Duration, l *logrus.Entry) {
	ticker := time.Tick(interval)
	go func() {
		for {
			//TODO is it right
			select {
			case <-ticker:
				if err := f(); err != nil {
					l.WithError(err).Error()
				}
			case <-t.stop:
				break
			}
		}
	}()
}

func (t *Timer) Stop() {
	t.stop <- true
}
