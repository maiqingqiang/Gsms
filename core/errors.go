package core

import (
	"errors"
	"fmt"
)

var (
	ErrInvalidIDDCode   = errors.New("invalid IDDCode")
	ErrGatewayNotFound  = errors.New("gateway not found")
	ErrMessageTypeError = errors.New("message type error")
)

type ErrGatewayFailed struct {
	Results []*Result
}

func NewErrGatewayFailed(results []*Result) *ErrGatewayFailed {
	return &ErrGatewayFailed{Results: results}
}

func (e *ErrGatewayFailed) Error() string {
	return fmt.Sprintf("all gateways failed to send message: %v", e.Results)
}
