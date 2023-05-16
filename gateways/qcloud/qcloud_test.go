package qcloud

import (
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestGateway_generateSign(t *testing.T) {
	g := &Gateway{
		SdkAppId:  "SdkAppId",
		SecretId:  "SecretId",
		SecretKey: "SecretKey",
		SignName:  "gsms",
	}

	params := &SendSmsRequest{
		PhoneNumberSet:   []string{"18888888888"},
		SmsSdkAppId:      g.SdkAppId,
		TemplateId:       "1111111",
		SignName:         g.SignName,
		TemplateParamSet: []string{"521410", "5"},
	}
	payload, _ := json.Marshal(params)

	sign := g.generateSign(string(payload), 1684225049)

	assert.Equal(
		t,
		"TC3-HMAC-SHA256 Credential=SecretId/2023-05-16/sms/tc3_request, SignedHeaders=content-type;host, Signature=c339e750c3b92ca97783de2c8d00434b12cf8f22c45d421fa8b0609d64e51358",
		sign,
	)
}
