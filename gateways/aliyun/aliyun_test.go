package aliyun

import (
	"github.com/stretchr/testify/assert"
	"net/url"
	"testing"
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
			assert.Equal(t, tt.want, generateSign(tt.args.httpMethod, tt.args.accessKeySecret, tt.args.query))
		})
	}
}
