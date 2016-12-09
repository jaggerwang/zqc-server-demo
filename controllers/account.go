package controllers

import (
	"net/http"

	valid "github.com/asaskevich/govalidator"
	"github.com/labstack/echo"
	"gopkg.in/mgo.v2/bson"

	"zqc/middlewares"
	"zqc/services"
)

type RegisterAccountParams struct {
	Mobile   string `valid:"matches(^[0-9]{11}$)"`
	Password string `valid:"stringlength(6|20)"`
	Code     string `valid:"stringlength(4)"`
}

func RegisterAccount(c echo.Context) (err error) {
	cc := c.(*middlewares.Context)
	params := RegisterAccountParams{
		Mobile:   cc.FormValue("mobile"),
		Password: cc.FormValue("password"),
		Code:     cc.FormValue("code"),
	}
	if ok, err := valid.ValidateStruct(params); !ok {
		return services.NewServiceError(services.ErrCodeInvalidParams, err.Error())
	}

	codeOk := true
	if params.Mobile != "" {
		codeOk = false
		verifyCode := cc.SessionVerifyCode()
		if verifyCode.By == "mobile" && verifyCode.Mobile == params.Mobile && verifyCode.Code == params.Code {
			codeOk = true
		}
	}
	if !codeOk {
		return services.NewServiceError(services.ErrCodeInvalidVerifyCode, "")
	}

	user, err := services.CreateUser(params.Mobile, params.Password)
	if err != nil {
		return err
	}

	cc.DeleteSessionItem("verifyCode")

	return ResponseJSON(http.StatusOK, Response{
		Data: map[string]interface{}{
			"user": user,
		},
	}, cc)
}

type LoginParams struct {
	Mobile   string `valid:"matches(^[0-9]{11}$),optional"`
	Email    string `valid:"email,optional"`
	Password string `valid:"stringlength(6|20)"`
}

func Login(c echo.Context) (err error) {
	cc := c.(*middlewares.Context)
	params := LoginParams{
		Mobile:   cc.FormValue("mobile"),
		Email:    cc.FormValue("email"),
		Password: cc.FormValue("password"),
	}
	if ok, err := valid.ValidateStruct(params); !ok {
		return services.NewServiceError(services.ErrCodeInvalidParams, err.Error())
	}

	username := ""
	if params.Mobile != "" {
		user, err := services.GetUserByMobile(params.Mobile)
		if err != nil {
			return services.NewServiceError(services.ErrCodeNotFound, "mobile not exist")
		}
		username = user.Username
	} else if params.Email != "" {
		user, err := services.GetUserByEmail(params.Email)
		if err != nil {
			return services.NewServiceError(services.ErrCodeNotFound, "email not exist")
		}
		username = user.Username
	} else {
		return services.NewServiceError(services.ErrCodeInvalidParams, err.Error())
	}

	user, err := services.VerifyUserPassword(username, params.Password)
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

type ResetPasswordParams struct {
	Mobile   string `valid:"matches(^[0-9]{11}$),optional"`
	Email    string `valid:"email,optional"`
	Password string `valid:"stringlength(6|20)"`
	Code     string `valid:"stringlength(4)"`
}

func ResetPassword(c echo.Context) (err error) {
	cc := c.(*middlewares.Context)
	params := ResetPasswordParams{
		Mobile:   cc.FormValue("mobile"),
		Email:    cc.FormValue("email"),
		Password: cc.FormValue("password"),
		Code:     cc.FormValue("code"),
	}
	if ok, err := valid.ValidateStruct(params); !ok {
		return services.NewServiceError(services.ErrCodeInvalidParams, err.Error())
	}

	var user *services.User
	if params.Mobile != "" {
		user, err = services.GetUserByMobile(params.Mobile)
		if err != nil {
			return services.NewServiceError(services.ErrCodeNotFound, "")
		}
	} else if params.Email != "" {
		user, err = services.GetUserByEmail(params.Email)
		if err != nil {
			return services.NewServiceError(services.ErrCodeNotFound, "")
		}
	} else {
		return services.NewServiceError(services.ErrCodeNotFound, "")
	}

	verifyCode := cc.SessionVerifyCode()
	if (params.Mobile != "" && !(verifyCode.By == "mobile" && verifyCode.Mobile == params.Mobile && verifyCode.Code == params.Code)) ||
		(params.Email != "" && !(verifyCode.By == "email" && verifyCode.Email == params.Email && verifyCode.Code == params.Code)) {
		return services.NewServiceError(services.ErrCodeInvalidVerifyCode, "")
	}

	user, err = services.UpdateUser(user.Id, bson.M{"password": params.Password})
	if err != nil {
		return err
	}

	cc.DeleteSessionItem("verifyCode")

	return ResponseJSON(http.StatusOK, Response{
		Data: map[string]interface{}{
			"user": user,
		},
	}, cc)
}

func IsLogined(c echo.Context) (err error) {
	cc := c.(*middlewares.Context)
	var user *services.User
	id := cc.SessionUserId()
	if id != "" {
		user, _ = services.GetUser(id)
	}

	return ResponseJSON(http.StatusOK, Response{
		Data: map[string]interface{}{
			"user": user,
		},
	}, cc)
}

func Logout(c echo.Context) (err error) {
	cc := c.(*middlewares.Context)
	cc.DeleteSession()

	return ResponseJSON(http.StatusOK, Response{}, cc)
}

type EditAccountParams struct {
	Nickname       string `valid:"stringlength(3|20),optional"`
	Gender         string `valid:"matches(^m|f$),optional"`
	Mobile         string `valid:"matches(^[0-9]{11}$),optional"`
	Code           string `valid:"stringlength(4),optional"`
	AvatarType     string `valid:"matches(^builtin|custom$),optional"`
	AvatarName     string `valid:"stringlength(1|50),optional"`
	AvatarId       string `valid:"objectidhex,optional"`
	Email          string `valid:"email,optional"`
	Intro          string `valid:"stringlength(1|50),optional"`
	BackgroundType string `valid:"matches(^builtin|custom$),optional"`
	BackgroundName string `valid:"stringlength(1|50),optional"`
	BackgroundId   string `valid:"objectidhex,optional"`
	Location       string `valid:"stringlength(3|30),optional"`
}

func EditAccount(c echo.Context) (err error) {
	cc := c.(*middlewares.Context)
	params := EditAccountParams{
		Nickname:       cc.FormValue("nickname"),
		Gender:         cc.FormValue("gender"),
		Mobile:         cc.FormValue("mobile"),
		Code:           cc.FormValue("code"),
		AvatarType:     cc.FormValue("avatarType"),
		AvatarName:     cc.FormValue("avatarName"),
		AvatarId:       cc.FormValue("avatarId"),
		Email:          cc.FormValue("email"),
		Intro:          cc.FormValue("intro"),
		BackgroundType: cc.FormValue("backgroundType"),
		BackgroundName: cc.FormValue("backgroundName"),
		BackgroundId:   cc.FormValue("backgroundId"),
		Location:       cc.FormValue("location"),
	}
	if ok, err := valid.ValidateStruct(&params); !ok {
		return services.NewServiceError(services.ErrCodeInvalidParams, err.Error())
	}
	if params.AvatarType == services.AvatarTypeBuiltin && params.AvatarName == "" {
		return services.NewServiceError(services.ErrCodeInvalidParams, "avatar path should not be empty")
	} else if params.AvatarType == services.AvatarTypeCustom && params.AvatarId == "" {
		return services.NewServiceError(services.ErrCodeInvalidParams, "avatar id should not be empty")
	}
	if params.BackgroundType == services.BackgroundTypeBuiltin && params.BackgroundName == "" {
		return services.NewServiceError(services.ErrCodeInvalidParams, "background path should not be empty")
	} else if params.BackgroundType == services.BackgroundTypeCustom && params.BackgroundId == "" {
		return services.NewServiceError(services.ErrCodeInvalidParams, "background id should not be empty")
	}

	codeOk := true
	if params.Mobile != "" {
		codeOk = false
		verifyCode := cc.SessionVerifyCode()
		if verifyCode.By == "mobile" && verifyCode.Mobile == params.Mobile && verifyCode.Code == params.Code {
			codeOk = true
		}
	}
	if !codeOk {
		return services.NewServiceError(services.ErrCodeInvalidVerifyCode, "")
	}

	if params.Email != "" {
		codeOk = false
		verifyCode := cc.SessionVerifyCode()
		if verifyCode.By == "email" && verifyCode.Email == params.Email && verifyCode.Code == params.Code {
			codeOk = true
		}
	}
	if !codeOk {
		return services.NewServiceError(services.ErrCodeInvalidVerifyCode, "")
	}

	formParams, err := cc.FormParams()
	if err != nil {
		return services.NewServiceError(services.ErrCodeInvalidParams, err.Error())
	}
	updateParams := bson.M{}
	if _, ok := formParams["nickname"]; ok {
		updateParams["nickname"] = params.Nickname
	}
	if _, ok := formParams["gender"]; ok {
		updateParams["gender"] = params.Gender
	}
	if _, ok := formParams["mobile"]; ok {
		updateParams["mobile"] = params.Mobile
	}
	if _, ok := formParams["avatarType"]; ok {
		updateParams["avatarType"] = params.AvatarType
	}
	if _, ok := formParams["avatarName"]; ok {
		updateParams["avatarName"] = params.AvatarName
	}
	if _, ok := formParams["avatarId"]; ok {
		avatarId, err := ParseObjectId(params.AvatarId)
		if err != nil {
			return services.NewServiceError(services.ErrCodeInvalidParams, err.Error())
		}
		updateParams["avatarId"] = avatarId
	}
	if _, ok := formParams["email"]; ok {
		updateParams["email"] = params.Email
	}
	if _, ok := formParams["intro"]; ok {
		updateParams["intro"] = params.Intro
	}
	if _, ok := formParams["backgroundType"]; ok {
		updateParams["backgroundType"] = params.BackgroundType
	}
	if _, ok := formParams["backgroundName"]; ok {
		updateParams["backgroundName"] = params.BackgroundName
	}
	if _, ok := formParams["backgroundId"]; ok {
		backgroundId, err := ParseObjectId(params.BackgroundId)
		if err != nil {
			return services.NewServiceError(services.ErrCodeInvalidParams, err.Error())
		}
		updateParams["backgroundId"] = backgroundId
	}
	if _, ok := formParams["location"]; ok {
		loc, err := ParseLocation(params.Location)
		if err != nil {
			return services.NewServiceError(services.ErrCodeInvalidParams, err.Error())
		}
		updateParams["location"] = loc
	}

	id := cc.SessionUserId()
	user, err := services.UpdateUser(id, updateParams)
	if err != nil {
		return err
	}

	cc.DeleteSessionItem("verifyCode")

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
