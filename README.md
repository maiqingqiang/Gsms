<a name="readme-top"></a>

<!-- PROJECT SHIELDS -->

[![go report card][go-report-card]][go-report-card-url]
[![Go.Dev reference][go.dev-reference]][go.dev-reference-url]
[![Go package][go-pacakge]][go-pacakge-url]
[![MIT License][license-shield]][license-url]
[![Contributors][contributors-shield]][contributors-url]
[![Forks][forks-shield]][forks-url]
[![Stargazers][stars-shield]][stars-url]
[![Issues][issues-shield]][issues-url]


<h1 align="center">Gsms :calling: </h1>

<p align="center">:calling: 一款满足你的多种发送需求的短信发送组件，此项目参考 <a href="https://github.com/overtrue/easy-sms">easy-sms</a> 实现的 Go 版本 </p>

> 正在开发中~~~~~~~~~~~

## 特点

1. 支持目前市面多家服务商
2. 一套写法兼容所有平台
3. 简单配置即可灵活增减服务商
4. 内置多种服务商轮询策略、支持自定义轮询策略
5. 统一的返回值格式，便于日志与监控
6. 自动轮询选择可用的服务商
7. 更多等你去发现与改进...

## 平台支持

- [云片](https://www.yunpian.com)
- [阿里云](https://www.aliyun.com/)

## 安装

```shell
go get -u github.com/maiqingqiang/gsms
```

## 使用

```go
package main

import (
	"github.com/maiqingqiang/gsms"
	"github.com/maiqingqiang/gsms/core"
	"github.com/maiqingqiang/gsms/gateways/yunpian"
	"log"
)

func main() {
	g := gsms.New(
		[]core.Gateway{
			&yunpian.Gateway{
				ApiKey:    "f4c1c41f48120eb311111111111097",
				Signature: "【默认签名】",
			},
		},
		gsms.WithDefaultGateway([]string{
			yunpian.NAME,
		}),
	)

	results, err := g.Send(18888888888, &core.Message{
		Template: "SMS_00000001",
		Data: map[string]string{
			"code": "521410",
		},
	})

	if err != nil {
		log.Fatalf("发送失败 %+v", err)
	}

	log.Printf("发送成功 %+v", results)
}
```

## 短信内容

由于使用多网关发送，所以一条短信要支持多平台发送，每家的发送方式不一样，但是我们抽象定义了以下公用属性：
- `Content` 文字内容，使用在像云片类似的以文字内容发送的平台
- `Template` 模板 ID，使用在以模板ID来发送短信的平台
- `Data`  模板变量，使用在以模板ID来发送短信的平台

所以，在使用过程中你可以根据所要使用的平台定义发送的内容。

```go
client.Send(18888888888, &core.Message{
    Template: "SMS_00000001",
    Data: map[string]string{
        "code": "521410",
    },
})

client.Send(18888888888, &core.Message{
    Content: "您的验证码为: 6379",
})
```

### 闭包方式
```go
client.Send(18888888888, &core.Message{
    Template: func(gateway core.GatewayInterface) string {
        if gateway.Name() == aliyun.NAME {
            return "SMS_271385117"
        }
        return "5532044"
    },
    Data: func(gateway core.GatewayInterface) map[string]string {
        if gateway.Name() == aliyun.NAME {
            return map[string]string{
                "code": "1111",
            }
        }
        return map[string]string{
            "code": "6379",
        }
    },
})
```

## 发送网关
```go
client.Send(18888888888, &core.Message{
    Template: "5532044",
    Data: map[string]string{
        "code": "6379",
    },
}, yunpian.NAME, aliyun.NAME)
```

## 自定义网关

只需要实现 `core.GatewayInterface` 接口即可，例如：

## 场景发送

```go

var _ core.MessageInterface = (*OrderPaidMessage)(nil)

type OrderPaidMessage struct {
	OrderNo string
}

func (o *OrderPaidMessage) Gateways() ([]string, error) {
	return []string{yunpian.NAME}, nil
}

func (o *OrderPaidMessage) Strategy() (core.StrategyInterface, error) {
	return nil, nil
}

func (o *OrderPaidMessage) GetContent(gateway core.GatewayInterface) (string, error) {
	return fmt.Sprintf("您的订单:%s, 已经完成付款", o.OrderNo), nil
}

func (o *OrderPaidMessage) GetTemplate(gateway core.GatewayInterface) (string, error) {
	return "5532044", nil
}

func (o *OrderPaidMessage) GetData(gateway core.GatewayInterface) (map[string]string, error) {
	return map[string]string{
		"code": "6379",
	}, nil
}

func (o *OrderPaidMessage) GetType(gateway core.GatewayInterface) (string, error) {
	return core.TextMessage, nil
}

client.Send(18888888888, &OrderPaidMessage{OrderNo: "1234"})

```

## 版权说明

该项目签署了 GNU 授权许可，详情请参阅 [LICENSE](LICENSE)

## 鸣谢

- [安正超](https://github.com/overtrue)
- [easy-sms](https://github.com/overtrue/easy-sms)

<!-- MARKDOWN LINKS & IMAGES -->

[contributors-shield]: https://img.shields.io/github/contributors/maiqingqiang/Gsms.svg
[contributors-url]: https://github.com/maiqingqiang/Gsms/graphs/contributors
[forks-shield]: https://img.shields.io/github/forks/maiqingqiang/Gsms.svg
[forks-url]: https://github.com/maiqingqiang/Gsms/network/members
[stars-shield]: https://img.shields.io/github/stars/maiqingqiang/Gsms.svg
[stars-url]: https://github.com/maiqingqiang/Gsms/stargazers
[issues-shield]: https://img.shields.io/github/issues/maiqingqiang/Gsms.svg
[issues-url]: https://github.com/maiqingqiang/Gsms/issues
[license-shield]: https://img.shields.io/github/license/maiqingqiang/Gsms.svg
[license-url]: https://github.com/maiqingqiang/Gsms/blob/master/LICENSE.txt
[go-report-card]: https://goreportcard.com/badge/github.com/maiqingqiang/Gsms
[go-report-card-url]: https://goreportcard.com/report/github.com/maiqingqiang/Gsms
[go.dev-reference]: https://img.shields.io/badge/go.dev-reference-blue?logo=go&logoColor=white
[go.dev-reference-url]: https://pkg.go.dev/github.com/maiqingqiang/gsms?tab=doc
[go-pacakge]: https://github.com/maiqingqiang/Gsms/actions/workflows/test.yml/badge.svg?branch=main
[go-pacakge-url]: https://github.com/maiqingqiang/Gsms/actions/workflows/test.yml