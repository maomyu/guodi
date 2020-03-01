package guodiemail

import (
	"fmt"
	"math/rand"
	"net/smtp"
	"strings"
	"time"

	"guodi/src/guodiredis"
)

func SendEmail(email string) bool {
	randnum := rand.New(rand.NewSource(time.Now().UnixNano())).Intn(10000)
	//代理参数
	auth := smtp.PlainAuth("", "nylgisc@163.com", "nylgisc0410", "smtp.163.com")
	to := []string{email}
	//邮件显示名
	nickname := "智慧校园"
	user := "nylgisc@163.com"
	//标题
	subject := "<noreply>裹递账号验证码"
	content_type := "Content-Type: text/plain; charset=UTF-8"
	body := fmt.Sprintf("邮箱验证码是： %v ,请不要将验证码展示给任何人", randnum)
	msg := []byte("To: " + strings.Join(to, ",") + "\r\nFrom: " + nickname +
		"<" + user + ">\r\nSubject: " + subject + "\r\n" + content_type + "\r\n\r\n" + body)
	//如果发送失败，发送邮件的账号会收到回执邮件，包含失败原因
	err := smtp.SendMail("smtp.163.com:25", auth, user, to, msg)
	if err != nil {
		return false
	}
	//guodizap.Debug()
	guodiredis.SaveAuthenticate(email, string(randnum))
	return true
}
