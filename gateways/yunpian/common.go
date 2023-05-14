package yunpian

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

type Response struct {
	Code   int    `json:"code"`   // 系统返回码
	Msg    string `json:"msg"`    // 例如""发送成功""，或者相应错误信息
	Detail string `json:"detail"` // 具体错误描述或解决方法
}
