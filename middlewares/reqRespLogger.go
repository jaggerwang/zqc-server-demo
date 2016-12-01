package middlewares

import (
	"os"
	"path/filepath"

	log "github.com/Sirupsen/logrus"
	"github.com/labstack/echo"
	"github.com/spf13/viper"
)

func ReqRespLogger() echo.MiddlewareFunc {
	logger := log.New()
	logger.Formatter = &log.JSONFormatter{}

	lvl, err := log.ParseLevel(viper.GetString("log.level"))
	if err != nil {
		panic(err)
	}
	logger.Level = lvl

	w, err := os.OpenFile(filepath.Join(viper.GetString("dir.data"), viper.GetString("log.reqresp.file")), os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0640)
	if err != nil {
		panic(err)
	}
	logger.Out = w

	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) (err error) {
			req := c.Request()
			resp := c.Response()
			if err = next(c); err != nil {
				c.Error(err)
			}
			params, err := c.FormParams()
			if err != nil {
				c.Error(err)
			}
			for _, vs := range params {
				for i, v := range vs {
					if len(v) > 100 {
						vs[i] = v[:100]
					}
				}
			}
			logger.WithFields(log.Fields{
				"req": log.Fields{
					"method": req.Method,
					"url":    req.URL,
					"header": req.Header,
					"params": params,
				},
				"resp": log.Fields{
					"status": resp.Status,
					"header": resp.Header,
					"body":   c.Get("respBody"),
				},
			}).Debug("request and response content")
			return err
		}
	}
}
