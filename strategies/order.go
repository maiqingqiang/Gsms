package strategies

import (
	"github.com/maiqingqiang/gsms/core"
	"sort"
)

var _ core.StrategyInterface = (*OrderStrategy)(nil)

type OrderStrategy struct {
}

func (o *OrderStrategy) Apply(gateways []string) []string {
	sort.Strings(gateways)
	return gateways
}
