package strategies

import (
	"github.com/maiqingqiang/gsms/core"
	"math/rand"
	"time"
)

var _ core.StrategyInterface = (*RandomStrategy)(nil)

type RandomStrategy struct {
}

func (o *RandomStrategy) Apply(gateways []string) []string {
	rand.Seed(time.Now().UnixNano())
	rand.Shuffle(
		len(gateways),
		func(i, j int) {
			gateways[i], gateways[j] = gateways[j], gateways[i]
		},
	)

	return gateways
}
