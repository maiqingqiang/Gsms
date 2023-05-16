package qcloud

import (
	"bytes"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"github.com/maiqingqiang/gsms"
	"github.com/maiqingqiang/gsms/utils"
	"github.com/maiqingqiang/gsms/utils/dove"
	"net/http"
	"time"
)

const NAME = "qcloud"

const EndpointUrl = "https://sms.tencentcloudapi.com"
const EndpointHost = "sms.tencentcloudapi.com"
const EndpointService = "sms"
const EndpointMethod = "SendSms"
const EndpointVersion = "2021-01-11"
const EndpointRegion = "ap-guangzhou"
const Ok = "Ok"

// SendSmsRequest 请求参数 https://cloud.tencent.com/document/api/382/55981
// https://github.com/TencentCloud/signature-process-demo/blob/main/sms/signature-v3/golang/demo.go
type SendSmsRequest struct {
	// 下发手机号码，采用 E.164 标准，格式为+[国家或地区码][手机号]，单次请求最多支持200个手机号且要求全为境内手机号或全为境外手机号。
	// 例如：+8613711112222， 其中前面有一个+号 ，86为国家码，13711112222为手机号。
	PhoneNumberSet []string `json:"PhoneNumberSet,omitempty"`

	// 短信 SdkAppId，在 [短信控制台](https://console.cloud.tencent.com/smsv2/app-manage)  添加应用后生成的实际 SdkAppId，示例如1400006666。
	SmsSdkAppId string `json:"SmsSdkAppId,omitempty"`

	// 模板 ID，必须填写已审核通过的模板 ID。模板 ID 可登录 [短信控制台](https://console.cloud.tencent.com/smsv2) 查看，若向境外手机号发送短信，仅支持使用国际/港澳台短信模板。
	TemplateId string `json:"TemplateId,omitempty"`

	// 短信签名内容，使用 UTF-8 编码，必须填写已审核通过的签名，例如：腾讯云，签名信息可登录 [短信控制台](https://console.cloud.tencent.com/smsv2)  查看。
	// 注：国内短信为必填参数。
	SignName string `json:"SignName,omitempty"`

	// 模板参数，若无模板参数，则设置为空。
	TemplateParamSet []string `json:"TemplateParamSet,omitempty"`

	// 短信码号扩展号，默认未开通，如需开通请联系 [sms helper](https://cloud.tencent.com/document/product/382/3773#.E6.8A.80.E6.9C.AF.E4.BA.A4.E6.B5.81)。
	ExtendCode string `json:"ExtendCode,omitempty"`

	// 用户的 session 内容，可以携带用户侧 ID 等上下文信息，server 会原样返回。
	SessionContext string `json:"SessionContext,omitempty"`

	// 国内短信无需填写该项；国际/港澳台短信已申请独立 SenderId 需要填写该字段，默认使用公共 SenderId，无需填写该字段。
	// 注：月度使用量达到指定量级可申请独立 SenderId 使用，详情请联系 [sms helper](https://cloud.tencent.com/document/product/382/3773#.E6.8A.80.E6.9C.AF.E4.BA.A4.E6.B5.81)。
	SenderId string `json:"SenderId,omitempty"`
}

type SendSmsResponse struct {
	Response *struct {
		SendStatusSet []*struct {
			SerialNo       string `json:"SerialNo"`
			PhoneNumber    string `json:"PhoneNumber"`
			Fee            int    `json:"Fee"`
			SessionContext string `json:"SessionContext"`
			Code           string `json:"Code"`
			Message        string `json:"Message"`
			IsoCode        string `json:"IsoCode"`
		} `json:"SendStatusSet"`

		Error *struct {
			Code    string `json:"Code"`
			Message string `json:"Message"`
		} `json:"Error"`

		RequestId string `json:"RequestId"`
	} `json:"SendSmsResponse"`
}

var _ gsms.Gateway = (*Gateway)(nil)

type Gateway struct {
	SdkAppId  string
	SecretId  string
	SecretKey string
	SignName  string
}

func (g *Gateway) Name() string {
	return NAME
}

func (g *Gateway) Send(to *gsms.PhoneNumber, message gsms.Message, config *gsms.Config) error {
	phone := fmt.Sprintf("%d", to.Number())
	if to.IDDCode() != 0 {
		phone = to.UniversalNumber()
	}

	template, err := message.GetTemplate(g)
	if err != nil {
		return err
	}

	data, err := message.GetData(g)
	if err != nil {
		return err
	}

	templateParamSet := make([]string, 0, len(data))

	for _, s := range data {
		templateParamSet = append(templateParamSet, s)
	}

	p := &SendSmsRequest{
		PhoneNumberSet: []string{
			phone,
		},
		SmsSdkAppId:      g.SdkAppId,
		SignName:         g.SignName,
		TemplateId:       template,
		TemplateParamSet: templateParamSet,
	}

	d := dove.New(dove.WithTimeout(config.Timeout), dove.WithLogger(config.Logger))

	timestamp := time.Now().Unix()

	payload, _ := json.Marshal(p)

	config.Logger.Infof("SendSmsRequest: %s", payload)

	header := http.Header{}
	header.Add("Authorization", g.generateSign(string(payload), timestamp))
	header.Add("Host", EndpointHost)
	header.Add("Content-Type", "application/json; charset=utf-8")
	header.Add("X-TC-Action", EndpointMethod)
	header.Add("X-TC-Region", EndpointRegion)
	header.Add("X-TC-Timestamp", fmt.Sprintf("%d", timestamp))
	header.Add("X-TC-Version", EndpointVersion)

	var response SendSmsResponse

	err = d.Request(http.MethodPost, EndpointUrl, header, bytes.NewReader(payload), &response)

	if err != nil {
		return err
	}

	if response.Response.Error != nil && response.Response.Error.Code != "" {
		return fmt.Errorf("send failed: %+v", response.Response.Error.Message)
	}

	for _, status := range response.Response.SendStatusSet {
		if status.Code != Ok {
			return fmt.Errorf("send failed, status message: %s", status.Message)
		}
	}

	return nil
}

// https://github.com/TencentCloud/signature-process-demo/blob/main/sms/signature-v3/golang/demo.go
func (g *Gateway) generateSign(params string, timestamp int64) string {
	// step 1: build canonical request string
	algorithm := "TC3-HMAC-SHA256"
	httpRequestMethod := "POST"
	canonicalURI := "/"
	canonicalQueryString := ""
	canonicalHeaders := "content-type:application/json; charset=utf-8\n" + "host:" + EndpointHost + "\n"
	signedHeaders := "content-type;host"
	hashedRequestPayload := utils.Sha256Hex(params)
	canonicalRequest := fmt.Sprintf("%s\n%s\n%s\n%s\n%s\n%s",
		httpRequestMethod,
		canonicalURI,
		canonicalQueryString,
		canonicalHeaders,
		signedHeaders,
		hashedRequestPayload)

	// step 2: build string to sign
	date := time.Unix(timestamp, 0).UTC().Format("2006-01-02")
	credentialScope := fmt.Sprintf("%s/%s/tc3_request", date, EndpointService)
	hashedCanonicalRequest := utils.Sha256Hex(canonicalRequest)
	string2sign := fmt.Sprintf("%s\n%d\n%s\n%s",
		algorithm,
		timestamp,
		credentialScope,
		hashedCanonicalRequest)

	// step 3: sign string
	secretDate := utils.HmacSha256(date, "TC3"+g.SecretKey)
	secretService := utils.HmacSha256(EndpointService, secretDate)
	secretSigning := utils.HmacSha256("tc3_request", secretService)
	signature := hex.EncodeToString([]byte(utils.HmacSha256(string2sign, secretSigning)))

	// step 4: build authorization
	authorization := fmt.Sprintf("%s Credential=%s/%s, SignedHeaders=%s, Signature=%s",
		algorithm,
		g.SecretId,
		credentialScope,
		signedHeaders,
		signature)

	return authorization
}
