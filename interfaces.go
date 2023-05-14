//go:generate mockgen -source=interfaces.go -destination=gsms_mock.go -package=gsms

package gsms

type Gateway interface {
	// Name Get gateway name
	Name() string
	// Send a short message
	Send(to *PhoneNumber, message Message, config *Config) error
}

// Message interface.
type Message interface {
	// Gateways Supported gateways.
	Gateways() ([]string, error)
	// Strategy Message strategy.
	Strategy() (Strategy, error)
	// GetContent Get message content.
	GetContent(gateway Gateway) (string, error)
	// GetTemplate Get message template.
	GetTemplate(gateway Gateway) (string, error)
	// GetData Get message data.
	GetData(gateway Gateway) (map[string]string, error)
	// GetType Get message type.
	GetType(gateway Gateway) (string, error)
}

type Strategy interface {
	// Apply the strategy and return result.
	Apply(gateways []string) []string
}

// LogLevel log level
type LogLevel int

type Logger interface {
	LogMode(LogLevel) Logger
	Info(...interface{})
	Infof(string, ...interface{})
	Warn(...interface{})
	Warnf(string, ...interface{})
	Error(...interface{})
	Errorf(string, ...interface{})
}
