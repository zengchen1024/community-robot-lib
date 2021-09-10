package utils

import "time"

type Timer struct {
	stop chan bool
}

// Start starts work. If the first attempt fails, then returns the error.
// It will trigger the work by the interval until recieving a stop signal.
func (t Timer) Start(f func(), interval time.Duration) {
	ticker := time.Tick(interval)
	go func() {
		for {
			//TODO is it right
			select {
			case <-ticker:
				f()
			case <-t.stop:
				break
			}
		}
	}()
}

func (t *Timer) Stop() {
	t.stop <- true
}
