package vatcheck

import (
	"time"

	"github.com/sony/gobreaker"
)

var (
	cb *gobreaker.CircuitBreaker
)

func init() {
	initCircuitBreaker()
}

func initCircuitBreaker() {
	cb = gobreaker.NewCircuitBreaker(gobreaker.Settings{
		Name:    "VIES",
		Timeout: time.Duration(20) * time.Second,
		ReadyToTrip: func(counts gobreaker.Counts) bool {
			failureRatio := float64(counts.TotalFailures) / float64(counts.Requests)
			return counts.Requests >= 3 && failureRatio >= 0.6
		},
	})
}
