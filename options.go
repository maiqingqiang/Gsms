package gsms

import (
	"time"
)

type Option func(*Gsms)

// WithTimeout set the timeout.
func WithTimeout(timeout time.Duration) func(*Gsms) {
	return func(gsms *Gsms) {
		gsms.config.Timeout = timeout
	}
}

// WithGateways set the gateways.
func WithGateways(gateways []string) func(*Gsms) {
	return func(gsms *Gsms) {
		gsms.defaultGateways = gateways
	}
}

// WithStrategy set the strategy.
func WithStrategy(strategy Strategy) func(*Gsms) {
	return func(gsms *Gsms) {
		gsms.strategy = strategy
	}
}
