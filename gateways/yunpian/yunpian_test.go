package yunpian

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/maiqingqiang/gsms/core"
	"github.com/stretchr/testify/assert"
	netUrl "net/url"
	"strings"
	"testing"
)

const NotSupportCountry = `{"http_status_code":400,"code":20,"msg":"暂不支持的国家地区","detail":"请确认号码归属地"}`
const Success = `{"code":0,"msg":"发送成功","count":1,"fee":0.05,"unit":"RMB","mobile":"18888888881","sid":74712264988}`

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
			assert.Equalf(t, tt.want, buildEndpoint(tt.args.product, tt.args.resource, tt.args.method), "buildEndpoint(%v, %v, %v)", tt.args.product, tt.args.resource, tt.args.method)
		})
	}
}

func TestGateway_Send(t *testing.T) {
	type fields struct {
		ApiKey    string
		Signature string
	}
	type args struct {
		to      int
		message core.MessageInterface
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    string
		wantErr string
	}{
		{
			name: "single send 1",
			fields: fields{
				ApiKey:    "ApiKey",
				Signature: "【云片】",
			},
			args: args{
				to: 18888888881,
				message: &core.Message{
					Content: "祝您万事如意，财源广进。",
				},
			},
			want:    "",
			wantErr: "暂不支持的国家地区",
		},
		{
			name: "single send 2",
			fields: fields{
				ApiKey:    "ApiKey",
				Signature: "【云片】",
			},
			args: args{
				to: 18888888882,
				message: &core.Message{
					Content: "祝您万事如意，财源广进。",
				},
			},
			want:    Success,
			wantErr: "",
		},
		{
			name: "tpl single send 1",
			fields: fields{
				ApiKey:    "ApiKey",
				Signature: "【云片】",
			},
			args: args{
				to: 18888888881,
				message: &core.Message{
					Template: "15320323",
					Data: map[string]string{
						"code": "6379",
					},
				},
			},
			want:    "",
			wantErr: "暂不支持的国家地区",
		}, {
			name: "tpl single send 2",
			fields: fields{
				ApiKey:    "ApiKey",
				Signature: "【云片】",
			},
			args: args{
				to: 18888888882,
				message: &core.Message{
					Template: "15320323",
					Data: map[string]string{
						"code": "6379",
					},
				},
			},
			want:    Success,
			wantErr: "",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gateway := &Gateway{
				ApiKey:    tt.fields.ApiKey,
				Signature: tt.fields.Signature,
			}

			phoneNumber := core.NewPhoneNumberWithoutIDDCode(tt.args.to)

			got, err := gateway.Send(phoneNumber, tt.args.message, &RequestTest{})
			if (err != nil) != (tt.wantErr != "") {
				assert.ErrorContainsf(t, err, tt.wantErr, "Send(%d, %v)", tt.args.to, tt.args.message)
			}

			assert.Equalf(t, tt.want, got, "Send(%v, %v)", tt.args.to, tt.args.message)
		})
	}
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
			assert.Equalf(t, tt.want, buildTplVal(tt.args.data), "buildTplVal(%v)", tt.args.data)
		})
	}
}

var _ core.RequestInterface = (*RequestTest)(nil)

type RequestTest struct {
}

func (r RequestTest) Request(method, url string, data interface{}, options ...core.Option) ([]byte, error) {
	return nil, nil
}

func (r RequestTest) Get(url string, data interface{}) ([]byte, error) {
	return nil, nil
}

func (r RequestTest) Post(url string, data interface{}) ([]byte, error) {
	if strings.Contains(url, "single_send") && data.(netUrl.Values).Get("mobile") == "18888888881" {
		return []byte(NotSupportCountry), nil
	}

	return []byte(Success), nil
}

func (r RequestTest) GetWithUnmarshal(url string, data interface{}, v interface{}) (string, error) {
	return "", nil
}

func (r RequestTest) PostWithUnmarshal(url string, data interface{}, v interface{}) (string, error) {
	body, err := r.Post(url, data)
	err = json.Unmarshal(body, v)
	if err != nil {
		return "", errors.New(fmt.Sprintf("json unmarshal error: %s body: %s", err.Error(), string(body)))
	}

	return string(body), nil
}
