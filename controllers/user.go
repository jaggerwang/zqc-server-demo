package controllers

import (
	"net/http"

	valid "github.com/asaskevich/govalidator"
	"github.com/labstack/echo"

	"jaggerwang.net/zqcserverdemo/middlewares"
	"jaggerwang.net/zqcserverdemo/services"
)

type UserInfoParams struct {
	Id string `valid:"objectidhex"`
}

func UserInfo(c echo.Context) (err error) {
	cc := c.(*middlewares.Context)
	params := UserInfoParams{
		Id: cc.FormValue("id"),
	}
	if ok, err := valid.ValidateStruct(params); !ok {
		return services.NewServiceError(services.ErrCodeInvalidParams, err.Error())
	}
	id, err := ParseObjectId(params.Id)
	if err != nil {
		return services.NewServiceError(services.ErrCodeInvalidParams, err.Error())
	}

	user, err := services.GetUser(id)
	if err != nil {
		return err
	}

	return ResponseJSON(http.StatusOK, Response{
		Data: map[string]interface{}{
			"user": user,
		},
	}, cc)
}

type NearbyUsersParams struct {
	Location string `valid:"stringlength(3|30)"`
	Distance string `valid:"int,optional"`
	Limit    string `valid:"int,optional"`
}

func NearbyUsers(c echo.Context) (err error) {
	cc := c.(*middlewares.Context)
	params := NearbyUsersParams{
		Location: cc.FormValue("location"),
		Distance: cc.FormValue("distance"),
		Limit:    cc.FormValue("limit"),
	}
	if ok, err := valid.ValidateStruct(params); !ok {
		return services.NewServiceError(services.ErrCodeInvalidParams, err.Error())
	}
	loc, err := ParseLocation(params.Location)
	if err != nil {
		return services.NewServiceError(services.ErrCodeInvalidParams, err.Error())
	}
	dist := 50000
	if params.Distance != "" {
		dist, err = ParseInt(params.Distance, 1000, 10000)
		if err != nil {
			return services.NewServiceError(services.ErrCodeInvalidParams, "dist must between 1000 to 10000")
		}
	}
	limit := 10
	if params.Limit != "" {
		limit, err = ParseInt(params.Limit, 1, 100)
		if err != nil {
			return services.NewServiceError(services.ErrCodeInvalidParams, "limit must between 1 to 100")
		}
	}

	users, err := services.NearbyUsers(loc, dist, limit)
	if err != nil {
		return err
	}

	return ResponseJSON(http.StatusOK, Response{
		Data: map[string]interface{}{
			"users": users,
		},
	}, cc)
}
