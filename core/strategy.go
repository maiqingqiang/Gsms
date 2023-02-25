package core

type StrategyInterface interface {
	// Apply the strategy and return result.
	Apply(gateways []string) []string
}
