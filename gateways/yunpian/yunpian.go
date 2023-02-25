package yunpian

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/maiqingqiang/gsms/core"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
	"time"
)

var _ core.GatewayInterface = (*Gateway)(nil)

type Gateway struct {
	ApiKey    string
	Signature string
}

// Send message.
func (g *Gateway) Send(to core.PhoneNumberInterface, message core.MessageInterface) (string, error) {

	p := url.Values{}
	method := MethodSingleSend
	p.Add("apikey", g.ApiKey)
	p.Add("mobile", to.UniversalNumber())

	template, err := message.GetTemplate(g)
	if err != nil {
		return "", err
	}

	data, err := message.GetData(g)
	if err != nil {
		return "", err
	}

	content, err := message.GetContent(g)
	if err != nil {
		return "", err
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

	req, err := http.NewRequest(http.MethodPost, endpoint, strings.NewReader(p.Encode()))
	if err != nil {
		return "", err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	client := &http.Client{
		Timeout: 5 * time.Second,
	}

	res, err := client.Do(req)
	if err != nil {
		return "", err
	}

	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return "", err
	}

	response := &Response{}

	err = json.Unmarshal(body, &response)
	if err != nil {
		return "", err
	}

	if !response.isSuccessful() {
		return "", errors.New(fmt.Sprintf("send failed code:%d msg:%s detail:%s", response.Code, response.Msg, response.Detail))
	}

	return string(body), nil
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
