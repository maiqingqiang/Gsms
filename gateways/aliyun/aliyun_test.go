package aliyun

import (
	"github.com/maiqingqiang/gsms/core"
	"github.com/stretchr/testify/assert"
	"net/url"
	netUrl "net/url"
	"strings"
	"testing"
)

const Success = `{"Message":"只能向已回复授权信息的手机号发送","RequestId":"F69545AD-66DC-53BE-B5BD-0E4D2E147AF7","Code":"OK"}`

func Test_generateSign(t *testing.T) {
	type args struct {
		httpMethod      string
		accessKeySecret string
		query           url.Values
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "test",
			args: args{
				httpMethod:      "GET",
				accessKeySecret: "access_key_secret",
				query: url.Values{
					"b": []string{"b"},
					"a": []string{"a"},
					"c": []string{"c"},
				},
			},
			want: "LQsZurDFe+jKVpTtyYcCDV+/y68=",
		},
		{
			// https://ram.aliyuncs.com/?UserName=test&SignatureVersion=1.0&Format=JSON&Timestamp=2015-08-18T03%3A15%3A45Z&AccessKeyId=testid&SignatureMethod=HMAC-SHA1&Version=2015-05-01&Action=CreateUser&SignatureNonce=6a6e0ca6-4557-11e5-86a2-b8e8563dc8d2
			// https://ram.aliyuncs.com/?UserName=test&SignatureVersion=1.0&Format=JSON&Timestamp=2015-08-18T03%3A15%3A45Z&AccessKeyId=testid&SignatureMethod=HMAC-SHA1&Version=2015-05-01&Signature=kRA2cnpJVacIhDMzXnoNZG9tDCI%3D&Action=CreateUser&SignatureNonce=6a6e0ca6-4557-11e5-86a2-b8e8563dc8d2
			name: "aliyun test",
			args: args{
				httpMethod:      "GET",
				accessKeySecret: "testsecret",
				query: url.Values{
					"UserName":         []string{"test"},
					"SignatureVersion": []string{"1.0"},
					"Format":           []string{"JSON"},
					"Timestamp":        []string{"2015-08-18T03:15:45Z"},
					"AccessKeyId":      []string{"testid"},
					"SignatureMethod":  []string{"HMAC-SHA1"},
					"Version":          []string{"2015-05-01"},
					"Action":           []string{"CreateUser"},
					"SignatureNonce":   []string{"6a6e0ca6-4557-11e5-86a2-b8e8563dc8d2"},
				},
			},
			want: "kRA2cnpJVacIhDMzXnoNZG9tDCI=",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, generateSign(tt.args.httpMethod, tt.args.accessKeySecret, tt.args.query))
		})
	}
}

func TestGateway_Send(t *testing.T) {
	type fields struct {
		AccessKeyId     string
		AccessKeySecret string
		SignName        string
	}
	type args struct {
		to      int
		message core.MessageInterface
		request core.ClientInterface
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    string
		wantErr string
	}{
		{
			name: "",
			fields: fields{
				AccessKeyId:     "AccessKeyId",
				AccessKeySecret: "AccessKeySecret",
				SignName:        "SignName",
			},
			args: args{
				to: 18888888881,
				message: &core.Message{
					Content: "祝您万事如意，财源广进。",
				},
				request: &ClientTest{},
			},
			want:    Success,
			wantErr: "",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			g := &Gateway{
				AccessKeyId:     tt.fields.AccessKeyId,
				AccessKeySecret: tt.fields.AccessKeySecret,
				SignName:        tt.fields.SignName,
			}

			phoneNumber := core.NewPhoneNumberWithoutIDDCode(tt.args.to)

			got, err := g.Send(phoneNumber, tt.args.message, tt.args.request)
			if (err != nil) != (tt.wantErr != "") {
				assert.ErrorContainsf(t, err, tt.wantErr, "Send(%d, %v)", tt.args.to, tt.args.message)
			}
			assert.Equalf(t, tt.want, got, "Send(%v, %v, %v)", tt.args.to, tt.args.message, tt.args.request)
		})
	}
}

var _ core.ClientInterface = (*ClientTest)(nil)

type ClientTest struct {
}

func (c ClientTest) GetWithUnmarshal(api string, data interface{}, v core.ResponseInterface) (string, error) {
	body := Success
	if strings.Contains(api, "single_send") && data.(netUrl.Values).Get("mobile") == "18888888881" {
		body = `{"Message":"OK","RequestId":"F69545AD-66DC-53BE-B5BD-0E4D2E147AF7","Code":"isv.SMS_TEST_NUMBER_LIMIT","BizId":"641921677331750484^0"}`
	}

	return body, v.Unmarshal([]byte(body))
}

func (c ClientTest) PostFormWithUnmarshal(api string, data string, v core.ResponseInterface) (string, error) {
	panic("implement me")
}
