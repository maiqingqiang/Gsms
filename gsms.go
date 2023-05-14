package gsms

import (
	"fmt"
	"github.com/maiqingqiang/gsms/strategies"
	"strconv"
	"time"
)

type Gsms struct {
	config          *Config
	defaultGateways []string
	strategy        Strategy
	gateways        map[string]Gateway
}

// Config gsms config.
type Config struct {
	Timeout time.Duration
	Logger  Logger
}

// StatusSuccess send message success.
const StatusSuccess = "success"

// StatusFailure send message failure.
const StatusFailure = "failure"

// Result Gateway send message result.
type Result struct {
	Gateway  string
	Status   string
	Template string
	Error    error
}

func (r *Result) String() string {
	return fmt.Sprintf("gateway: %s, status: %s, template: %s, error: %v", r.Gateway, r.Status, r.Template, r.Error)
}

// New a gsms instance.
func New(gateways []Gateway, options ...Option) *Gsms {
	gatewaysMap := make(map[string]Gateway, len(gateways))

	for _, gateway := range gateways {
		gatewaysMap[gateway.Name()] = gateway
	}

	gsms := &Gsms{
		config: &Config{
			Timeout: 5 * time.Second,
			Logger:  NewLogger(),
		},
		gateways: gatewaysMap,
		strategy: &strategies.OrderStrategy{},
	}

	for _, option := range options {
		option(gsms)
	}

	return gsms
}

// Send a message.
func (g *Gsms) Send(to interface{}, message Message, gateways ...string) ([]*Result, error) {

	if len(gateways) == 0 {
		var err error
		gateways, err = message.Gateways()
		if err != nil {
			return nil, err
		}
	}

	if len(gateways) == 0 {
		gateways = g.defaultGateways
	}

	gateways = g.formatGateways(message, gateways)

	var phoneNumber *PhoneNumber

	switch to.(type) {
	case *PhoneNumber:
		phoneNumber = to.(*PhoneNumber)
	case int:
		phoneNumber = NewPhoneNumberWithoutIDDCode(to.(int))
	case string:
		number, err := strconv.Atoi(to.(string))
		if err != nil {
			return nil, err
		}

		phoneNumber = NewPhoneNumberWithoutIDDCode(number)
	default:
		return nil, ErrInvalidPhoneNumber
	}

	var results []*Result
	isSuccessful := false

	for _, gateway := range gateways {

		result := &Result{
			Gateway: gateway,
			Status:  StatusSuccess,
		}

		var gw Gateway

		gw, result.Error = g.Gateway(gateway)

		if result.Error != nil {
			result.Status = StatusFailure
			results = append(results, result)
			continue
		}

		result.Template, result.Error = message.GetTemplate(gw)
		if result.Error != nil {
			result.Status = StatusFailure
			results = append(results, result)
			continue
		}

		g.config.Logger.Infof("[%s] start send [template: %s] message", gateway, result.Template)

		result.Error = gw.Send(phoneNumber, message, g.config)

		if result.Error != nil {
			result.Status = StatusFailure
			g.config.Logger.Warnf("[%s] send [template: %s] message failed: %+v", gateway, result.Template, result.Error)
		} else {
			g.config.Logger.Infof("[%s] send [template: %s] message success", gateway, result.Template)
		}

		g.config.Logger.Infof("[%s] end send [template: %s] message\n", gateway, result.Template)

		results = append(results, result)

		if result.Status == StatusSuccess {
			isSuccessful = true
			break
		}
	}

	if !isSuccessful {
		return nil, NewErrGatewayFailed(results)
	}

	return results, nil
}

// Gateway Get gateway by name
func (g *Gsms) Gateway(name string) (Gateway, error) {
	if gateway, ok := g.gateways[name]; ok {
		return gateway, nil
	}

	return nil, ErrGatewayNotFound
}

// formatGateways format gateways
func (g *Gsms) formatGateways(message Message, gateways []string) []string {

	if strategy, err := message.Strategy(); err == nil && strategy != nil {
		return strategy.Apply(gateways)
	}

	if g.strategy != nil {
		return g.strategy.Apply(gateways)
	}

	return gateways
}

// Debug set logger level to info
func (g *Gsms) Debug() *Gsms {
	newGsms := *g
	newGsms.config.Logger = newGsms.config.Logger.LogMode(Info)
	return &newGsms
}
