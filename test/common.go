package test

import (
	"zqc/models"
	"zqc/services"
)

func emptyDB() {
	models.EmptyDB("zqc", "zqc", "")
}

func createUser() services.User {
	user, err := services.CreateUser("18600000000", "123456", "jag", "m")
	if err != nil {
		panic(err)
	}

	return user
}
