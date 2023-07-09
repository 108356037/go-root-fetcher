package queue

import (
	"context"
	"time"
)

func IntervalTrigger(intervalTime time.Duration, triggerChan chan string, quitChan chan struct{}, done context.CancelFunc) {
	for {
		select {
		case <-quitChan:
			done()

		case t := <-triggerChan:
			time.Sleep(intervalTime)
			triggerChan <- t
		}

	}
}
