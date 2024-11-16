package utils

import (
	"errors"
	"log"
	"time"
)

func Retry(maxAttempts int, initialDelay time.Duration, fn func() error) error {
	var err error
	delay := initialDelay

	for attempt := 1; attempt <= maxAttempts; attempt++ {
		err = fn()
		if err == nil {
			return nil
		}
		log.Printf("attemp %v/%v failed: %v", attempt, maxAttempts, err)

		// exponential backoff
		time.Sleep(delay)
		delay *= 2
	}
	return errors.New("max retry attempts exceeded")
}
