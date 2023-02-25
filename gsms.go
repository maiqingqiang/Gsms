package gsms

import (
	"github.com/maiqingqiang/gsms/core"
	"github.com/maiqingqiang/gsms/strategies"
	"strconv"
	"time"
)

type Gsms struct {
	DefaultGateways []string
	Timeout         time.Duration
	Strategy        core.StrategyInterface
	Gateways        map[string]core.GatewayInterface
}

func New(gateways []core.GatewayInterface, options ...Option) *Gsms {
	gatewaysMap := make(map[string]core.GatewayInterface, len(gateways))

	for _, gateway := range gateways {
		gatewaysMap[gateway.Name()] = gateway
	}

	gsms := &Gsms{
		Gateways: gatewaysMap,
		Strategy: &strategies.OrderStrategy{},
	}

	for _, option := range options {
		option(gsms)
	}

	return gsms
}

// Send a message.
func (g *Gsms) Send(to interface{}, message core.MessageInterface, gateways ...string) ([]*core.Result, error) {

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

	var phoneNumber *core.PhoneNumber

	switch to.(type) {
	case *core.PhoneNumber:
		phoneNumber = to.(*core.PhoneNumber)
	case int:
		phoneNumber = core.NewPhoneNumberWithoutIDDCode(to.(int))
	case string:
		number, err := strconv.Atoi(to.(string))
		if err != nil {
			return nil, err
		}

		phoneNumber = core.NewPhoneNumberWithoutIDDCode(number)
	default:
		return nil, core.ErrInvalidPhoneNumber
	}

	var results []*core.Result
	isSuccessful := false

	request := &core.Request{
		Timeout: g.Timeout,
	}

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

		result.Result, result.Error = gw.Send(phoneNumber, message, request)

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
