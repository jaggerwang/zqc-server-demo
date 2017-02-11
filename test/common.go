package test

import (
	"zqc/models"
	"zqc/services"
)

func emptyDb() {
	models.EmptyDb("zqc", "zqc", "")
}

func createUser() services.User {
	user, err := services.CreateUser("18600000000", "123456", "jag", "m")
	if err != nil {
		panic(err)
	}

	return user
}
