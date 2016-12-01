package controllers

import (
	"fmt"
	"net/http"

	log "github.com/Sirupsen/logrus"
	valid "github.com/asaskevich/govalidator"
	"github.com/labstack/echo"

	"jaggerwang.net/zqcserverdemo/middlewares"
	"jaggerwang.net/zqcserverdemo/services"
)

type SendVerifyCodeParams struct {
	By     string `valid:"matches(^mobile|email$)"`
	Mobile string `valid:"matches(^[0-9]{11}$),optional"`
	Email  string `valid:"email,optional"`
}

func SendVerifyCode(c echo.Context) (err error) {
	cc := c.(*middlewares.Context)
	params := SendVerifyCodeParams{
		By:     cc.FormValue("by"),
		Mobile: cc.FormValue("mobile"),
		Email:  cc.FormValue("email"),
	}
	if ok, err := valid.ValidateStruct(params); !ok {
		return services.NewServiceError(services.ErrCodeInvalidParams, err.Error())
	}
	if params.Mobile == "" && params.Email == "" {
		return services.NewServiceError(services.ErrCodeInvalidParams, "either username or mobile is needed")
	}

	var verifyCode *services.VerifyCode
	if params.By == "mobile" {
		verifyCode, err = services.SendVerifyCodeByMobile(params.Mobile)
		if err != nil {
			return err
		}
	} else if params.By == "email" {
		verifyCode, err = services.SendVerifyCodeByEmail(params.Email)
		if err != nil {
			return err
		}
	} else {
		return services.NewServiceError(services.ErrCodeInvalidParams, fmt.Sprintf("can not send by %s", params.By))
	}

	cc.SetSessionItem("verifyCode", verifyCode)
	log.WithField("verifyCode", verifyCode).Info("send verify code")

	return ResponseJSON(http.StatusOK, Response{}, cc)
}
