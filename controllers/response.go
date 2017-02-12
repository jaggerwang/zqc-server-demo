package controllers

import (
	"net/http"

	log "github.com/Sirupsen/logrus"
	"github.com/gorilla/sessions"
	"github.com/labstack/echo"

	"zqc/services"
)

type Response struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Context interface{} `json:"context"`
	Data    interface{} `json:"data"`
}

func ResponseJSON(status int, resp Response, c echo.Context) (err error) {
	if session, ok := c.Get("session").(*sessions.Session); ok {
		if session.IsNew || c.Get("sessionModified") != nil {
			if err := session.Save(c.Request(), c.Response()); err != nil {
				log.Error(err)
			}
		}
	}

	return c.JSON(status, resp)
}

func ErrorHandler(err error, c echo.Context) {
	status := http.StatusOK
	code := services.ErrCodeFail
	var message string
	var context interface{}

	if c.Echo().Debug {
		message = err.Error()
	}

	if he, ok := err.(*echo.HTTPError); ok {
		status = he.Code
		code = services.ErrCodeHttp
		message = he.Error()
	} else if se, ok := err.(*services.Error); ok {
		code = se.Code
		if c.Echo().Debug {
			message = se.Message
			context = se.Context
		} else {
			message = services.ErrMessages[code]
		}
	}

	if !c.Response().Committed {
		ResponseJSON(status, Response{
			Code:    code,
			Message: message,
			Context: context,
		}, c)
	}
}
