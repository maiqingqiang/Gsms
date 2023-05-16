package aliyun

import (
	"github.com/jarcoal/httpmock"
	"github.com/maiqingqiang/gsms"
	"github.com/maiqingqiang/gsms/message"
	"github.com/stretchr/testify/assert"
	"net/url"
	"testing"
	"time"
)

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
			g := &Gateway{AccessKeySecret: tt.args.accessKeySecret}

			assert.Equal(t, tt.want, g.generateSign(tt.args.query))
		})
	}
}

func TestGateway_Send(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	httpmock.RegisterResponder("GET", `http://dysmsapi.aliyuncs.com`,
		httpmock.NewStringResponder(200, `{"Message":"发送成功","RequestId":"F69545AD-66DC-53BE-B5BD-0E4D2E147AF1","Code":"OK"}`))

	g := &Gateway{
		AccessKeyId:     "AccessKeyId",
		AccessKeySecret: "AccessKeySecret",
		SignName:        "SignName",
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

	httpmock.RegisterResponder("GET", `http://dysmsapi.aliyuncs.com`,
		httpmock.NewStringResponder(200, `{"Message":"发送失败","RequestId":"F69545AD-66DC-53BE-B5BD-0E4D2E147AF1","Code":"FAIL"}`))

	g := &Gateway{
		AccessKeyId:     "AccessKeyId",
		AccessKeySecret: "AccessKeySecret",
		SignName:        "SignName",
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
