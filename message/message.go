package message

import (
	"github.com/maiqingqiang/gsms"
)

// TextMessage Text message type.
const TextMessage = "text"

// VoiceMessage Voice message type.
const VoiceMessage = "voice"

// ContentFunc Content function.
type ContentFunc func(gatewayName string) string

// TemplateFunc Template function.
type TemplateFunc func(gatewayName string) string

// DataFunc Data function.
type DataFunc func(gatewayName string) map[string]string

// TypeFunc Type function.
type TypeFunc func(gatewayName string) string

//var _ gsms.Message = (*Message)(nil)

type Message struct {
	Content  interface{}
	Template interface{}
	Data     interface{}
	Type     interface{}
}

// Gateways Supported gateways.
func (m *Message) Gateways() ([]string, error) {
	return nil, nil
}

// Strategy Message strategy.
func (m *Message) Strategy() (gsms.Strategy, error) {
	return nil, nil
}

// GetContent Get message content.
func (m *Message) GetContent(gateway gsms.Gateway) (string, error) {
	switch content := m.Content.(type) {
	case string:
		return content, nil
	case func(gateway gsms.Gateway) string:
		return content(gateway), nil
	}
	return "", nil
}

// GetTemplate Get message template.
func (m *Message) GetTemplate(gateway gsms.Gateway) (string, error) {
	switch template := m.Template.(type) {
	case string:
		return template, nil
	case func(gateway gsms.Gateway) string:
		return template(gateway), nil
	}

	return "", nil
}

// GetData Get message data.
func (m *Message) GetData(gateway gsms.Gateway) (map[string]string, error) {
	switch data := m.Data.(type) {
	case map[string]string:
		return data, nil
	case func(gateway gsms.Gateway) map[string]string:
		return data(gateway), nil
	}

	return nil, nil
}

// GetType Get message type.
func (m *Message) GetType(gateway gsms.Gateway) (string, error) {
	switch messageType := m.Type.(type) {
	case string:
		return messageType, nil
	case func(gateway gsms.Gateway) string:
		return messageType(gateway), nil
	}

	return "", nil
}
