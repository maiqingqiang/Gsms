package core

import "time"

type GatewayInterface interface {
	// Name Get gateway name
	Name() string
	// Send a short message
	Send(to PhoneNumberInterface, message MessageInterface, request RequestInterface) (string, error)
}

type GatewayBase struct {
	Timeout time.Duration
}
