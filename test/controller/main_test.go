package controller

import (
	"os"
	"testing"

	"github.com/spf13/viper"

	"jaggerwang.net/zqcserverdemo/models"
)

func TestMain(m *testing.M) {
	viper.Set("mongodb.zqc", map[string]interface{}{
		"addrs": "127.0.0.1:27019",
	})

	models.EmptyDb("zqc", "zqc", "")

	createDbIndexes()

	result := m.Run()

	os.Exit(result)
}
