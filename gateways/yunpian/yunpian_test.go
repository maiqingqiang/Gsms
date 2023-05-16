package yunpian

import (
	"github.com/jarcoal/httpmock"
	"github.com/maiqingqiang/gsms"
	"github.com/maiqingqiang/gsms/message"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func Test_buildEndpoint(t *testing.T) {
	type args struct {
		product  string
		resource string
		method   string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "sms",
			args: args{
				product:  "sms",
				resource: "sms",
				method:   "single_send",
			},
			want: "https://sms.yunpian.com/v2/sms/single_send.json",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			g := &Gateway{}
			assert.Equalf(t, tt.want, g.buildEndpoint(tt.args.product, tt.args.resource, tt.args.method), "buildEndpoint(%v, %v, %v)", tt.args.product, tt.args.resource, tt.args.method)
		})
	}
}

func TestGateway_Send(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	httpmock.RegisterResponder("POST", `https://sms.yunpian.com/v2/sms/single_send.json`,
		httpmock.NewStringResponder(200, `{"code":0,"msg":"发送成功","count":1,"fee":0.05,"unit":"RMB","mobile":"18888888888","sid":74712264988}`))

	httpmock.RegisterResponder("POST", `https://sms.yunpian.com/v2/sms/tpl_single_send.json`,
		httpmock.NewStringResponder(200, `{"code":0,"msg":"发送成功","count":1,"fee":0.05,"unit":"RMB","mobile":"18888888888","sid":74712264988}`))

	g := &Gateway{
		ApiKey:    "ApiKey",
		Signature: "Signature",
	}

	config := &gsms.Config{
		Timeout: 5 * time.Second,
		Logger:  gsms.NewLogger().LogMode(gsms.Info),
	}

	phoneNumber := gsms.NewPhoneNumberWithoutIDDCode(188888888888)

	err := g.Send(
		phoneNumber,
		&message.Message{
			Template: "SMS_00000001",
			Data: map[string]string{
				"code": "9527",
			},
		},
		config,
	)

	assert.NoError(t, err)

	err = g.Send(
		phoneNumber,
		&message.Message{
			Template: func(gateway gsms.Gateway) string {
				if gateway.Name() == NAME {
					return "SMS_271311117"
				}
				return "5532011"
			},
			Data: func(gateway gsms.Gateway) map[string]string {
				if gateway.Name() == NAME {
					return map[string]string{
						"code": "1111",
					}
				}
				return map[string]string{
					"code": "6379",
				}
			},
		},
		config,
	)

	assert.NoError(t, err)

	err = g.Send(
		phoneNumber,
		&message.Message{
			Content: "【Gsms】您的验证码是521410",
		},
		config,
	)

	assert.NoError(t, err)

	err = g.Send(
		phoneNumber,
		&message.Message{
			Content: func(gateway gsms.Gateway) string {
				if gateway.Name() == NAME {
					return "【Gsms】您的验证码是521410"
				}
				return "【Gsms】活动验证码是111"
			},
		},
		config,
	)

	assert.NoError(t, err)
}

func TestGateway_Send_Failed(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	httpmock.RegisterResponder("POST", `https://sms.yunpian.com/v2/sms/single_send.json`,
		httpmock.NewStringResponder(200, `{"http_status_code":400,"code":20,"msg":"暂不支持的国家地区","detail":"请确认号码归属地"}`))

	httpmock.RegisterResponder("POST", `https://sms.yunpian.com/v2/sms/tpl_single_send.json`,
		httpmock.NewStringResponder(200, `{"http_status_code":400,"code":20,"msg":"暂不支持的国家地区","detail":"请确认号码归属地"}`))

	g := &Gateway{
		ApiKey:    "ApiKey",
		Signature: "Signature",
	}

	config := &gsms.Config{
		Timeout: 5 * time.Second,
		Logger:  gsms.NewLogger().LogMode(gsms.Info),
	}

	phoneNumber := gsms.NewPhoneNumberWithoutIDDCode(188888888888)

	err := g.Send(
		phoneNumber,
		&message.Message{
			Template: "SMS_00000001",
			Data: map[string]string{
				"code": "9527",
			},
		},
		config,
	)

	assert.Error(t, err)

	err = g.Send(
		phoneNumber,
		&message.Message{
			Template: func(gateway gsms.Gateway) string {
				if gateway.Name() == NAME {
					return "SMS_271311117"
				}
				return "5532011"
			},
			Data: func(gateway gsms.Gateway) map[string]string {
				if gateway.Name() == NAME {
					return map[string]string{
						"code": "1111",
					}
				}
				return map[string]string{
					"code": "6379",
				}
			},
		},
		config,
	)

	assert.Error(t, err)

	err = g.Send(
		phoneNumber,
		&message.Message{
			Content: "【Gsms】您的验证码是521410",
		},
		config,
	)

	assert.Error(t, err)

	err = g.Send(
		phoneNumber,
		&message.Message{
			Content: func(gateway gsms.Gateway) string {
				if gateway.Name() == NAME {
					return "【Gsms】您的验证码是521410"
				}
				return "【Gsms】活动验证码是111"
			},
		},
		config,
	)

	assert.Error(t, err)
}

func Test_buildTplVal(t *testing.T) {
	type args struct {
		data map[string]string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "build template value",
			args: args{
				data: map[string]string{
					"code": "6379",
				},
			},
			want: "#code#=6379",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			g := &Gateway{}
			assert.Equalf(t, tt.want, g.buildTplVal(tt.args.data), "buildTplVal(%v)", tt.args.data)
		})
	}
}
