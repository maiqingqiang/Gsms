package yunpian

import (
	"fmt"
	"github.com/maiqingqiang/gsms"
	"github.com/maiqingqiang/gsms/utils/dove"
	"net/url"
	"strings"
)

const NAME = "yunpian"

const EndpointTemplate = "https://%s.yunpian.com/%s/%s/%s.%s"
const EndpointVersion = "v2"
const EndpointFormat = "json"
const ProductSms = "sms"
const ResourceSms = "sms"
const SuccessCode = 0

// MethodSingleSend https://www.yunpian.com/official/document/sms/zh_cn/domestic_single_send
const MethodSingleSend = "single_send"

// MethodTplSingleSend https://www.yunpian.com/official/document/sms/zh_CN/domestic_tpl_single_send
const MethodTplSingleSend = "tpl_single_send"

var _ gsms.Gateway = (*Gateway)(nil)

type Gateway struct {
	ApiKey    string
	Signature string
}

type SendSmsResponse struct {
	Code   int    `json:"code"`   // 系统返回码
	Msg    string `json:"msg"`    // 例如""发送成功""，或者相应错误信息
	Detail string `json:"detail"` // 具体错误描述或解决方法
}

// Send message.
func (g *Gateway) Send(to *gsms.PhoneNumber, message gsms.Message, config *gsms.Config) (err error) {
	p := url.Values{}
	method := MethodSingleSend
	p.Add("apikey", g.ApiKey)
	p.Add("mobile", to.UniversalNumber())

	template, err := message.GetTemplate(g)
	if err != nil {
		return
	}

	data, err := message.GetData(g)
	if err != nil {
		return
	}

	content, err := message.GetContent(g)
	if err != nil {
		return
	}

	if template != "" {
		method = MethodTplSingleSend
		p.Add("tpl_id", template)
		p.Add("tpl_value", g.buildTplVal(data))
	} else {
		if !strings.HasPrefix(content, "【") {
			content = g.Signature + content
		}

		p.Add("text", content)
	}

	endpoint := g.buildEndpoint(ProductSms, ResourceSms, method)

	var response SendSmsResponse

	d := dove.New(dove.WithTimeout(config.Timeout), dove.WithLogger(config.Logger))

	err = d.PostForm(endpoint, strings.NewReader(p.Encode()), &response)
	if err != nil {
		return err
	}

	if response.Code != SuccessCode {
		return fmt.Errorf("send failed code:%d msg:%s detail:%s", response.Code, response.Msg, response.Detail)
	}

	return
}

// buildEndpoint Build endpoint url.
func (g *Gateway) buildEndpoint(product, resource, method string) string {
	return fmt.Sprintf(EndpointTemplate, product, EndpointVersion, resource, method, EndpointFormat)
}

// buildTplVal Build template value.
func (g *Gateway) buildTplVal(data map[string]string) string {
	tplVals := make([]string, 0, len(data))
	for k, v := range data {
		tplVals = append(tplVals, fmt.Sprintf("#%s#=%s", k, v))
	}

	return strings.Join(tplVals, "&")
}

// Name Get gateway name.
func (g *Gateway) Name() string {
	return NAME
}
