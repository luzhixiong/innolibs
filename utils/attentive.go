package utils

import (
	"encoding/json"
	"fmt"
)

const (
	attentiveMobileApiUrl = "https://api.attentivemobile.com/v1/text/send"
)

type attentiveMessage struct {
	To string `json:"to"` // 接收者号码，To和subscriberExternalId至少必填一个
	//SubscriberExternalId int64  `json:"subscriberExternalId"` //
	SubscriptionType string `json:"subscriptionType"` // 必填 MARKETING/TRANSACTIONAL
	Body             string `json:"body"`             // 内容 必填
	//MediaUrl         string `json:"mediaUrl"`         // 彩信图片地址
	UseShortLinks bool   `json:"useShortLinks"` // 自动缩短正文中的链接 必填
	MessageName   string `json:"messageName"`   // 为了跟踪目的而创建的消息的名称。 每个公司 24 小时内最多只能发送 20 条消息。
	SkipFatigue   bool   `json:"skipFatigue"`   // 强制向疲劳的订阅者发送消息。 适用于所有订阅类型。 默认true
}

func attentiveHeader(apiKey string) map[string]string {
	return map[string]string{
		"Authorization": fmt.Sprintf(`Bearer %s`, apiKey),
		"Content-Type":  `application/json`,
	}
}

func AttentiveSendSms(apiKey, to, content string) (resp []byte, err error) {
	body := &attentiveMessage{
		To: to,
		//SubscriberExternalId: to,
		Body:             content,
		SubscriptionType: "TRANSACTIONAL",
		UseShortLinks:    true,
		MessageName:      "FlyDelivery Notification",
		SkipFatigue:      true,
	}
	data, _ := json.Marshal(body)
	res, err := HttpPostWithHeader(attentiveMobileApiUrl, data, attentiveHeader(apiKey))
	if err != nil {
		return
	}
	resp = res.Body()
	return
}
