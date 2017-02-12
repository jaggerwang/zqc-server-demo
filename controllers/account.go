package controllers

import (
	"net/http"

	valid "github.com/asaskevich/govalidator"
	"github.com/labstack/echo"
	"gopkg.in/mgo.v2/bson"

	"zqc/middlewares"
	"zqc/services"
)

func RegisterAccount(c echo.Context) (err error) {
	cc := c.(*middlewares.Context)
	params := struct {
		Mobile   string `valid:"matches(^[0-9]{11}$)"`
		Password string `valid:"stringlength(6|20)"`
		Nickname string `valid:"stringlength(3|20)"`
		Gender   string `valid:"stringlength(1|1)"`
	}{
		Mobile:   cc.FormValue("mobile"),
		Password: cc.FormValue("password"),
		Nickname: cc.FormValue("nickname"),
		Gender:   cc.FormValue("gender"),
	}
	if ok, err := valid.ValidateStruct(params); !ok {
		return services.NewError(services.ErrCodeInvalidParams, err.Error())
	}

	user, err := services.CreateUser(params.Mobile, params.Password, params.Nickname, params.Gender)
	if err != nil {
		return err
	}

	return ResponseJSON(http.StatusOK, Response{
		Data: map[string]interface{}{
			"user": user,
		},
	}, cc)
}

func Login(c echo.Context) (err error) {
	cc := c.(*middlewares.Context)
	params := struct {
		Mobile   string `valid:"matches(^[0-9]{11}$)"`
		Password string `valid:"stringlength(6|20)"`
	}{
		Mobile:   cc.FormValue("mobile"),
		Password: cc.FormValue("password"),
	}
	if ok, err := valid.ValidateStruct(params); !ok {
		return services.NewError(services.ErrCodeInvalidParams, err.Error())
	}

	user, err := services.GetUserByMobile(params.Mobile)
	if err != nil {
		return services.NewError(services.ErrCodeNotFound, "mobile not exist")
	}

	user, err = services.VerifyUserPassword(user.Id, params.Password)
	if err != nil {
		return err
	}

	cc.SetSessionItem("userId", user.Id)

	return ResponseJSON(http.StatusOK, Response{
		Data: map[string]interface{}{
			"user": user,
		},
	}, cc)
}

func IsLogined(c echo.Context) (err error) {
	cc := c.(*middlewares.Context)

	id := cc.SessionUserId()
	if id != "" {
		user, err := services.GetUser(id)
		if err != nil {
			return err
		}
		return ResponseJSON(http.StatusOK, Response{
			Data: map[string]interface{}{
				"user": user,
			},
		}, cc)
	} else {
		return ResponseJSON(http.StatusOK, Response{}, cc)
	}

}

func Logout(c echo.Context) (err error) {
	cc := c.(*middlewares.Context)
	cc.DeleteSession()

	return ResponseJSON(http.StatusOK, Response{}, cc)
}

func EditAccount(c echo.Context) (err error) {
	cc := c.(*middlewares.Context)
	params := struct {
		Mobile   string `valid:"matches(^[0-9]{11}$),optional"`
		Nickname string `valid:"stringlength(3|20),optional"`
		Gender   string `valid:"matches(^m|f$),optional"`
	}{
		Mobile:   cc.FormValue("mobile"),
		Nickname: cc.FormValue("nickname"),
		Gender:   cc.FormValue("gender"),
	}
	if ok, err := valid.ValidateStruct(&params); !ok {
		return services.NewError(services.ErrCodeInvalidParams, err.Error())
	}

	formParams, err := cc.FormParams()
	if err != nil {
		return services.NewError(services.ErrCodeInvalidParams, err.Error())
	}
	updateParams := bson.M{}
	if _, ok := formParams["mobile"]; ok {
		updateParams["mobile"] = params.Mobile
	}
	if _, ok := formParams["nickname"]; ok {
		updateParams["nickname"] = params.Nickname
	}
	if _, ok := formParams["gender"]; ok {
		updateParams["gender"] = params.Gender
	}

	id := cc.SessionUserId()
	user, err := services.UpdateUser(id, updateParams)
	if err != nil {
		return err
	}

	return ResponseJSON(http.StatusOK, Response{
		Data: map[string]interface{}{
			"user": user,
		},
	}, cc)
}

func AccountInfo(c echo.Context) (err error) {
	cc := c.(*middlewares.Context)

	id := cc.SessionUserId()
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
