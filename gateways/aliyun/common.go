package aliyun

const NAME = "aliyun"

const EndpointUrl = "http://dysmsapi.aliyuncs.com"

const EndpointMethod = "SendSms"

const EndpointVersion = "2017-05-25"

const EndpointFormat = "JSON"

const EndpointRegionId = "cn-hangzhou"

const EndpointSignatureMethod = "HMAC-SHA1"

const EndpointSignatureVersion = "1.0"

const OK = "OK"

// Response
// https://help.aliyun.com/document_detail/419273.htm?spm=a2c4g.11186623.0.0.4a0879bebUrJyq#api-detail-40
type Response struct {
	Code      string `json:"Code"`
	Message   string `json:"Message"`
	BizId     string `json:"BizId"`
	RequestId string `json:"RequestId"`
}
