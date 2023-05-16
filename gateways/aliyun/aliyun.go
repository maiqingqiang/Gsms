package aliyun

import (
	"crypto/hmac"
	"crypto/sha1"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"github.com/maiqingqiang/gsms"
	"github.com/maiqingqiang/gsms/utils/dove"
	"net/http"
	"net/url"
	"strings"
	"time"
)

const NAME = "aliyun"

const EndpointUrl = "http://dysmsapi.aliyuncs.com"
const EndpointMethod = "SendSms"
const EndpointVersion = "2017-05-25"
const EndpointFormat = "JSON"
const EndpointRegionId = "cn-hangzhou"
const EndpointSignatureMethod = "HMAC-SHA1"
const EndpointSignatureVersion = "1.0"
const OK = "OK"

var _ gsms.Gateway = (*Gateway)(nil)

type Gateway struct {
	AccessKeyId     string
	AccessKeySecret string
	SignName        string
}

// SendSmsResponse
// https://help.aliyun.com/document_detail/419273.htm?spm=a2c4g.11186623.0.0.4a0879bebUrJyq#api-detail-40
type SendSmsResponse struct {
	Code      string `json:"Code"`
	Message   string `json:"Message"`
	BizId     string `json:"BizId"`
	RequestId string `json:"RequestId"`
}

func (g *Gateway) Name() string {
	return NAME
}

func (g *Gateway) Send(to *gsms.PhoneNumber, message gsms.Message, config *gsms.Config) error {
	query := url.Values{}

	query.Add("RegionId", EndpointRegionId)
	query.Add("AccessKeyId", g.AccessKeyId)
	query.Add("Format", EndpointFormat)
	query.Add("SignatureMethod", EndpointSignatureMethod)
	query.Add("SignatureVersion", EndpointSignatureVersion)
	query.Add("SignatureNonce", uuid.New().String())
	query.Add("Timestamp", time.Now().UTC().Format("2006-01-02T15:04:05Z"))
	query.Add("Action", EndpointMethod)
	query.Add("Version", EndpointVersion)
	query.Add("PhoneNumbers", to.UniversalNumber())
	query.Add("SignName", g.SignName)

	template, err := message.GetTemplate(g)
	if err != nil {
		return err
	}
	query.Add("TemplateCode", template)

	data, err := message.GetData(g)
	if err != nil {
		return err
	}
	marshal, err := json.Marshal(data)
	if err != nil {
		return err
	}
	query.Add("TemplateParam", string(marshal))

	query.Add("Signature", g.generateSign(query))

	var response SendSmsResponse

	d := dove.New(dove.WithTimeout(config.Timeout), dove.WithLogger(config.Logger))

	err = d.Get(EndpointUrl, strings.NewReader(query.Encode()), &response)
	if err != nil {
		return err
	}

	if response.Code != OK {
		return fmt.Errorf("send failed: %+v", response)
	}

	return nil
}

// generateSign Generate sign.
// https://help.aliyun.com/document_detail/101343.html
func (g *Gateway) generateSign(query url.Values) string {
	encode := url.QueryEscape(query.Encode())

	encode = strings.Replace(encode, "+", "%20", -1)
	encode = strings.Replace(encode, "*", "%2A", -1)
	encode = strings.Replace(encode, "%7E", "~", -1)

	h := hmac.New(sha1.New, []byte(g.AccessKeySecret+"&"))

	h.Write([]byte(fmt.Sprintf("%s&%%2F&%s", http.MethodGet, encode)))

	return base64.StdEncoding.EncodeToString(h.Sum(nil))
}
