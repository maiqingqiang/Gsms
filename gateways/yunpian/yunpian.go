package yunpian

import (
	"fmt"
	"github.com/maiqingqiang/gsms"
	"github.com/maiqingqiang/gsms/utils/dove"
	"net/url"
	"strings"
)

var _ gsms.Gateway = (*Gateway)(nil)

type Gateway struct {
	ApiKey    string
	Signature string
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
		p.Add("tpl_value", buildTplVal(data))
	} else {
		if !strings.HasPrefix(content, "【") {
			content = g.Signature + content
		}

		p.Add("text", content)
	}

	endpoint := buildEndpoint(ProductSms, ResourceSms, method)

	var response Response

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
func buildEndpoint(product, resource, method string) string {
	return fmt.Sprintf(EndpointTemplate, product, EndpointVersion, resource, method, EndpointFormat)
}

// buildTplVal Build template value.
func buildTplVal(data map[string]string) string {
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
