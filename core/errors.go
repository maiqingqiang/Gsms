package core

import (
	"errors"
	"fmt"
)

var (
	ErrInvalidIDDCode       = errors.New("invalid IDDCode")
	ErrGatewayNotFound      = errors.New("gateway not found")
	ErrMessageTypeError     = errors.New("message type error")
	ErrInvalidPhoneNumber   = errors.New("invalid phone number")
	ErrRequestDataTypeError = errors.New("request data type error")
)

type ErrGatewaysFailed struct {
	Results []*Result
}

func NewErrGatewayFailed(results []*Result) *ErrGatewaysFailed {
	return &ErrGatewaysFailed{Results: results}
}

func (e *ErrGatewaysFailed) Error() string {
	return fmt.Sprintf("all gateways failed to send message: %v", e.Results)
}

type ErrRequestFailed struct {
	StatusCode int
	Body       string
	Message    string
}

func (e *ErrRequestFailed) Error() string {
	return e.Message
}
