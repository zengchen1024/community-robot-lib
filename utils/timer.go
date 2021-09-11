package utils

import "time"

type Timer interface {
	Start(f func(), interval time.Duration)
	Stop()
}

func NewTimer() Timer {
	return timer{
		stop:    make(chan struct{}),
		stopped: make(chan struct{}),
	}
}

type timer struct {
	stop    chan struct{}
	stopped chan struct{}
}

// Start triggers the work by the interval until recieving a stop signal.
func (t timer) Start(f func(), interval time.Duration) {
	ticker := time.Tick(interval)

	go func() {
		for {
			select {
			case <-ticker:
				f()
			case <-t.stop:
				close(t.stopped)
				return
			}
		}
	}()
}

func (t timer) Stop() {
	close(t.stop)

	ticker := time.Tick(1 * time.Millisecond)

	for range ticker {
		select {
		case <-t.stopped:
			return
		}
	}
}
