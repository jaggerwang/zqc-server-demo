package middlewares

import (
	"github.com/gorilla/sessions"
	"github.com/labstack/echo"
	"gopkg.in/mgo.v2/bson"
)

type Context struct {
	echo.Context
}

func (c *Context) Session() *sessions.Session {
	return c.Get("session").(*sessions.Session)
}

func (c *Context) DeleteSession() {
	c.Session().Options.MaxAge = -1
	c.Set("sessionModified", true)
}

func (c *Context) SetSessionItem(key string, value interface{}) {
	c.Session().Values[key] = value
	c.Set("sessionModified", true)
}

func (c *Context) DeleteSessionItem(key string) {
	delete(c.Session().Values, key)
	c.Set("sessionModified", true)
}

func (c *Context) SessionUserId() (userId bson.ObjectId) {
	if v, ok := c.Session().Values["userId"]; ok {
		return v.(bson.ObjectId)
	} else {
		return userId
	}
}

func MiddlewareContext() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) (err error) {
			cc := &Context{c}
			return next(cc)
		}
	}
}
