package controllers

import (
	"net/http"

	log "github.com/Sirupsen/logrus"
	"github.com/gorilla/sessions"
	"github.com/labstack/echo"

	"zqc/services"
)

type Response struct {
	Code    int
	Message string
	Data    map[string]interface{}
	Context interface{}
}

func ResponseJSON(status int, resp Response, c echo.Context) (err error) {
	respBody := map[string]interface{}{
		"code":    resp.Code,
		"message": resp.Message,
		"data":    resp.Data,
		"context": resp.Context,
	}
	c.Set("respBody", respBody)

	if session, ok := c.Get("session").(*sessions.Session); ok {
		if session.IsNew || c.Get("sessionModified") != nil {
			if err := session.Save(c.Request(), c.Response()); err != nil {
				log.Error(err)
			}
		}
	}

	return c.JSON(status, respBody)
}

func HttpErrorHandler(err error, c echo.Context) {
	status := http.StatusOK
	code := services.ErrCodeFail
	var message string
	if c.Echo().Debug {
		message = err.Error()
	} else {
		message = ""
	}
	var context interface{}
	if he, ok := err.(*echo.HTTPError); ok {
		status = he.Code
		code = services.ErrCodeHttp
		message = he.Error()
	} else if se, ok := err.(*services.ServiceError); ok {
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
