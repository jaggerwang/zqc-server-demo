package test

import (
	"os"
	"testing"

	"github.com/spf13/viper"
)

func TestMain(m *testing.M) {
	initConfig()

	result := m.Run()

	os.Exit(result)
}

func initConfig() {
	viper.SetConfigFile("../config.json")
	err := viper.ReadInConfig()
	if err != nil {
		panic(err)
	}
}
