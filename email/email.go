package email

import (
	"crypto/tls"
	"fmt"
	"github.com/goinggo/mapstructure"
	"github.com/jordan-wright/email"
	"net/smtp"
)

// 邮箱发送(配置参数)
func SendByConfig(emailMap map[string]interface{}, to []string, subject string, body string) (err error) {
	emailConfig := Email{}
	err = mapstructure.Decode(emailMap, &emailConfig)
	if err != nil {
		return err
	}
	from := emailConfig.From
	nickname := emailConfig.Nickname
	secret := emailConfig.Secret
	host := emailConfig.Host
	port := emailConfig.Port
	isSSL := emailConfig.IsSSL

	auth := smtp.PlainAuth("", from, secret, host)
	e := email.NewEmail()
	if nickname != "" {
		e.From = fmt.Sprintf("%s <%s>", nickname, from)
	} else {
		e.From = from
	}
	e.To = to
	e.Subject = subject
	e.HTML = []byte(body)
	hostAddr := fmt.Sprintf("%s:%d", host, port)
	if isSSL {
		err = e.SendWithTLS(hostAddr, auth, &tls.Config{ServerName: host})
	} else {
		err = e.Send(hostAddr, auth)
	}
	return err
}
