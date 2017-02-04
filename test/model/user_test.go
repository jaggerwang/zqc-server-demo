package model

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"

	"zqc/models"
)

func TestEmptyUserColl(t *testing.T) {
	Convey("Insert a user to collection", t, func() {
		uc, err := models.NewUserColl()
		So(err, ShouldBeNil)
		err = uc.Insert(models.User{
			Mobile:   "18683420507",
			Password: "123456",
			Nickname: "jag",
			Gender:   "m",
		})
		So(err, ShouldBeNil)

		Convey("Empty collection", func() {
			info, err := uc.RemoveAll(nil)
			So(err, ShouldBeNil)
			So(info.Removed, ShouldEqual, 1)
		})

		Reset(func() {
			models.EmptyDb("zqc", "zqc", "")
		})
	})
}
