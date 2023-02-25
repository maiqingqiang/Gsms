package gsms

import (
	"github.com/maiqingqiang/gsms/core"
	"time"
)

type Gsms struct {
	DefaultGateways []string
	Timeout         time.Duration
	Strategy        core.StrategyInterface
	Gateways        map[string]core.GatewayInterface
}

func WithTimeout(timeout time.Duration) func(*Gsms) {
	return func(gsms *Gsms) {
		gsms.Timeout = timeout
	}
}

func WithDefaultGateway(gateways []string) func(*Gsms) {
	return func(gsms *Gsms) {
		gsms.DefaultGateways = gateways
	}
}

func WithStrategy(strategy core.StrategyInterface) func(*Gsms) {
	return func(gsms *Gsms) {
		gsms.Strategy = strategy
	}
}

func New(gateways []core.GatewayInterface, options ...func(*Gsms)) *Gsms {
	gatewaysMap := make(map[string]core.GatewayInterface, len(gateways))

	for _, gateway := range gateways {
		gatewaysMap[gateway.Name()] = gateway
	}

	gsms := &Gsms{
		Gateways: gatewaysMap,
	}

	for _, option := range options {
		option(gsms)
	}

	return gsms
}

func (g *Gsms) Send(to int64, message core.MessageInterface, gateways ...string) ([]*core.Result, error) {

	if len(gateways) == 0 {
		var err error
		gateways, err = message.Gateways()
		if err != nil {
			return nil, err
		}
	}

	if len(gateways) == 0 {
		gateways = g.DefaultGateways
	}

	gateways = g.formatGateways(message, gateways)

	phoneNumber := core.NewNewPhoneNumberWithoutIDDCode(to)

	var results []*core.Result
	isSuccessful := false

	for _, gateway := range gateways {

		result := &core.Result{
			Gateway: gateway,
			Status:  core.StatusSuccess,
		}

		var gw core.GatewayInterface

		gw, result.Error = g.Gateway(gateway)

		if result.Error != nil {
			result.Status = core.StatusFailure
			results = append(results, result)
			continue
		}

		result.Template, result.Error = message.GetTemplate(gw)
		if result.Error != nil {
			result.Status = core.StatusFailure
			results = append(results, result)
			continue
		}

		result.Result, result.Error = gw.Send(phoneNumber, message)

		if result.Error != nil {
			result.Status = core.StatusFailure
		}

		results = append(results, result)

		if result.Status == core.StatusSuccess {
			isSuccessful = true
			break
		}
	}

	if !isSuccessful {
		return nil, core.NewErrGatewayFailed(results)
	}

	return results, nil
}

// Gateway Get gateway by name
func (g *Gsms) Gateway(name string) (core.GatewayInterface, error) {
	if gateway, ok := g.Gateways[name]; ok {
		return gateway, nil
	}

	return nil, core.ErrGatewayNotFound
}

// formatGateways format gateways
func (g *Gsms) formatGateways(message core.MessageInterface, gateways []string) []string {

	if strategy, err := message.Strategy(); err == nil && strategy != nil {
		return strategy.Apply(gateways)
	}

	if g.Strategy != nil {
		return g.Strategy.Apply(gateways)
	}

	return gateways
}
