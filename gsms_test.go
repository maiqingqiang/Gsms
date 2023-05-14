package gsms

import (
	"fmt"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"testing"
)

type Test1Gateway struct {
}

func (t Test1Gateway) Name() string {
	return "Test1"
}

func (t Test1Gateway) Send(to *PhoneNumber, message Message, config *Config) error {
	return nil
}

var _ Gateway = (*Test1Gateway)(nil)

type Test2Gateway struct {
}

func (t Test2Gateway) Name() string {
	return "Test2"
}

func (t Test2Gateway) Send(to *PhoneNumber, message Message, config *Config) error {
	return nil
}

var _ Gateway = (*Test2Gateway)(nil)

func TestGsms_Gateway(t *testing.T) {
	type fields struct {
		config          *Config
		DefaultGateways []string
		Strategy        Strategy
		Gateways        map[string]Gateway
	}
	type args struct {
		name string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    Gateway
		wantErr assert.ErrorAssertionFunc
	}{
		{
			name:   "gsms.Gateway Return Test1",
			fields: fields{},
			args: args{
				name: "Test1",
			},
			want:    &Test1Gateway{},
			wantErr: assert.NoError,
		},
		{
			name:   "gsms.Gateway Return Test2",
			fields: fields{},
			args: args{
				name: "Test2",
			},
			want:    &Test2Gateway{},
			wantErr: assert.NoError,
		},
		{
			name:   "gsms.Gateway Return Test3",
			fields: fields{},
			args: args{
				name: "Test3",
			},
			want:    nil,
			wantErr: assert.Error,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			g := New([]Gateway{
				&Test1Gateway{},
				&Test2Gateway{},
			})

			got, err := g.Gateway(tt.args.name)
			if !tt.wantErr(t, err, fmt.Sprintf("Gateway(%v)", tt.args.name)) {
				return
			}
			assert.Equalf(t, tt.want, got, "Gateway(%v)", tt.args.name)
		})
	}
}

func TestGsms_Send(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockMessage := NewMockMessage(ctrl)
	mockGateway := NewMockGateway(ctrl)

	mockMessage.EXPECT().Gateways().Return(nil, nil)
	mockMessage.EXPECT().Strategy().Return(nil, nil)
	mockMessage.EXPECT().GetTemplate(mockGateway).Return("SMS_00000001", nil)

	mockGateway.EXPECT().Name().Return("mockGateway")
	mockGateway.EXPECT().Send(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil)

	g := New([]Gateway{
		mockGateway,
	}, WithGateways([]string{
		"mockGateway",
	}))

	result, err := g.Debug().Send(18888888888, mockMessage)
	if !assert.NoError(t, err) {
		return
	}

	assert.Equal(t, result, []*Result{
		{
			Gateway:  "mockGateway",
			Status:   "success",
			Template: "SMS_00000001",
			Error:    nil,
		},
	})
}

func TestGsms_Send_Failed(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockMessage := NewMockMessage(ctrl)
	mockGateway := NewMockGateway(ctrl)

	mockMessage.EXPECT().Gateways().Return(nil, nil)
	mockMessage.EXPECT().Strategy().Return(nil, nil)
	mockMessage.EXPECT().GetTemplate(gomock.Any()).Return("SMS_00000001", nil)

	mockGateway.EXPECT().Name().Return("mockGateway")
	mockGateway.EXPECT().Send(gomock.Any(), gomock.Any(), gomock.Any()).Return(fmt.Errorf("send failed"))

	g := New([]Gateway{
		mockGateway,
	}, WithGateways([]string{
		"mockGateway",
	}))
	_, err := g.Debug().Send(18888888888, mockMessage)
	assert.Error(t, err)
}
