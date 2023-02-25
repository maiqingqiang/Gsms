package core

import "time"

type GatewayInterface interface {
	// Name Get gateway name
	Name() string
	// Send a short message
	Send(to PhoneNumberInterface, message MessageInterface) (string, error)
}

type base struct {
	timeout time.Duration
}
