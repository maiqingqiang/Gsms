package dove

import (
	"github.com/maiqingqiang/gsms"
	"time"
)

type Option func(dove *Dove)

func WithTimeout(timeout time.Duration) Option {
	return func(dove *Dove) {
		dove.client.Timeout = timeout
	}
}

func WithLogger(logger gsms.Logger) Option {
	return func(dove *Dove) {
		dove.logger = logger
	}
}

func WithUnmarshal(unmarshal Unmarshal) Option {
	return func(dove *Dove) {
		dove.unmarshal = unmarshal
	}
}
