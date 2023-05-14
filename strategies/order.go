package strategies

import (
	"sort"
)

type OrderStrategy struct {
}

func (o *OrderStrategy) Apply(gateways []string) []string {
	sort.Strings(gateways)
	return gateways
}
