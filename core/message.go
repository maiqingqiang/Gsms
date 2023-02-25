package core

// TextMessage Text message type.
const TextMessage = "text"

// VoiceMessage Voice message type.
const VoiceMessage = "voice"

// MessageInterface Message interface.
type MessageInterface interface {
	// Gateways Supported gateways.
	Gateways() ([]string, error)
	// Strategy Message strategy.
	Strategy() (StrategyInterface, error)
	// GetContent Get message content.
	GetContent(gateway GatewayInterface) (string, error)
	// GetTemplate Get message template.
	GetTemplate(gateway GatewayInterface) (string, error)
	// GetData Get message data.
	GetData(gateway GatewayInterface) (map[string]string, error)
	// GetType Get message type.
	GetType(gateway GatewayInterface) (string, error)
}

// ContentFunc Content function.
type ContentFunc func(gatewayName string) string

// TemplateFunc Template function.
type TemplateFunc func(gatewayName string) string

// DataFunc Data function.
type DataFunc func(gatewayName string) map[string]string

// TypeFunc Type function.
type TypeFunc func(gatewayName string) string

var _ MessageInterface = (*Message)(nil)

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
func (m *Message) Strategy() (StrategyInterface, error) {
	return nil, nil
}

// GetContent Get message content.
func (m *Message) GetContent(gateway GatewayInterface) (string, error) {
	switch content := m.Content.(type) {
	case string:
		return content, nil
	case func(gateway GatewayInterface) string:
		return content(gateway), nil
	}
	return "", nil
}

// GetTemplate Get message template.
func (m *Message) GetTemplate(gateway GatewayInterface) (string, error) {
	switch template := m.Template.(type) {
	case string:
		return template, nil
	case func(gateway GatewayInterface) string:
		return template(gateway), nil
	}

	return "", nil
}

// GetData Get message data.
func (m *Message) GetData(gateway GatewayInterface) (map[string]string, error) {
	switch data := m.Data.(type) {
	case map[string]string:
		return data, nil
	case func(gateway GatewayInterface) map[string]string:
		return data(gateway), nil
	}

	return nil, nil
}

// GetType Get message type.
func (m *Message) GetType(gateway GatewayInterface) (string, error) {
	switch messageType := m.Type.(type) {
	case string:
		return messageType, nil
	case func(gateway GatewayInterface) string:
		return messageType(gateway), nil
	}

	return "", nil
}
