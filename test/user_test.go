package test

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"

	"zqc/services"
)

func TestCreateUser(t *testing.T) {
	Convey("Given an empty db", t, func() {
		emptyDb()

		Convey("Create a user", func() {
			_, err := services.CreateUser("18600000000", "123456", "jag", "m")
			So(err, ShouldBeNil)
		})
	})
}

func TestEditUser(t *testing.T) {
	Convey("Given an exist user", t, func() {
		emptyDb()
		user := createUser()

		Convey("Update user", func() {
			nickname := "jag1"
			gender := "f"
			user, err := services.UpdateUser(user.Id, map[string]interface{}{
				"nickname": nickname,
				"gender":   gender,
			})
			So(err, ShouldBeNil)
			So(user.Nickname, ShouldEqual, nickname)
			So(user.Gender, ShouldEqual, gender)
		})
	})
}
