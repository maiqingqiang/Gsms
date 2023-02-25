package yunpian

import (
	"github.com/jarcoal/httpmock"
	"github.com/maiqingqiang/gsms/core"
	"github.com/stretchr/testify/assert"
	"net/http"
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
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	httpmock.RegisterResponder(
		http.MethodPost,
		"https://sms.yunpian.com/v2/sms/single_send.json",
		func(request *http.Request) (*http.Response, error) {
			if request.FormValue("mobile") == "18888888881" {
				return httpmock.NewStringResponse(400, NotSupportCountry), nil
			}

			return httpmock.NewStringResponse(200, Success), nil
		},
	)

	httpmock.RegisterResponder(
		http.MethodPost,
		"https://sms.yunpian.com/v2/sms/tpl_single_send.json",
		func(request *http.Request) (*http.Response, error) {
			if request.FormValue("mobile") == "18888888881" {
				return httpmock.NewStringResponse(400, NotSupportCountry), nil
			}

			return httpmock.NewStringResponse(200, Success), nil
		},
	)

	type fields struct {
		ApiKey    string
		Signature string
	}
	type args struct {
		to      int64
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

			phoneNumber := core.NewNewPhoneNumberWithoutIDDCode(tt.args.to)

			got, err := gateway.Send(phoneNumber, tt.args.message)
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
