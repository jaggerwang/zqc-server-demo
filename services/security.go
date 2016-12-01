package services

import (
	log "github.com/Sirupsen/logrus"

	"jaggerwang.net/zqcserverdemo/utils"
)

type VerifyCode struct {
	By     string
	Code   string
	Mobile string
	Email  string
}

func SendVerifyCodeByMobile(mobile string) (verifyCode *VerifyCode, err error) {
	code := utils.RandString(4, []rune("0123456789"))
	// TODO need to send really
	log.WithFields(log.Fields{
		"mobile": mobile,
		"code":   code,
	}).Info("send verify code by mobile")
	return &VerifyCode{
		By:     "mobile",
		Code:   code,
		Mobile: mobile,
	}, nil
}

func SendVerifyCodeByEmail(email string) (verifyCode *VerifyCode, err error) {
	code := utils.RandString(4, []rune("0123456789"))
	// TODO need to send really
	log.WithFields(log.Fields{
		"email": email,
		"code":  code,
	}).Info("send verify code by email")
	return &VerifyCode{
		By:    "email",
		Code:  code,
		Email: email,
	}, nil
}
