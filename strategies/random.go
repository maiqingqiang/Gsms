package strategies

import (
	"math/rand"
	"time"
)

type RandomStrategy struct {
}

func (o *RandomStrategy) Apply(gateways []string) []string {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	r.Shuffle(
		len(gateways),
		func(i, j int) {
			gateways[i], gateways[j] = gateways[j], gateways[i]
		},
	)

	return gateways
}
