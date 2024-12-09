package utils

import (
	"github.com/jordan-wright/email"
	"im/config"
	"log"
	"math/rand"
	"net/smtp"
	"strconv"
)

const (
	CodeSubject = "Verification Code"
)

func SendEmail(sbuject, html string, tos []string) error {
	email := email.NewEmail()
	email.From = "IM <" + config.EmailForm + ">"
	email.To = tos
	email.Subject = sbuject
	email.HTML = []byte(html)
	err := email.Send(config.EmailServer+":"+strconv.Itoa(config.EmailPort), smtp.PlainAuth("", config.EmailForm, config.EmailPassword, config.EmailServer))
	return err
}

func SendCode(code, to string) error {
	err := SendEmail(CodeSubject, "您的验证码是：<b>"+code+"</b>,验证码5分钟内有效,请妥善保管.", []string{to})
	if err != nil {
		log.Printf("send email fail %v", err)
	}
	return err
}

func GetNumCode(len int) string {
	s := ""
	for i := 0; i < len; i++ {
		s += strconv.Itoa(rand.Intn(10))
	}

	return s
}
