package middlewares

import (
	"github.com/labstack/echo"
	"github.com/spf13/viper"
	"gopkg.in/boj/redistore.v1"

	"zqc/models"
)

func Session() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) (err error) {
			store, err := redistore.NewRediStoreWithPool(
				models.RedisPool("zqc"), []byte(viper.GetString("secretkey")))
			if err != nil {
				panic(err)
			}
			store.SetMaxAge(viper.GetInt("session.maxAge"))
			store.SetMaxLength(viper.GetInt("session.maxLength"))
			store.SetKeyPrefix(viper.GetString("session.keyPrefix"))

			session, err := store.Get(c.Request(), viper.GetString("session.name"))
			if err != nil {
				panic(err)
			}
			c.Set("session", session)

			err = next(c)

			return err
		}
	}
}
