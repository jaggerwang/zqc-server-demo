package controller

import (
	"net/url"
	"testing"

	. "github.com/smartystreets/goconvey/convey"

	"jaggerwang.net/zqcserverdemo/common"
	"jaggerwang.net/zqcserverdemo/models"
	"jaggerwang.net/zqcserverdemo/services"
)

func TestRegisterUser(t *testing.T) {
	Convey("Given some user info", t, func() {
		username := "jaggerwang"
		password := "198157"
		nickname := "jag"
		gender := "m"
		mobile := "18683420507"

		Convey("Register a user", func() {
			f := make(url.Values)
			f.Set("username", username)
			f.Set("password", password)
			f.Set("nickname", nickname)
			f.Set("gender", gender)
			f.Set("mobile", mobile)
			result := postResult("/account/register", f, "")

			So(result["code"].(float64), ShouldEqual, common.ErrCodeOk)
			user := result["data"].(map[string]interface{})["user"].(map[string]interface{})
			So(user["username"].(string), ShouldEqual, username)
			So(user["nickname"].(string), ShouldEqual, nickname)
			So(user["gender"].(string), ShouldEqual, gender)
			So(user["mobile"].(string), ShouldEqual, mobile)
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

		nickname := "jagger"
		mobile := "01234567890"

		Convey("Edit user", func() {
			f := make(url.Values)
			f.Set("nickname", nickname)
			f.Set("mobile", mobile)
			result := postResult("/account/edit", f, token)

			So(result["code"].(float64), ShouldEqual, common.ErrCodeOk)
			user := result["data"].(map[string]interface{})["user"].(map[string]interface{})
			So(user["nickname"].(string), ShouldEqual, nickname)
			So(user["mobile"].(string), ShouldEqual, mobile)
			So(user["username"].(string), ShouldEqual, u["username"])
		})

		Reset(func() {
			models.EmptyDb("zqc", "zqc", "")
		})
	})
}

func createUser() map[string]interface{} {
	f := make(url.Values)
	f.Set("username", "jaggerwang")
	f.Set("password", "198157")
	f.Set("nickname", "jag")
	f.Set("gender", "m")
	f.Set("mobile", "18683420507")
	result := postResult("/account/register", f, "")
	So(result["code"].(float64), ShouldEqual, common.ErrCodeOk)
	return result["data"].(map[string]interface{})["user"].(map[string]interface{})
}

func login() string {
	f := make(url.Values)
	f.Set("username", "jaggerwang")
	f.Set("password", "198157")
	result := postResult("/account/login", f, "")
	So(result["code"].(float64), ShouldEqual, common.ErrCodeOk)
	return result["data"].(map[string]interface{})["token"].(string)
}
