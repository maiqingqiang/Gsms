package gsms

import (
	"github.com/maiqingqiang/gsms/core"
	"time"
)

type Option func(*Gsms)

// WithTimeout set the timeout.
func WithTimeout(timeout time.Duration) func(*Gsms) {
	return func(gsms *Gsms) {
		gsms.Timeout = timeout
	}
}

// WithGateways set the gateways.
func WithGateways(gateways []string) func(*Gsms) {
	return func(gsms *Gsms) {
		gsms.DefaultGateways = gateways
	}
}

// WithStrategy set the strategy.
func WithStrategy(strategy core.StrategyInterface) func(*Gsms) {
	return func(gsms *Gsms) {
		gsms.Strategy = strategy
	}
}
