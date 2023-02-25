package aliyun

import (
	"crypto/hmac"
	"crypto/sha1"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"github.com/maiqingqiang/gsms/core"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
	"time"
)

var _ core.GatewayInterface = (*Gateway)(nil)

type Gateway struct {
	AccessKeyId     string
	AccessKeySecret string
	SignName        string
}

func (g *Gateway) Name() string {
	return NAME
}

func (g *Gateway) Send(to core.PhoneNumberInterface, message core.MessageInterface) (string, error) {
	req, err := http.NewRequest(http.MethodGet, EndpointUrl, nil)
	if err != nil {
		return "", err
	}
	query := req.URL.Query()

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
		return "", err
	}
	query.Add("TemplateCode", template)

	data, err := message.GetData(g)
	if err != nil {
		return "", err
	}
	marshal, err := json.Marshal(data)
	if err != nil {
		return "", err
	}
	query.Add("TemplateParam", string(marshal))

	query.Add("Signature", generateSign(http.MethodGet, g.AccessKeySecret, query))

	req.URL.RawQuery = query.Encode()

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
		return "", errors.New(fmt.Sprintf("send failed code:%s msg:%s", response.Code, response.Message))
	}

	return string(body), nil
}

// generateSign Generate sign.
// https://help.aliyun.com/document_detail/101343.html
func generateSign(httpMethod, accessKeySecret string, query url.Values) string {
	httpMethod = strings.ToUpper(httpMethod)

	encode := url.QueryEscape(query.Encode())

	encode = strings.Replace(encode, "+", "%20", -1)
	encode = strings.Replace(encode, "*", "%2A", -1)
	encode = strings.Replace(encode, "%7E", "~", -1)

	h := hmac.New(sha1.New, []byte(accessKeySecret+"&"))

	h.Write([]byte(fmt.Sprintf("%s&%%2F&%s", httpMethod, encode)))

	return base64.StdEncoding.EncodeToString(h.Sum(nil))
}
