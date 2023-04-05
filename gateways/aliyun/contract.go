package aliyun

import (
	"encoding/json"
	"fmt"
	"github.com/maiqingqiang/gsms/core"
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

var _ core.ResponseInterface = (*Response)(nil)

// Response
// https://help.aliyun.com/document_detail/419273.htm?spm=a2c4g.11186623.0.0.4a0879bebUrJyq#api-detail-40
type Response struct {
	Code      string `json:"Code"`
	Message   string `json:"Message"`
	BizId     string `json:"BizId"`
	RequestId string `json:"RequestId"`
}

func (r *Response) Unmarshal(data []byte) error {
	if err := json.Unmarshal(data, r); err != nil {
		return fmt.Errorf("unmarshal response failed: %w data: %s", err, data)
	}

	if r.Code != OK {
		return fmt.Errorf("send failed code:%s msg:%s", r.Code, r.Message)
	}

	return nil
}
