package controller

import (
	"net/url"
	"testing"

	. "github.com/smartystreets/goconvey/convey"

	"zqc/common"
	"zqc/models"
	"zqc/services"
)

func TestRegisterUser(t *testing.T) {
	Convey("Given some user info", t, func() {
		mobile := "18683420507"
		password := "123456"
		nickname := "jag"
		gender := "m"

		Convey("Register a user", func() {
			f := make(url.Values)
			f.Set("mobile", mobile)
			f.Set("password", password)
			f.Set("nickname", nickname)
			f.Set("gender", gender)
			result := postResult("/account/register", f, "")

			So(result["code"].(float64), ShouldEqual, common.ErrCodeOk)
			user := result["data"].(map[string]interface{})["user"].(map[string]interface{})
			So(user["mobile"].(string), ShouldEqual, mobile)
			So(user["nickname"].(string), ShouldEqual, nickname)
			So(user["gender"].(string), ShouldEqual, gender)
			So(user, ShouldNotContainKey, "password")
			So(user, ShouldNotContainKey, "salt")
		})

		Reset(func() {
			models.EmptyDb("zqc", "zqc", "")
		})
	})
}

func TestEditUser(t *testing.T) {
	Convey("Given a register user and some info edit to", t, func() {
		u := createUser()

		token := login()

		mobile := "01234567890"
		nickname := "jagger"

		Convey("Edit user", func() {
			f := make(url.Values)
			f.Set("mobile", mobile)
			f.Set("nickname", nickname)
			result := postResult("/account/edit", f, token)

			So(result["code"].(float64), ShouldEqual, common.ErrCodeOk)
			user := result["data"].(map[string]interface{})["user"].(map[string]interface{})
			So(user["mobile"].(string), ShouldEqual, mobile)
			So(user["nickname"].(string), ShouldEqual, nickname)
		})

		Reset(func() {
			models.EmptyDb("zqc", "zqc", "")
		})
	})
}

func createUser() map[string]interface{} {
	f := make(url.Values)
	f.Set("mobile", "18683420507")
	f.Set("password", "123456")
	f.Set("nickname", "jag")
	f.Set("gender", "m")
	result := postResult("/account/register", f, "")
	So(result["code"].(float64), ShouldEqual, common.ErrCodeOk)
	return result["data"].(map[string]interface{})["user"].(map[string]interface{})
}

func login() string {
	f := make(url.Values)
	f.Set("mobile", "18683420507")
	f.Set("password", "123456")
	result := postResult("/account/login", f, "")
	So(result["code"].(float64), ShouldEqual, common.ErrCodeOk)
	return result["data"].(map[string]interface{})["token"].(string)
}
